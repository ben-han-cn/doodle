package controller

//Deploy ---> replicaset ---> pod

type ReplicaSetController struct {
	kubeClient    clientset.Interface
	podController controller.PodControlInterface
	rsLister      appslisters.ReplicaSetLister
	podLister     corelisters.PodLister
	queue         workqueue.RateLimitingInterface
}
