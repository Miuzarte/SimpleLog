# SimpleLog

简单到有点简陋的日志库

## 安装

```sh
go get -u github.com/yourusername/SimpleLog
```

## 使用示例

```go
package main

import (
    slog "github.com/yourusername/SimpleLog"
    "os"
)

func main() {
    logger := slog.New("[module1]", true, true)
    logger.SetLevel(slog.DebugLevel)
    logger.Info("...")

    logger2 := slog.New("[module2]", true, true)
    // debug level too
    logger2.Infof("...")
}
```

## 功能

- 所有日志实例共享相同的 `level` 和 `output`, 统一控制
- 除此之外 `banner`, `color`, `escapeNewline` 可单独设置
- 日志级别：Trace, Debug, Info, Warn, Error, Fatal, Panic
- 较短的日期格式化, 彩色日志级别标题
- 日志换行符转义

## 接口

### New

创建一个新的日志实例

```go
func New(banner string, color, escapeNewline bool) *Logger
```

### SetOutput

设置日志输出

```go
func (l *Logger) SetOutput(w io.Writer) *Logger
func (l *Logger) AddOutput(w io.Writer) *Logger
```

### SetLevel

设置日志级别

```go
func (l *Logger) SetLevel(level Level) *Logger
```

### SetBanner

设置日志前缀

```go
func (l *Logger) SetBanner(banner string) *Logger
```

### SetEscapeNewline

设置是否转义换行符

```go
func (l *Logger) SetEscapeNewline(escape bool) *Logger
```

### 日志方法

- `Trace(a ...any)`
- `Debug(a ...any)`
- `Info(a ...any)`
- `Warn(a ...any)`
- `Error(a ...any)`
- `Fatal(a ...any)`
- `Panic(a ...any)`
- `FakePanic(a ...any)`

每个方法都有对应的格式化版本，如 `Tracef(format string, a ...any)`
