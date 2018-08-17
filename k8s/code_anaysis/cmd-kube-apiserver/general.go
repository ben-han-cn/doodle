package apiserver

func Run() {
	CreateKubeAPIServerConfig()
	createAPIExtensionsConfig()
	apiExtensionsServer := createAPIExtensionsServer()
	kubeAPIServer := CreateKubeAPIServer()

	kubeAPIServer.GenericAPIServer.PrepareRun()
	apiExtensionsServer.GenericAPIServer.PrepareRun()
	createAggregatorConfig()
	aggregatorServer := createAggregatorServer()
	aggregatorServer.GenericAPIServer.PrepareRun().Run(stopCh)
}
