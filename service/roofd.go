package service

import (
	"context"

	"github.com/jacobalberty/roofwater/service/utils"
	"go.uber.org/zap"
)

type RoofD struct {
	valve *Valve
}

func (r *RoofD) CheckWeatherAndCool(ctx context.Context, w *ExpiringWeather, cfg Config) {
	var (
		err error
		t   float64
	)

	ctx, span := utils.Tracer.Start(ctx, "checkWeatherAndCool")
	defer span.End()

	t, err = w.CurrentTempByZip(ctx)
	if err != nil {
		utils.Logger.Ctx(ctx).Error("Failed to get weather", zap.Error(err))
	}
	if t > cfg.MinTemp {
		if r.valve == nil {
			if cfg.MQTTConfig.URL != "" {
				r.valve = &Valve{
					Addr:       cfg.ValveConfig.Topic,
					MQTTConfig: &cfg.MQTTConfig,
				}
			} else {
				r.valve = &Valve{Addr: cfg.ValveConfig.Addr}
			}
		}
		utils.Logger.Ctx(ctx).Info("Temperature is too high", zap.Float64("temp", t))
		r.valve.RWPulse(ctx, cfg.PulseWidth, cfg.PulseInterval)
	}
}
