package config

import (
	"fmt"

	"github.com/marzeq/mconf"
	"github.com/marzeq/mconf/mconf_values"
	"github.com/marzeq/windigo/curve"
	"github.com/marzeq/windigo/device"
	"github.com/marzeq/windigo/fan"
	"github.com/marzeq/windigo/sensor"
)

type ConfigFile = map[string]mconf_values.MconfValue

type Config = struct {
	Devices device.Devices
	Sensors sensor.Sensors
	Curves  curve.Curves
	Fans    fan.Fans
	Path    string
}

func ReadConfig(configFile string) (Config, error) {
	config, _, err := mconf.ParseFromFile(configFile)
	if err != nil {
		return Config{}, fmt.Errorf("failed to parse config file: %w", err)
	}

	devices, err := GetDevices(config)
	if err != nil {
		return Config{}, fmt.Errorf("failed to get devices: %w", err)
	}

	sensors, err := GetSensors(config, devices)
	if err != nil {
		return Config{}, fmt.Errorf("failed to get sensors: %w", err)
	}

	curves, err := GetCurves(config, sensors)
	if err != nil {
		return Config{}, fmt.Errorf("failed to get curves: %w", err)
	}

	fans, err := GetFans(config, devices, curves)
	if err != nil {
		return Config{}, fmt.Errorf("failed to get fans: %w", err)
	}

	return Config{
		Devices: devices,
		Sensors: sensors,
		Curves:  curves,
		Fans:    fans,
		Path:    configFile,
	}, nil
}
