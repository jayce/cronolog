package cronolog

import "testing"

func TestUnixParse(t *testing.T) {
	data := []struct {
		test string
		s    string
	}{
		{"%%%Y", "%2006"},
		{"%Y%m%d", "20060102"},
		{"%Y/%m/%d", "2006/01/02"},
		{"%Y/%m/%d %H:%M:%S", "2006/01/02 15:04:05"},
		{"%T", "15:04:05"},
		{"%D", "01/02/06"},
		{"%Y%m%d%H%M%S", "20060102150405"},
	}

	for _, v := range data {
		s, err := UnixToGolang(v.test)
		if err != nil {
			t.Fatal(err)
		}

		if s != v.s {
			t.Errorf("expect '%s', but got '%s'", v.s, s)
		}
	}
}
