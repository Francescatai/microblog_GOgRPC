// Copyright 2023 Francesca <https://github.com/Francescatai>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/Francescatai/microblog_GOgRPC.

package microblog

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"

	"microblog/internal/microblog/controller/v1/user"
	"microblog/internal/microblog/store"
	"microblog/internal/pkg/known"
	"microblog/internal/pkg/log"
	mwRequestId "microblog/internal/pkg/middleware"
	pb "microblog/pkg/proto/microblog/v1"
	"microblog/pkg/token"
	"microblog/pkg/version/verflag"
)

var cfgFile string

// 創建一個 *cobra.Command 對象之後，可以使用 Command 對象的 Execute 方法來啟動應用
func NewMicroBlogCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "microblog",
		Short: "A Go practical project",
		Long: `A Go practical project, used to create user with basic information.

				Find more information at:
				https://github.com/Francescatai/microblog_GOgRPC`,

		// true -> 可以保持命令出錯時顯示錯誤資訊
		SilenceUsage: true,
		// 指定調用 cmd.Execute() 時，執行Run function，執行失敗會返回錯誤資訊
		RunE: func(cmd *cobra.Command, args []string) error {
			log.Init(logOptions())
			defer log.Sync() // 將緩存中的log存到文件中

			// if `--version=true`，print version info and exit
			verflag.PrintAndExitIfRequested()

			return run()
		},
		// 不需要指定命令行參數
		Args: func(cmd *cobra.Command, args []string) error {
			for _, arg := range args {
				if len(arg) > 0 {
					return fmt.Errorf("%q does not take any arguments, got %q", cmd.CommandPath(), args)
				}
			}

			return nil
		},
	}

	// initConfig 函數在每個命令行都會被調用讀取配置
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

	// 初始化 store
	if err := initStore(); err != nil {
		return err
	}

	token.Init(viper.GetString("jwt-secret"), known.XUsernameKey)

	// Gin mode
	gin.SetMode(viper.GetString("runmode"))

	g := gin.New()

	mws := []gin.HandlerFunc{gin.Recovery(), mwRequestId.NoCache, mwRequestId.Cors, mwRequestId.Secure, mwRequestId.RequestID()}

	g.Use(mws...)

	if err := installRouters(g); err != nil {
		return err
	}

	// Server instantiation
	// HTTP
	httpsrv := startInsecureServer(g)
	// HTTPS
	httpssrv := startSecureServer(g)
	// grpc
	grpcsrv := startGRPCServer()

	// 等待中斷signal關閉server（10 秒超時)
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT -> Ctrl +C
	// kill -9 is syscall.SIGKILL but can't be caught, so don't need to add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) // 接收到以上兩種singal才會繼續執行
	<-quit
	log.Infow("Shutting down server ...")

	// 新增 ctx 用於通知server goroutine, 有 10 秒時間完成當前正在處理的請求
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 10 秒内將未完成的請求處理完再關閉server，超時直接退出
	if err := httpsrv.Shutdown(ctx); err != nil {
		log.Errorw("Insecure Server forced to shutdown", "err", err)
		return err
	}
	if err := httpssrv.Shutdown(ctx); err != nil {
		log.Errorw("Secure Server forced to shutdown", "err", err)
		return err
	}

	grpcsrv.GracefulStop()

	log.Infow("Server exiting")

	return nil
}

// HTTP server
func startInsecureServer(g *gin.Engine) *http.Server {

	httpsrv := &http.Server{Addr: viper.GetString("addr"), Handler: g}

	log.Infow("Start to listening the incoming requests on http address", "addr", viper.GetString("addr"))
	go func() {
		if err := httpsrv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalw(err.Error())
		}
	}()

	return httpsrv
}

// HTTPS server
func startSecureServer(g *gin.Engine) *http.Server {
	httpssrv := &http.Server{Addr: viper.GetString("tls.addr"), Handler: g}

	log.Infow("Start to listening the incoming requests on https address", "addr", viper.GetString("tls.addr"))
	cert, key := viper.GetString("tls.cert"), viper.GetString("tls.key")
	if cert != "" && key != "" {
		go func() {
			if err := httpssrv.ListenAndServeTLS(cert, key); err != nil && !errors.Is(err, http.ErrServerClosed) {
				log.Fatalw(err.Error())
			}
		}()
	}

	return httpssrv
}

func startGRPCServer() *grpc.Server {
	lis, err := net.Listen("tcp", viper.GetString("grpc.addr"))
	if err != nil {
		log.Fatalw("Failed to listen", "err", err)
	}

	grpcsrv := grpc.NewServer()
	pb.RegisterMicroblogServer(grpcsrv, user.New(store.S, nil))

	log.Infow("Start to listening the incoming requests on grpc address", "addr", viper.GetString("grpc.addr"))
	go func() {
		if err := grpcsrv.Serve(lis); err != nil {
			log.Fatalw(err.Error())
		}
	}()

	return grpcsrv
}
