package dd

//meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
//types.go
//store all the meta info about the abstract resource/object
type TypeMeta struct {
	Kind       string
	APIVersion string
}

type ObjectMeta struct {
	Name                       string
	GenerateName               string
	Namespace                  string
	SelfLink                   string
	UID                        types.UID
	ResourceVersion            string
	Generation                 int64
	CreationTimestamp          Time
	DeletionTimestamp          *Time
	DeletionGracePeriodSeconds *int64
	Labels                     map[string]string
	Annotations                map[string]string
	OwnerReferences            []OwnerReference
	Initializers               *Initializers
	Finalizers                 []string
	ClusterName                string
}

type OwnerReferences struct {
	APIVersion string
	Kind       string
	Name       string
	UID        types.UID
	// If true, this reference points to the managing controller.
	Controller *bool
	// If true, AND if the owner has the "foregroundDeletion" finalizer, then
	// the owner cannot be deleted from the key-value store until this
	// reference is removed.
	BlockOwnerDeletion *bool
}

//"k8s.io/api/core/v1"
//store all the basic struct/definition of basic resources
//types.go
type Pod struct {
	metav1.TypeMeta
	metav1.ObjectMeta
	Spec   PodSpec
	Status PodStatus
}

//k8s.io/apimachinery/pkg/runtime/schema
//protobuf generate mod, which is mainly used to serialization
type GroupVersionKind struct {
	Group   string
	Version string
	Kind    string
}

//v1core "k8s.io/client-go/kubernetes/typed/core/v1"
//clients which can control all the resources
//by conversion, xxxxGetter will return xxxxInterface
type CoreV1Interface interface {
	RESTClient() rest.Interface
	ComponentStatusesGetter
	ConfigMapsGetter
	EndpointsGetter
	EventsGetter
	LimitRangesGetter
	NamespacesGetter
	NodesGetter
	PersistentVolumesGetter
	PersistentVolumeClaimsGetter
	PodsGetter
	PodTemplatesGetter
	ReplicationControllersGetter
	ResourceQuotasGetter
	SecretsGetter
	ServicesGetter
	ServiceAccountsGetter
}

//core v1 interface contructor
func NewForConfig(c *rest.Config) (*CoreV1Client, error)
func New(c rest.Interface) *CoreV1Client

//work with pod resource
type PodInterface interface {
	Create(*v1.Pod) (*v1.Pod, error)
	Update(*v1.Pod) (*v1.Pod, error)
	UpdateStatus(*v1.Pod) (*v1.Pod, error)
	Delete(name string, options *meta_v1.DeleteOptions) error
	DeleteCollection(options *meta_v1.DeleteOptions, listOptions meta_v1.ListOptions) error
	Get(name string, options meta_v1.GetOptions) (*v1.Pod, error)
	List(opts meta_v1.ListOptions) (*v1.PodList, error)
	Watch(opts meta_v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.Pod, err error)
	PodExpansion
}

//work with replication controller
type ReplicationControllerInterface interface {
	Create(*v1.ReplicationController) (*v1.ReplicationController, error)
	Update(*v1.ReplicationController) (*v1.ReplicationController, error)
	UpdateStatus(*v1.ReplicationController) (*v1.ReplicationController, error)
	Delete(name string, options *meta_v1.DeleteOptions) error
	DeleteCollection(options *meta_v1.DeleteOptions, listOptions meta_v1.ListOptions) error
	Get(name string, options meta_v1.GetOptions) (*v1.ReplicationController, error)
	List(opts meta_v1.ListOptions) (*v1.ReplicationControllerList, error)
	Watch(opts meta_v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.ReplicationController, err error)
	GetScale(replicationControllerName string, options meta_v1.GetOptions) (*v1beta1.Scale, error)
	UpdateScale(replicationControllerName string, scale *v1beta1.Scale) (*v1beta1.Scale, error)

	ReplicationControllerExpansion //interface{}
}

//v1core "k8s.io/client-go/kubernetes/typed/core/v1"
//implement CoreV1Interface
type CoreV1Client struct {
	restClient rest.Interface
}

//k8s.io/client-go/rest/client.go
type Interface interface {
	GetRateLimiter() flowcontrol.RateLimiter
	Verb(verb string) *Request
	Post() *Request
	Put() *Request
	Patch(pt types.PatchType) *Request
	Get() *Request
	Delete() *Request
	APIVersion() schema.GroupVersion
}

