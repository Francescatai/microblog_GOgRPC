// Copyright 2023 Innkeeper Francesca <https://github.com/Francescatai>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/Francescatai/microblog_GOgRPC.

/*
cobra setting
*/
package microblog

import (
    "encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// NewMicroBlogCommand 創建一個 *cobra.Command 對象之後，可以使用 Command 對象的 Execute 方法來啟動應用
func NewMicroBlogCommand() *cobra.Command {
    cmd := &cobra.Command{
        // 名字會出現在幫助資訊中
        Use: "microblog",
        // 命令敘述
        Short: "A Go practical project",
        // 命令詳情
        Long: `A Go practical project, used to create user with basic information.

				Find more information at:
				https://github.com/Francescatai/microblog_GOgRPC`,

        // true-> 可以保持命令出錯時顯示錯誤資訊
        SilenceUsage: true,
        // 指定調用 cmd.Execute() 時，執行Run function，執行失敗會返回錯誤資訊
        RunE: func(cmd *cobra.Command, args []string) error {
            return run()
        },
        // 命令運行時，不需要指定命令行參數
        Args: func(cmd *cobra.Command, args []string) error {
            for _, arg := range args {
                if len(arg) > 0 {
                    return fmt.Errorf("%q does not take any arguments, got %q", cmd.CommandPath(), args)
                }
            }

            return nil
        },
    }

	// 以下設定，使得 initConfig 函數在每個命令行都會被調用讀取配置
	cobra.OnInitialize(initConfig)

	// 定義flag與配置設定

	// Cobra 支持PersistentFlag，此flag可用於它所分配的命令以及該命令下的每個子命令
	cmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "The path to the microblog configuration file. Empty string for no configuration file.")

	// Cobra 也支持本地flag，本地flag只能在其所绑定的命令上使用
	cmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

    return cmd
}

// 實際的業務程式碼入口函數
func run() error {
    fmt.Println("start CLI service")
	settings, _ := json.Marshal(viper.AllSettings())
    fmt.Println(string(settings))
    // print for test: db -> username 
    fmt.Println(viper.GetString("db.username"))
    return nil
}