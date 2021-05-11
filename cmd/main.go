package main

import (
	"clover/controller"
	"clover/model/mysql"
	"clover/model/redis"
	"clover/pkg/log"
	"clover/pkg/snowflake"
	"clover/setting"
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	// parse the flag
	flag.Parse()

	// init logger
	log.InitLogger()

	log.WithCategory("main").Info("Starting server...")

	// parse the application setting
	setting.InitAppSettings()

	// init db
	log.WithCategory("main").Info("Starting DB...")
	mysql.InitDB(setting.GetMysqlConfig(), setting.GetTimeZone())

	// init redis
	log.WithCategory("main").Info("Starting redis...")
	redis.InitRedis(setting.GetRedisConfig())

	// init snowflake
	snowflake.Init(setting.GetMachineID())

	// Init gin router
	log.WithCategory("main").Info("Starting router...")
	router := controller.InitRouter()
	server := http.Server{
		Addr:    setting.GetEndpoint(),
		Handler: router,
	}

	sig := make(chan os.Signal, 0)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	go func() {
		sg := <-sig
		log.WithCategory("main").Info("Catched the signal", sg)

		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		err := server.Shutdown(ctx)
		if err != nil {
			log.WithCategory("main").WithError(err).Error("Shutdown http server error")
		}

		mysql.CloseMysql()
		redis.CloseRedis()
	}()

	err := server.ListenAndServe()
	if err != nil {
		log.WithCategory("main").WithError(err).Error("Server stopping...")
	}

	log.WithCategory("main").Info("Server Stopped...")
}
