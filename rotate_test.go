package cronolog

import (
	"log"
	"os"
	"path"
	"testing"
	"time"
)

var tmpdir = "./tmp"

func removeTMPDir() error {
	return os.RemoveAll(tmpdir)
}

func checkFiles(files int, layout, period string) bool {
	perd, _ := time.ParseDuration(period)
	date, _ := alignTime(perd)

	for i := files; i > 0; i-- {
		date = date.Add(-perd)
		name := date.Format(layout)
		_, err := os.Stat(name)
		if os.IsNotExist(err) {
			return false
		}
	}

	return true
}

func TestRotate(t *testing.T) {
	defer removeTMPDir()

	tmpFile := path.Join(tmpdir, "/t-20060102150405.log")
	period := "2s"

	rt, err := NewRotater(tmpFile, period)
	if err != nil {
		t.Fatal(err)
	}
	defer rt.Close()
	log.SetOutput(rt)

	files := 3
	for i := 0; i < files*2; i++ {
		log.Println("rotate", i)
		time.Sleep(1 * time.Second)
	}

	if !checkFiles(files, tmpFile, period) {
		t.Fatal("generate files not ok")
	}
}
