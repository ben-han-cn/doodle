package storage

//how api group is registered
//k8s.io/apiserver/pkg/server/genericapiserver.go
type APIGroupInfo struct {
	GroupMeta                    apimachinery.GroupMeta
	VersionedResourcesStorageMap map[string]map[string]rest.Storage
	//....
}

//apiGroupInfo.VersionedResourcesStorageMap ---> rest.Storage --> rest.Creater --> Create
func (s *GenericAPIServer) InstallAPIGroup(apiGroupInfo *APIGroupInfo) error {
	s.installAPIResources(APIGroupPrefix, apiGroupInfo)
	apiGroup := metav1.APIGroup{
		Name:             apiGroupInfo.GroupMeta.GroupVersion.Group,
		Versions:         apiVersionsForDiscovery,
		PreferredVersion: preferredVersionForDiscovery,
	}

	s.DiscoveryGroupManager.AddGroup(apiGroup)
	s.Handler.GoRestfulContainer.Add(discovery.NewAPIGroupHandler(s.Serializer, apiGroup, s.requestContextMapper).WebService())
}

func (s *GenericAPIServer) installAPIResources(apiPrefix string, apiGroupInfo *APIGroupInfo) error {
	for _, groupVersion := range apiGroupInfo.GroupMeta.GroupVersions {
		apiGroupVersion.InstallREST(s.Handler.GoRestfulContainer)
	}
}

func (g *APIGroupVersion) InstallREST(container *restful.Container) error {
	prefix := path.Join(g.Root, g.GroupVersion.Group, g.GroupVersion.Version)
	installer := &APIInstaller{
		group:                        g,
		prefix:                       prefix,
		minRequestTimeout:            g.MinRequestTimeout,
		enableAPIResponseCompression: g.EnableAPIResponseCompression,
	}

	apiResources, ws, registrationErrors := installer.Install()
	versionDiscoveryHandler := discovery.NewAPIVersionHandler(g.Serializer, g.GroupVersion, staticLister{apiResources}, g.Context)
	versionDiscoveryHandler.AddToWebService(ws)
	container.Add(ws)
	return utilerrors.NewAggregate(registrationErrors)
}

//k8s.io/apiserver/pkg/endpoints/installer.go
//restfulCreateResource
//restfulListResource
type APIInstaller struct {
	group                        *APIGroupVersion
	prefix                       string // Path prefix where API resources are to be registered.
	minRequestTimeout            time.Duration
	enableAPIResponseCompression bool
}

func (a *APIInstaller) Install() ([]metav1.APIResource, *restful.WebService, []error) {
	ws := a.newWebService()
	paths := make([]string, len(a.group.Storage))
	for _, path := range paths {
		apiResource, err := a.registerResourceHandlers(path, a.group.Storage[path], ws)
	}
}

//registerResourceHandlers ---> very big function to map path to handler
func restfulCreateResource(r rest.Creater, scope handlers.RequestScope, typer runtime.ObjectTyper, admit admission.Interface) restful.RouteFunction {
	return func(req *restful.Request, res *restful.Response) {
		handlers.CreateResource(r, scope, typer, admit)(res.ResponseWriter, req.Request)
	}
}

//k8s.io/apiserver/pkg/endpoints/handlers/create.go
//CreateResource --> createHandler
func createHandler(r rest.NamedCreater, scope RequestScope, typer runtime.ObjectTyper, admit admission.Interface, includeName bool) http.HandlerFunc {
	ctx := scope.ContextFunc(req)
	ctx = request.WithNamespace(ctx, namespace)

	gv := scope.Kind.GroupVersion()
	s, err := negotiation.NegotiateInputSerializer(req, false, scope.Serializer)
	decoder := scope.Serializer.DecoderToVersion(s.Serializer, schema.GroupVersion{Group: gv.Group, Version: runtime.APIVersionInternal})
	body, err := readBody(req)
	defaultGVK := scope.Kind
	original := r.New()
	obj, gvk, err := decoder.Decode(body, &defaultGVK, original)
	ae := request.AuditEventFrom(ctx)
	userInfo, _ := request.UserFrom(ctx)
	result, err := finishRequest(timeout, func() (runtime.Object, error) {
		return r.Create(
			ctx,
			name,
			obj,
			rest.AdmissionToValidateObjectFunc(admit, admissionAttributes),
			includeUninitialized,
		)
	})
	transformResponseObject(ctx, scope, req, w, code, result)
}

