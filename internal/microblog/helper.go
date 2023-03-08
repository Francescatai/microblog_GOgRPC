// Copyright 2023 Francesca <https://github.com/Francescatai>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/Francescatai/microblog_GOgRPC.

package microblog

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"microblog/internal/microblog/store"
	"microblog/internal/pkg/log"
	"microblog/pkg/db"
)

const (
	// recommendedHomeDir 定義放置 microblog 服務配置的默認目錄
	recommendedHomeDir = ".microblog"

	// defaultConfigName 指定 microblog 服務默認配置文件名
	defaultConfigName = "./configs/microblog.yaml"
)

// initConfig 設置需要讀取的默認配置文件、環境變量，病毒取配置文件内容到 viper 中
func initConfig() {
	if cfgFile != "" {
		// 從命令行選項指定的配置文件讀取
		viper.SetConfigFile(cfgFile)
	} else {
		// 查找用户主目錄
		home, err := os.UserHomeDir()
		// 如果取得用戶主目錄失敗，印 `'Error: xxx` 錯誤並退出
		cobra.CheckErr(err)

		// 將 `$HOME/<recommendedHomeDir>` 目錄加入到配置文件的查找路徑
		viper.AddConfigPath(filepath.Join(home, recommendedHomeDir))

		// 把當前目錄加入到配置文件的查找路徑中
		viper.AddConfigPath(".")

		// 設置配置文件格式為 YAML
		viper.SetConfigType("yaml")

		// 配置文件名稱
		viper.SetConfigName(defaultConfigName)
	}

	// 讀取匹配的環境變量
	viper.AutomaticEnv()

	// 讀取環境變量的前綴為 microblog，如果是 microblog，自動轉大寫
	viper.SetEnvPrefix("MICROBLOG")

	// 將 viper.Get(key) key 字符串中 '.' 和 '-' 替換為 '_'
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)

	// 如果指定了配置文件名，則使用指定的配置文件，否则在註冊的查找路徑中搜尋
	if err := viper.ReadInConfig(); err != nil {
		log.Errorw("Failed to read viper configuration file", "err", err)
	}

	// viper 當前使用的配置文件
	log.Infow("Using config file", "file", viper.ConfigFileUsed())
}

// 從 viper 中讀取log配置，創建 `*log.Options` 並返回
func logOptions() *log.Options {
	return &log.Options{
		DisableCaller:     viper.GetBool("log.disable-caller"),
		DisableStacktrace: viper.GetBool("log.disable-stacktrace"),
		Level:             viper.GetString("log.level"),
		Format:            viper.GetString("log.format"),
		OutputPaths:       viper.GetStringSlice("log.output-paths"),
	}
}

// 創建 gorm.DB 實例，初始化 store
func initStore() error {
	dbOptions := &db.MySQLOptions{
		Host:                  viper.GetString("db.host"),
		Username:              viper.GetString("db.username"),
		Password:              viper.GetString("db.password"),
		Database:              viper.GetString("db.database"),
		MaxIdleConnections:    viper.GetInt("db.max-idle-connections"),
		MaxOpenConnections:    viper.GetInt("db.max-open-connections"),
		MaxConnectionLifeTime: viper.GetDuration("db.max-connection-life-time"),
		LogLevel:              viper.GetInt("db.log-level"),
	}

	ins, err := db.NewMySQL(dbOptions)
	if err != nil {
		return err
	}

	_ = store.NewStore(ins)

	return nil
}
