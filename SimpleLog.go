package SimpleLog

import (
	"fmt"
	"io"
	"os"
	"runtime/debug"
	"strings"
	"sync"
	"time"
)

type Level int

const (
	TraceLevel Level = iota
	DebugLevel
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
	PanicLevel
)

// 全局唯一的日志实例, 统一控制 level
type logger struct {
	*sync.Mutex
	Out   io.Writer
	Level Level
}

// 外部接口, 自定义某些选项
type Logger struct {
	*logger
	banner        string
	color         bool
	escapeNewline bool
}

var (
	LevelBannerN = map[Level]string{
		TraceLevel: "[TRACE]",
		DebugLevel: "[DEBUG]",
		InfoLevel:  " [INFO]",
		WarnLevel:  " [WARN]",
		ErrorLevel: "[ERROR]",
		FatalLevel: "[FATAL]",
		PanicLevel: "[PANIC]",
	}
	LevelBannerC = map[Level]string{
		TraceLevel: "\x1b[94m[TRACE]\x1b[m",
		DebugLevel: "\x1b[92m[DEBUG]\x1b[m",
		InfoLevel:  "\x1b[97m [INFO]\x1b[m",
		WarnLevel:  "\x1b[93m [WARN]\x1b[m",
		ErrorLevel: "\x1b[91m[ERROR]\x1b[m",
		FatalLevel: "\x1b[91;5m[FATAL]\x1b[m",
		PanicLevel: "\x1b[91;5;7m[PANIC]\x1b[m",
	}
)

var defaultLogger = &logger{
	Out:   os.Stderr,
	Mutex: new(sync.Mutex),
}

func New(banner string, color, escapeNewline bool) *Logger {
	return &Logger{defaultLogger, banner, color, escapeNewline}
}

func (l *Logger) AddOutput(w io.Writer) *Logger {
	l.Out = io.MultiWriter(l.Out, w)
	return l
}

func (l *Logger) SetOutput(w io.Writer) *Logger {
	l.Out = w
	return l
}

func (l *Logger) SetLevel(level Level) *Logger {
	l.Level = level
	return l
}

func (l *Logger) SetBanner(banner string) *Logger {
	if len(banner) > 0 && banner[0] != '[' {
		banner = "[" + banner
	}
	if len(banner) > 0 && banner[len(banner)-1] != ']' {
		banner = banner + "]"
	}
	l.banner = banner
	return l
}

func (l *Logger) SetEscapeNewline(escape bool) *Logger {
	l.escapeNewline = escape
	return l
}

var (
	lastLogoutMonth int // 新的一月时输出一次带月份的日志
	lastLogoutDay   int // 新的一天时输出一次带日期的日志
)

func (l *logger) formatTime() string {
	t := time.Now()
	month, day := int(t.Month()), t.Day()
	defer func() {
		lastLogoutMonth, lastLogoutDay = month, day
	}()
	if month != lastLogoutMonth {
		return t.Format("[15:04-|01/02]")
	} else if day != lastLogoutDay {
		return t.Format("[15:04:05-|02]")
	} else {
		return t.Format("[15:04:05.000]")
	}
}

var newLineReplacer = strings.NewReplacer("\n", "\x1b[97m\\n\x1b[m")

func (l *Logger) Format(level Level, s string) string {
	if l.escapeNewline {
		s = newLineReplacer.Replace(s)
	}
	var lvl string
	if l.color {
		lvl = LevelBannerC[level]
	} else {
		lvl = LevelBannerN[level]
	}
	t := l.formatTime()
	sb := new(strings.Builder)
	sb.Grow(len(lvl) + len(t) + len(l.banner) + len(s) + 2)
	sb.WriteString(lvl)
	sb.WriteString(t)
	sb.WriteString(l.banner)
	sb.WriteByte(' ')
	sb.WriteString(s)
	sb.WriteByte('\n')
	return sb.String()
}

func (l *Logger) Output(s string) {
	l.Lock()
	defer l.Unlock()
	l.Out.Write([]byte(s))
}

func (l *Logger) Print(level Level, a ...any) {
	l.Output(l.Format(level, fmt.Sprint(a...)))
}

func (l *Logger) Printf(level Level, format string, a ...any) {
	l.Output(l.Format(level, fmt.Sprintf(format, a...)))
}

func (l *Logger) levelOk(level Level) bool {
	return level >= l.Level // 大于等于则输出
}

func (l *Logger) Trace(a ...any) {
	if !l.levelOk(TraceLevel) {
		return
	}
	l.Print(TraceLevel, a...)
}

func (l *Logger) Tracef(format string, a ...any) {
	if !l.levelOk(TraceLevel) {
		return
	}
	l.Printf(TraceLevel, format, a...)
}

func (l *Logger) Debug(a ...any) {
	if !l.levelOk(DebugLevel) {
		return
	}
	l.Print(DebugLevel, a...)
}

func (l *Logger) Debugf(format string, a ...any) {
	if !l.levelOk(DebugLevel) {
		return
	}
	l.Printf(DebugLevel, format, a...)
}

func (l *Logger) Info(a ...any) {
	if !l.levelOk(InfoLevel) {
		return
	}
	l.Print(InfoLevel, a...)
}

func (l *Logger) Infof(format string, a ...any) {
	if !l.levelOk(InfoLevel) {
		return
	}
	l.Printf(InfoLevel, format, a...)
}

func (l *Logger) Warn(a ...any) {
	if !l.levelOk(WarnLevel) {
		return
	}
	l.Print(WarnLevel, a...)
}

func (l *Logger) Warnf(format string, a ...any) {
	if !l.levelOk(WarnLevel) {
		return
	}
	l.Printf(WarnLevel, format, a...)
}

func (l *Logger) Error(a ...any) {
	if !l.levelOk(ErrorLevel) {
		return
	}
	l.Print(ErrorLevel, a...)
}

func (l *Logger) Errorf(format string, a ...any) {
	if !l.levelOk(ErrorLevel) {
		return
	}
	l.Printf(ErrorLevel, format, a...)
}

func (l *Logger) Fatal(a ...any) {
	if !l.levelOk(FatalLevel) {
		return
	}
	l.Print(FatalLevel, a...)
	os.Exit(1)
}

func (l *Logger) Fatalf(format string, a ...any) {
	if !l.levelOk(FatalLevel) {
		return
	}
	l.Printf(FatalLevel, format, a...)
	os.Exit(1)
}

func (l *Logger) Panic(a ...any) {
	if !l.levelOk(PanicLevel) {
		return
	}
	l.Print(PanicLevel, a...)
	panic(fmt.Sprint(a...))
}

func (l *Logger) Panicf(format string, a ...any) {
	if !l.levelOk(PanicLevel) {
		return
	}
	l.Printf(PanicLevel, format, a...)
	panic(fmt.Sprintf(format, a...))
}

// FakePanic only print stack
func (l *Logger) FakePanic(a ...any) {
	if !l.levelOk(PanicLevel) {
		return
	}
	l.Print(PanicLevel, a...)
	l.Output(string(debug.Stack()))
}

// FakePanic only print stack
func (l *Logger) FakePanicf(format string, a ...any) {
	if !l.levelOk(PanicLevel) {
		return
	}
	l.Printf(PanicLevel, format, a...)
	l.Output(string(debug.Stack()))
}
