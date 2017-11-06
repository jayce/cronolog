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
	LDebug = iota + Level(0)
	Linfo
	LWarn
	LError
	LAlert
)

func StringToLevel(level string) Level {
	switch level {
	case "alert":
		return LAlert
	case "error":
		return LError
	case "warn":
		return LWarn
	case "info":
		return Linfo
	case "debug":
		return LDebug
	}
	return LDebug
}

func LevelToString(level Level) string {
	switch level {
	case LDebug:
		return "debug"
	case Linfo:
		return "info"
	case LWarn:
		return "warn"
	case LError:
		return "error"
	case LAlert:
		return "alert"
	}
	return "debug"
}

const calldep = 2

type logger struct {
	l     *log.Logger
	mux   sync.Mutex
	level Level
}

func NewLogger(out io.Writer, flags int) *logger {
	return &logger{
		l:     log.New(out, "", flags),
		level: Linfo,
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

func (l *logger) SetFlags(flags int) {
	l.l.SetFlags(flags)
}

func (l *logger) Debug(v ...interface{}) {
	l.log(calldep, LDebug, v...)
}

func (l *logger) Info(v ...interface{}) {
	l.log(calldep, Linfo, v...)
}

func (l *logger) Warn(v ...interface{}) {
	l.log(calldep, LWarn, v...)
}

func (l *logger) Error(v ...interface{}) {
	l.log(calldep, LError, v...)
}

func (l *logger) Alert(v ...interface{}) {
	l.log(calldep, LAlert, v...)
}

func (l *logger) isOutput(level Level) bool {
	return l.level > level
}

func (l *logger) log(calldep int, level Level, v ...interface{}) {
	if l.isOutput(level) {
		return
	}

	buf := fmt.Sprintf("[%s] ", LevelToString(level)) + fmt.Sprintln(v...)
	l.l.Output(calldep+1, buf)
}

func (l *logger) logf(calldep int, level Level, format string, v ...interface{}) {
	if l.isOutput(level) {
		return
	}

	l.log(calldep+1, level, fmt.Sprintf(format, v...))
}

func (l *logger) Debugf(format string, v ...interface{}) {
	l.logf(calldep, LDebug, format, v...)
}

func (l *logger) Infof(format string, v ...interface{}) {
	l.logf(calldep, Linfo, format, v...)
}

func (l *logger) Warnf(format string, v ...interface{}) {
	l.logf(calldep, LWarn, format, v...)
}

func (l *logger) Errorf(format string, v ...interface{}) {
	l.logf(calldep, LError, format, v...)
}

func (l *logger) Alertf(format string, v ...interface{}) {
	l.logf(calldep, LAlert, format, v...)
}

func (l *logger) NewScope(s string) *scope {
	return &scope{
		l: l,
		s: s,
	}
}

type scope struct {
	l *logger
	s string
}

func (s *scope) log(level Level, v ...interface{}) {
	if s.l.isOutput(level) {
		return
	}

	buf := fmt.Sprintf("[%s] %s: ", LevelToString(level), s.s) + fmt.Sprintln(v...)
	s.l.l.Output(calldep+1, buf)
}

func (s *scope) Debug(v ...interface{}) {
	s.log(LDebug, fmt.Sprint(v...))
}

func (s *scope) Info(v ...interface{}) {
	s.log(Linfo, fmt.Sprint(v...))
}

func (s *scope) Warn(v ...interface{}) {
	s.log(LWarn, fmt.Sprint(v...))
}

func (s *scope) Error(v ...interface{}) {
	s.log(LError, fmt.Sprint(v...))
}

func (s *scope) Alert(v ...interface{}) {
	s.log(LAlert, fmt.Sprint(v...))
}

func (s *scope) Debugf(format string, v ...interface{}) {
	s.log(LDebug, fmt.Sprintf(format, v...))
}

func (s *scope) Infof(format string, v ...interface{}) {
	s.log(Linfo, fmt.Sprintf(format, v...))
}

func (s *scope) Warnf(format string, v ...interface{}) {
	s.log(LWarn, fmt.Sprintf(format, v...))
}

func (s *scope) Errorf(format string, v ...interface{}) {
	s.log(LError, fmt.Sprintf(format, v...))
}

func (s *scope) Alertf(format string, v ...interface{}) {
	s.log(LAlert, fmt.Sprintf(format, v...))
}

var stderr = NewLogger(os.Stderr, log.LstdFlags)

func NewScope(scope string) *scope {
	return stderr.NewScope(scope)
}

func SetOutput(out io.Writer) {
	stderr.SetOutput(out)
}

func SetLevel(level Level) {
	stderr.SetLevel(level)
}

func SetFlags(flags int) {
	stderr.SetFlags(flags)
}

func Debug(v ...interface{}) {
	stderr.log(calldep, LDebug, v...)
}

func Info(v ...interface{}) {
	stderr.log(calldep, Linfo, v...)
}

func Warn(v ...interface{}) {
	stderr.log(calldep, LWarn, v...)
}

func Error(v ...interface{}) {
	stderr.log(calldep, LError, v...)
}

func Alert(v ...interface{}) {
	stderr.log(calldep, LAlert, v...)
}

func Debugf(format string, v ...interface{}) {
	stderr.logf(calldep, LDebug, format, v...)
}

func Infof(format string, v ...interface{}) {
	stderr.logf(calldep, Linfo, format, v...)
}

func Warnf(format string, v ...interface{}) {
	stderr.logf(calldep, LWarn, format, v...)
}

func Errorf(format string, v ...interface{}) {
	stderr.logf(calldep, LError, format, v...)
}

func Alertf(format string, v ...interface{}) {
	stderr.logf(calldep, LAlert, format, v...)
}
