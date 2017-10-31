package cronolog

import "fmt"

var flagsMap = map[string]string{
	"%": "%",
	"a": "Mon", "A": "Monday",
	"b": "Jan", "B": "January",
	"d": "02", "D": "01/02/06",
	"e": "02",
	"F": "2006-01-02",
	// "g": "", "G": "",
	"H": "15", "h": "Jan",
	"I": "03",
	// "j": "",
	"k": "15",
	"l": "03",
	"m": "01", "M": "04",
	//"s": "",
	"S": "05",
	"t": "\t", "T": "15:04:05",
	"Y": "2006", "y": "06",
}

func UnixToGolang(f string) (string, error) {
	newf := ""
	w := len(f)
	flag := false

	for i := 0; i < w; i++ {
		c := string(f[i])

		switch {
		case flag == true:
			if v, ok := flagsMap[c]; ok {
				newf += v
			} else {
				return "", fmt.Errorf("'%s: not support '%s'", f, "%"+c)
			}
			flag = false
		case c == "%":
			flag = true
		default:
			newf += c
		}
	}

	return newf, nil
}
