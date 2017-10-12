package cronolog

import (
	"fmt"
	"io"
	"log"
	"os"
	"sync"
)

const (
	Ldate         = log.Ldate
	Ltime         = log.Ltime
	Lmicroseconds = log.Lmicroseconds
	Llongfile     = log.Llongfile
	Lshortfile    = log.Lshortfile
	LUTC          = log.LUTC
	LstdFlags     = log.LstdFlags
)

type Level int

const (
	LDebug Level = 0x0
	LWarn  Level = 0x1
	LError Level = 0x2
	LAlert Level = 0x3
)

func StringToLevel(level string) Level {
	switch level {
	case "alert":
		return LAlert
	case "error":
		return LError
	case "warn":
		return LWarn
	case "debug":
		return LDebug
	}
	return LDebug
}

func LevelToString(level Level) string {
	switch level {
	case LDebug:
		return "debug"
	case LWarn:
		return "warn"
	case LError:
		return "error"
	case LAlert:
		return "alert"
	}
	return "debug"
}

type logger struct {
	l     *log.Logger
	mux   sync.Mutex
	level Level
}

func NewLogger(out io.Writer, prefix string, flag int) *logger {
	return &logger{
		l:     log.New(out, prefix, flag),
		level: LDebug,
	}
}

func (l *logger) SetLevel(level Level) {
	l.mux.Lock()
	defer l.mux.Unlock()
	l.level = level
}

func (l *logger) SetOutput(w io.Writer) {
	l.l.SetOutput(w)
}

func (l *logger) Debug(v ...interface{}) {
	l.log(LDebug, v...)
}

func (l *logger) Warn(v ...interface{}) {
	l.log(LWarn, v...)
}

func (l *logger) Error(v ...interface{}) {
	l.log(LError, v...)
}

func (l *logger) Alert(v ...interface{}) {
	l.log(LAlert, v...)
}

func (l *logger) isOutput(level Level) bool {
	return l.level > level
}

func (l *logger) log(level Level, v ...interface{}) {
	if l.isOutput(level) {
		return
	}

	s := "[" + LevelToString(l.level) + "] " + fmt.Sprintln(v...)
	l.l.Output(4, s)
}

func (l *logger) Debugf(format string, v ...interface{}) {
	l.logf(LDebug, format, v...)
}

func (l *logger) Warnf(format string, v ...interface{}) {
	l.logf(LWarn, format, v...)
}

func (l *logger) Errorf(format string, v ...interface{}) {
	l.logf(LError, format, v...)
}

func (l *logger) Alertf(format string, v ...interface{}) {
	l.logf(LAlert, format, v...)
}

func (l *logger) logf(level Level, format string, v ...interface{}) {
	if l.isOutput(level) {
		return
	}

	s := "[" + LevelToString(l.level) + "] " + fmt.Sprintf(format, v...)
	l.l.Output(4, s)
}

var stderr = NewLogger(os.Stderr, "", log.LstdFlags)

func SetOutput(out io.Writer) {
	stderr.SetOutput(out)
}

func SetLevel(level Level) {
	stderr.SetLevel(level)
}

func Debug(v ...interface{}) {
	stderr.Debug(v...)
}

func Warn(v ...interface{}) {
	stderr.Warn(v...)
}

func Error(v ...interface{}) {
	stderr.Error(v...)
}

func Alert(v ...interface{}) {
	stderr.Alert(v...)
}

func Debugf(format string, v ...interface{}) {
	stderr.Debugf(format, v...)
}

func Warnf(format string, v ...interface{}) {
	stderr.Warnf(format, v...)
}

func Errorf(format string, v ...interface{}) {
	stderr.Errorf(format, v...)
}

func Alertf(format string, v ...interface{}) {
	stderr.Alertf(format, v...)
}
