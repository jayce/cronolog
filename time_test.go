package cronolog

import (
	"testing"
	"time"
)

func TestUnixToGolang(t *testing.T) {
	date := time.Now()

	data := []struct {
		test string
		s    string
	}{
		{"02%%%Y", "02%" + date.Format("2006")},
		{"%Y%m%d", date.Format("20060102")},
		{"%Y/%m/%d", date.Format("2006/01/02")},
		{"%Y/%m/%d %H:%M:%S", date.Format("2006/01/02 15:04:05")},
		{"%T", date.Format("15:04:05")},
		{"%D", date.Format("01/02/06")},
		{"%Y%m%d%H%M%S", date.Format("20060102150405")},
	}

	for _, v := range data {
		s, err := UnixToGolang(v.test, date)
		if err != nil {
			t.Fatal(err)
		}

		if s != v.s {
			t.Errorf("expect '%s', but got '%s'", v.s, s)
		}
	}
}
