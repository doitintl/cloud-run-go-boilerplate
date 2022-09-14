package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	sdlog "github.com/Mattel/logrus-stackdriver-formatter"
	"github.com/doitintl/cloud-run-go-boilerplate/api"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func initLogger() *logrus.Entry {
	switch viper.GetString("LOG_LEVEL") {
	case "debug", "DEBUG":
		logrus.SetLevel(logrus.DebugLevel)
	case "info", "INFO":
		logrus.SetLevel(logrus.InfoLevel)
	case "warning", "WARNING":
		logrus.SetLevel(logrus.WarnLevel)
	case "error", "ERROR":
		logrus.SetLevel(logrus.ErrorLevel)
	case "fatal", "FATAL":
		logrus.SetLevel(logrus.FatalLevel)
	case "panic", "PANIC":
		logrus.SetLevel(logrus.PanicLevel)
	default:
		logrus.SetLevel(logrus.InfoLevel)
	}

	if viper.GetBool("LOG_JSON") {
		formatter := sdlog.NewFormatter(
			sdlog.WithProjectID(viper.GetString("PROJECT_ID")),
			sdlog.WithService("ServiceName"),
		)
		logrus.SetFormatter(formatter)
	}

	logrus.SetReportCaller(true)

	return logrus.WithFields(logrus.Fields{})
}

func initSettings() {
	viper.AutomaticEnv()

	viper.SetDefault("LOG_JSON", false)
	viper.SetDefault("LOG_LEVEL", "INFO")
	viper.SetDefault("PORT", "8080")
}

func main() {
	initSettings()

	log := initLogger()

	serviceAPI, err := api.New(context.Background(), log)
	if err != nil {
		log.Panicf("could not create api server: %s", err)
	}

	log.Infof("Listening on %s", viper.GetString("PORT"))

	server := &http.Server{
		Addr:              ":" + viper.GetString("PORT"),
		Handler:           serviceAPI.Mux,
		ReadTimeout:       1 * time.Second,
		WriteTimeout:      1 * time.Second,
		IdleTimeout:       30 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
	}

	go func() {
		log.Println("application started")

		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		} else {
			log.Println("application stopped gracefully")
		}
	}()

	stopCh := make(chan os.Signal, 1)
	signal.Notify(stopCh, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	defer func() {
		close(stopCh)
	}()

	log.Println("notified:", <-stopCh)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		panic(err)
	} else {
		log.Println("application shutdowned")
	}
}
