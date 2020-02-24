package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"math/rand"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	extenderv1 "k8s.io/kubernetes/pkg/scheduler/apis/extender/v1"
)

const (
	Percent = 0.3
	RSKind = "ReplicaSet"
)

type GetPodsByNodeNameFunc func(ns string, nodeName string) ([]v1.Pod, error)
type PerNodeControl struct  {
	threshold int
	current int
	skiprest bool
}

func Filter(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var buff bytes.Buffer
	var extenderArgs extenderv1.ExtenderArgs
	var extenderFilterResult *extenderv1.ExtenderFilterResult

	body := io.TeeReader(r.Body, &buff)

	if err := json.NewDecoder(body).Decode(&extenderArgs); err != nil {
		extenderFilterResult = &extenderv1.ExtenderFilterResult{Error: err.Error()}
	} else {
		extenderFilterResult = filter(extenderArgs)
	}

	log.Printf("%+v", extenderFilterResult)

	if response, err := json.Marshal(extenderFilterResult); err != nil {
		log.Fatalln(err)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(response)
	}
}

func Prioritize(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var buf bytes.Buffer
	var extenderArgs extenderv1.ExtenderArgs
	var hostPriorityList *extenderv1.HostPriorityList

	body := io.TeeReader(r.Body, &buf)

	if err := json.NewDecoder(body).Decode(&extenderArgs); err != nil {
		log.Println(err)
		hostPriorityList = &extenderv1.HostPriorityList{}
	} else {
		hostPriorityList = prioritize(extenderArgs)
	}

	if response, err := json.Marshal(hostPriorityList); err != nil {
		log.Fatalln(err)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(response)
	}
}

func prioritize(args extenderv1.ExtenderArgs) *extenderv1.HostPriorityList {
	nodes := args.Nodes.Items

	hostPriorityList := make(extenderv1.HostPriorityList, len(nodes))
	for i, node := range nodes {
		score := rand.Int63n(extenderv1.MaxExtenderPriority)
		hostPriorityList[i] = extenderv1.HostPriority{
			Host:  node.Name,
			Score: score,
		}
	}

	return &hostPriorityList
}

func isReplicaSet(ref *metav1.OwnerReference) bool {
	if ref != nil && ref.Kind == RSKind {
		return true
	}
	return false
}

func getReplicaSetOwnerRef(refs *[]metav1.OwnerReference) *metav1.OwnerReference {
	for _, owner := range *refs {
		if isReplicaSet(&owner) {
			return &owner
		}
	}
	return nil
}

func filter(args extenderv1.ExtenderArgs) *extenderv1.ExtenderFilterResult {
	pod := args.Pod
	var filteredNodes []v1.Node
	failedNodes := make(extenderv1.FailedNodesMap)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	podInfo, err := CS.CoreV1().Pods(pod.Namespace).Get(ctx, pod.Name, metav1.GetOptions{})
	if err != nil {
		panic(err.Error())
	}

	ownerReference := getReplicaSetOwnerRef(&podInfo.OwnerReferences)
	if ownerReference != nil {
		limiter, ok := M[ownerReference.Name]
		if !ok {
			rs, err := CS.AppsV1().ReplicaSets(pod.Namespace).Get(ctx, ownerReference.Name, metav1.GetOptions{})
			if err != nil {
				panic(err.Error())
			}
			var m int
			switch rs.Spec.Replicas {
				case nil: m = 1
				default: m  = int(*rs.Spec.Replicas)
			}
			M[ownerReference.Name] = PerNodeControl{int(math.Round(float64(m) * Percent)), 0, false}
		}

		for _, node := range args.Nodes.Items {
			if limiter.skiprest {
				filteredNodes = append(filteredNodes, node)
			} else {
				podList, err := getPodsAssignedToNode()(node.Namespace, node.Name)
				if err != nil {
					failedNodes[node.Name] = err.Error()
				} else {
					for _, podOnNode := range podList {
						if strings.HasPrefix(podOnNode.Name, ownerReference.Name) {
							limiter.current++
						}
					}
					if limiter.current >= limiter.threshold {
						limiter.skiprest = true
						failedNodes[node.Name] = fmt.Sprintf("Node[%s] could not accept more pod from replica set[%s]", node.Name, ownerReference.Name)
					} else {
						filteredNodes = append(filteredNodes, node)
					}
				}
			}
		}
	} else {
		filteredNodes = args.Nodes.Items
	}

	result := extenderv1.ExtenderFilterResult{
		Nodes: &v1.NodeList{
			Items: filteredNodes,
		},
		FailedNodes: failedNodes,
		Error:       "",
	}
	log.Printf("%+v", M)

	return &result
}

func getPodsAssignedToNode() GetPodsByNodeNameFunc {
	return func(ns string, nodeName string) ([]v1.Pod, error) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		selector := fields.SelectorFromSet(fields.Set{"spec.nodeName": nodeName})
		pods, err := CS.CoreV1().Pods(ns).List(ctx, metav1.ListOptions{
			FieldSelector: selector.String(),
		})
		if err != nil {
			return []v1.Pod{}, fmt.Errorf("failed to get Pods assigned to node %v", nodeName)
		}
		return pods.Items, nil
	}
}

var CS *kubernetes.Clientset
var M map[string]PerNodeControl

func init() {
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}

	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	CS = clientSet
	M = make(map[string]PerNodeControl)
}

func main() {
	log.Print("Scheduler extender init...")
	router := httprouter.New()
	router.POST("/filter", Filter)
	router.POST("/prioritize", Prioritize)
	log.Fatal(http.ListenAndServe(":8888", router))
}
