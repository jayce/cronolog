package cronolog

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const (
	perm     os.FileMode = 0664
	fileflag int         = os.O_CREATE | os.O_WRONLY | os.O_APPEND
)

type Rotater struct {
	lock   sync.Mutex
	period time.Duration
	timer  *time.Timer
	fd     *os.File

	layout   string
	lastname string
	backlogs int
	isRotate bool // Whether to rotate
}

// NewRotater create a rotater for time period
func NewRotater(layout, period string, backlogs int) (r *Rotater, err error) {
	var perd time.Duration

	perd, err = time.ParseDuration(period)
	if err != nil {
		return
	}

	if perd < time.Second {
		err = fmt.Errorf("period must be great then %s", time.Second)
		return
	}

	if backlogs < 0 {
		err = fmt.Errorf("backlogs must be great then 0")
		return
	}

	date, since := alignTime(perd)

	file, err := UnixToGolang(layout, date)
	if err != nil {
		return
	}

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
		backlogs: backlogs,
	}

	return
}

func (r *Rotater) rotate() (err error) {
	r.isRotate = false

	date, _ := alignTime(r.period)
	name, _ := UnixToGolang(r.layout, date)
	if r.lastname == name {
		return
	}

	err = os.MkdirAll(filepath.Dir(name), perm)
	if err != nil {
		return
	}

	var (
		newfd *os.File
		flag  = os.O_CREATE | os.O_WRONLY | os.O_APPEND
	)

	newfd, err = os.OpenFile(name, flag, perm)
	if err != nil {
		return
	}

	if r.fd != nil {
		r.fd.Close()
	}

	r.fd = newfd
	r.lastname = name

	return r.deleteBacklog(r.backlogs)
}

func (r *Rotater) deleteBacklog(backlog int) error {
	if backlog == 0 {
		return nil
	}

	date, _ := alignTime(r.period)
	since := r.period * time.Duration(backlog)
	past := date.Add(-since)
	name, _ := UnixToGolang(r.layout, past)

	if r.lastname == name {
		return nil
	}

	if _, err := os.Stat(name); !os.IsNotExist(err) {
		return os.Remove(name)
	}
	return nil
}

// Write: implements io.Writer
func (r *Rotater) Write(p []byte) (n int, err error) {
	r.lock.Lock()
	defer r.lock.Unlock()

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
	_, offset := date.Zone()
	since := time.Duration(date.Unix()+int64(offset)) * time.Second % period
	return date.Add(-since), period - since
}
