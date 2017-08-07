package cronolog

import (
	"errors"
	"fmt"
	"os"
	"sync"
)

type RotaterSize struct {
	lock      sync.Mutex
	fd        *os.File
	filename  string
	duplicate int
	wsize     int64
	rsize     int64
}

func fileIsExist(name string) bool {
	info, err := os.Stat(name)
	if err != nil {
		return os.IsExist(err)
	}

	if info.IsDir() {
		return false
	}

	return true
}

func fileSize(name string) int64 {
	info, err := os.Stat(name)
	if err != nil {
		return 0
	}

	if info.IsDir() {
		return 0
	}

	return info.Size()
}

func NewRotaterSize(name, size string, duplicate int) (r *RotaterSize, err error) {
	var (
		rsize Size
		fsize int64
		fd    *os.File
	)

	if duplicate < 0 {
		return nil, errors.New("duplicate can't less than 0")
	}

	rsize, err = ParseSize(size)
	if err != nil {
		return nil, err
	}

	fsize = fileSize(name)
	fd, err = os.OpenFile(name, fileflag, perm)
	if err != nil {
		return nil, err
	}

	return &RotaterSize{
		rsize:     int64(rsize),
		wsize:     fsize,
		filename:  name,
		duplicate: duplicate,
		fd:        fd,
	}, nil
}

func (r *RotaterSize) Write(p []byte) (n int, err error) {
	r.lock.Lock()
	defer r.lock.Unlock()

	n, err = r.fd.Write(p)
	if err != nil {
		return
	}

	r.wsize += int64(n)
	if r.wsize >= r.rsize {
		err = r.rotate()
		r.wsize = 0
	}

	return
}

func (r *RotaterSize) rotate() (err error) {
	origin := r.filename
	first := origin + ".1"

	for i := r.duplicate; i > 1; i-- {
		last := fmt.Sprintf("%s.%d", origin, i)
		previous := fmt.Sprintf("%s.%d", origin, i-1)
		if !fileIsExist(previous) {
			continue
		}

		if err := os.Rename(previous, last); err != nil {
			return err
		}
	}

	r.fd.Close()
	os.Rename(origin, first)

	var fd *os.File
	fd, err = os.OpenFile(origin, fileflag, perm)
	if err != nil {
		return err
	}

	r.fd = fd
	return err
}

func (r *RotaterSize) Close() (err error) {
	return r.fd.Close()
}
