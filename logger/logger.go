package logger

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"sync/atomic"
	"time"
)

var (
	colors map[string]string
	logNo  uint64
)

const (
	Black = (iota + 30)
	Red
	Green
	Yellow
	Blue
	Magenta
	Cyan
	White
)

type Worker struct {
	Minion *log.Logger
	Color  int
}

type Info struct {
	Id      uint64
	Time    string
	File    string
	Line    int
	Module  string
	Level   string
	Message string
	format  string
}

type Logger struct {
	Module string
	worker *Worker
}

func (r *Info) Output() string {
	msg := fmt.Sprintf(r.format, r.Id, r.Time, r.File, r.Line, r.Level, r.Message)
	return msg
}

func NewWorker(prefix string, flag int, color int, out io.Writer) *Worker {
	return &Worker{Minion: log.New(out, prefix, flag), Color: color}
}

func (w *Worker) Log(level string, calldepth int, info *Info) error {
	if w.Color != 0 {
		buf := &bytes.Buffer{}
		buf.Write([]byte(colors[level]))
		buf.Write([]byte(info.Output()))
		buf.Write([]byte("\033[0m"))
		return w.Minion.Output(calldepth+1, buf.String())
	} else {
		return w.Minion.Output(calldepth+1, info.Output())
	}
}

func colorString(color int) string {
	return fmt.Sprintf("\033[%dm", int(color))
}

func initColors() {
	colors = map[string]string{
		"PANIC":   colorString(Magenta),
		"ERROR":   colorString(Red),
		"WARNING": colorString(Yellow),
		"INFO":    colorString(Green),
		"NOTICE":  colorString(White),
		"DEBUG":   colorString(Cyan),
		"TRACE":   colorString(Blue),
	}
}

func NewLogger(args ...interface{}) *Logger {
	initColors()

	var module string = "DEFAULT"
	var color int = 1
	var out io.Writer = os.Stderr

	for _, arg := range args {
		switch t := arg.(type) {
		case string:
			module = t
		case int:
			color = t
		case io.Writer:
			out = t
		default:
			panic("logger: Unknown argument")
		}
	}
	newWorker := NewWorker("", 0, color, out)
	return &Logger{Module: module, worker: newWorker}
}

func (l *Logger) Log(lvl string, message string) {
	_, file, line, _ := runtime.Caller(2)

	var formatString string = "#%d %s %s:%d ▶ %.5s %s"
	info := &Info{
		Id:      atomic.AddUint64(&logNo, 1),
		Time:    time.Now().Format("2006-01-02 15:04:05"),
		File:    file,
		Line:    line,
		Module:  l.Module,
		Level:   lvl,
		Message: message,
		format:  formatString,
	}
	l.worker.Log(lvl, 2, info)
}

func (l *Logger) Panic(args ...interface{}) {
	l.Log("PANIC", fmt.Sprint(args...))
	panic(args)
	os.Exit(1)
}

func (l *Logger) Error(args ...interface{}) {
	l.Log("ERROR", fmt.Sprint(args...))
}

func (l *Logger) Warn(args ...interface{}) {
	l.Log("WARNING", fmt.Sprint(args...))
}

func (l *Logger) Info(args ...interface{}) {
	l.Log("INFO", fmt.Sprint(args...))
}

func (l *Logger) Debug(args ...interface{}) {
	l.Log("DEBUG", fmt.Sprint(args...))
}

func (l *Logger) Trace(args ...interface{}) {
	l.Log("TRACE", fmt.Sprint(args...))
}

func (l *Logger) Notice(args ...interface{}) {
	l.Log("NOTICE", fmt.Sprint(args...))
}

func (l *Logger) Panicf(format string, args ...interface{}) {
	l.Log("PANIC", fmt.Sprintf(format, args...))
	panic(fmt.Sprintf(format, args...))
	os.Exit(1)
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	l.Log("ERROR", fmt.Sprintf(format, args...))
}

func (l *Logger) Warnf(format string, args ...interface{}) {
	l.Log("WARNING", fmt.Sprintf(format, args...))
}

func (l *Logger) Infof(format string, args ...interface{}) {
	l.Log("INFO", fmt.Sprintf(format, args...))
}

func (l *Logger) Debugf(format string, args ...interface{}) {
	l.Log("DEBUG", fmt.Sprintf(format, args...))
}

func (l *Logger) Tracef(format string, args ...interface{}) {
	l.Log("TRACE", fmt.Sprintf(format, args...))
}

func (l *Logger) Noticef(format string, args ...interface{}) {
	l.Log("NOTICE", fmt.Sprintf(format, args...))
}
