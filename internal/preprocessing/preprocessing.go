package preprocessing

import (
	"regexp"
	"sync"

	"github.com/tylitianrui/file-clickhouse-exporter/pkg/aggregation"
)

const (
	PreprocessorAggregation = "aggregation"
	PreprocessorStatic      = "static"
	PreprocessorDynamic     = "dynamic"
	PreprocessorRaw         = "raw"
)

var (
	regexColumns      = regexp.MustCompile(`((\$\d+)|((aggregation)|(static)|(dynamic))\.(\w*))(\((\w*)\))?`)
	columnsAggregator = aggregation.NewAggregation(regexColumns)
)

type ProcessorLogic map[string]string

type Preprocessor interface {
	SetReadColumns(columns []string)
}

/*
  columns:
    time: $1(time)
    time_utc: $1(time_utc)
    name: aggregation.key1
    tags: static.a
    action: $2
    duration: $4(int32)
  Preprocessing:
    aggregation:
      key1: $2+" "+$3
      $4: $4[3:]
    static:
      a: 1
    dynamic:
      id: gen_uuid()

*/
// Preprocessor
type Preprocessing struct {
	rawColumns []string                     // raw configuration. eg:$1(time_utc),aggregation.key1
	readIndex  map[string]map[string]string // map[raw] = { map[$1] = time_utc} ,map[aggregation] = { map[key1] = string}
	processors map[string]ProcessorLogic
	mu         sync.Mutex
}

func NewPreprocessor() *Preprocessing {
	preprocessor := &Preprocessing{
		rawColumns: []string{},
		readIndex: map[string]map[string]string{
			PreprocessorRaw:         {},
			PreprocessorAggregation: {},
			PreprocessorStatic:      {},
			PreprocessorDynamic:     {},
		},
		processors: map[string]ProcessorLogic{},
	}
	return preprocessor
}

func (p *Preprocessing) SetColumns(columns map[string]string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.rawColumns = p.rawColumns[:0]
	for _, v := range columns {
		p.rawColumns = append(p.rawColumns, v)
	}
}

func (p *Preprocessing) SetReadColumns(columns []string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.rawColumns = append(p.rawColumns[:0], columns...)
}

func (p *Preprocessing) SetProcessorLogic(name string, logic map[string]string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.processors[name] = logic
}

func (p *Preprocessing) Load() {
	p.mu.Lock()
	defer p.mu.Unlock()
	for _, v := range p.rawColumns {
		res := columnsAggregator.Parse(v)
		if len(res) == 0 {
			panic("configure err:" + v)
		}
		configure := res[0]
		if len(configure) != 10 {
			panic("configure err:" + v)
		}
		rawKey := configure[2]
		preprocessingSource := configure[3]
		preprocessingSourceKey := configure[7]
		columnType := configure[9]
		if len(columnType) == 0 {
			columnType = "string"
		}
		if len(rawKey) > 0 {
			p.readIndex[PreprocessorRaw][rawKey] = columnType

		} else if len(preprocessingSource) > 0 {
			p.readIndex[preprocessingSource][preprocessingSourceKey] = columnType
		}
	}
}
