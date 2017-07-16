package cronolog

import (
	"os"
	"time"
)

const (
	TimeMinuteFormat = "200601021504"
	TimeHourFormat   = "2006010215"
	TimeDayFormat    = "20060102"
	DefaultFormat    = TimeDayFormat
)

const (
	rotateFiles = 3
)

type Rotater struct {
	duration *time.Duration
	fd       *os.File
}

func NewRotate(name, duration string) (r *Rotater, err error) {
	r = new(Rotater)

	r.duration, err = time.ParseDuration(duration)
	if err != nil {
		return nil, err
	}

	flag := os.O_CREATE | os.O_WRONLY | os.O_APPEND
	perm := 0664

	r.fd, err = os.OpenFile(name, flag, perm)
	return r, err
}

func (r *Rotater) rotate() {

}

func (r *Rotater) Write(p []byte) (n int, err error) {
	return r.fd.Write(p)
}

func (r *Rotate) Close() error {
	return r.fd.Close()
}
