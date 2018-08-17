package sharedinformer

//k8s.io/client-go/tools/cache/shared_informer.go
type SharedInformer interface {
	AddEventHandler(handler ResourceEventHandler)
	AddEventHandlerWithResyncPeriod(handler ResourceEventHandler, resyncPeriod time.Duration)
	GetStore() Store
	GetController() Controller
	Run(stopCh <-chan struct{})
	HasSynced() bool
	LastSyncResourceVersion() string
}
