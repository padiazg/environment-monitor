package settings

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

type SensorType string

const (
	SensorNone   SensorType = "none"
	SensorSPS30  SensorType = "SPS30"
	SensorPMS007 SensorType = "PMS007"
	SensorZH07   SensorType = "ZH07"
)

type Settings struct {
	Url         string
	Sensor      string
	ApiKey      string
	Source      string
	Description string
	Lat         float64
	Lon         float64
	Interval    time.Duration
	ZH07        *ZH07Settings
	Dummy       *DummySettings
}

func (s *Settings) Show() {
	b, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		fmt.Println("error:", err)
	}

	fmt.Printf("Settings:\n%s\n", string(b))
}

func (s *Settings) ShowKeyValuePairs() {
	fmt.Printf("Key/Value pairs:\n")
	// print all the keys
	for _, key := range viper.AllKeys() {
		val := viper.Get(key)
		fmt.Printf("  %s: %v\n", key, val)
	}
}

func (s *Settings) Save(name string) error {
	// if err := viper.WriteConfigAs(name); err != nil {
	// 	return fmt.Errorf("writing config file: %+v", err)
	// }

	file, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("error opening/creating file: %v", err)
	}
	defer file.Close()

	enc := yaml.NewEncoder(file)

	err = enc.Encode(s)
	if err != nil {
		log.Fatalf("error encoding: %v", err)
	}

	return nil
}

// validate checks configuration values
func (s *Settings) Validate() error {

	if s.Interval <= 0 {
		return fmt.Errorf("invalid interval '%s': must be positive", s.Interval)
	}

	if s.Sensor == "" {
		return fmt.Errorf("sensor can't be empty")
	}

	return nil
}
