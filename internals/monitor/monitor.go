package monitor

import (
	"strings"
	"sync"
	"time"

	"github.com/padiazg/environment-monitor-daemon/internals/models/settings"
	"github.com/padiazg/environment-monitor-daemon/internals/sensor"
)

type data struct {
	reading sensor.Reading
	err     error
}

type Monitor struct {
	settings    *settings.Settings
	readingChan chan data
	doneChan    chan bool
	sensor      sensor.SensorInterface
	ticker      *time.Ticker
	mu          sync.RWMutex
	stopped     bool
}

func NewMonitor(settings *settings.Settings) *Monitor {
	var s sensor.SensorInterface
	switch strings.ToLower(settings.Sensor) {
	case "zh07":
		s = sensor.NewZH07(settings.ZH07)
	case "dummy":
		s = sensor.NewDummy(settings.Dummy)
	default:
		s = nil
	}

	return &Monitor{
		settings: settings,
		sensor:   s,
	}
}

func (m *Monitor) Init() (chan data, func(), error) {
	if err := m.sensor.Init(); err != nil {
		return nil, nil, err
	}

	m.readingChan = make(chan data, 1)
	m.doneChan = make(chan bool)

	return m.readingChan,
		func() { m.Stop() },
		nil
}

func (m *Monitor) Stop() {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.stopped {
		return
	}
	m.stopped = true

	select {
	case m.doneChan <- true:
	default:
	}

	close(m.doneChan)
	close(m.readingChan)
	m.sensor.Close()
}

func (m *Monitor) Run() {
	go func() {
		m.ticker = time.NewTicker(m.settings.Interval)
		defer m.ticker.Stop()

		for {
			select {
			case <-m.ticker.C:
				r, err := m.sensor.Read()
				if err != nil {
					m.mu.RLock()
					if !m.stopped {
						m.readingChan <- data{err: err}
					}
					m.mu.RUnlock()
					continue
				}

				m.mu.RLock()
				if !m.stopped {
					m.readingChan <- data{reading: *r}
				}
				m.mu.RUnlock()

			case <-m.doneChan:
				return
			}
		}
	}()
}
