package cronolog

import (
	"fmt"
	"testing"
)

func TestRotaterSize(t *testing.T) {
	initTMPDir()
	defer removeTMPDir()

	testFile := tmpdir + "/test.log"

	duplicate := 2
	fileSize := 1024
	buf := make([]byte, fileSize)

	rSize, err := NewRotaterSize(testFile, "1kb", duplicate)
	if err != nil {
		t.Fatal(err)
	}
	defer rSize.Close()

	for i := 0; i < duplicate; i++ {
		_, err = rSize.Write(buf)
		if err != nil {
			t.Error(err)
		}
	}

	for i := 1; i <= duplicate; i++ {
		name := fmt.Sprintf("%s.%d", testFile, i)
		if !fileIsExist(name) {
			t.Errorf("%s should be created, but isn't exist.", name)
		}
	}
}
