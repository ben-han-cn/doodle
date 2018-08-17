package runtime

//resource/object is split into groups
//resource/object will evolve with different version
//kind the the concrete type of resource
//url/endpoint for resource is
// /apis/<group>/<version>/namespaces/<namespace>/<kind-plural>
type Object interface {
	GetObjectKind() ObjectKind
	DeepCopyObject() Object
}

// Unstructured allows objects that do not have Golang structs
// registered to be manipulated generically. This can be used
// to deal with the API objects from a plug-in.
type Unstructured interface {
	Object

	UnstructuredContent() map[string]interface{}
	SetUnstructuredContent(map[string]interface{})
	IsList() bool
	EachListItem(func(Object) error) error
}

type ObjectKind interface {
	SetGroupVersionKind(kind GroupVersionKind)
	GroupVersionKind() GroupVersionKind
}

type GroupResource struct {
	Group    string
	Resource string
}

type GroupVersionResource struct {
	Group    string
	Version  string
	Resource string
}

type GroupKind struct {
	Group string
	Kind  string
}

type GroupVersionKind struct {
	Group   string
	Version string
	Kind    string
}

type Encoder interface {
	Encode(obj Object, w io.Writer) error
}

type Decoder interface {
	Decode(data []byte, defaults *GroupVersionKind, into Object) (Object, *GroupVersionKind, error)
}

type Serializer interface {
	Encoder
	Decoder
}
type Codec Serializer

type StorageSerializer interface {
	SupportedMediaTypes() []SerializerInfo
	UniversalDeserializer() Decoder
	EncoderForVersion(Encoder, GroupVersioner) Encoder
	DecoderForVersion(Decoder, GroupVersioner) Decoder
}

type GroupVersioner interface {
	KindForGroupVersionKinds(kinds []GroupVersionKind) (GroupVersionKind, ok)
}

type ObjectDefaulter interface {
	Default(Object)
}

type ObjectVersioner interface {
	ConvertToVersion(Object, GroupVersioner) (Object, error)
}

type ObjectConverter interface {
	Convert(in, out, context interface{}) error
	ConvertToVersion(Object, GroupVersioner) (Object, error)
	ConvertFieldLabel(version, kind, label, value string) (string, string, error)
}

type ObjectTyper interface {
	ObjectKinds(Object) ([]GroupVersionKind, bool, error)
	Recognizes(GroupVersionKind) bool
}

type ObjectCreater interface {
	New(GroupVersionKind) (Object, error)
}

type ResourceVersioner interface {
	SetResourceVersion(obj Object, version string) error
	ResourceVersion(obj Object) (string, error)
}
