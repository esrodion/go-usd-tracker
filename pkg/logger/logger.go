package logger

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	LevelDebug   = "DEBUG"
	LevelInfo    = "INFO"
	LevelWarning = "WARNING"
	LevelError   = "ERROR"
)

var (
	logger      *zap.Logger
	once        sync.Once
	atomicLevel zap.AtomicLevel
	logOpenTime time.Time
	logLifetime time.Duration
	outFile     *os.File
	outFileM    sync.Mutex
)

func BuildLogger(logLevel string) {
	once.Do(func() {
		atomicLevel = zap.NewAtomicLevel()
		err := SetLevel(logLevel)

		var encoderCfg zapcore.EncoderConfig
		if os.Getenv("LOG_EXTENDED") == "true" {
			encoderCfg = zap.NewProductionEncoderConfig()
		} else {
			encoderCfg = zap.NewDevelopmentEncoderConfig()
			encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
		}

		var encoder zapcore.Encoder
		if os.Getenv("LOG_TO_JSON") == "true" {
			encoder = zapcore.NewJSONEncoder(encoderCfg)
		} else {
			encoder = zapcore.NewConsoleEncoder(encoderCfg)
		}

		var opts []zap.Option
		if os.Getenv("LOG_TO_FILE") != "true" {
			outFile = os.Stdout
		} else if outFile, err = openLogFile(); err == nil {
			n, _ := strconv.Atoi(os.Getenv("LOG_FILE_LIFETIME"))
			if n > 0 {
				logLifetime = time.Duration(n) * time.Second
				fileReopener := func(core zapcore.Core) zapcore.Core {
					return &zapcoreFileReopen{Core: core}
				}
				opts = append(opts, zap.WrapCore(fileReopener))
			}
		}

		logger = zap.New(zapcore.NewCore(encoder, outFile, atomicLevel), opts...)

		if err != nil {
			logger.Debug(err.Error())
		}
	})
}

func SetLevel(logLevel string) error {
	switch strings.ToUpper(logLevel) {
	case LevelDebug:
		atomicLevel.SetLevel(zapcore.DebugLevel)
	case LevelInfo:
		atomicLevel.SetLevel(zapcore.InfoLevel)
	case LevelWarning:
		atomicLevel.SetLevel(zapcore.WarnLevel)
	case LevelError:
		atomicLevel.SetLevel(zapcore.ErrorLevel)
	default:
		atomicLevel.SetLevel(zapcore.DebugLevel)
	}

	return nil
}

func CurrentLevel() string {
	return atomicLevel.String()
}

// Logger returns a global logger defined in this package.
// If logger is nil function returns a logger with level specified in LOBBY_LOG_LEVEL env variable.
func Logger() *zap.Logger {
	if logger == nil {
		BuildLogger(os.Getenv("LOG_LEVEL"))
	}
	return logger
}

// // zapcore.Core wrapper for outdated file reopenning
type zapcoreFileReopen struct {
	zapcore.Core
}

func (z *zapcoreFileReopen) CheckFile() {
	if time.Since(logOpenTime) > logLifetime {
		outFileM.Lock()
		defer outFileM.Unlock()
		if time.Since(logOpenTime) < logLifetime {
			return
		}
		var err error
		var file *os.File
		outFile.Close()
		file, err = openLogFile()
		if err != nil {
			logOpenTime = time.Now()
			file = os.Stdout
			logger.Error("Could not open file specified in LOG_FILE environment varibale, exporting logs to os.Stdout: " + err.Error())
		}

		*outFile = *file
	}
}

func (z *zapcoreFileReopen) Check(e zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	z.CheckFile()
	return z.Core.Check(e, ce)
}

func openLogFile() (*os.File, error) {
	var out *os.File = os.Stdout
	var file *os.File
	var err error

	logOpenTime = time.Now()
	fileName := os.Getenv("LOG_FILE")
	if len(fileName) > 0 {
		fileName = fmt.Sprintf(fileName, logOpenTime.Format(os.Getenv("LOG_FILE_TIME_FORMAT")))
		file, err = os.OpenFile(fileName, os.O_APPEND|os.O_CREATE, os.ModePerm)
		if err == nil {
			out = file
		}
	}

	return out, err
}
