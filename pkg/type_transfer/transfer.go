package type_transfer

import (
	"time"

	"github.com/spf13/cast"
)

func String2Int64(s string) int64 {
	return cast.ToInt64(s)
}

func String2UInt64(s string) uint64 {
	return cast.ToUint64(s)
}
func String2Time(s string) time.Time {
	return cast.ToTime(s)

}
