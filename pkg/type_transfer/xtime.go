package type_transfer

import (
	"fmt"
	"time"
)

var timeFormatList = []string{
	time.RFC3339,
	time.RFC1123Z,
	time.RFC1123,
	time.RFC822Z,
	time.RFC822,
	time.RFC850,
	time.ANSIC,
	time.UnixDate,
	time.RubyDate,
	time.Kitchen,
	time.Stamp,
	time.StampMilli,
	time.StampMicro,
	time.StampNano,
	"2006-01-02T15:04:05",
	"2006-01-02 15:04:05.999999999 -0700 MST",
	"2006-01-02",
	"02 Jan 2006",
	"02/Jan/2006:15:04:05",
	"2006-01-02T15:04:05-0700",
	"2006-01-02T15:04:05.999999-0700",
	"2006-01-02 15:04:05 -07:00",
	"2006-01-02 15:04:05 -0700",
	"2006-01-02 15:04:05Z07:00",
	"2006-01-02 15:04:05Z0700",
	"2006-01-02 15:04:05",
}

func StringToTimeWithLocation(s string, loc *time.Location) (time.Time, error) {
	return parseWithLocation(s, loc, timeFormatList)
}

func parseWithLocation(s string, loc *time.Location, dates []string) (d time.Time, e error) {
	for _, dateType := range dates {
		if d, e = time.ParseInLocation(dateType, s, loc); e == nil {
			return
		}
	}
	return d, fmt.Errorf("unable to parse date: %s", s)
}