//k8s.io/apiserver/pkg/registry/rest/rest.go
// NamedCreater is an object that can create an instance of a RESTful object using a name parameter.
type NamedCreater interface {
	New() runtime.Object
	Create(ctx genericapirequest.Context, name string, obj runtime.Object, createValidation ValidateObjectFunc, includeUninitialized bool) (runtime.Object, error)
}

//kubernetes/pkg/master/master.go
func (m *Master) InstallAPIs(apiResourceConfigSource serverstorage.APIResourceConfigSource, restOptionsGetter generic.RESTOptionsGetter, restStorageProviders ...RESTStorageProvider) {
	apiGroupsInfo := []genericapiserver.APIGroupInfo{}

	for _, restStorageBuilder := range restStorageProviders {
		groupName := restStorageBuilder.GroupName()
		apiGroupInfo, _ := restStorageBuilder.NewRESTStorage(apiResourceConfigSource, restOptionsGetter)
		if postHookProvider, ok := restStorageBuilder.(genericapiserver.PostStartHookProvider); ok {
			name, hook, err := postHookProvider.PostStartHook()
			m.GenericAPIServer.AddPostStartHookOrDie(name, hook)
		}
		apiGroupsInfo = append(apiGroupsInfo, apiGroupInfo)
	}

	for i := range apiGroupsInfo {
		m.GenericAPIServer.InstallAPIGroup(&apiGroupsInfo[i])
	}
}

//NewRESTStorage ----> NewDefaultAPIGroupInfo [k8s.io/apiserver/pkg/server/genericapiserver.go]
func NewDefaultAPIGroupInfo(group string, registry *registered.APIRegistrationManager, scheme *runtime.Scheme, parameterCodec runtime.ParameterCodec, codecs serializer.CodecFactory) APIGroupInfo {
	groupMeta := registry.GroupOrDie(group)

	return APIGroupInfo{
		GroupMeta:                    *groupMeta,
		VersionedResourcesStorageMap: map[string]map[string]rest.Storage{},
		// TODO unhardcode this.  It was hardcoded before, but we need to re-evaluate
		OptionsExternalVersion: &schema.GroupVersion{Version: "v1"},
		Scheme:                 scheme,
		ParameterCodec:         parameterCodec,
		NegotiatedSerializer:   codecs,
	}
}

//k8s.io/apiserver/pkg/registry/rest/rest.go
type Storage interface {
	// New returns an empty object that can be used with Create and Update after request data has been put into it.
	// This object must be a pointer type for use with Codec.DecodeInto([]byte, runtime.Object)
	New() runtime.Object
}

//kubernetes/pkg/registry/storage/rest/storage_storage.go
func (p RESTStorageProvider) NewRESTStorage(apiResourceConfigSource serverstorage.APIResourceConfigSource, restOptionsGetter generic.RESTOptionsGetter) (genericapiserver.APIGroupInfo, bool) {
	apiGroupInfo.VersionedResourcesStorageMap[storageapiv1.SchemeGroupVersion.Version] = p.v1Storage(apiResourceConfigSource, restOptionsGetter)
}

func (p RESTStorageProvider) v1Storage(apiResourceConfigSource serverstorage.APIResourceConfigSource, restOptionsGetter generic.RESTOptionsGetter) map[string]rest.Storage {
	storage := map[string]rest.Storage{}
	// storageclasses
	storageClassStorage := storageclassstore.NewREST(restOptionsGetter)
	storage["storageclasses"] = storageClassStorage
	return storage
}

