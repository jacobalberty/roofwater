package main

import (
	"context"
	"errors"
	"flag"
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
			envconfig.Usage("RW", &cfg)
		}
		os.Exit(1)
	}

	if err = owm.ValidAPIKey(cfg.Weather.APIKey); err != nil {
		utils.Logger.Fatal("Invalid OpenWeatherMap API key", zap.Error(err))
	}

	utils.Logger.Info("Pulse interval", zap.Duration("duration", cfg.PulseInterval))

	w, err = service.NewExpiringWeather(cfg)
	if err != nil {
		utils.Logger.Fatal("Failed to create weather client", zap.Error(err))
	}

	if daemon {
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)

		go func() {
			utils.Logger.Info("Roof water loop successfully started")

			for {
				checkWeatherAndCool(context.Background(), w, cfg)
				time.Sleep(cfg.PulseInterval)
			}
		}()

		<-c
		utils.Logger.Info("Received interrupt, exiting")
	} else {
		checkWeatherAndCool(context.Background(), w, cfg)
	}
}

func checkWeatherAndCool(ctx context.Context, w *service.ExpiringWeather, cfg service.Config) {
	var (
		err error
		t   float64
	)

	t, err = w.CurrentTempByZip(ctx)
	if err != nil {
		utils.Logger.Ctx(ctx).Error("Failed to get weather", zap.Error(err))
	}
	if t > cfg.MinTemp {
		utils.Logger.Ctx(ctx).Info("Temperature is too high", zap.Float64("temp", t))
		service.Valve{IP: cfg.Valve}.RWPulse(ctx, cfg.PulseWidth)
	}
}