//Request support build factory style
func (c *pods) Update(pod *v1.Pod) (result *v1.Pod, err error) {
	result = &v1.Pod{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("pods").
		Name(pod.Name).
		Body(pod).
		Do().
		Into(result)
	return
}

//k8s.io/client-go/kubernetes/clientset.go
//note:In 1.9, apps/v1 is introduced, and extensions/v1beta1, apps/v1beta1 and apps/v1beta2 are deprecated.
type Interface interface {
	Discovery() discovery.DiscoveryInterface
	AdmissionregistrationV1beta1() admissionregistrationv1beta1.AdmissionregistrationV1beta1Interface
	AppsV1() appsv1.AppsV1Interface
	AuthenticationV1() authenticationv1.AuthenticationV1Interface
	AuthorizationV1() authorizationv1.AuthorizationV1Interface
	AutoscalingV1() autoscalingv1.AutoscalingV1Interface
	BatchV1() batchv1.BatchV1Interface
	CertificatesV1beta1() certificatesv1beta1.CertificatesV1beta1Interface
	CoreV1() corev1.CoreV1Interface
	EventsV1beta1() eventsv1beta1.EventsV1beta1Interface
	ExtensionsV1beta1() extensionsv1beta1.ExtensionsV1beta1Interface
	NetworkingV1() networkingv1.NetworkingV1Interface
	PolicyV1beta1() policyv1beta1.PolicyV1beta1Interface
	RbacV1() rbacv1.RbacV1Interface
	SchedulingV1alpha1() schedulingv1alpha1.SchedulingV1alpha1Interface
	SettingsV1alpha1() settingsv1alpha1.SettingsV1alpha1Interface
	StorageV1() storagev1.StorageV1Interface
}

//kubeconfig ---> ~/.kube/config
//config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
//clientset, err := kubernetes.NewForConfig(config)

//k8s.io/client-go/discovery/discovery_client.go
type DiscoveryInterface interface {
	RESTClient() restclient.Interface
	ServerGroupsInterface
	ServerResourcesInterface
	ServerVersionInterface
	OpenAPISchemaInterface
}

//k8s.io/client-go/kubernetes/typed/admissionregistration/v1beta1/admissionregistration_client.go
type AdmissionregistrationV1beta1Interface interface {
	RESTClient() rest.Interface
	MutatingWebhookConfigurationsGetter
	ValidatingWebhookConfigurationsGetter
}

//k8s.io/client-go/kubernetes/typed/apps/v1/apps_client.go
type AppsV1Interface interface {
	RESTClient() rest.Interface
	ControllerRevisionsGetter
	DaemonSetsGetter
	DeploymentsGetter
	ReplicaSetsGetter
	StatefulSetsGetter
}

//k8s.io/client-go/kubernetes/typed/authentication/v1/authentication_client.go
type AuthenticationV1Interface interface {
	RESTClient() rest.Interface
	TokenReviewsGetter
}

//k8s.io/client-go/kubernetes/typed/authorization/v1/authorization_client.go
type AuthorizationV1Interface interface {
	RESTClient() rest.Interface
	LocalSubjectAccessReviewsGetter
	SelfSubjectAccessReviewsGetter
	SelfSubjectRulesReviewsGetter
	SubjectAccessReviewsGetter
}

//k8s.io/client-go/kubernetes/typed/autoscaling/v1/autoscaling_client.go
type AutoscalingV1Interface interface {
	RESTClient() rest.Interface
	HorizontalPodAutoscalersGetter
}

//k8s.io/client-go/kubernetes/typed/batch/v1/batch_client.go
type BatchV1Interface interface {
	RESTClient() rest.Interface
	JobsGetter
}

//k8s.io/client-go/kubernetes/typed/certificates/v1beta1/certificates_client.go
type CertificatesV1beta1Interface interface {
	RESTClient() rest.Interface
	CertificateSigningRequestsGetter
}

//k8s.io/client-go/kubernetes/typed/events/v1beta1/events_client.go
type EventsV1beta1Interface interface {
	RESTClient() rest.Interface
	EventsGetter
}

//k8s.io/client-go/kubernetes/typed/extensions/v1beta1/extensions_client.go
type ExtensionsV1beta1Interface interface {
	RESTClient() rest.Interface
	DaemonSetsGetter
	DeploymentsGetter
	IngressesGetter
	PodSecurityPoliciesGetter
	ReplicaSetsGetter
	ScalesGetter
}

//k8s.io/client-go/kubernetes/typed/networking/v1/networking_client.go
type NetworkingV1Interface interface {
	RESTClient() rest.Interface
	NetworkPoliciesGetter
}

//k8s.io/client-go/kubernetes/typed/policy/v1beta1/policy_client.go
type PolicyV1beta1Interface interface {
	RESTClient() rest.Interface
	EvictionsGetter
	PodDisruptionBudgetsGetter
	PodSecurityPoliciesGetter
}

//k8s.io/client-go/kubernetes/typed/rbac/v1/rbac_client.go
type RbacV1Interface interface {
	RESTClient() rest.Interface
	ClusterRolesGetter
	ClusterRoleBindingsGetter
	RolesGetter
	RoleBindingsGetter
}

//k8s.io/client-go/kubernetes/typed/scheduling/v1alpha1/scheduling_client.go
type SchedulingV1alpha1Interface interface {
	RESTClient() rest.Interface
	PriorityClassesGetter
}

//k8s.io/client-go/kubernetes/typed/settings/v1alpha1/settings_client.go
type SettingsV1alpha1Interface interface {
	RESTClient() rest.Interface
	PodPresetsGetter
}

//k8s.io/client-go/kubernetes/typed/storage/v1/storage_client.go
type StorageV1Interface interface {
	RESTClient() rest.Interface
	StorageClassesGetter
}
