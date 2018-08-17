package label

//selector in
//	equality-based
//	set-based

type Labels interface {
	Has(label string) (exists bool)
	Get(label string) (value string)
}

type Selector interface {
	Matches(Labels) bool
	Empty() bool
	String() string
	Add(r ...Requirement) Selector
	Requirements() (requirements Requirements, selectable bool)
	DeepCopySelector() Selector
}

type Requirement struct {
	key       string
	operator  selection.Operator
	strValues []string
}

func NewRequirement(key string, op selection.Operator, vals []string) (*Requirement, error) {
	if err := validateLabelKey(key); err != nil {
		return nil, err
	}
	switch op {
	case selection.In, selection.NotIn:
		assert(len(vars) > 0)
	case selection.Equals, selection.DoubleEquals, selection.NotEquals:
		assert(len(vars) == 1)
	case selection.Exists, selection.DoesNotExist:
		assert(len(vars) == 0)
	case selection.GreaterThan, selection.LessThan:
		assert(len(vars) == 1)
		for i := range vals {
			assert(strconv.ParseInt(vals[i], 10, 64) == nil)
		}
	default:
	}

	for i := range vals {
		if err := validateLabelValue(vals[i]); err != nil {
			return nil, err
		}
	}
	sort.Strings(vals)
	return &Requirement{key: key, operator: op, strValues: vals}, nil
}
