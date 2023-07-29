package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_config_GetCnfDefault(t *testing.T) {
	a := assert.New(t)
	configParser := NewConfig()
	a.IsType(&config{}, configParser)

	configParser.SetCnfFile("yaml", "test", ".")
	configParser.Load()

	name := configParser.c.Get("test.name")
	a.EqualValues(name, "tyltr")

	hello := configParser.c.Get("hello")
	a.EqualValues(hello, nil)

	defkey := configParser.GetCnfDefault("heelo", 0)
	a.EqualValues(defkey, 0)

}
