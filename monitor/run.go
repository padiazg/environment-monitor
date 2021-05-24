package monitor

import (
	"context"
	"os"
	"time"

	"github.com/d2r2/go-logger"
	"github.com/padiazg/environment-monitor-daemon/config"
)

var (
	lg = logger.NewPackageLogger("monitor", logger.InfoLevel)
)

// Run monitoring main loop
func Run(ctx context.Context, c *config.Config) error {
	c.Init(os.Args)
	lg.Info("Start monitoring")

	bus := GetAQISensor(c)
	defer bus.Close()

	lg.Infof("Starting ticker. Triggering every %s", c.Tick)
	var ticker *time.Ticker = time.NewTicker(c.Tick)

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			if err := aqiMeasurement(c); err != nil {
				lg.Errorf("Reading measurement.", err)
			}
		} // select
	} // for ...
} // Run ...
