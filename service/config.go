package service

import (
	"time"

	"github.com/jacobalberty/roofwater/service/utils"
)

type Config struct {
	PulseWidth  float64            `envconfig:"PULSEWIDTH" default:"0.05" desc:"Percent of the Interval expressed as a decimal to turn valve on"`
	PulsePeriod time.Duration      `envconfig:"RW_PULSEPERIOD" default:"5m" desc:"Period of pulses"`
	MinTemp     float64            `envconfig:"MINTEMP" default:"90" desc:"Minimum temperature to run the valve"`
	Weather     WeatherConfig      `envconfig:"OWM" required:"true" desc:"OpenWeatherMap configuration"`
	Tracing     utils.TracerConfig `envconfig:"TRACING" required:"true" desc:"Tracing configuration"`
	ValveConfig ValveConfig        `envconfig:"VALVE" required:"true" desc:"Valve configuration"`
	MQTTConfig  MQTTConfig         `envconfig:"MQTT" required:"true" desc:"MQTT configuration"`
	PIDConfig   PIDConfig          `envconfig:"PID" required:"true" desc:"PID configuration"`
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

type PIDConfig struct {
	Kp float64 `envconfig:"KP" default:"1.0" desc:"Proportional gain"`
	Ki float64 `envconfig:"KI" default:"0.0" desc:"Integral gain"`
	Kd float64 `envconfig:"KD" default:"0.0" desc:"Derivative gain"`
}

type MQTTConfig struct {
	URL               string        `envconfig:"URL" desc:"MQTT Broker URL"`
	User              string        `envconfig:"USER" desc:"MQTT username"`
	Pass              string        `envconfig:"PASS" desc:"MQTT password"`
	KeepAlive         uint16        `envconfig:"KEEPALIVE" default:"5" desc:"MQTT keepalive in seconds"`
	ConnectRetryDelay time.Duration `envconfig:"CONNECT_RETRY_DELAY" default:"10s" desc:"MQTT connect retry delay"`
	Timeout           time.Duration `envconfig:"TIMEOUT" default:"10s" desc:"MQTT timeout"`
}
