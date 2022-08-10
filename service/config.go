package service

import (
	"net"
	"time"
)

type Config struct {
	Valve         net.IP        `envconfig:"VALVE_IP" required:"true" desc:"IP address of the valve"`
	PulseWidth    time.Duration `envconfig:"PULSEWIDTH" default:"15s" desc:"Duration of the time to turn valve on"`
	PulseInterval time.Duration `envconfig:"PULSEINTERVAL" default:"5m" desc:"Interval between pulses"`
	MinTemp       float64       `envconfig:"MINTEMP" default:"90" desc:"Minimum temperature to run the valve"`
	Weather       WeatherConfig `envconfig:"OWM" required:"true" desc:"OpenWeatherMap configuration"`
}

type WeatherConfig struct {
	APIKey        string        `envconfig:"API_KEY" required:"true" desc:"OpenWeatherMap API key"`
	Zip           string        `envconfig:"ZIP" required:"true" desc:"Zip code to get weather for"`
	Country       string        `envconfig:"COUNTRY" required:"true" desc:"Country to get weather for"`
	CacheDuration time.Duration `envconfig:"CACHE_DURATION" default:"1h" desc:"Duration to cache weather data"`
}
