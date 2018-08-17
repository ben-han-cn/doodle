package runtime

type TypeMeta struct {
	APIVersion string
	Kind       string
}

type RawExtension struct {
	Raw    []byte
	Object Object
}

type Unknown struct {
	TypeMeta

	Raw             []byte
	ContentEncoding string
	ContentType     string
}
