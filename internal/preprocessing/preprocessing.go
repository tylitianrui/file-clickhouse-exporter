package preprocessing

import (
	"fmt"
	"regexp"
	"sync"

	"github.com/tylitianrui/file-clickhouse-exporter/pkg/aggregation"
	"github.com/tylitianrui/file-clickhouse-exporter/pkg/collector"
)

const (
	PreprocessorAggregation = "aggregation"
	PreprocessorStatic      = "static"
	PreprocessorDynamic     = "dynamic"
	PreprocessorRaw         = "raw"
)

var (
	preprocessingTyps = []string{PreprocessorAggregation, PreprocessorStatic, PreprocessorDynamic}
	regexColumns      = regexp.MustCompile(`((\$\d+)|((aggregation)|(static)|(dynamic))\.(\w*))(\((\w*)\))?`)
	columnsAggregator = aggregation.NewAggregation(regexColumns)
)

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
	rawColumns               []string                     // raw configuration. eg:$1(time_utc),aggregation.key1
	readInxFromPreprocessing map[string]map[string]string // map[raw] = { map[$1] = time_utc} ,map[aggregation] = { map[key1] = string}
	processingConfiguration  map[string]map[string]string // processing configuration
	readIdxFromFile          []string
	mu                       sync.Mutex
}

func NewPreprocessor() *Preprocessing {
	preprocessor := &Preprocessing{
		rawColumns: []string{},
		readInxFromPreprocessing: map[string]map[string]string{
			PreprocessorRaw:         {},
			PreprocessorAggregation: {},
			PreprocessorStatic:      {},
			PreprocessorDynamic:     {},
		},
		processingConfiguration: map[string]map[string]string{},
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

func (p *Preprocessing) SetProcessing(processing map[string]map[string]string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.processingConfiguration = processing
}

func (p *Preprocessing) Load() {
	strJoinAggregation := aggregation.NewStrJoinAggregation()
	p.mu.Lock()
	defer p.mu.Unlock()
	set := collector.NewSet()
	// 加载数据索引
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
			p.readInxFromPreprocessing[PreprocessorRaw][rawKey] = columnType
			set.Add(rawKey)

		} else if len(preprocessingSource) > 0 {
			p.readInxFromPreprocessing[preprocessingSource][preprocessingSourceKey] = columnType
		}
	}

	// 检查 期望读取的预处理数据是否在配置中
	for _, preprocessingTyp := range preprocessingTyps {
		for expectKey, _ := range p.readInxFromPreprocessing[preprocessingTyp] {
			var contain bool
			for containKey, _ := range p.readInxFromPreprocessing[preprocessingTyp] {
				if containKey == expectKey {
					contain = true
				}

			}
			if !contain {
				msg := fmt.Sprintf("configuration[%s] does not  contain key[%s]", PreprocessorAggregation, expectKey)
				panic(msg)
			}
		}
	}

	// 配置中存在aggregation 并且会在aggregation中读取数据
	aggregations, ok := p.processingConfiguration[PreprocessorAggregation]
	if ok && len(p.readInxFromPreprocessing[PreprocessorAggregation]) > 0 {
		for k, v := range aggregations {
			if _, exist := p.readInxFromPreprocessing[PreprocessorAggregation][k]; exist {
				configs := strJoinAggregation.Parse(v)
				for _, config := range configs {
					idx := config[2]
					set.Add(idx)
				}
			}
		}
	}

	p.readIdxFromFile = p.readIdxFromFile[:0]
	for _, item := range set.AllItems() {
		p.readIdxFromFile = append(p.readIdxFromFile, item.(string))
	}
}
