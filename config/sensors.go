package config

import (
	"fmt"
	"os"

	"github.com/marzeq/mconf/mconf_values"
	"github.com/marzeq/windigo/device"
	"github.com/marzeq/windigo/sensor"
)

func GetSensors(config ConfigFile, devices device.Devices) (sensor.Sensors, error) {
	snsrs, exists, typeOk := mconf_values.ObjGetObject(config, "sensors")
	if !exists {
		return nil, fmt.Errorf("section 'sensors' not found")
	} else if !typeOk {
		return nil, fmt.Errorf("section 'sensors' must be an object")
	}
	sensors := make(sensor.Sensors)

	for device, sensorsForDevice := range snsrs {
		if _, ok := devices[device]; !ok {
			return nil, fmt.Errorf("device '%s' not found", device)
		}

		sensorsForDeviceVal, ok := mconf_values.UnwrapObject(sensorsForDevice)
		if !ok {
			return nil, fmt.Errorf("sensors for device '%s' must be an object", device)
		}

		for sensorName, sens := range sensorsForDeviceVal {
			sensorVal, ok := mconf_values.UnwrapObject(sens)
			if !ok {
				return nil, fmt.Errorf("sensor '%s' for device '%s' must be an object", sensorName, device)
			}

			readFile, exists, typeOk := mconf_values.ObjGetString(sensorVal, "readfile")
			if !exists {
				return nil, fmt.Errorf("sensor '%s' for device '%s' must have 'readfile' property", sensorName, device)
			} else if !typeOk {
				return nil, fmt.Errorf("sensor '%s' for device '%s' 'readfile' property must be a string", sensorName, device)
			}

			readFile = fmt.Sprintf("%s/%s", devices[device].Path, readFile)
			if _, err := os.Stat(readFile); os.IsNotExist(err) {
				return nil, fmt.Errorf("sensor '%s' for device '%s' 'readfile' property must be a valid file", sensorName, device)
			}

			offset, exists, typeOk := mconf_values.ObjGetFloat(sensorVal, "offset")
			if !exists {
				offset = 0.0
			} else if !typeOk {
				offsetInt, _, typeOk := mconf_values.ObjGetInt(sensorVal, "offset")
				if !typeOk {
					return nil, fmt.Errorf("sensor '%s' for device '%s' 'offset' property must be a number", sensorName, device)
				}
				offset = float64(offsetInt)
			}

			minValue, exists, typeOk := mconf_values.ObjGetFloat(sensorVal, "min")
			if !exists {
				minValue = -10000.0
			} else if !typeOk {
				minValueInt, _, typeOk := mconf_values.ObjGetInt(sensorVal, "min")
				if !typeOk {
					return nil, fmt.Errorf("sensor '%s' for device '%s' 'min' property must be a number", sensorName, device)
				}
				minValue = float64(minValueInt)
			}

			maxValue, exists, typeOk := mconf_values.ObjGetFloat(sensorVal, "max")
			if !exists {
				maxValue = 10000.0
			} else if !typeOk {
				maxValueInt, _, typeOk := mconf_values.ObjGetInt(sensorVal, "max")
				if !typeOk {
					return nil, fmt.Errorf("sensor '%s' for device '%s' 'max' property must be a number", sensorName, device)
				}
				maxValue = float64(maxValueInt)
			}

			sensors[sensorName] = sensor.Sensor{
				Name:     sensorName,
				Path:     readFile,
				Offset:   offset,
				MinValue: minValue,
				MaxValue: maxValue,
			}
		}
	}

	return sensors, nil
}
