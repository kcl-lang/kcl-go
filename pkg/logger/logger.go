// Copyright The KCL Authors. All rights reserved.

package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"sync/atomic"
)

var logger Logger = NewStdLogger(os.Stderr, "", "", 0)

// NewStdLogger create new logger based on std log.
// If level is empty string, use WARN as the default level.
// If flag is zore, use 'log.LstdFlags|log.Lshortfile' as the default flag.
// Level: DEBUG < INFO < WARN < ERROR < PANIC < FATAL
func NewStdLogger(out io.Writer, prefix, level string, flag int) Logger {
	return newStdLogger(out, prefix, level, flag)
}

func GetLogger() Logger {
	return logger
}

func SetLogger(new Logger) (old Logger) {
	old, logger = logger, new
	return
}

// Logger interface
//
// See https://github.com/chai2010/logger
type Logger interface {
	Debug(v ...any)
	Debugln(v ...any)
	Debugf(format string, v ...any)
	Info(v ...any)
	Infoln(v ...any)
	Infof(format string, v ...any)
	Warning(v ...any)
	Warningln(v ...any)
	Warningf(format string, v ...any)
	Error(v ...any)
	Errorln(v ...any)
	Errorf(format string, v ...any)
	Panic(v ...any)
	Panicln(v ...any)
	Panicf(format string, v ...any)
	Fatal(v ...any)
	Fatalln(v ...any)
	Fatalf(format string, v ...any)

	// Level: DEBUG < INFO < WARN < ERROR < PANIC < FATAL
	GetLevel() string
	SetLevel(new string) (old string)
}

type logLevelType uint32

const (
	logInvalidLevel logLevelType = iota // invalid
	logDebugLevel
	logInfoLevel
	logWarnLevel
	logErrorLevel
	logPanicLevel
	logFatalLevel
)

func (level logLevelType) Valid() bool {
	return level >= logDebugLevel && level <= logFatalLevel
}

func newLogLevel(name string) logLevelType {
	switch name {
	case "DEBUG":
		return logDebugLevel
	case "INFO":
		return logInfoLevel
	case "WARN":
		return logWarnLevel
	case "ERROR":
		return logErrorLevel
	case "PANIC":
		return logPanicLevel
	case "FATAL":
		return logFatalLevel
	}
	return logInvalidLevel
}

func (level logLevelType) String() string {
	switch level {
	case logDebugLevel:
		return "DEBUG"
	case logInfoLevel:
		return "INFO"
	case logWarnLevel:
		return "WARN"
	case logErrorLevel:
		return "ERROR"
	case logPanicLevel:
		return "PANIC"
	case logFatalLevel:
		return "FATAL"
	}
	return "INVALID"
}

type stdLogger struct {
	level logLevelType
	*log.Logger
}

func newStdLogger(out io.Writer, prefix, level string, flag int) *stdLogger {
	if flag == 0 {
		flag = log.LstdFlags | log.Lshortfile
	}
	if level == "" {
		level = "WARN"
	}

	p := &stdLogger{Logger: log.New(out, prefix, flag)}
	p.SetLevel(level)
	return p
}

func (p *stdLogger) getLevel() logLevelType {
	return logLevelType(atomic.LoadUint32((*uint32)(&p.level)))
}
func (p *stdLogger) setLevel(level logLevelType) logLevelType {
	return logLevelType(atomic.SwapUint32((*uint32)(&p.level), uint32(level)))
}

func (p *stdLogger) getLevelName() string {
	return p.getLevel().String()
}
func (p *stdLogger) setLevelByName(levelName string) string {
	level := newLogLevel(levelName)
	if !level.Valid() {
		panic("invalid level: " + levelName)
	}
	return p.setLevel(level).String()
}

func (p *stdLogger) GetLevel() string {
	return p.getLevel().String()
}
func (p *stdLogger) SetLevel(new string) (old string) {
	return p.setLevelByName(new)
}

