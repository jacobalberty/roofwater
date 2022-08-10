package service

import (
	"net"
	"time"
)

type Config struct {
	Valve         net.IP        `envconfig:"VALVE_IP" required:"true"`
	PulseWidth    time.Duration `envconfig:"PULSEWIDTH" default:"15s"`
	PulseInterval time.Duration `envconfig:"PULSEINTERVAL" default:"5m"`
	MinTemp       float64       `envconfig:"MINTEMP" default:"90"`
	Weather       WeatherConfig `envconfig:"OWM"`
}

type WeatherConfig struct {
	APIKey        string        `envconfig:"API_KEY" required:"true"`
	Zip           string        `envconfig:"ZIP" required:"true"`
	Country       string        `envconfig:"COUNTRY" required:"true"`
	CacheDuration time.Duration `envconfig:"CACHE_DURATION" default:"1h"`
}
