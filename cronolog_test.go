package cronolog

import (
	"bytes"
	"io"
	"io/ioutil"
	"testing"
)

var buf = bytes.NewBuffer(make([]byte, 4096))
var blog = NewLogger(buf, 0)

func checkBuf(n int) (bool, int) {
	for {
		_, err := buf.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
		}
		n--
	}

	return n == 0, n
}

func TestLogger(t *testing.T) {
	buf.Reset()
	blog.SetLevel(LError)

	blog.Debug("test debug")
	blog.Warn("test warn")
	blog.Error("test error")

	if ok, n := checkBuf(1); !ok {
		t.Errorf("There shold be 1 messages, but found %d", n)
	}

	buf.Reset()
	blog.SetLevel(LDebug)

	blog.Debug("test debug")
	blog.Warn("test warn")
	blog.Error("test error")

	if ok, n := checkBuf(3); !ok {
		t.Errorf("There shold be 3 messages, but found %d", n)
	}
}

func BenchmarkLogger(b *testing.B) {
	clog := NewLogger(ioutil.Discard, Lshortfile|LstdFlags)
	for i := 0; i < b.N; i++ {
		clog.Debug("test", "teste", "test1")
	}
}
