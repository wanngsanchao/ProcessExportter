package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

// 定义日志级别
type Level logrus.Level

var (
	accessinstance *logrus.Logger
	errorinstance  *logrus.Logger
	once           sync.Once
)

var (
	accesslog string = "app-access.log"
	errorlog  string = "app-error.log"
)

const (
	DebugLevel Level = iota // 0
	InfoLevel               // 1
	WarnLevel               // 2
	ErrorLevel              // 3
	FatalLevel              // 4
	PanicLevel              // 5
)

// 颜色代码
const (
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorPurple = "\033[35m"
	colorCyan   = "\033[36m"
	colorReset  = "\033[0m"
)

// 初始化日志配置
func InitLogger(logFile string, maxAge time.Duration, rotationTime time.Duration, level Level) (*Logrus.Logger, error) {

	once.Do(func() {
		// 创建日志实例
		instance = logrus.New()

		// 设置日志级别
		instance.SetLevel(logrus.Level(level))

		// 设置日志格式
		instance.SetFormatter(&logrus.TextFormatter{
			ForceColors:     true,
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05",
			CallerPrettyfier: func(f *runtime.Frame) (string, string) {
				filename := filepath.Base(f.File)
				return "", fmt.Sprintf(" [%s:%d]", filename, f.Line)
			},
		})

		// 同时输出到终端和文件
		multiWriter := io.MultiWriter(os.Stdout, getRotatedLogWriter(logFile, maxAge, rotationTime))
		instance.SetOutput(multiWriter)

		// 显示调用者信息
		instance.SetReportCaller(true)

		// 设置颜色钩子
		setColorHook()
	})

	if instance == nil {
		err = errors.New("初始化日志失败")
	}

	return instance, nil
}

// 获取日志轮转写入器
func getRotatedLogWriter(logFile string, maxAge time.Duration, rotationTime time.Duration) io.Writer {
	// 日志文件格式: app.log.20230101
	writer, err := rotatelogs.New(
		logFile+".%Y%m%d",
		rotatelogs.WithLinkName(logFile),          // 生成软链，指向最新日志文件
		rotatelogs.WithMaxAge(maxAge),             // 文件最大保存时间
		rotatelogs.WithRotationTime(rotationTime), // 日志切割时间间隔
	)

	if err != nil {
		fmt.Printf("日志轮转初始化失败: %v\n", err)
		return os.Stdout
	}

	return writer
}

// 设置颜色钩子
func setColorHook() {
	instance.AddHook(lfshook.NewHook(lfshook.WriterMap{
		logrus.DebugLevel: os.Stdout,
		logrus.InfoLevel:  os.Stdout,
		logrus.WarnLevel:  os.Stdout,
		logrus.ErrorLevel: os.Stdout,
		logrus.FatalLevel: os.Stdout,
		logrus.PanicLevel: os.Stdout,
	}, &logrus.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			filename := filepath.Base(f.File)
			return "", fmt.Sprintf(" [%s:%d]", filename, f.Line)
		},
		// 自定义日志级别颜色
		DisableLevelTruncation: true,
		Formatter: func(entry *logrus.Entry) ([]byte, error) {
			var levelColor string
			switch entry.Level {
			case logrus.DebugLevel:
				levelColor = colorBlue
			case logrus.InfoLevel:
				levelColor = colorGreen
			case logrus.WarnLevel:
				levelColor = colorYellow
			case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
				levelColor = colorRed
			}

			levelText := strings.ToUpper(entry.Level.String())
			msg := fmt.Sprintf(
				"%s %s%s%s %s%s\n",
				entry.Time.Format("2006-01-02 15:04:05"),
				levelColor, levelText, colorReset,
				entry.Message,
				entry.Caller,
			)
			return []byte(msg), nil
		},
	}))
}

func Logger(Level logrus.Level, msg string, data interface{}) {

}
