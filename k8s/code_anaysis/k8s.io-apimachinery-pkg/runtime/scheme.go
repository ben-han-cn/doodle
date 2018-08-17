package runtime

type Scheme struct {
	gvkToType map[GroupVerisonKind]reflect.Type
	typeToGVK map[reflect.Type][]GroupVerisonKind

	unversionedTypes map[reflect.Type]GroupVerisonKind
	unversionedKinds map[string]reflect.Type

	fieldLabelConversionFuncs map[string]map[string]FieldLabelConversionFunc
	defaultFuncs              map[reflect.Type]func(interface{})
	converter                 *Converter
}

type Converter struct {
	conversionFuncs          ConversionFuncs
	generatedConversionFuncs ConversionFuncs

	genericConversions []GenericConversionFunc
	ignoredConversions map[typePair]struct{}

	structFieldDests   map[typeNamePair][]typeNamePair
	structFieldSources map[typeNamePair][]typeNamePair

	inputFieldMappingFuncs map[reflect.Type]FieldMappingFunc
	inputDefaultFlags      map[reflect.Type]FieldMatchingFlags

	Debug    DebugLogger
	nameFunc func(t reflect.Type) string
}

type GenericConversionFunc func(a, b interface{}, scope Scope) (bool, error)

type Scope interface {
	Convert(src, dest interface{}, flags FieldMatchingFlags) error
	DefaultConvert(src, dest interface{}, flags FieldMatchingFlags) error
	SrcTag() reflect.StructTag
	DestTag() reflect.StructTag
	Flags() FieldMatchingFlags
	Meta() *Meta
}

type FieldMatchingFlags int

//these constants could be combined
const (
	//only care dest, ignore field only exists in source
	DestFromSource FieldMatchingFlags = 0
	//only care source, ignore field only exists in dest
	SourceToDest FieldMatchingFlags = 1 << iota
	IgnoreMissingFields
	AllowDifferentFieldTypeNames
)
