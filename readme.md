# Environment monitor

An environment monitoring application that can run as a command line application or as a Linux daemon.

This application is meant to be run in a Raspberry Pi

This project is intended to be used to colaborate with the [Aire Libre](http://airelib.re/) project

# Build

### Cross compile
Maybe you would prefer to build the binary on your desktop and then copy the binary to the pi to be runned.

Make sure the target folder `/usr/local/environment-monitor` exists in the target RPi, create if not
```bash
$ ssh root@pizerow.local 
$ mkdir -p /usr/local/environment-monitor
```

Go back to your desktop, build and copy the binary to the target
```bash
OOS=linux GOARCH=arm GOARM=6 go build
scp environment-monitor-daemon root@pizerow.local:/usr/local/environment-monitor
```