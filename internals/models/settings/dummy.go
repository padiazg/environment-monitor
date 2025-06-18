package settings

import "fmt"

type DummySettings struct {
	PM1  float32 // Mass Concentration PM1.0 [μg/m3]
	PM25 float32 // Mass Concentration PM2.5 [μg/m3]
	PM10 float32 // Mass Concentration PM10 [μg/m3]
}

func (d *DummySettings) Validate() error {
	if d == nil {
		return fmt.Errorf("Dummy settings not set")
	}

	if d.PM1 == 0 {
		return fmt.Errorf("PM1 cannot be empty")
	}

	if d.PM25 == 0 {
		return fmt.Errorf("PM25 cannot be empty")
	}

	if d.PM10 == 0 {
		return fmt.Errorf("PM10 cannot be empty")
	}

	return nil
}
