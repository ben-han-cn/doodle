package v1

type ObjectMeta struct {
	Name                       string
	GenerateName               string
	Namesapce                  string
	SelfLink                   string
	UID                        types.UID
	ResourceVersion            string
	Generation                 int64 //?
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

type OwnerReference struct {
	APIVersion         string
	Kind               string
	Name               string
	UID                types.UID
	Controller         *bool
	BlockOwnerDeletion *bool
}

type ListMeta struct {
	SelfLink        string
	ResourceVersion string
	Continue        string
}
