package op

type PredicateCondition string

const (
	Always       PredicateCondition = "AL"
	Equal        PredicateCondition = "EQ"
	NotEqual     PredicateCondition = "NE"
	Greater      PredicateCondition = "GT"
	GreaterEqual PredicateCondition = "GE"
	Less         PredicateCondition = "LT"
	LessEqual    PredicateCondition = "LE"
)

var flagMap = map[PredicateCondition]uint8{
	Always:       0b1110,
	Equal:        0b0000,
	NotEqual:     0b0001,
	Greater:      0b1100,
	GreaterEqual: 0b1010,
	Less:         0b1011,
	LessEqual:    0b1101,
}

// Predicate represents ARM's conditional execution
type Predicate struct {
	Condition PredicateCondition // EQ, NE, GT, GE, LT, LE, etc.
	Flags     uint8              // NZCV flags
}

func NewPredicate(condition PredicateCondition) Predicate {
	return Predicate{
		Condition: condition,
		Flags:     flagMap[condition],
	}
}

func (p Predicate) String() string {
	return string(p.Condition)
}
