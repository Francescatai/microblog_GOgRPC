// Copyright 2023 Francesca <https://github.com/Francescatai>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/Francescatai/microblog_GOgRPC.

package log

import (
	"go.uber.org/zap/zapcore"
)

// log相關配置
type Options struct {
	// 是否開啟 caller
	DisableCaller bool
	// 是否禁止在 panic 及以上级别打印堆栈資訊
	DisableStacktrace bool
	// 指定log level：debug, info, warn, error, dpanic, panic, fatal
	Level string
	// log fomat：console, json
	Format string
	// log輸出位置
	OutputPaths []string
}

// 創建一個帶有默認參數的 Options 對象
func NewOptions() *Options {
	return &Options{
		DisableCaller:     false,
		DisableStacktrace: false,
		Level:             zapcore.InfoLevel.String(),
		Format:            "console",
		OutputPaths:       []string{"stdout"},
	}
}