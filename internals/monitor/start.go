package monitor

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/padiazg/environment-monitor-daemon/internals/models/settings"
)

func Start(settings *settings.Settings) {
	var (
		ctx, cancel = context.WithCancel(context.Background())
		signalChan  = make(chan os.Signal, 1)
	)

	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	defer func() {
		cancel()
		signal.Stop(signalChan)
	}()

	monitor := NewMonitor(settings)
	rc, done, err := monitor.Init()
	if err != nil {
		log.Fatalf("Initializing driver: %+v\n", err)
	}
	defer done()

	go func() {
		for r := range rc {
			if r.err != nil {
				log.Printf("error: %s\n", r.err.Error())
				continue
			}

			r.reading.Show()
		}
	}()

	monitor.Run()

	select {
	case s := <-signalChan:
		switch s {
		case syscall.SIGINT, syscall.SIGTERM:
			fmt.Printf("\nGot SIGINT/SIGTERM, exiting.\n")
			cancel()
		case syscall.SIGHUP:
			fmt.Println("Got SIGHUP, reloading configuration.")
			// if err := c.Init(os.Args); err != nil {
			// 	log.Printf("Error reloading config: %+v\n", err)
			// } else {
			// 	c.Show()
			// }
		}
	case <-ctx.Done():
		fmt.Println("Done")
		monitor.Stop()
	}
}
