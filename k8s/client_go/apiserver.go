apiserver:
  k8s.io/kubernetes/pkg/master/master.go 

completedConfig.New(genericapiserver.DelegationTarget)
  -> genericAPIServer := genericConfig.New()
  -> m.InstallAPIs(RESTStorageProvider)
  -> genericapiserver.InstallAPIGroup(APIGroupInfo)


registry:
  k8s.io/kubernetes/pkg/registry
  implements the storage and system logic for the core 
of api server.
