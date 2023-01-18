package main

import (
	"context"
	"fmt"
	"metalflow/initialize"
	"metalflow/pkg/global"
	"net/http"

	"gopkg.in/alecthomas/kingpin.v2"

	// 加入pprof性能分析
	_ "net/http/pprof" //nolint:gosec
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"
	"time"
)

const version = "1.4.0"

var (
	app        = kingpin.New("metalflow", "Metal Flow").Version(version)
	configFile = app.Flag("config-file", "Config file (.yml)").String()
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			// 将异常写入日志
			global.Log.Error(fmt.Sprintf("项目启动失败: %v\n堆栈信息: %v", err, string(debug.Stack())))
		}
	}()

	// 解析命令行
	kingpin.MustParse(app.Parse(os.Args[1:]))
	// 初始化配置
	initialize.Config(*configFile)

	// 初始化日志
	initialize.Logger()

	// 初始化数据库
	initialize.Mysql()

	// 初始化redis
	initialize.Redis()

	// 初始化casbin策略管理器
	initialize.CasbinEnforcer()

	// 初始校验器
	initialize.Validate()

	// 初始化路由
	r := initialize.Routers()

	// 初始化数据
	initialize.Data()

	// 初始化异步任务Machinery，其要在redis之后,需在consul之前
	initialize.Async()

	// 初始化consul watch, consul初始化要在mysql初始化之后
	initialize.Consul()

	// 初始化定时任务
	initialize.Cron()

	host := "0.0.0.0"
	port := global.Conf.System.Port
	// 服务器启动以及优雅的关闭
	// 参考地址https://github.com/gin-gonic/examples/blob/master/graceful-shutdown/graceful-shutdown/server.go
	srv := &http.Server{
		Addr:              fmt.Sprintf("%s:%d", host, port),
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		// 加入pprof性能分析
		global.Log.Info("Debug pprof is running at ", fmt.Sprintf("%s:%d", host, global.Conf.System.PprofPort))
		PprofSrv := &http.Server{
			Addr:              fmt.Sprintf("%s:%d", host, global.Conf.System.PprofPort),
			Handler:           nil,
			ReadHeaderTimeout: 5 * time.Second,
		}
		if err := PprofSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			global.Log.Error("listen pprof error: ", err)
		}
	}()

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			global.Log.Error("listen error: ", err)
		}
	}()

	global.Log.Info(fmt.Sprintf("Server is running at %s:%d/%s", host, port, global.Conf.System.UrlPathPrefix))

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	global.Log.Info("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) //nolint:gomnd
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		global.Log.Error("Server forced to shutdown: ", err)
	}

	global.Log.Info("Server exiting")
}
