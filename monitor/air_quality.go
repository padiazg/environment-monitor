package monitor

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/d2r2/go-i2c"
	"github.com/d2r2/go-logger"
	"github.com/padiazg/environment-monitor-daemon/config"
	"github.com/padiazg/go-sps30"
)

var (
	bus    *i2c.I2C
	sensor *sps30.SPS30 // create an interface for sensor
)

func GetAQISensor(c *config.Config) *i2c.I2C {
	if bus != nil {
		return bus
	}

	var err error
	bus, err = i2c.NewI2C(0x69, 1)
	if err != nil {
		lg.Error("Creating I2C bus connection. ", err)
	}
	logger.ChangePackageLogLevel("i2c", logger.InfoLevel)

	return bus
} // GetAQISensor ...

func GetSensor() *sps30.SPS30 {
	if sensor != nil {
		return sensor
	}
	sensor = sps30.NewSPS30(bus)
	return sensor
} // GetSensor ...

func aqiMeasurement(c *config.Config) error {
	// fmt.Println("Started measurement ", t)
	if err := GetSensor().StartMeasurement(); err != nil {
		return err
	}
	time.Sleep(2 * time.Second)

	// Read if there is new data available
	var dataReady int = GetSensor().ReadDataReady()
	lg.Infof("Data ready: %v", dataReady)

	if dataReady == 1 {
		// Read measurements
		m, err := GetSensor().ReadMeasurement()
		if err != nil {
			lg.Errorf("read-measurement: %v", err)
			return err
		}
		// print to console
		formatMeasurementHuman(m)
		postMeasurement(m, c)
		// lg.Infof("data: %v", data)
	}

	// Stop measurement, go to idle-mode again
	if err := GetSensor().StopMeasurement(); err != nil {
		return err
	}
	lg.Infof("Done measurement. Sleeping.")

	return nil
} // aqiMeasurement ...

// Formats data for human reading
func formatMeasurementHuman(m *sps30.AirQualityReading) {
	lg.Infof("pm0.5 count: %8s", fmt.Sprintf("%4.3f", m.NumberPM05))
	lg.Infof("pm1   count: %8s ug: %6s", fmt.Sprintf("%4.3f", m.NumberPM1), fmt.Sprintf("%2.3f", m.MassPM1))
	lg.Infof("pm2.5 count: %8s ug: %6s", fmt.Sprintf("%4.3f", m.NumberPM25), fmt.Sprintf("%2.3f", m.MassPM25))
	lg.Infof("pm4   count: %8s ug: %6s", fmt.Sprintf("%4.3f", m.NumberPM4), fmt.Sprintf("%2.3f", m.MassPM4))
	lg.Infof("pm10  count: %8s ug: %6s", fmt.Sprintf("%4.3f", m.NumberPM10), fmt.Sprintf("%2.3f", m.MassPM10))
	lg.Infof("pm_typ: %4.3f", m.TypicalParticleSize)
} // formatMeasurementHuman ...

func postMeasurement(m *sps30.AirQualityReading, c *config.Config) error {
	recorded := time.Now()
	r0, _ := json.Marshal([]struct {
		Sensor      string  `json:"sensor"`      // Model of the device
		Source      string  `json:"source"`      // Name used to identify the device
		Description string  `json:"description"` // User friendly name to identify the device
		Pm1Dot0     int     `json:"pm1dot0"`     // Concentration of PM1.0 inhalable particles per ug/m3
		Pm2Dot5     int     `json:"pm2dot5"`     // Concentration of PM2.5 inhalable particles per ug/m3
		Pm10        int     `json:"pm10"`        // Concentration of PM10 inhalable particles per ug/m3
		Longitude   float64 `json:"longitude"`   // Physical longitude coordinate of the device
		Latitude    float64 `json:"latitude"`    // Physical latitude coordinate of the device
		Recorded    string  `json:"recorded"`    // Date and time for when these values were measured
	}{{
		Sensor:      c.AQISensor,
		Source:      c.Source,
		Description: c.Description,
		Latitude:    c.Latitude,
		Longitude:   c.Longitude,
		Pm1Dot0:     int(m.MassPM10),
		Pm2Dot5:     int(m.MassPM25),
		Pm10:        int(m.MassPM10),
		Recorded:    recorded.Format(time.RFC3339Nano), // "2021-01-18T22:06:54.673
	}})

	lg.Info("r0 =>", string(r0))

	// method := "POST"
	payload := strings.NewReader(string(r0))

	client := &http.Client{}
	req, err := http.NewRequest("POST", c.Url, payload)

	if err != nil {
		lg.Error(err)
		return err
	}

	req.Header.Add("X-API-Key", c.ApiKey)
	req.Header.Add("Content-Type", "application/javascript")

	res, err := client.Do(req)
	if err != nil {
		lg.Error(err)
		return err
	}

	defer res.Body.Close()
	// body, err := ioutil.ReadAll(res.Body)
	// fmt.Println(string(body))

	return nil
} // postMeasurement ...
