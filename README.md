Cronolog
========

a simple logging lib for Go, with rotating file.

Feature
-------

- [x] Rotater: Spport part of the unix time-format, like: '%Y%m%d%H%M%S'.

Usage
-----

```
package main

import (
	"fmt"
	"os"

	clog "github.com/jayce/cronolog"
)

func main() {
	rotater, err := clog.NewRotater("cronolog/2006010215.log", "24h", 0)
	// or
	// rotater, err := clog.NewRotater("cronolog/%Y%m%d%H.log", "24h", 0)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	logger := clog.NewLogger(rotater, "", clog.LstdFlags)
	logger.SetLevel(clog.LWarn)

	logger.Debug("debug infomation")
	logger.Warn("Warn infomation")
	logger.Error("Error infomation")

	// output:
	// 2017/10/24 01:00:31 main.go:19 [warn] "Warn infomation"
	// 2017/10/24 01:00:31 main.go:20 [error] "Warn infomation"

	plog := clog.NewScope("Cache")

	plog.Warn("Warn infomation")

	// output:
	// 2017/10/24 01:00:31 main.go:20 [warn] Cache: Warn infomation
}
```