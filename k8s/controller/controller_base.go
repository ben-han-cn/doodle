package controller

import (
	appsinformers "k8s.io/client-go/informers/apps/v1"
	coreinformers "k8s.io/client-go/informers/core/v1"

	appslisters "k8s.io/client-go/listers/apps/v1"
	corelisters "k8s.io/client-go/listers/core/v1"

	"k8s.io/client-go/tools/record"
)

type GeneralController struct {
	//resource manipulator
	//actually all api resource operation is through kubeclient
	//PodControlInterface and RSControlInterface is based on kubeclient
	podControl controller.PodControlInterface
	rsContol   controller.RSControlInterface
	kubeclient clientset.Interface

	//lister for watch event
	rsLister  appslisters.ReplicaSetLister
	podLister corelisters.PodLister

	//recorder for broadcast event
	recorder record.EventRecorder

	//queue
	queue workqueue.RateLimitingInterface
}

func newGeneralController(rsInformer appsinformers.ReplicaSetInformer,
	podInformer coreinformers.PodInformer,
	kubeClient clientset.Interface) *GeneralController {

	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartLogging(glog.Infof)
	eventBroadcaster.StartRecordingToSink(&v1core.EventSinkImpl{Interface: kubeClient.CoreV1().Events("")})
	recorder := eventBroadcaster.NewRecorder(scheme.Scheme, v1.EventSource{Component: "controller_name"})

	rsInformer.Informer().AddEventHandler()
	rsLister = rsInformer.Lister()
	podInformer.Informer().AddEventHandler()
	podLister = podInformer.Lister()
}

func (c *GeneralController) Run(workers int, stopCh <-chan struct{}) {
	defer c.queue.ShutDown()
	//wait for event sync
	if !controller.WaitForCacheSync(rsc.Kind, stopCh, rsc.podListerSynced, rsc.rsListerSynced) {
		return
	}
	for i := 0; i < workers; i++ {
		go wait.Until(rsc.worker, time.Second, stopCh)
	}
	<-stopCh
}

func (c *GeneralController) worker() {
	for rsc.processNextWorkItem() {
	}
}

func (c *GeneralController) processNextWorkItem() bool {
	key, quit := rsc.queue.Get()
	if quit {
		return false
	}
	defer c.queue.Done(key)

	err := c.syncHandler(key.(string))
	if err == nil {
		rsc.queue.Forget(key)
		return true
	}

	utilruntime.HandleError(fmt.Errorf("Sync %q failed with %v", key, err))
	rsc.queue.AddRateLimited(key)
	return true
}

func (rsc *ReplicaSetController) syncHandler(key string) error {
	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	rs, err := rsLister.ReplicaSets(namespace).Get(name)
	//....
}
