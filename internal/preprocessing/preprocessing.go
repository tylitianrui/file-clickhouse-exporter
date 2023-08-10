package preprocessing

import (
	"errors"
	"fmt"
	"regexp"
	"sync"

	"github.com/tylitianrui/file-clickhouse-exporter/pkg/aggregation"
	"github.com/tylitianrui/file-clickhouse-exporter/pkg/collector"
)

const (
	BuffSize = 1 << 6
)

const (
	PreprocessorAggregation = "aggregation"
	PreprocessorStatic      = "static"
	PreprocessorDynamic     = "dynamic"
	PreprocessorRaw         = "raw"
)

var (
	preprocessingTyps  = []string{PreprocessorAggregation, PreprocessorStatic, PreprocessorDynamic}
	regexColumns       = regexp.MustCompile(`((\$\d+)|((aggregation)|(static)|(dynamic))\.(\w*))(\((\w*)\))?`)
	columnsAggregator  = aggregation.NewAggregation(regexColumns)
	strJoinAggregation = aggregation.NewStrJoinAggregation()
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
	mu                  sync.Mutex
	preprocessingConfig map[string]map[string]string      // config.yaml
	columnsConfig       map[string]string                 // config.yaml
	aggregations        map[string]aggregation.Aggregator //
	readIndexInFile     []string
}

func NewPreprocessor() *Preprocessing {
	preprocessor := &Preprocessing{
		preprocessingConfig: map[string]map[string]string{},
		columnsConfig:       map[string]string{},
		aggregations:        make(map[string]aggregation.Aggregator),
		readIndexInFile:     make([]string, 0),
	}
	return preprocessor
}

// config.yaml  option preprocessing
func (p *Preprocessing) SetPreprocessingConfig(cnf map[string]map[string]string) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	for k, _ := range cnf {
		if collector.IndexInStringArray(k, preprocessingTyps) < 0 {
			msg := fmt.Sprintf("config option err:%s is unknown", k)
			return errors.New(msg)
		}
	}
	p.preprocessingConfig = cnf
	return nil
}

func (p *Preprocessing) SetColumns(columns map[string]string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.columnsConfig = columns
}

func (p *Preprocessing) LoadConfig() error {
	p.mu.Lock()
	defer p.mu.Unlock()
	readColumnsSet := collector.NewSet() // 在文件中读取的索引
	if len(p.columnsConfig) == 0 {
		return errors.New("err: have not  run `SetColumns` yet")
	}
	if len(p.preprocessingConfig) == 0 {
		return errors.New("err: have not  run `SetPreprocessingConfig` yet")
	}
	for _, cnf := range p.columnsConfig {
		res, err := columnsAggregator.ParseRule(cnf)
		if err != nil {
			return err
		}
		if len(res) == 0 {
			panic("configure err:" + cnf)
		}
		configure := res[0]
		if len(configure) != 10 {
			panic("configure err:" + cnf)
		}

		// columns[time]: $1(time)
		// rawIndex $1
		rawIndex := configure[2]
		if rawIndex != "" {
			readColumnsSet.Add(rawIndex)
		}

		//  name: aggregation.key1
		//  preprocessingSource= aggregation
		//  preprocessingSourceKey=key1
		preprocessingSource := configure[3]
		preprocessingSourceKey := configure[7]
		columnType := configure[9]
		if len(columnType) == 0 {
			columnType = "string"
		}
		if preprocessingSource == PreprocessorAggregation {
			aggregationRule := p.preprocessingConfig[preprocessingSource][preprocessingSourceKey]
			strJoinAggregation := aggregation.NewStrJoinAggregation()
			rules, _ := strJoinAggregation.ParseRule(aggregationRule)
			p.aggregations[preprocessingSourceKey] = strJoinAggregation
			for _, rule := range rules {
				if len(rule) == 0 {
					return errors.New("")
				}
				idx := rule[0]
				readColumnsSet.Add(idx)
			}
		}
	}
	setItems := readColumnsSet.AllItems()
	for _, item := range setItems {
		i := item.(string)
		p.readIndexInFile = append(p.readIndexInFile, i)
	}
	return nil
}

func (p *Preprocessing) GetIndexOfFile() []string {
	return p.readIndexInFile
}

func (p *Preprocessing) Do(data map[string]string) map[string]string {
	result := make(map[string]string)
	for k, v := range p.aggregations {
		resItem, _ := v.Aggregate(data)
		result[k] = resItem
	}
	return result
}
