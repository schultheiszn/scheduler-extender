package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"

	"github.com/julienschmidt/httprouter"

	v1 "k8s.io/api/core/v1"
	schedulerapi "k8s.io/kubernetes/pkg/scheduler/apis/extender/v1"
)

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome to sample-scheduler-extender!\n")
	log.Print("myscheduler-extender Index")
}

func Filter(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	log.Print("myscheduler-extender Filter")
	var buf bytes.Buffer
	body := io.TeeReader(r.Body, &buf)
	var extenderArgs schedulerapi.ExtenderArgs
	var extenderFilterResult *schedulerapi.ExtenderFilterResult
	if err := json.NewDecoder(body).Decode(&extenderArgs); err != nil {
		extenderFilterResult = &schedulerapi.ExtenderFilterResult{
			Error: err.Error(),
		}
	} else {
		extenderFilterResult = filter(extenderArgs)
	}

	if response, err := json.Marshal(extenderFilterResult); err != nil {
		log.Fatalln(err)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(response)
	}
}

func Prioritize(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	log.Print("myscheduler-extender Prioritize")
	var buf bytes.Buffer
	body := io.TeeReader(r.Body, &buf)
	var extenderArgs schedulerapi.ExtenderArgs
	var hostPriorityList *schedulerapi.HostPriorityList
	if err := json.NewDecoder(body).Decode(&extenderArgs); err != nil {
		log.Println(err)
		hostPriorityList = &schedulerapi.HostPriorityList{}
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

const (
	// lucky priority gives a random [0, schedulerapi.MaxPriority] score
	// currently schedulerapi.MaxPriority is 10
	luckyPrioMsg = "pod %v/%v is lucky to get score %v\n"
)

func prioritize(args schedulerapi.ExtenderArgs) *schedulerapi.HostPriorityList {
	pod := args.Pod
	nodes := args.Nodes.Items

	hostPriorityList := make(schedulerapi.HostPriorityList, len(nodes))
	for i, node := range nodes {
		// score := rand.Intn(math.MaxInt32) //schedulerapi.MaxExtenderPriority)
		score := rand.Int63n(schedulerapi.MaxExtenderPriority)
		log.Printf(luckyPrioMsg, pod.Name, pod.Namespace, score)
		hostPriorityList[i] = schedulerapi.HostPriority{
			Host:  node.Name,
			Score: score,
		}
	}

	return &hostPriorityList
}

func filter(args schedulerapi.ExtenderArgs) *schedulerapi.ExtenderFilterResult {
	// var filteredNodes []v1.Node
	failedNodes := make(schedulerapi.FailedNodesMap)
	// pod := args.Pod

	// for _, node := range args.Nodes.Items {
	//fits, failReasons, _ := true, nil, nil //podFitsOnNode(pod, node)
	//if fits {
	// filteredNodes = append(filteredNodes, node)
	//} else {
	//	failedNodes[node.Name] = strings.Join(failReasons, ",")
	//}
	// }
	result := schedulerapi.ExtenderFilterResult{
		Nodes: &v1.NodeList{
			Items: args.Nodes.Items,
		},
		FailedNodes: failedNodes,
		Error:       "",
	}

	return &result
}

func main() {
	log.Print("myscheduler-extender init")
	router := httprouter.New()
	router.GET("/", Index)
	router.POST("/filter", Filter)
	router.POST("/prioritize", Prioritize)

	log.Fatal(http.ListenAndServe("0.0.0.0:8888", router))
}
