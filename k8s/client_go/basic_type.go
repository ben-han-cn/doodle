//k8s.io/apimachinery/pkg/runtime/interfaces.go
type Object interface {
  GetObjectKind() schema.ObjectKind
  DeepCopyObject() Object
}


 //k8s.io/apimachinery/pkg/api/meta
 func Accessor(obj interface{}) (metav1.Object, error)
 func TypeAccessor(obj interface{}) (metav1.Type, error)


 //k8s.io/apimachinery/pkg/apis/meta/v1/meta.go
 type Object interface {
  GetNamespace() string
  SetNamespace(namespace string)
  GetName() string
  SetName(name string)
  GetGenerateName() string
  SetGenerateName(name string)
  GetUID() types.UID
  SetUID(uid types.UID)
  GetResourceVersion() string
  SetResourceVersion(version string)
  GetGeneration() int64
  SetGeneration(generation int64)
  GetSelfLink() string
  SetSelfLink(selfLink string)
  GetCreationTimestamp() Time
  SetCreationTimestamp(timestamp Time)
  GetDeletionTimestamp() *Time
  SetDeletionTimestamp(timestamp *Time)
  GetDeletionGracePeriodSeconds() *int64
  SetDeletionGracePeriodSeconds(*int64)
  GetLabels() map[string]string
  SetLabels(labels map[string]string)
  GetAnnotations() map[string]string
  SetAnnotations(annotations map[string]string)
  GetInitializers() *Initializers
  SetInitializers(initializers *Initializers)
  GetFinalizers() []string
  SetFinalizers(finalizers []string)
  GetOwnerReferences() []OwnerReference
  SetOwnerReferences([]OwnerReference)
  GetClusterName() string
  SetClusterName(clusterName string)
}


type Type interface {
  GetAPIVersion() string
  SetAPIVersion(version string)
  GetKind() string
  SetKind(kind string)
}

