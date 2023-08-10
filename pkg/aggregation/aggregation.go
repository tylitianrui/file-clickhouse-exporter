// aggregation
package aggregation

import (
	"bytes"
	"errors"
	"math"
	"regexp"
	"strconv"
)

var (
	aggregationStrJoinRegex = regexp.MustCompile(`((\$\d+)(\[-?(\d*):(-?\d*)\])?(\+\"(.*?)\"\+)?)`)
)

type Aggregator interface {
	ParseRule(str string) ([][]string, error)
	Aggregate(data map[string]string) (string, error)
}

type Aggregation struct {
	re *regexp.Regexp
}

func NewAggregation(re *regexp.Regexp) Aggregator {
	return &Aggregation{
		re: re,
	}
}

func (a *Aggregation) ParseRule(str string) ([][]string, error) {
	res := a.re.FindAllStringSubmatch(str, -1)
	return res, nil
}
func (a *Aggregation) Aggregate(data map[string]string) (string, error) {
	return "", nil
}

// $3[1:3]+"hello"
//
//	idx=$3
//	from=1
//	to=3
//	join=hello
type aggregationItem struct {
	idx  string
	from int
	to   int
	join string
}

func (a *aggregationItem) frommAndToIndex(val string) (from int, to int, err error) {
	length := len(val)
	if a.from == math.MinInt && a.to == math.MaxInt {
		return 0, length, nil
	} else if a.from == math.MinInt {
		if a.to < 0 {
			to = length + a.to
		} else {
			to = a.to
		}
		from = 0

	} else if a.to == math.MaxInt {
		if a.from < 0 {
			from = length + a.from
		} else {
			from = a.from
		}
		to = length

	} else {
		if a.from < 0 {
			from = length + a.from
		} else {
			from = a.from
		}
		if a.to < 0 {
			to = length + a.to
		} else {
			to = a.to
		}
	}
	if from < 0 {
		return 0, 0, errors.New("")
	}
	if from > length {
		return 0, 0, errors.New("")
	}
	if to < 0 {
		return 0, 0, errors.New("")
	}
	if to > length {
		return 0, 0, errors.New("")
	}
	if from > to {
		return 0, 0, errors.New("")
	}
	return from, to, nil

}

// StrJoinAggregation string join and split.
// StrJoinAggregation handles one rule only.
type StrJoinAggregation struct {
	aggregation      *Aggregation
	aggregationRules []aggregationItem
}

// Parse implements Aggregator.
// eg: $3[3:3]+\"   \"+$4[:4]+$5[:5]+\"he\"+$6[6:6]
func (sa *StrJoinAggregation) ParseRule(ruleStr string) ([][]string, error) {

	rawRules, err := sa.aggregation.ParseRule(ruleStr)
	if err != nil {
		return nil, err
	}
	rules := make([][]string, len(rawRules))
	rules = rules[:0]
	aggregationRules := make([]aggregationItem, len(rawRules))
	aggregationRules = aggregationRules[:0]
	for _, ruleItems := range rawRules {
		fi := math.MinInt
		ti := math.MaxInt
		raw := ruleItems[0]
		if len(ruleItems) < 8 {
			return nil, errors.New("configure syntax error:" + raw)
		}
		idx := ruleItems[2]
		f := ruleItems[4]
		if f != "" {
			fi, err = strconv.Atoi(f)
			if err != nil {
				return nil, errors.New("configure syntax error:" + raw)
			}
		}
		t := ruleItems[5]
		if t != "" {
			ti, err = strconv.Atoi(t)
			if err != nil {
				return nil, errors.New("configure syntax error:" + raw)
			}
		}
		joinstr := ruleItems[7]
		aggregationRulesItem := aggregationItem{
			idx:  idx,
			from: fi,
			to:   ti,
			join: joinstr,
		}
		aggregationRules = append(aggregationRules, aggregationRulesItem)
		item := []string{idx, f, t, joinstr}
		rules = append(rules, item)
	}
	sa.aggregationRules = aggregationRules

	return rules, nil
}

func NewStrJoinAggregation() Aggregator {
	aggregation := NewAggregation(aggregationStrJoinRegex)
	return &StrJoinAggregation{
		aggregation:      aggregation.(*Aggregation),
		aggregationRules: []aggregationItem{},
	}
}

// according aggregation rules ,Aggregate data to a string
// aggregation rules: $3[3:3]+\"   \"+$4[:4]+$5[:5]+\"he\"+$6[6:6]
// data={"$3":"hello","$4":"world",$5:"yourname",$6:"tyltrli"}
// return `   worlyournhe`
func (sa *StrJoinAggregation) Aggregate(data map[string]string) (string, error) {
	var buffer bytes.Buffer
	for i := 0; i < len(sa.aggregationRules); i++ {
		rule := sa.aggregationRules[i]
		key := rule.idx
		val, exist := data[key]
		if !exist {
			return "", errors.New("error:" + key + " does not exist in data")
		}
		f, t, err := rule.frommAndToIndex(val)
		if err != nil {
			return "", err
		}

		val = val[f:t]
		buffer.WriteString(val)
		join := rule.join
		buffer.WriteString(join)

	}
	return buffer.String(), nil

}
