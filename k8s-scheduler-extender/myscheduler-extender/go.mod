module banzaicloud.com/myscheduler-extender

go 1.13

require (
	github.com/julienschmidt/httprouter v1.2.0
	k8s.io/api v0.0.0
	k8s.io/kubernetes v1.17.3
)

replace (
	k8s.io/api => k8s.io/api v0.0.0-20200131193051-d9adff57e763
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.0.0-20200131201446-6910daba737d
	k8s.io/apimachinery => k8s.io/apimachinery v0.17.4-beta.0
	k8s.io/apiserver => k8s.io/apiserver v0.0.0-20200131195721-b64b0ef70370
	k8s.io/cli-runtime => k8s.io/cli-runtime v0.0.0-20200131202043-1dc23f43cc94
	k8s.io/client-go => k8s.io/client-go v0.0.0-20200210225353-0ff5a65499e6
	k8s.io/cloud-provider => k8s.io/cloud-provider v0.0.0-20200131203830-fe5589c708de
	k8s.io/cluster-bootstrap => k8s.io/cluster-bootstrap v0.0.0-20200131203557-3c6746d7c617
	k8s.io/code-generator => k8s.io/code-generator v0.17.4-beta.0
	k8s.io/component-base => k8s.io/component-base v0.0.0-20200131194811-85b325a9731b
	k8s.io/cri-api => k8s.io/cri-api v0.17.4-beta.0
	k8s.io/csi-api => k8s.io/csi-api v0.0.0-20190313123203-94ac839bf26c
	k8s.io/csi-translation-lib => k8s.io/csi-translation-lib v0.0.0-20200131204100-4311b557c8ce
	k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.0.0-20200131200134-d62c64b672cc
	k8s.io/kube-controller-manager => k8s.io/kube-controller-manager v0.0.0-20200131203333-c935c9222556
	k8s.io/kube-proxy => k8s.io/kube-proxy v0.0.0-20200131202556-6b094e7591d1
	k8s.io/kube-scheduler => k8s.io/kube-scheduler v0.0.0-20200131203102-8e9ee8fa0785
	k8s.io/kubectl => k8s.io/kubectl v0.0.0-20200131205129-9ef1401eb3ec
	k8s.io/kubelet => k8s.io/kubelet v0.0.0-20200131202828-eb1b5c1ce7fb
	k8s.io/legacy-cloud-providers => k8s.io/legacy-cloud-providers v0.0.0-20200131204342-ef4bac7ed518
	k8s.io/metrics => k8s.io/metrics v0.0.0-20200131201757-ffbb7a48f604
	k8s.io/sample-apiserver => k8s.io/sample-apiserver v0.0.0-20200131200511-51b2302b2589
	k8s.io/sample-cli-plugin => k8s.io/sample-cli-plugin v0.0.0-20200131202323-14126e90c844
	k8s.io/sample-controller => k8s.io/sample-controller v0.0.0-20200131200932-3fd12213be16
)
