package service

import (
	"context"
	"time"

	owm "github.com/briandowns/openweathermap"
	"github.com/jacobalberty/roofwater/service/utils"
)

type ExpiringWeather struct {
	w          *owm.CurrentWeatherData
	cfg        Config
	lastUpdate time.Time
}

func (e *ExpiringWeather) CurrentTempByZip(ctx context.Context) (float64, error) {
	var err error
	if time.Since(e.lastUpdate) > e.cfg.Weather.CacheDuration {
		utils.Logger.Ctx(ctx).Info("Updating weather cache")
		err = e.w.CurrentByZipcode(e.cfg.Weather.Zip, e.cfg.Weather.Country)
		e.lastUpdate = time.Now()
	}
	return e.w.Main.FeelsLike, err
}

func NewExpiringWeather(cfg Config) (*ExpiringWeather, error) {
	var (
		err error
		w   *owm.CurrentWeatherData
	)
	w, err = owm.NewCurrent("F", "EN", cfg.Weather.APIKey)
	if err != nil {
		return nil, err
	}

	return &ExpiringWeather{
		w:   w,
		cfg: cfg,
	}, nil
}
