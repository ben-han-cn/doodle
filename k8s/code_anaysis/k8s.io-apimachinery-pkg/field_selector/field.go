package field

type Fields interface {
	Has(string) bool
	Get(string) string
}

type Requirement struct {
	Operator Operator
	Field    string
	Value    string
}

// field selector.
type Selector interface {
	Matches(Fields) bool
	Empty() bool
	RequiresExactMatch(field string) (value string, found bool)
	Transform(fn TransformFunc) (Selector, error)
	Requirements() Requirements
	String() string
	DeepCopySelector() Selector
}
type TransformFunc func(field, value string) (newField, newValue string, err error)

func OneTermEqualSelector(k, v string) Selector {
	return &hasTerm{field: k, value: v}
}

func AndSelectors(selectors ...Selector) Selector {
	return andTerm(selectors)
}

type Operator string

const (
	DoesNotExist Operator = "!"
	Equals       Operator = "="
	DoubleEquals Operator = "=="
	In           Operator = "in"
	NotEquals    Operator = "!="
	NotIn        Operator = "notin"
	Exists       Operator = "exists"
	GreaterThan  Operator = "gt"
	LessThan     Operator = "lt"
)
