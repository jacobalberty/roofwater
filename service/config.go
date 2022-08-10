package service

import (
	"net"
	"time"
)

type Config struct {
	Valve         net.IP        `envconfig:"RW_VALVE_IP" required:"true"`
	PulseWidth    time.Duration `envconfig:"RW_PULSEWIDTH" default:"15s"`
	PulseInterval time.Duration `envconfig:"RW_PULSEINTERVAL" default:"5m"`
	MinTemp       float64       `envconfig:"RW_MINTEMP" default:"80"`
	Weather       WeatherConfig
}

type WeatherConfig struct {
	APIKey        string        `envconfig:"RW_OWM_API_KEY" required:"true"`
	Zip           string        `envconfig:"RW_OWM_ZIP" required:"true"`
	Country       string        `envconfig:"RW_OWM_COUNTRY" required:"true"`
	CacheDuration time.Duration `envconfig:"RW_OWM_CACHE_DURATION" default:"1h"`
}
