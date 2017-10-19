Cronolog
========

a simple logging lib for Go, with rotating file.

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

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	logger := clog.NewLogger(rotater, "", clog.LstdFlags)
	logger.SetLevel(clog.LError)

	logger.Debug("debug infomation")
	logger.Warn("Warn infomation")
	logger.Error("Error infomation")
}
```