//kubernetes/pkg/registry/storage/storageclass/storage/storage.go
// NewREST returns a RESTStorage object that will work against persistent volumes.
func NewREST(optsGetter generic.RESTOptionsGetter) *REST {
	store := &genericregistry.Store{
		NewFunc:                  func() runtime.Object { return &storageapi.StorageClass{} },
		NewListFunc:              func() runtime.Object { return &storageapi.StorageClassList{} },
		DefaultQualifiedResource: storageapi.Resource("storageclasses"),

		CreateStrategy:      storageclass.Strategy,
		UpdateStrategy:      storageclass.Strategy,
		DeleteStrategy:      storageclass.Strategy,
		ReturnDeletedObject: true,

		TableConvertor: printerstorage.TableConvertor{TablePrinter: printers.NewTablePrinter().With(printersinternal.AddHandlers)},
	}
	options := &generic.StoreOptions{RESTOptions: optsGetter}
	if err := store.CompleteWithOptions(options); err != nil {
		panic(err) // TODO: Propagate error up
	}

	return &REST{store}
}

//k8s.io/apiserver/pkg/registry/generic/registry/store.go
type Store struct {
	Storage storage.Interface //k8s.io/apiserver/pkg/storage
}

func (e *Store) CompleteWithOptions(options *generic.StoreOptions) error {
	opts, err := options.RESTOptions.GetRESTOptions(e.DefaultQualifiedResource)
	e.Storage, e.DestroyFunc = opts.Decorator(
		opts.StorageConfig,
		e.NewFunc(),
		prefix,
		keyFunc,
		e.NewListFunc,
		attrFunc,
		triggerFunc,
	)
}

//kubernetes/cmd/kube-apiserver/app/aggregator.go
func createAggregatorConfig() {
	genericConfig.RESTOptionsGetter = &genericoptions.SimpleRestOptionsFactory{Options: etcdOptions}
}

//k8s.io/apiserver/pkg/server/options/etcd.go
func (f *SimpleRestOptionsFactory) GetRESTOptions(resource schema.GroupResource) (generic.RESTOptions, error) {
	ret := generic.RESTOptions{
		StorageConfig: &f.Options.StorageConfig,
		Decorator:     generic.UndecoratedStorage,
	}
	return ret, nil
}

//k8s.io/apiserver/pkg/registry/generic/storage_decorator.go
func UndecoratedStorage() (storage.Interface, factory.DestroyFunc) {
	return NewRawStorage(config)
}

// NewRawStorage creates the low level kv storage. This is a work-around for current
// two layer of same storage interface.
// TODO: Once cacher is enabled on all registries (event registry is special), we will remove this method.
func NewRawStorage(config *storagebackend.Config) (storage.Interface, factory.DestroyFunc) {
	s, d, err := factory.Create(*config)
	if err != nil {
		glog.Fatalf("Unable to create storage backend: config (%v), err (%v)", config, err)
	}
	return s, d
}

//k8s.io/apiserver/pkg/storage/storagebackend/factory/factory.go
func Create(c storagebackend.Config) (storage.Interface, DestroyFunc, error) {
	return newETCD3Storage(c)
}
func newETCD3Storage(c storagebackend.Config) (storage.Interface, DestroyFunc, error) {
	return etcd3.New(client, c.Codec, c.Prefix, transformer, c.Paging), destroyFunc, nil
}

//k8s.io/apiserver/pkg/storage/etcd3/store.go
func newStore(c *clientv3.Client, quorumRead, pagingEnabled bool, codec runtime.Codec, prefix string, transformer value.Transformer) *store {
	versioner := etcd.APIObjectVersioner{}
	result := &store{
		client:        c,
		codec:         codec,
		versioner:     versioner,
		transformer:   transformer,
		pagingEnabled: pagingEnabled,
		// for compatibility with etcd2 impl.
		// no-op for default prefix of '/registry'.
		// keeps compatibility with etcd2 impl for custom prefixes that don't start with '/'
		pathPrefix: path.Join("/", prefix),
		watcher:    newWatcher(c, codec, versioner, transformer),
	}
	if !quorumRead {
		// In case of non-quorum reads, we can set WithSerializable()
		// options for all Get operations.
		result.getOps = append(result.getOps, clientv3.WithSerializable())
	}
	return result
}
