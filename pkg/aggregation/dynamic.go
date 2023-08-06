package aggregation

import (
	"sync"

	uuid "github.com/satori/go.uuid"
)

var DefaultDynamicGenerator *Dynamic

func init() {
	DefaultDynamicGenerator = NewDynamicGenerator(4)
	DefaultDynamicGenerator.Register("gen_uuid()", GenUUID)
}

type DynamicGenerator func() string

type Dynamic struct {
	generators map[string]DynamicGenerator
	mu         sync.Mutex
}

func NewDynamicGenerator(size int) *Dynamic {
	return &Dynamic{
		generators: make(map[string]DynamicGenerator, size),
	}
}

func (d *Dynamic) Register(name string, gen DynamicGenerator) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.generators[name] = gen
}
func (d *Dynamic) GetDynamicGenerator(name string) (DynamicGenerator, bool) {
	d.mu.Lock()
	defer d.mu.Unlock()
	generator, exist := d.generators[name]
	return generator, exist
}

func GenUUID() string {
	u4 := uuid.NewV4()
	return u4.String()
}
