package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	owm "github.com/briandowns/openweathermap"
	"github.com/jacobalberty/roofwater/service"
	"github.com/jacobalberty/roofwater/service/utils"
	"github.com/kelseyhightower/envconfig"
	"go.uber.org/zap"
)

func main() {
	var (
		err    error
		cfg    service.Config
		c      = make(chan os.Signal, 1)
		w      *service.ExpiringWeather
		daemon bool
	)

	flag.BoolVar(&daemon, "d", false, "Run as a daemon")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Fprint(os.Stderr, "\n")
		err = envconfig.Usage("RW", &cfg)
		if err != nil {
			utils.Logger.Fatal("Failed to print config usage", zap.Error(err))
		}
	}

	flag.Parse()

	utils.InitializeLogger(daemon)

	defer func() {
		err := utils.Logger.Sync()
		if err != nil && !errors.Is(err, syscall.ENOTTY) && !errors.Is(err, syscall.EINVAL) {
			log.Println("error syncing logs", err)
		}
	}()

	err = envconfig.Process("RW", &cfg)
	if err != nil {
		if daemon {
			utils.Logger.Fatal("Failed to process config", zap.Error(err))
		} else {
			err = envconfig.Usage("RW", &cfg)
			if err != nil {
				utils.Logger.Fatal("Failed to print config usage", zap.Error(err))
			}
		}
		os.Exit(1)
	}

	if err = owm.ValidAPIKey(cfg.Weather.APIKey); err != nil {
		utils.Logger.Fatal("Invalid OpenWeatherMap API key", zap.Error(err))
	}

	utils.Logger.Info("Pulse period", zap.Duration("period", cfg.PulsePeriod))

	w, err = service.NewExpiringWeather(cfg)
	if err != nil {
		utils.Logger.Fatal("Failed to create weather client", zap.Error(err))
	}

	if daemon {
		ctx := context.Background()
		shutdown := utils.InitializeTracer(ctx, cfg.Tracing)
		defer shutdown()

		signal.Notify(c, os.Interrupt, syscall.SIGTERM)

		go func() {
			var (
				r = service.RoofD{}
			)
			utils.Logger.Info("Roof water loop successfully started")

			for {
				func() {
					r.CheckWeatherAndCool(ctx, w, cfg)
					time.Sleep(cfg.PulsePeriod)
				}()
			}
		}()

		<-c
		utils.Logger.Info("Received interrupt, exiting")
	} else {
		(&service.RoofD{}).CheckWeatherAndCool(context.Background(), w, cfg)
	}
}
