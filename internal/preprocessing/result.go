package preprocessing

import "github.com/tylitianrui/file-clickhouse-exporter/pkg/type_transfer"

type Result struct {
	raw         map[string]string
	aggregation map[string]string
	dynamic     map[string]string
	static      map[string]string
}

func NewResult() *Result {
	return &Result{}
}

func (r *Result) SetRaw(raw map[string]string) {
	r.raw = raw
}

func (r *Result) SetAggregation(aggregation map[string]string) {
	r.aggregation = aggregation
}
func (r *Result) SetDynamic(dynamic map[string]string) {
	r.dynamic = dynamic
}
func (r *Result) SetStatic(static map[string]string) {
	r.static = static
}

func (r *Result) Result(index []ResultIdx) map[string]interface{} {
	res := make(map[string]interface{})
	for _, resultIdx := range index {
		var dataSource map[string]string
		switch resultIdx.Source {
		case PreprocessorRaw:
			dataSource = r.raw
		case PreprocessorAggregation:
			dataSource = r.aggregation
		case PreprocessorDynamic:
			dataSource = r.dynamic
		case PreprocessorStatic:
			dataSource = r.static
		default:
			dataSource = r.raw
		}
		val, exist := dataSource[resultIdx.SourceKey]
		if exist {
			switch resultIdx.SourceType {
			case "int16":
				res[resultIdx.Key] = type_transfer.String2Int16(val)
			case "int32":
				res[resultIdx.Key] = type_transfer.String2Int32(val)
			case "int64":
				res[resultIdx.Key] = type_transfer.String2Int64(val)
			case "time":
				res[resultIdx.Key] = type_transfer.String2Time(val)
			case "time_utc":
				res[resultIdx.Key] = type_transfer.String2TimeUTC(val)
			default:
				res[resultIdx.Key] = val
			}

		}
	}
	return res

}
