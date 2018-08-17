package meta

/*
group: defines a logical names for a set of REST
resources under the root.

version: REST resource endpoints that evolve together within a group

type/kind: a named entity (Pod, Deployment)

resource: endpoint that handles REST request an represented as:
/apis/<grouop>/<version>/namespaces/<namespace>/<kind-plural>

*/

//k8s.io/apimachinery/pkg/apimachinery/types.go
// GroupMeta stores the metadata of a group.
type GroupMeta struct {
	GroupVersion  schema.GroupVersion
	GroupVersions []schema.GroupVersion

	// mapping between REST paths and the objects
	// declared in a Scheme and all known versions.
	RESTMapper meta.RESTMapper

	// InterfacesFor returns the default Codec and ResourceVersioner
	// for a given version string, or an error if the version is not known.
	InterfacesFor func(version schema.GroupVersion) (*meta.VersionInterfaces, error)
	// InterfacesByVersion stores the per-version interfaces.
	InterfacesByVersion map[schema.GroupVersion]*meta.VersionInterfaces
}

//k8s.io/apimachinery/pkg/api/meta/interfaces.go
type VersionInterfaces struct {
	runtime.ObjectConvertor
	MetadataAccessor
}

type MetadataAccessor interface {
	APIVersion(obj runtime.Object) (string, error)
	SetAPIVersion(obj runtime.Object, version string) error

	Kind(obj runtime.Object) (string, error)
	SetKind(obj runtime.Object, kind string) error

	Namespace(obj runtime.Object) (string, error)
	SetNamespace(obj runtime.Object, namespace string) error

	Name(obj runtime.Object) (string, error)
	SetName(obj runtime.Object, name string) error

	GenerateName(obj runtime.Object) (string, error)
	SetGenerateName(obj runtime.Object, name string) error

	UID(obj runtime.Object) (types.UID, error)
	SetUID(obj runtime.Object, uid types.UID) error

	SelfLink(obj runtime.Object) (string, error)
	SetSelfLink(obj runtime.Object, selfLink string) error

	Labels(obj runtime.Object) (map[string]string, error)
	SetLabels(obj runtime.Object, labels map[string]string) error

	Annotations(obj runtime.Object) (map[string]string, error)
	SetAnnotations(obj runtime.Object, annotations map[string]string) error

	Continue(obj runtime.Object) (string, error)
	SetContinue(obj runtime.Object, c string) error

	runtime.ResourceVersioner
}

//k8s.io/apimachinery/pkg/runtime/interfaces.go
// ObjectConvertor converts an object to a different version.
type ObjectConvertor interface {
	Convert(in, out, context interface{}) error
	ConvertToVersion(in Object, gv GroupVersioner) (out Object, err error)
	ConvertFieldLabel(version, kind, label, value string) (string, string, error)
}

type Object interface {
	GetObjectKind() schema.ObjectKind
	DeepCopyObject() Object
}

type ResourceVersioner interface {
	SetResourceVersion(obj Object, version string) error
	ResourceVersion(obj Object) (string, error)
}

//k8s.io/apimachinery/pkg/runtime/scheme.go
// In a Scheme, a Type is a particular Go struct, a Version is a point-in-time
// identifier for a particular representation of that Type (typically backwards
// compatible), a Kind is the unique name for that Type within the Version, and a
// Group identifies a set of Versions, Kinds, and Types that evolve over time. An
// Unversioned Type is one that is not yet formally bound to a type and is promised
// to be backwards compatible (effectively a "v1" of a Type that does not expect
// to break in the future).
//
// Schemes are not expected to change at runtime and are only threadsafe after
// registration is complete.
type Scheme struct {
	gvkToType        map[schema.GroupVersionKind]reflect.Type
	typeToGVK        map[reflect.Type][]schema.GroupVersionKind
	unversionedTypes map[reflect.Type]schema.GroupVersionKind
	unversionedKinds map[string]reflect.Type

	// Map from version and resource to the corresponding func to convert
	// resource field labels in that version to internal version.
	fieldLabelConversionFuncs map[string]map[string]FieldLabelConversionFunc

	// defaulterFuncs is an array of interfaces to be called with an object to provide defaulting
	// the provided object must be a pointer.
	defaulterFuncs map[reflect.Type]func(interface{})

	// converter stores all registered conversion functions. It also has
	// default coverting behavior.
	converter *conversion.Converter
}

//serialization
//k8s.io/apimachinery/pkg/runtime/interfaces.go
type Encoder interface {
	Encode(obj Object, w io.Writer) error
}

type Decoder interface {
	Decode(data []byte, defaults *schema.GroupVersionKind, into Object) (Object, *schema.GroupVersionKind, error)
}

type Serializer interface {
	Encoder
	Decoder
}

type Codec Serializer

// SerializerInfo contains information about a specific serialization format
type SerializerInfo struct {
	// MediaType is the value that represents this serializer over the wire.
	MediaType string
	// EncodesAsText indicates this serializer can be encoded to UTF-8 safely.
	EncodesAsText    bool
	Serializer       Serializer
	PrettySerializer Serializer
	StreamSerializer *StreamSerializerInfo
}

type NegotiatedSerializer interface {
	SupportedMediaTypes() []SerializerInfo
	EncoderForVersion(serializer Encoder, gv GroupVersioner) Encoder
	DecoderToVersion(serializer Decoder, gv GroupVersioner) Decoder
}

//discovery addresses
//"k8s.io/apiserver/pkg/endpoints/discovery/addresses.go"
// return server address based on client ip who sending the request
type Addresses interface {
	ServerAddressByClientCIDRs(net.IP) []metav1.ServerAddressByClientCIDR
}

type ServerAddressByClientCIDR struct {
	ClientCIDR    string `json:"clientCIDR" protobuf:"bytes,1,opt,name=clientCIDR"`
	ServerAddress string `json:"serverAddress" protobuf:"bytes,2,opt,name=serverAddress"`
}
