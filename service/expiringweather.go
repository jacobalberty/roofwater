package service

import (
	"context"
	"strings"
	"time"

	owm "github.com/briandowns/openweathermap"
	"github.com/jacobalberty/roofwater/service/utils"
	"go.uber.org/zap"
)

type ExpiringWeather struct {
	w          *owm.CurrentWeatherData
	cfg        Config
	lastUpdate time.Time
}

func (e *ExpiringWeather) CurrentTempByZip(ctx context.Context) (float64, error) {
	var err error

	ctx, span := utils.Tracer.Start(ctx, "CurrentTempByZip")
	defer span.End()

	if time.Since(e.lastUpdate) > e.cfg.Weather.CacheDuration {
		err = e.w.CurrentByZipcode(e.cfg.Weather.Zip, e.cfg.Weather.Country)
		utils.Logger.Ctx(ctx).Info("Updated weather cache",
			zap.Float64("current.feels_like", e.w.Main.FeelsLike),
			zap.Float64("current.temp", e.w.Main.Temp),
		)
		e.lastUpdate = time.Now()
	}
	return e.w.Main.FeelsLike, err
}

func NewExpiringWeather(cfg Config) (*ExpiringWeather, error) {
	var (
		err error
		w   *owm.CurrentWeatherData
	)

	w, err = owm.NewCurrent(
		strings.ToUpper(cfg.Weather.Unit),
		strings.ToUpper(cfg.Weather.Language),
		cfg.Weather.APIKey,
	)

	if err != nil {
		return nil, err
	}

	return &ExpiringWeather{
		w:   w,
		cfg: cfg,
	}, nil
}
