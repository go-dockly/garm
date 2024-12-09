package mapper

import "fmt"

// LabelGenerator handles creation of unique labels
type LabelGenerator struct {
	counter int
	prefix  string
}

func NewLabelGenerator(prefix string) *LabelGenerator {
	return &LabelGenerator{
		counter: 0,
		prefix:  prefix,
	}
}

func (lg *LabelGenerator) Next() string {
	label := fmt.Sprintf("%s_%d", lg.prefix, lg.counter)
	lg.counter++
	return label
}

// LabelManager handles creation and tracking of assembly labels
type LabelManager struct {
	counter  int
	used     map[string]bool
	prefixes map[string]string
}

func NewLabelManager() *LabelManager {
	return &LabelManager{
		used: make(map[string]bool),
		prefixes: map[string]string{
			"loop":     "L",
			"cond":     "C",
			"body":     "B",
			"post":     "P",
			"end":      "E",
			"continue": "CONT",
			"break":    "BRK",
		},
	}
}

func (lm *LabelManager) Generate(prefix string) string {
	label := fmt.Sprintf("%s%d", lm.prefixes[prefix], lm.counter)
	lm.counter++
	lm.used[label] = false
	return label
}

func (lm *LabelManager) MarkUsed(label string) {
	lm.used[label] = true
}
