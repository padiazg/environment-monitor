package config

import (
	"time"

	"github.com/namsral/flag"
)

type AQISensorType string

const (
	defaultTick = 120 * time.Second

	AQISensorNone   AQISensorType = "none"
	AQISensorSPS30  AQISensorType = "SPS30"
	AQISensorPMS007 AQISensorType = "PMS007"
)

// Config configuration structure
type Config struct {
	Tick        time.Duration
	Url         string
	AQISensor   string
	ApiKey      string
	Source      string
	Description string
	Latitude    float64
	Longitude   float64
}

// Init reads configuration
func (c *Config) Init(args []string) error {
	flags := flag.NewFlagSet(args[0], flag.ExitOnError)
	flags.String(flag.DefaultConfigFlagname, "", "Path to config file")

	var (
		tick        = flags.Duration("tick", defaultTick, "Ticking interval")
		url         = flags.String("url", "", "Request URL")
		apiKey      = flags.String("api_key", "", "API key to post readings")
		aqiSensor   = flags.String("aqi_sensor", "", "Model of the device")
		source      = flags.String("source", "", "Name used to identify the device")
		description = flags.String("description", "", "User friendly name to identify the device")
		latitude    = flags.Float64("latitude", 0.0, "Physical latitude coordinate of the device")
		longitude   = flags.Float64("longitude", 0.0, "Physical longitude coordinate of the device")
	)

	if err := flags.Parse(args[1:]); err != nil {
		return err
	}

	c.Tick = *tick
	c.Url = *url
	c.ApiKey = *apiKey
	c.AQISensor = *aqiSensor
	c.Source = *source
	c.Description = *description
	c.Latitude = *latitude
	c.Longitude = *longitude

	return nil
}
