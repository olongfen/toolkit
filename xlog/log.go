package xlog

import (
	"context"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
	"os"
	"time"
)

type Config struct {
	InfoFile   string
	ErrorFile  string
	MaxSize    int
	MaxBackups int
	MaxAge     int
	Compress   bool
}

func encodeJSON() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	return zapcore.NewJSONEncoder(encoderConfig)
}

func writer(isInfo bool, conf Config) zapcore.WriteSyncer {
	var (
		lofFile string
	)
	if isInfo {
		lofFile = conf.InfoFile
	} else {
		lofFile = conf.ErrorFile
	}
	lumberJackLogger := &lumberjack.Logger{
		Filename:   lofFile,
		MaxSize:    conf.MaxSize,    // 在进行切割之前，日志文件的最大大小（以MB为单位）
		MaxBackups: conf.MaxBackups, // 保留旧文件的最大个数
		MaxAge:     conf.MaxAge,     // 保留旧文件的最大天数
		Compress:   conf.Compress,   // 是否压缩/归档旧文件
	}
	return zapcore.AddSync(lumberJackLogger)
}

func NewProduceLogger(configs ...Config) *zap.Logger {
	var (
		conf Config
	)
	if len(configs) != 0 {
		conf = configs[0]
	}
	if conf.InfoFile == "" {
		conf.InfoFile = "./logs/server.log"
	}
	if conf.ErrorFile == "" {
		conf.ErrorFile = "./logs/error.log"
	}
	if conf.MaxSize == 0 {
		conf.MaxSize = 200
	}
	if conf.MaxAge == 0 {
		conf.MaxAge = 60
	}
	if conf.MaxBackups == 0 {
		conf.MaxBackups = 50
	}
	core := zapcore.NewCore(encodeJSON(), zapcore.NewMultiWriteSyncer(writer(true, conf), os.Stdout), zap.LevelEnablerFunc(func(level zapcore.Level) bool {
		return level < zap.ErrorLevel && level >= zap.DebugLevel
	}))
	coreErr := zapcore.NewCore(encodeJSON(), zapcore.NewMultiWriteSyncer(writer(false, conf), os.Stdout), zap.LevelEnablerFunc(func(level zapcore.Level) bool {
		return level >= zap.ErrorLevel
	}))
	return zap.New(zapcore.NewTee(core, coreErr), zap.AddCaller())
}

func NewDevelopment() *zap.Logger {

	cfg := zap.NewProductionEncoderConfig()
	cfg.EncodeLevel = zapcore.LowercaseColorLevelEncoder
	cfg.EncodeTime = zapcore.ISO8601TimeEncoder
	cfg1 := zap.NewProductionConfig()
	cfg1.EncoderConfig = cfg
	cfg1.Encoding = "console"

	log, _ := cfg1.Build()
	return log
}

type DBLog struct {
	*zap.Logger
	LogLevel                            zapcore.Level
	SlowThreshold                       time.Duration
	IgnoreRecordNotFoundError           bool
	traceStr, traceErrStr, traceWarnStr string
}

func (l *DBLog) LogMode(level logger.LogLevel) logger.Interface {
	newlogger := *l
	newlogger.LogLevel = zapcore.Level(level)
	return &newlogger
}

func (l *DBLog) Info(ctx context.Context, s string, i ...interface{}) {
	l.Sugar().Info(s, i)
}

func (l *DBLog) Warn(ctx context.Context, s string, i ...interface{}) {
	l.Sugar().Warn(s, i)
}

func (l *DBLog) Error(ctx context.Context, s string, i ...interface{}) {
	l.Sugar().Error()
}

func (l *DBLog) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if l.LogLevel <= zapcore.DebugLevel {
		return
	}

	elapsed := time.Since(begin)
	switch {
	case err != nil && l.LogLevel >= zapcore.ErrorLevel && (!errors.Is(err, gorm.ErrRecordNotFound) || !l.IgnoreRecordNotFoundError):
		sql, rows := fc()
		if rows == -1 {
			l.Sugar().Infof(l.traceErrStr, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			l.Sugar().Infof(l.traceErrStr, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	case elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.LogLevel >= zap.WarnLevel:
		sql, rows := fc()
		slowLog := fmt.Sprintf("SLOW SQL >= %v", l.SlowThreshold)
		if rows == -1 {
			l.Sugar().Infof(l.traceWarnStr, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			l.Sugar().Infof(l.traceWarnStr, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	case l.LogLevel == zap.InfoLevel:
		sql, rows := fc()
		if rows == -1 {
			l.Sugar().Infof(l.traceStr, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			l.Sugar().Infof(l.traceStr, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	}
}

func NewDBLog(zapLog *zap.Logger) logger.Interface {
	var (
		traceStr     = `%s [%.3fms] [rows:%v] %s`
		traceWarnStr = `%s %s[%.3fms] [rows:%v] %s`
		traceErrStr  = `%s %s[%.3fms] [rows:%v] %s`
	)
	return &DBLog{
		Logger:                    zapLog,
		IgnoreRecordNotFoundError: false,
		traceStr:                  traceStr,
		traceWarnStr:              traceWarnStr,
		traceErrStr:               traceErrStr,
		SlowThreshold:             200 * time.Millisecond,
		LogLevel:                  zapcore.WarnLevel,
	}
}
