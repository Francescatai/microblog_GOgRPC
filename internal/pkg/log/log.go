// Copyright 2023 Francesca <https://github.com/Francescatai>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/Francescatai/microblog_GOgRPC.

package log

import (
	"context"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"microblog/internal/pkg/known"
)

// project的log interface. 只包含了支持的log method
type Logger interface {
	Debugw(msg string, keysAndValues ...interface{})
	Infow(msg string, keysAndValues ...interface{})
	Warnw(msg string, keysAndValues ...interface{})
	Errorw(msg string, keysAndValues ...interface{})
	Panicw(msg string, keysAndValues ...interface{})
	Fatalw(msg string, keysAndValues ...interface{})
	Sync()
}

type zapLogger struct {
	z *zap.Logger
}

// 確認 zapLogger 實現了 Logger interface
var _ Logger = &zapLogger{}

var (
	mu sync.Mutex

	// 默認的全局 Logger
	std = NewLogger(NewOptions())
)

// 使用指定選項初始化 Logger
func Init(opts *Options) {
	mu.Lock()
	defer mu.Unlock()

	std = NewLogger(opts)
}

// 根據傳入的 opts 創建 Logger
func NewLogger(opts *Options) *zapLogger {
	if opts == nil {
		opts = NewOptions()
	}

	// 將log level轉為 zapcore.Level
	var zapLevel zapcore.Level
	if err := zapLevel.UnmarshalText([]byte(opts.Level)); err != nil {
		// 默認使用 info 级别
		zapLevel = zapcore.InfoLevel
	}

	// 創建一個默認的 encoder 配置
	encoderConfig := zap.NewProductionEncoderConfig()
	// 自定義 MessageKey 為 message
	encoderConfig.MessageKey = "message"
	// 自定義 TimeKey 為 timestamp
	encoderConfig.TimeKey = "timestamp"
	// 將時間序列化為 `2006-01-02 15:04:05.000` 格式
	encoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
	}
	// 將 time.Duration 序列化為經過的毫秒數的浮點數
	encoderConfig.EncodeDuration = func(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendFloat64(float64(d) / float64(time.Millisecond))
	}

	// 創建 zap.Logger 所需配置
	cfg := &zap.Config{
		// 在log中顯示調用log所在的位置，例如：`"caller":"microblog/microblog.go:75"`
		DisableCaller: opts.DisableCaller,
		// 是否禁止在 panic 及以上级别打印堆栈資訊
		DisableStacktrace: opts.DisableStacktrace,
		// 指定kog level
		Level: zap.NewAtomicLevelAt(zapLevel),
		// 指定log格式
		Encoding:      opts.Format,
		EncoderConfig: encoderConfig,
		// 指定log存檔位置
		OutputPaths: opts.OutputPaths,
		// 設定 zap 内部錯誤輸出位置
		ErrorOutputPaths: []string{"stderr"},
	}

	// 使用 cfg 創建 *zap.Logger 對象
	z, err := cfg.Build(zap.AddStacktrace(zapcore.PanicLevel), zap.AddCallerSkip(1))
	if err != nil {
		panic(err)
	}
	logger := &zapLogger{z: z}

	// 把標準庫的 log.Logger 的 info 级别的輸出重定向到 zap.Logger
	zap.RedirectStdLog(z)

	return logger
}

// 調用 zap.Logger 的 Sync 方法，將緩存中的log刷新到文件中. 主程序需要在退出前調用 Sync
func Sync() { std.Sync() }

func (l *zapLogger) Sync() {
	_ = l.z.Sync()
}

// Debugw 輸出 debug level的log
func Debugw(msg string, keysAndValues ...interface{}) {
	std.z.Sugar().Debugw(msg, keysAndValues...)
}

func (l *zapLogger) Debugw(msg string, keysAndValues ...interface{}) {
	l.z.Sugar().Debugw(msg, keysAndValues...)
}

// Infow 輸出 info level的log
func Infow(msg string, keysAndValues ...interface{}) {
	std.z.Sugar().Infow(msg, keysAndValues...)
}

func (l *zapLogger) Infow(msg string, keysAndValues ...interface{}) {
	l.z.Sugar().Infow(msg, keysAndValues...)
}

// Warnw 輸出 warning level的log
func Warnw(msg string, keysAndValues ...interface{}) {
	std.z.Sugar().Warnw(msg, keysAndValues...)
}

func (l *zapLogger) Warnw(msg string, keysAndValues ...interface{}) {
	l.z.Sugar().Warnw(msg, keysAndValues...)
}

// Errorw 輸出 error level的log
func Errorw(msg string, keysAndValues ...interface{}) {
	std.z.Sugar().Errorw(msg, keysAndValues...)
}

func (l *zapLogger) Errorw(msg string, keysAndValues ...interface{}) {
	l.z.Sugar().Errorw(msg, keysAndValues...)
}

// Panicw 輸出 panic level的log
func Panicw(msg string, keysAndValues ...interface{}) {
	std.z.Sugar().Panicw(msg, keysAndValues...)
}

func (l *zapLogger) Panicw(msg string, keysAndValues ...interface{}) {
	l.z.Sugar().Panicw(msg, keysAndValues...)
}

// Fatalw 輸出 fatal level的log
func Fatalw(msg string, keysAndValues ...interface{}) {
	std.z.Sugar().Fatalw(msg, keysAndValues...)
}

func (l *zapLogger) Fatalw(msg string, keysAndValues ...interface{}) {
	l.z.Sugar().Fatalw(msg, keysAndValues...)
}

/* write requestID to log */
// C 解析傳入的 context，取出key-value，並新增到 zap.Logger log中
func C(ctx context.Context) *zapLogger {
	return std.C(ctx)
}

func (l *zapLogger) C(ctx context.Context) *zapLogger {
	lc := l.clone()

	if requestID := ctx.Value(known.XRequestIDKey); requestID != nil {
		lc.z = lc.z.With(zap.Any(known.XRequestIDKey, requestID))
	}

	if userID := ctx.Value(known.XUsernameKey); userID != nil {
		lc.z = lc.z.With(zap.Any(known.XUsernameKey, userID))
	}

	return lc
}

// clone deep copy zapLogger
func (l *zapLogger) clone() *zapLogger {
	lc := *l
	return &lc
}
