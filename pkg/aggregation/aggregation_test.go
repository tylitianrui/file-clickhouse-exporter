// preprocessing
package aggregation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAggregation_ALL(t *testing.T) {
	a := assert.New(t)
	strJoinAggregation := NewStrJoinAggregation()
	demo1 := "$3[3:3]+\"   \"+$4[:4]+$5[:5]+\"he\"+$6[6:6]"
	res1 := strJoinAggregation.Parse(demo1)
	expected1 := [][]string{
		{"$3[3:3]+\"   \"+", "$3[3:3]+\"   \"+", "$3", "[3:3]", "3", "3", "+\"   \"+", "   "},
		{"$4[:4]", "$4[:4]", "$4", "[:4]", "", "4", "", ""},
		{"$5[:5]+\"he\"+", "$5[:5]+\"he\"+", "$5", "[:5]", "", "5", "+\"he\"+", "he"},
		{"$6[6:6]", "$6[6:6]", "$6", "[6:6]", "6", "6", "", ""},
	}
	a.Equal(expected1, res1)
}
