package aggregation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDynamic(t *testing.T) {
	a := assert.New(t)
	gen, exist := DefaultDynamicGenerator.GetDynamicGenerator("uuid")
	a.Equal(exist, false)
	a.Nil(gen)

	gen, exist = DefaultDynamicGenerator.GetDynamicGenerator("gen_uuid()")
	a.Equal(exist, true)
	a.NotNil(gen)
}
