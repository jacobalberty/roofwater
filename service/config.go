package service

import (
	"time"

	"github.com/jacobalberty/roofwater/service/utils"
)

type Config struct {
	PulseWidth    time.Duration      `envconfig:"PULSEWIDTH" default:"15s" desc:"Duration of the time to turn valve on"`
	PulseInterval time.Duration      `envconfig:"PULSEINTERVAL" default:"5m" desc:"Interval between pulses"`
	MinTemp       float64            `envconfig:"MINTEMP" default:"90" desc:"Minimum temperature to run the valve"`
	Weather       WeatherConfig      `envconfig:"OWM" required:"true" desc:"OpenWeatherMap configuration"`
	Tracing       utils.TracerConfig `envconfig:"TRACING" required:"true" desc:"Tracing configuration"`
	ValveConfig   ValveConfig        `envconfig:"VALVE" required:"true" desc:"Valve configuration"`
	MQTTConfig    MQTTConfig         `envconfig:"MQTT" required:"true" desc:"MQTT configuration"`
}

type WeatherConfig struct {
	APIKey        string        `envconfig:"API_KEY" required:"true" desc:"OpenWeatherMap API key"`
	Zip           string        `envconfig:"ZIP" required:"true" desc:"Zip code to get weather for"`
	Country       string        `envconfig:"COUNTRY" required:"true" desc:"Country to get weather for"`
	CacheDuration time.Duration `envconfig:"CACHE_DURATION" default:"1h" desc:"Duration to cache weather data"`
	Unit          string        `envconfig:"UNIT" default:"F" desc:"Unit to use for weather data"`
	Language      string        `envconfig:"LANGUAGE" default:"EN" desc:"Language to use for weather data"`
}

type ValveConfig struct {
	Addr  string `envconfig:"HTTP_ADDR" desc:"HTTP Address of the valve"`
	Topic string `envconfig:"MQTT_TOPIC" desc:"MQTT topic to publish to"`
}

type MQTTConfig struct {
	URL  string `envconfig:"URL" desc:"MQTT Broker URL"`
	User string `envconfig:"USER" desc:"MQTT username"`
	Pass string `envconfig:"PASS" desc:"MQTT password"`
}
