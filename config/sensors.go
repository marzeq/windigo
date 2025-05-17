package config

import (
	"fmt"
	"os"

	"github.com/marzeq/mconf/mconf_values"
	"github.com/marzeq/windigo/device"
	"github.com/marzeq/windigo/sensor"
)

func GetSensors(config ConfigFile, devices device.Devices) (sensor.Sensors, error) {
	sensorsGot, ok := config["sensors"]
	if !ok {
		return nil, fmt.Errorf("section 'sensors' not found")
	}

	sensorsVal, ok := mconf_values.MconfUnwrapObject(sensorsGot)
	if !ok {
		return nil, fmt.Errorf("section 'sensors' must be an object")
	}
	sensors := make(sensor.Sensors)

	for device, sensorsForDevice := range sensorsVal {
		if _, ok := devices[device]; !ok {
			return nil, fmt.Errorf("device '%s' not found", device)
		}

		sensorsForDeviceVal, ok := mconf_values.MconfUnwrapObject(sensorsForDevice)
		if !ok {
			return nil, fmt.Errorf("sensors for device '%s' must be an object", device)
		}

		for sensorFile, sensorName := range sensorsForDeviceVal {
			sensorNameVal, ok := mconf_values.MconfUnwrapString(sensorName)
			if !ok {
				return nil, fmt.Errorf("sensor '%s' must be a string", sensorFile)
			}

			sensorPath := fmt.Sprintf("%s/%s", devices[device].Path, sensorFile)
			if _, err := os.Stat(sensorPath); os.IsNotExist(err) {
				return nil, fmt.Errorf("sensor '%s' not found in device '%s'", sensorFile, device)
			}

			if _, ok := sensors[sensorNameVal]; ok {
				return nil, fmt.Errorf("sensor '%s' already exists", sensorNameVal)
			}

			sensors[sensorNameVal] = sensor.Sensor{
				Name: sensorNameVal,
				Path: sensorPath,
			}
		}
	}

	return sensors, nil
}
