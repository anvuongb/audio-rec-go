package main

import (
	"audio-rec-go/src/config"
	"audio-rec-go/src/httpserver"
	"audio-rec-go/src/services/voice"
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"audio-rec-go/src/sql"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"gorm.io/gorm"
)

func main() {
	config.GetConfigStruct(&config.GlobalConfig)

	var httpAddr = flag.String("http", ":"+config.GlobalConfig.HTTP.Port, "http listen address")
	var logger log.Logger
	logger = log.NewLogfmtLogger(os.Stderr)
	logger = log.NewSyncLogger(logger)
	logger = log.With(
		logger,
		"service", "voice",
		"time:", log.DefaultTimestampUTC,
		"caller", log.DefaultCaller,
	)
	level.Info(logger).Log("msg", "service started")
	defer level.Info(logger).Log("msg", "service ended")

	var db *gorm.DB = sql.InitDB(
		config.GlobalConfig.DB.Username, config.GlobalConfig.DB.Password,
		config.GlobalConfig.DB.Host, config.GlobalConfig.DB.Port,
		config.GlobalConfig.DB.Database,
	)

	flag.Parse()

	ctx := context.Background()
	repo := voice.NewRepository(db, logger)
	service := voice.NewService(repo, logger)

	errs := make(chan error)

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	serviceEndpoints := voice.MakeEndpoints(service)

	go func() {
		fmt.Println("listening on port", *httpAddr)
		handler := httpserver.NewHTTPServer(ctx, logger, config.GlobalConfig.HTTP.Cors, serviceEndpoints)
		errs <- http.ListenAndServe(*httpAddr, handler)
	}()

	level.Error(logger).Log("exit", <-errs)
}
