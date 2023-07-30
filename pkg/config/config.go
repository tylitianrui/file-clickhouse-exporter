package config

import (
	"regexp"
)

var C Config

type Config struct {
	Clickhouse Clickhouse
	File       File
}

var regexstr = regexp.MustCompile(`(\$\d+)\((\w+)\)`)

type Clickhouse struct {
	DB          string            `json:"db,omitempty" yaml:"db" gorm:"db" mapstructure:"db"`
	Table       string            `json:"table,omitempty" yaml:"table" gorm:"table" mapstructure:"table"`
	Host        string            `json:"host,omitempty" yaml:"host" gorm:"host" mapstructure:"host"`
	Port        int               `json:"port,omitempty" yaml:"port" gorm:"port" mapstructure:"port"`
	Credentials Credentials       `json:"credentials,omitempty" yaml:"credentials" gorm:"credentials" mapstructure:"credentials"`
	Columns     map[string]string `json:"columns,omitempty" yaml:"columns" gorm:"columns" mapstructure:"columns"`
}

func (c *Clickhouse) BuildColumns() (columns []string, index []string, types []string) {
	for column, idx := range c.Columns {
		columns = append(columns, column)
		substr := regexstr.FindStringSubmatch(idx)
		if len(substr) == 3 {
			index = append(index, substr[1])
			types = append(types, substr[2])
		} else {
			index = append(index, idx)
			types = append(types, "string")
		}

	}
	return

}

type Credentials struct {
	User     string `json:"user,omitempty" yaml:"user" gorm:"user" mapstructure:"user"`
	Password string `json:"password,omitempty" yaml:"password" gorm:"password" mapstructure:"password"`
}

type File struct {
	Path string `json:"path,omitempty" yaml:"path" gorm:"path" mapstructure:"path"`
}
