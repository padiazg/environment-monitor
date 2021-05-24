package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/d2r2/go-logger"
	"github.com/padiazg/environment-monitor-daemon/config"
	"github.com/padiazg/environment-monitor-daemon/monitor"
)

var lg = logger.NewPackageLogger("main", logger.InfoLevel)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	c := &config.Config{}

	defer func() {
		signal.Stop(signalChan)
		logger.FinalizeLogger()
		cancel()
	}() // defer func...

	go func() {
		for {
			select {
			case s := <-signalChan:
				switch s {
				case syscall.SIGINT, syscall.SIGTERM:
					lg.Info("Got SIGINT/SIGTERM, exiting.")
					cancel()
					os.Exit(1)
				case syscall.SIGHUP:
					lg.Info("Got SIGHUP, reloading configuration.")
					c.Init(os.Args)
				} // switch ...
			case <-ctx.Done():
				lg.Info("Done")
				os.Exit(1)
			} // select ...
		} // for ...
	}() // go func ...

	if err := monitor.Run(ctx, c); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
} // main ...
