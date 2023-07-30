package file_parser

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

var (
	reg = regexp.MustCompile(`\s+`)
)

func init() {
	DefaultParserController.RegisterParser("file", &FileParser{})
}

type FileParser struct {
	formatString []string
	idx          []int
}

func (fp *FileParser) SetFormatString(s string) error {
	s = strings.ReplaceAll(s, " ", "")
	l := strings.Split(s, "$")
	l = l[1:]
	formatString := make([]string, len(l))
	idx := make([]int, len(l))
	for i := 0; i < len(l); i++ {
		formatString[i] = "$" + l[i]
		index, err := strconv.Atoi(l[i])
		if err != nil {
			return err
		}
		if index < 0 {
			return errors.New("$i:i >= 0")
		}
		idx[i] = index

	}
	fp.formatString = formatString
	fp.idx = idx
	return nil
}

func (fp *FileParser) SetFormat(format_list []string) error {
	idx := make([]int, len(format_list))
	for i := 0; i < len(format_list); i++ {
		format := strings.ReplaceAll(format_list[i], " ", "")
		if len(format) <= 1 {
			return errors.New("format:err " + format)
		}
		formatidx := format[1:]
		index, err := strconv.Atoi(formatidx)
		if err != nil {
			return err
		}
		if index < 0 {
			return errors.New("$i:i >= 0")
		}
		idx[i] = index
		format_list[i] = format

	}
	fp.formatString = format_list
	fp.idx = idx
	return nil
}

func (fp *FileParser) Parse(s string) map[string]string {
	s = strings.Trim(s, " ")
	result := reg.Split(s, -1)
	res := make(map[string]string)
	for i, v := range fp.idx {
		if v == 0 {
			res[fp.formatString[i]] = s
			continue
		}
		if v > len(result) {
			res[fp.formatString[i]] = ""
			continue
		}
		res[fp.formatString[i]] = result[v-1]

	}
	return res
}

func (fp *FileParser) ParseColumns(idx []string, s string) []string {
	panic("")

}
