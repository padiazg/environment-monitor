# Environment monitor

An environment monitoring application that can run as a command line application or as a Linux daemon.

This application is meant to be run in a Raspberry Pi

This project is intended to be used to colaborate with the [Aire Libre](http://airelib.re/) project

## Build
You can build the project either in the sbc or a desktop
```shell
# for arm32
$ make build-32

# for arm64
$ make build-64

# for desktop
$ make build
```

**Copy the binary to the target folder**
```shell
$ scp environment-monitor-daemon-arm32 root@pizerow.local:/usr/local/bin/
```

## Run
```shell
# create a sample configuration file
$ environment-monitor save-default

# rename
$ mv .environment-monitor.default.yaml .environment-monitor.yaml
```

Update the config file accordingly. The Dummy sensor should work fine for testing if you don't have a phisical sensor, it will respond with the values set in the config.
```yaml
url: http://api/path
sensor: dummy
apikey: abcdef0123456789
source: MySensor/1
description: My first sensor
lat: -25.376054
lon: -57.545052
interval: 5s
zh07:
    mode: qa
    port: /dev/serial0
dummy:
    pm1: 1
    pm25: 2.5
    pm10: 10
```

```shell
$ environment-monitor-daemon-arm32
```

## Parameters