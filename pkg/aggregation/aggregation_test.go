// preprocessing
package aggregation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAggregation_ALL(t *testing.T) {
	a := assert.New(t)
	strJoinAggregation := NewStrJoinAggregation()
	demo1 := "$3[:-1]+\"   \"+$4[:4]+$5[:5]+\"he\"+$6[2:-1]"
	res1, _ := strJoinAggregation.ParseRule(demo1)
	expected1 := [][]string{
		{"$3", "", "-1", "   "},
		{"$4", "", "4", ""},
		{"$5", "", "5", "he"},
		{"$6", "2", "-1", ""},
	}
	a.Equal(expected1, res1)
	data := map[string]string{
		"$3": "hello",
		"$4": "world",
		"$5": "myname",
		"$6": "tyltrli",
	}
	res, _ := strJoinAggregation.Aggregate(data)
	expeact2 := "hell   worlmynamheltrl"
	a.Equal(expeact2, res)

}
