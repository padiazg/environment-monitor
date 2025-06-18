package sensor

type SensorInterface interface {
	Init() error
	Read() (*Reading, error)
	Close() error
}
