package sensor

import (
	"bufio"

	"github.com/padiazg/environment-monitor-daemon/internals/models/settings"
	"github.com/padiazg/go-zh07"
	"github.com/tarm/serial"
)

const (
	zh07Name = "ZH07"
)

type ZH07 struct {
	settings   *settings.ZH07Settings
	sensor     zh07.SensorInterface
	serialPort *serial.Port
}

var _ SensorInterface = (*ZH07)(nil)

func NewZH07(settings *settings.ZH07Settings) *ZH07 {
	return &ZH07{settings: settings}
}

func (z *ZH07) Init() error {
	if err := z.settings.Validate(); err != nil {
		return &SensorInvalidSettingsError{Name: zh07Name, Message: err.Error()}
	}

	// open TTY port
	s, err := serial.OpenPort(&serial.Config{
		Name:     z.settings.Port,
		Baud:     9600,
		Parity:   serial.ParityNone,
		StopBits: serial.Stop1,
	})

	if err != nil {
		return &SensorCommunicationError{Name: zh07Name, Message: err.Error()}
	}

	z.serialPort = s

	rw := bufio.NewReadWriter(
		bufio.NewReader(s),
		bufio.NewWriter(s),
	)

	// create a sensor instance
	switch z.settings.Mode {
	case "qa":
		z.sensor = zh07.NewZH07q(&zh07.Config{RW: rw})
	case "initiative":
		z.sensor = zh07.NewZH07i(&zh07.Config{RW: rw})
	default:
		return &SensorUnknownModeError{Name: zh07Name, Message: z.settings.Mode}
	}

	return nil
}

func (z *ZH07) Read() (*Reading, error) {
	r, err := z.sensor.Read()
	if err != nil {
		return nil, &SensorReadError{Name: zh07Name, Message: err.Error()}
	}

	return &Reading{
		MassPM1:  float32(r.PM1),
		MassPM25: float32(r.PM25),
		MassPM10: float32(r.PM10),
	}, nil
}

func (z *ZH07) Close() error {
	if z.serialPort != nil {
		return z.serialPort.Close()
	}

	return nil
}
