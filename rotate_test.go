package cronolog

import (
	"log"
	"os"
	"path"
	"path/filepath"
	"testing"
	"time"
)

var tmpdir = "./tmp"

func initTMPDir() {
	os.MkdirAll(tmpdir, 0666)
}

func removeTMPDir() error {
	return os.RemoveAll(tmpdir)
}

func countFiles(dir string) int {
	total := 0
	filepath.Walk(tmpdir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		total++
		return nil
	})

	return total
}

func TestRotate(t *testing.T) {
	initTMPDir()
	defer removeTMPDir()

	tmpFile := path.Join(tmpdir, "/t-20060102150405.log")
	period := "1s"

	rt, err := NewRotater(tmpFile, period, 0)
	if err != nil {
		t.Fatal(err)
	}
	defer rt.Close()
	log.SetOutput(rt)

	files := 3
	for i := 0; i < files; i++ {
		log.Println("rotate", i)
		time.Sleep(time.Second)
	}

	n := countFiles(tmpdir)
	if n != files {
		t.Fatalf("should be gernate %d files, but got %d files", files, n)
	}
}

func TestBacklogs(t *testing.T) {
	initTMPDir()
	defer removeTMPDir()

	tmpFile := path.Join(tmpdir, "/backlog-20060102150405.log")
	backlogs := 3
	period := "1s"
	r, err := NewRotater(tmpFile, period, backlogs)
	if err != nil {
		t.Fatal(err)
	}

	defer r.Close()
	log.SetOutput(r)

	times := backlogs * 2
	for i := 0; i < times; i++ {
		log.Println(i)
		time.Sleep(time.Second)
	}

	n := countFiles(tmpdir)
	if n != backlogs {
		t.Fatalf("should be gernate %d files, but got %d files", backlogs, n)
	}
}
