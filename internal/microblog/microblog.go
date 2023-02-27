// Copyright 2023 Francesca <https://github.com/Francescatai>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/Francescatai/microblog_GOgRPC.

package microblog

import (
	"encoding/json"
	"fmt"
	"net/http"
	"errors"
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/gin-gonic/gin"

	"microblog/internal/pkg/log"
	"microblog/pkg/version/verflag"
	mwRequestId "microblog/internal/pkg/middleware"
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
			log.Init(logOptions())
 			defer log.Sync() // 將緩存中的log存到文件中

			// if `--version=true`，print version info and exit
			verflag.PrintAndExitIfRequested()

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

	// add --version flag
	verflag.AddFlags(cmd.PersistentFlags())

	return cmd
}

// 實際的業務程式碼入口函數
func run() error {
	fmt.Println("start CLI service")
	settings, _ := json.Marshal(viper.AllSettings())
	log.Infow(string(settings))

	// Gin mode
	gin.SetMode(viper.GetString("runmode"))


	g := gin.New()

	mws := []gin.HandlerFunc{gin.Recovery(), mwRequestId.NoCache, mwRequestId.Cors, mwRequestId.Secure, mwRequestId.RequestID()}

	g.Use(mws...)

	// 404 Handler
	g.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "Page not found."})
	})

	// healthz handler
	g.GET("/healthz", func(c *gin.Context) {
		log.C(c).Infow("Healthz function called")
		
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// HTTP Server instantiation
	httpsrv := &http.Server{Addr: viper.GetString("addr"), Handler: g}

	log.Infow("Start to listening the incoming requests on http address", "addr", viper.GetString("addr"))
	go func() {
        if err := httpsrv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
            log.Fatalw(err.Error())
        }
    }()

	// 等待中斷signal關閉server（10 秒超時)
    quit := make(chan os.Signal, 1)
    // kill (no param) default send syscall.SIGTERM
    // kill -2 is syscall.SIGINT -> Ctrl +C
    // kill -9 is syscall.SIGKILL but can't be caught, so don't need to add it
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) 
    <-quit                                               // 接收到以上兩種singal才會繼續執行
    log.Infow("Shutting down server ...")

    // 新增 ctx 用於通知server goroutine, 有 10 秒時間完成當前正在處理的請求
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    // 10 秒内將未完成的請求處理完再關閉server，超時直接退出
    if err := httpsrv.Shutdown(ctx); err != nil {
        log.Errorw("Insecure Server forced to shutdown", "err", err)
        return err
    }

    log.Infow("Server exiting")

	return nil
}
