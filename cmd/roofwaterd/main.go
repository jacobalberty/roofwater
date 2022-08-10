package main

import (
	"context"
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
		err error
		cfg service.Config
		c   = make(chan os.Signal, 1)
		w   *service.ExpiringWeather
	)
	utils.InitializeLogger()

	defer func() {
		err := utils.Logger.Sync()
		if err != nil {
			panic(err)
		}
	}()

	utils.Logger.Info("Roof Water started")

	err = envconfig.Process("roofwaterd", &cfg)
	if err != nil {
		utils.Logger.Fatal("Failed to process config", zap.Error(err))
	}

	if err = owm.ValidAPIKey(cfg.Weather.APIKey); err != nil {
		utils.Logger.Fatal("Invalid OpenWeatherMap API key", zap.Error(err))
	}

	utils.Logger.Info("Pulse interval", zap.Duration("duration", cfg.PulseInterval))

	w, err = service.NewExpiringWeather(cfg)
	if err != nil {
		utils.Logger.Fatal("Failed to create weather client", zap.Error(err))
	}

	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		var (
			err error
			t   float64
		)

		for {
			ctx := context.Background()
			t, err = w.CurrentTempByZip(ctx)
			if err != nil {
				utils.Logger.Ctx(ctx).Error("Failed to get weather", zap.Error(err))
			}
			if t > cfg.MinTemp {

				utils.Logger.Ctx(ctx).Info("Temperature is too high", zap.Float64("temp", t))
				service.Valve{IP: cfg.Valve}.RWPulse(ctx, cfg.PulseWidth)
			}
			time.Sleep(cfg.PulseInterval)
		}
	}()

	<-c
	utils.Logger.Info("Received interrupt, exiting")
}
