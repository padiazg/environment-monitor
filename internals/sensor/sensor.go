package sensor

import (
	"fmt"
)

type SensorInterface interface {
	Init() error
	Read() (*Reading, error)
	Close() error
}

type SensorInvalidSettingsError struct {
	Name    string
	Message string
}

func (e *SensorInvalidSettingsError) Error() string {
	return fmt.Sprintf("%s invalid settings: %s", e.Name, e.Message)
}

type SensorCommunicationError struct {
	Name    string
	Message string
}

func (e *SensorCommunicationError) Error() string {
	return fmt.Sprintf("%s communication error: %s", e.Name, e.Message)
}

type SensorUnknownModeError struct {
	Name    string
	Message string
}

func (e *SensorUnknownModeError) Error() string {
	return fmt.Sprintf("%s unknown mode: %s", e.Name, e.Message)
}

type SensorReadError struct {
	Name    string
	Message string
}

func (e *SensorReadError) Error() string {
	return fmt.Sprintf("%s reading: %s", e.Name, e.Message)
}
