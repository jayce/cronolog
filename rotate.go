package cronolog

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const (
	perm os.FileMode = 0664
)

type Rotater struct {
	lock   sync.Mutex
	period time.Duration
	timer  *time.Timer
	fd     *os.File

	layout   string
	lastname string // Last file name
	rotates  int    // Number of rotations
	isRotate bool   // Whether to rotate
}

// NewRotater create a rotater for time period
func NewRotater(layout, period string) (r *Rotater, err error) {
	var perd time.Duration

	perd, err = time.ParseDuration(period)
	if err != nil {
		return
	}

	if perd < 1*time.Second {
		err = fmt.Errorf("period less then %s", time.Second)
		return
	}

	date, since := alignTime(perd)

	file := date.Format(layout)
	if len(file) > 0 &&
		(file[len(file)-1] == '.' ||
			file[len(file)-1] == '/') {
		err = fmt.Errorf("'%s': Is a directory", file)
		return
	}

	since = perd - since

	r = &Rotater{
		timer: time.AfterFunc(since, func() {
			r.timer.Reset(r.period)

			r.lock.Lock()
			r.isRotate = true
			r.lock.Unlock()
		}),
		period:   perd,
		layout:   layout,
		isRotate: true,
	}

	return
}

func (r *Rotater) rotate() (err error) {
	date, _ := alignTime(r.period)
	name := date.Format(r.layout)
	if r.lastname == name {
		return
	}

	err = os.MkdirAll(filepath.Dir(name), perm)
	fmt.Println(filepath.Dir(name), err)
	if err != nil {
		return
	}

	var (
		newfd *os.File
		flag  = os.O_CREATE | os.O_WRONLY | os.O_APPEND
	)

	newfd, err = os.OpenFile(name, flag, perm)
	fmt.Println(err)
	if err != nil {
		return
	}

	r.lock.Lock()
	defer r.lock.Unlock()

	if r.fd != nil {
		r.fd.Close()
	}

	r.fd = newfd
	r.lastname = name
	r.isRotate = false
	r.rotates++

	return err
}

// Write: implements io.Writer
func (r *Rotater) Write(p []byte) (n int, err error) {
	if r.isRotate {
		if err = r.rotate(); err != nil {
			return
		}
	}
	return r.fd.Write(p)
}

// Close: implements io.Closer
func (r *Rotater) Close() error {
	if r.timer != nil {
		r.timer.Stop()
	}

	if r.fd != nil {
		return r.fd.Close()
	}

	return nil
}

func alignTime(period time.Duration) (time.Time, time.Duration) {
	date := time.Now()
	since := time.Duration(date.UnixNano()) % period
	return date.Add(-since), since
}
