package config

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
)

var C Config

type Config struct {
	ClickHouse ClickHouse
	Setting    Setting
}

var regexstr = regexp.MustCompile(`(\$\d+)(\[(\d*):(\d*)\])?(\((\w+)\))?`)

type PreprocessingConfig map[string]string

type ClickHouse struct {
	DB            string                         `json:"db,omitempty" yaml:"db" gorm:"db" mapstructure:"db"`
	Table         string                         `json:"table,omitempty" yaml:"table" gorm:"table" mapstructure:"table"`
	Host          string                         `json:"host,omitempty" yaml:"host" gorm:"host" mapstructure:"host"`
	Port          int                            `json:"port,omitempty" yaml:"port" gorm:"port" mapstructure:"port"`
	Credentials   Credentials                    `json:"credentials,omitempty" yaml:"credentials" gorm:"credentials" mapstructure:"credentials"`
	Columns       map[string]string              `json:"columns,omitempty" yaml:"columns" gorm:"columns" mapstructure:"columns"`
	Preprocessing map[string]PreprocessingConfig `json:"preprocessing,omitempty" yaml:"preprocessing" gorm:"preprocessing" mapstructure:"preprocessing"`
}

type Preprocessing struct {
	Columns []string
	Index   []string
	Split   [][]int
	Types   []string
}

func (c *ClickHouse) BuildColumns() (preprocessing *Preprocessing, err error) {
	var (
		columns   []string = make([]string, len(c.Columns))
		index     []string = make([]string, len(c.Columns))
		str_split [][]int  = make([][]int, len(c.Columns))
		types     []string = make([]string, len(c.Columns))
	)
	columns = columns[:0]
	index = index[:0]
	str_split = str_split[:0]
	types = types[:0]

	for column, format_str := range c.Columns {
		columns = append(columns, column)
		substr := regexstr.FindStringSubmatch(format_str)
		if len(substr) != 7 {
			return nil, errors.New("Syntax Error: column: " + column + " format:" + format_str)

		}
		fmt.Println("parse column: " + column + " format:" + format_str)
		idx := substr[1]
		split := substr[2]
		split_from := substr[3]
		split_to := substr[4]
		trans := substr[5]
		trans_type := substr[6]
		from, to := 0, -1
		if len(split) > 0 {
			if len(split_from) > 0 {
				from, err = strconv.Atoi(split_from)
				if err != nil {
					msg := fmt.Sprintf("Syntax Error [%s :%s] `[int1:int2]`: expect int ,but got %s", column, format_str, split_from)
					return nil, errors.New(msg)
				}

			}
			if len(split_to) > 0 {
				to, err = strconv.Atoi(split_to)
				if err != nil {
					msg := fmt.Sprintf("Syntax Error [%s :%s] `[int1:int2]`: expect int ,but got %s", column, format_str, split_from)
					return nil, errors.New(msg)
				}

			}
		}
		if len(trans) > 0 {
			if len(trans_type) == 0 {
				msg := fmt.Sprintf("Syntax Error [%s :%s] `(type)`: expect a type ,but got %s", column, format_str, trans_type)
				return nil, errors.New(msg)
			}
		} else {
			trans_type = "string"
		}

		index = append(index, idx)
		types = append(types, trans_type)
		str_split = append(str_split, []int{from, to})
		fmt.Println("")

	}

	preprocessing = &Preprocessing{
		Columns: columns,
		Index:   index,
		Split:   str_split,
		Types:   types,
	}
	return preprocessing, nil

}

type Credentials struct {
	User     string `json:"user,omitempty" yaml:"user" gorm:"user" mapstructure:"user"`
	Password string `json:"password,omitempty" yaml:"password" gorm:"password" mapstructure:"password"`
}

type Setting struct {
	FilePath         string `yaml:"file_path" mapstructure:"file_path"`
	MaxlineEveryRead int    `yaml:"max_line_every_read" mapstructure:"max_line_every_read"`
	Interval         int    `yaml:"interval" mapstructure:"interval"`
}