func (p *stdLogger) Debug(v ...any) {
	if l := logDebugLevel; p.getLevel() <= l {
		p.Output(2, "["+l.String()+"] "+fmt.Sprint(v...))
	}
}
func (p *stdLogger) Debugln(v ...any) {
	if l := logDebugLevel; p.getLevel() <= l {
		p.Output(2, "["+l.String()+"] "+fmt.Sprintln(v...))
	}
}
func (p *stdLogger) Debugf(format string, v ...any) {
	if l := logDebugLevel; p.getLevel() <= l {
		p.Output(2, "["+l.String()+"] "+fmt.Sprintf(format, v...))
	}
}

func (p *stdLogger) Info(v ...any) {
	if l := logInfoLevel; p.getLevel() <= l {
		p.Output(2, "["+l.String()+"] "+fmt.Sprint(v...))
	}
}
func (p *stdLogger) Infoln(v ...any) {
	if l := logInfoLevel; p.getLevel() <= l {
		p.Output(2, "["+l.String()+"] "+fmt.Sprintln(v...))
	}
}
func (p *stdLogger) Infof(format string, v ...any) {
	if l := logInfoLevel; p.getLevel() <= l {
		p.Output(2, "["+l.String()+"] "+fmt.Sprintf(format, v...))
	}
}

func (p *stdLogger) Warning(v ...any) {
	if l := logWarnLevel; p.getLevel() <= l {
		p.Output(2, "["+l.String()+"] "+fmt.Sprint(v...))
	}
}
func (p *stdLogger) Warningln(v ...any) {
	if l := logWarnLevel; p.getLevel() <= l {
		p.Output(2, "["+l.String()+"] "+fmt.Sprintln(v...))
	}
}
func (p *stdLogger) Warningf(format string, v ...any) {
	if l := logWarnLevel; p.getLevel() <= l {
		p.Output(2, "["+l.String()+"] "+fmt.Sprintf(format, v...))
	}
}

func (p *stdLogger) Error(v ...any) {
	if l := logErrorLevel; p.getLevel() <= l {
		p.Output(2, "["+l.String()+"] "+fmt.Sprint(v...))
	}
}
func (p *stdLogger) Errorln(v ...any) {
	if l := logErrorLevel; p.getLevel() <= l {
		p.Output(2, "["+l.String()+"] "+fmt.Sprintln(v...))
	}
}
func (p *stdLogger) Errorf(format string, v ...any) {
	if l := logErrorLevel; p.getLevel() <= l {
		p.Output(2, "["+l.String()+"] "+fmt.Sprintf(format, v...))
	}
}

func (p *stdLogger) Panic(v ...any) {
	s := fmt.Sprint(v...)
	if l := logPanicLevel; p.getLevel() <= l {
		p.Output(2, "["+l.String()+"] "+s)
	}
	panic(s)
}
func (p *stdLogger) Panicln(v ...any) {
	s := fmt.Sprintln(v...)
	if l := logPanicLevel; p.getLevel() <= l {
		p.Output(2, "["+l.String()+"] "+s)
	}
	panic(s)
}
func (p *stdLogger) Panicf(format string, v ...any) {
	s := fmt.Sprintf(format, v...)
	if l := logPanicLevel; p.getLevel() <= l {
		p.Output(2, "["+l.String()+"] "+s)
	}
	panic(s)
}

func (p *stdLogger) Fatal(v ...any) {
	const l = logFatalLevel
	p.Output(2, "["+l.String()+"] "+fmt.Sprint(v...))
	os.Exit(1)
}
func (p *stdLogger) Fatalln(v ...any) {
	const l = logFatalLevel
	p.Output(2, "["+l.String()+"] "+fmt.Sprintln(v...))
	os.Exit(1)
}
func (p *stdLogger) Fatalf(format string, v ...any) {
	const l = logFatalLevel
	p.Output(2, "["+l.String()+"] "+fmt.Sprintf(format, v...))
	os.Exit(1)
}
