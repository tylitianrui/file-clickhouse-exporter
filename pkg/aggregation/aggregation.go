// aggregation
package aggregation

import "regexp"

var (
	aggregationStrJoinRegex = regexp.MustCompile(`((\$\d+)(\[(\d*):(\d*)\])?(\+\"(.*?)\"\+)?)`)
)

type Aggregator interface {
	Parse(str string) [][]string
}

type Aggregation struct {
	re *regexp.Regexp
}

func NewAggregation(re *regexp.Regexp) Aggregator {
	return &Aggregation{
		re: re,
	}
}

func (a *Aggregation) Parse(str string) [][]string {
	res := a.re.FindAllStringSubmatch(str, -1)
	return res
}

// StrJoinAggregation string join and split.
type StrJoinAggregation struct {
	aggregation *Aggregation
}

// Parse implements Aggregator.
func (sa *StrJoinAggregation) Parse(str string) [][]string {
	return sa.aggregation.Parse(str)
}

func NewStrJoinAggregation() Aggregator {
	aggregation := NewAggregation(aggregationStrJoinRegex)
	return &StrJoinAggregation{
		aggregation: aggregation.(*Aggregation),
	}
}
