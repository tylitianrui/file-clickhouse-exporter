package type_transfer

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestString2Time(t *testing.T) {
	a := assert.New(t)
	tt := "2023-07-24T09:00:01.355626+0000"
	ttime := String2Time(tt)
	etime, _ := time.ParseInLocation("2006-01-02T15:04:05.99999", "2023-07-24T09:00:01.355626", time.UTC)
	ok := etime.Equal(ttime)
	a.Equal(true, ok)
}
