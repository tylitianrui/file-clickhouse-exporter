package preprocessing

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPreprocessor(t *testing.T) {
	a := assert.New(t)
	demo1 := "$1"
	res1 := regexColumns.FindAllStringSubmatch(demo1, -1)
	expect1 := [][]string{{"$1", "$1", "$1", "", "", "", "", "", "", ""}}
	a.Equal(expect1, res1)

	demo2 := "$2(string)"
	res2 := regexColumns.FindAllStringSubmatch(demo2, -1)
	expect2 := [][]string{{"$2(string)", "$2", "$2", "", "", "", "", "", "(string)", "string"}}
	a.Equal(expect2, res2)

	demo3 := "aggregation.key1(string)"
	res3 := regexColumns.FindAllStringSubmatch(demo3, -1)
	expect3 := [][]string{{"aggregation.key1(string)", "aggregation.key1", "", "aggregation", "aggregation", "", "", "key1", "(string)", "string"}}
	a.Equal(expect3, res3)

}
