package file_parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileParser_Parse(t *testing.T) {

	a := assert.New(t)
	s := "    1 2    你 3 4  \n\r"
	fp := &FileParser{}

	err := fp.SetFormatString("$1   $2$3")
	a.NoError(err)
	res := fp.Parse(s)
	expect := map[string]string{
		"$1": "1",
		"$2": "2",
		"$3": "你",
	}
	a.Equal(expect, res)

	err = fp.SetFormatString("$1   $2$100")
	a.NoError(err)
	res = fp.Parse(s)
	expect = map[string]string{
		"$1":   "1",
		"$2":   "2",
		"$100": "",
	}
	a.Equal(expect, res)

	err = fp.SetFormat([]string{"$2", "$4"})
	a.NoError(err)
	res = fp.Parse(s)
	expect = map[string]string{
		"$4": "3",
		"$2": "2",
	}
	a.Equal(expect, res)

}
