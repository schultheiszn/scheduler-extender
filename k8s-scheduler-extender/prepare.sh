#kubectl edit clusterrole cluster-autoscaler
kubectl edit clusterrole system:kube-scheduler
kubectl create clusterrolebinding myscheduler --clusterrole cluster-admin --serviceaccount=kube-system:myscheduler
