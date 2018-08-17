package httpserver

//how http server is started and which handler to run
//k8s.io/apiserver/pkg/server/genericapiserver.go

type GenericAPIServer struct {
	//build cluster IPs for discovery.
	discoveryAddresses discovery.Addresses

	// LoopbackClientConfig is a config for a privileged loopback connection to the API server
	LoopbackClientConfig *restclient.Config

	// minRequestTimeout is how short the request timeout can be.  This is used to build the RESTHandler
	minRequestTimeout time.Duration

	// ShutdownTimeout is the timeout used for server shutdown. This specifies the timeout before server
	// gracefully shutdown returns.
	ShutdownTimeout time.Duration

	// legacyAPIGroupPrefixes is used to set up URL parsing for authorization and for validating requests
	// to InstallLegacyAPIGroup
	legacyAPIGroupPrefixes sets.String

	// admissionControl is used to build the RESTStorage that backs an API Group.
	admissionControl admission.Interface

	// requestContextMapper provides a way to get the context for a request.  It may be nil.
	requestContextMapper apirequest.RequestContextMapper

	SecureServingInfo *SecureServingInfo

	// ExternalAddress is the address (hostname or IP and port) that should be used in
	// external (public internet) URLs for this GenericAPIServer.
	ExternalAddress string

	// Serializer controls how common API objects not in a group/version prefix are serialized for this server.
	// Individual APIGroups may define their own serializers.
	Serializer runtime.NegotiatedSerializer

	// "Outputs"
	// Handler holds the handlers being used by this API server
	Handler *APIServerHandler

	// listedPathProvider is a lister which provides the set of paths to show at /
	listedPathProvider routes.ListedPathProvider

	// DiscoveryGroupManager serves /apis
	DiscoveryGroupManager discovery.GroupManager

	// Enable swagger and/or OpenAPI if these configs are non-nil.
	swaggerConfig *swagger.Config
	openAPIConfig *openapicommon.Config

	// PostStartHooks are each called after the server has started listening, in a separate go func for each
	// with no guarantee of ordering between them.  The map key is a name used for error reporting.
	// It may kill the process with a panic if it wishes to by returning an error.
	postStartHookLock      sync.Mutex
	postStartHooks         map[string]postStartHookEntry
	postStartHooksCalled   bool
	disabledPostStartHooks sets.String

	preShutdownHookLock    sync.Mutex
	preShutdownHooks       map[string]preShutdownHookEntry
	preShutdownHooksCalled bool

	// healthz checks
	healthzLock    sync.Mutex
	healthzChecks  []healthz.HealthzChecker
	healthzCreated bool

	// auditing. The backend is started after the server starts listening.
	AuditBackend audit.Backend

	// enableAPIResponseCompression indicates whether API Responses should support compression
	// if the client requests it via Accept-Encoding
	enableAPIResponseCompression bool

	// delegationTarget is the next delegate in the chain or nil
	delegationTarget DelegationTarget

	// HandlerChainWaitGroup allows you to wait for all chain handlers finish after the server shutdown.
	HandlerChainWaitGroup *utilwaitgroup.SafeWaitGroup
}

// FullHandlerChain -> Director -> {GoRestfulContainer,NonGoRestfulMux} based on inspection of registered web services
type APIServerHandler struct {
	FullHandlerChain   http.Handler
	GoRestfulContainer *restful.Container   //web service container
	NonGoRestfulMux    *mux.PathRecorderMux //path -> handler map, doesn't support dynamic url
	Director           http.Handler
}

type director struct {
	name               string
	goRestfulContainer *restful.Container
	nonGoRestfulMux    *mux.PathRecorderMux
}

//Director.goRestfulContainer == APIServerHandler.goRestfulContainer
//Director.nonGoRestfulMux == APIServerHandler.nonGoRestfulMux
func (d director) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	d.goRestfulContainer.Dispatch(w, req)
	d.nonGoRestfulMux.ServeHTTP(w, req)
}

//a.FullHandlerChain = handlerChainBuilder(director)
func (a *APIServerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.FullHandlerChain.ServeHTTP(w, r)
}

//k8s.io/apiserver/pkg/server/config.go
//handlerChainBuilder == DefaultBuildHandlerChain
func DefaultBuildHandlerChain(apiHandler http.Handler, c *Config) http.Handler {
	handler := genericapifilters.WithAuthorization(apiHandler, c.RequestContextMapper, c.Authorization.Authorizer, c.Serializer)
	handler = genericfilters.WithMaxInFlightLimit(handler, c.MaxRequestsInFlight, c.MaxMutatingRequestsInFlight, c.RequestContextMapper, c.LongRunningFunc)
	handler = genericapifilters.WithImpersonation(handler, c.RequestContextMapper, c.Authorization.Authorizer, c.Serializer)
	if utilfeature.DefaultFeatureGate.Enabled(features.AdvancedAuditing) {
		handler = genericapifilters.WithAudit(handler, c.RequestContextMapper, c.AuditBackend, c.AuditPolicyChecker, c.LongRunningFunc)
	} else {
		handler = genericapifilters.WithLegacyAudit(handler, c.RequestContextMapper, c.LegacyAuditWriter)
	}
	failedHandler := genericapifilters.Unauthorized(c.RequestContextMapper, c.Serializer, c.Authentication.SupportsBasicAuth)
	if utilfeature.DefaultFeatureGate.Enabled(features.AdvancedAuditing) {
		failedHandler = genericapifilters.WithFailedAuthenticationAudit(failedHandler, c.RequestContextMapper, c.AuditBackend, c.AuditPolicyChecker)
	}
	handler = genericapifilters.WithAuthentication(handler, c.RequestContextMapper, c.Authentication.Authenticator, failedHandler)
	handler = genericfilters.WithCORS(handler, c.CorsAllowedOriginList, nil, nil, nil, "true")
	handler = genericfilters.WithTimeoutForNonLongRunningRequests(handler, c.RequestContextMapper, c.LongRunningFunc, c.RequestTimeout)
	handler = genericfilters.WithWaitGroup(handler, c.RequestContextMapper, c.LongRunningFunc, c.HandlerChainWaitGroup)
	handler = genericapifilters.WithRequestInfo(handler, c.RequestInfoResolver, c.RequestContextMapper)
	handler = apirequest.WithRequestContext(handler, c.RequestContextMapper)
	handler = genericfilters.WithPanicRecovery(handler)
	return handler
}

//k8s.io/apiserver/pkg/server/genericapiserver.go
func (s *GenericAPIServer) PrepareRun() preparedGenericAPIServer
func (s preparedGenericAPIServer) Run(stopCh <-chan struct{}) error
func (s preparedGenericAPIServer) NonBlockingRun(stopCh <-chan struct{}) error {
	s.SecureServingInfo.Serve(s.Handler, s.ShutdownTimeout, internalStopCh)
}

//k8s.io/apiserver/pkg/server/serve.go
func (s *SecureServingInfo) Serve(handler http.Handler, shutdownTimeout time.Duration, stopCh <-chan struct{}) error {
	secureServer := &http.Server{
		Addr:           s.Listener.Addr().String(),
		Handler:        handler,
		MaxHeaderBytes: 1 << 20,
		TLSConfig: &tls.Config{
			NameToCertificate: s.SNICerts,
			MinVersion:        tls.VersionTLS12,
			NextProtos:        []string{"h2", "http/1.1"},
		},
	}
	return RunServer(secureServer, s.Listener, shutdownTimeout, stopCh)
}

func RunServer(server *http.Server, ln net.Listener, shutDownTimeout time.Duration, stopCh <-chan struct{}) error {
	err := server.Serve(listener)
}
