package type_transfer

import (
	"time"

	"github.com/spf13/cast"
)

func String2Int64(s string) int64 {
	return cast.ToInt64(s)
}

func String2Int32(s string) int32 {
	return cast.ToInt32(s)
}

func String2Int16(s string) int16 {
	return cast.ToInt16(s)
}
func String2UInt64(s string) uint64 {
	return cast.ToUint64(s)
}

func String2Float64(s string) float64 {
	return cast.ToFloat64(s)
}

func String2Float32(s string) float32 {
	return cast.ToFloat32(s)
}

func String2Time(s string) time.Time {
	t, _ := StringToTimeWithLocation(s, time.Local)
	return t
}

func String2TimeUTC(s string) time.Time {
	t, _ := StringToTimeWithLocation(s, time.UTC)
	return t
}
