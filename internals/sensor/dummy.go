package sensor

import (
	"github.com/padiazg/environment-monitor-daemon/internals/models/settings"
)

const (
	dummyName = "Dummy"
)

type Dummy struct {
	settings *settings.DummySettings
}

var _ SensorInterface = (*Dummy)(nil)

func NewDummy(settings *settings.DummySettings) *Dummy {
	return &Dummy{settings: settings}
}

func (d *Dummy) Init() error {
	if err := d.settings.Validate(); err != nil {
		return &SensorInvalidSettingsError{Name: dummyName, Message: err.Error()}
	}

	return nil
}

func (d *Dummy) Read() (*Reading, error) {
	return &Reading{
		MassPM1:  d.settings.PM1,
		MassPM25: d.settings.PM25,
		MassPM10: d.settings.PM10,
	}, nil
}

func (d *Dummy) Close() error {
	return nil
}
