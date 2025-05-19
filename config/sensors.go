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

		for sensorName, sens := range sensorsForDeviceVal {
			sensorVal, ok := mconf_values.MconfUnwrapObject(sens)
			if !ok {
				return nil, fmt.Errorf("sensor '%s' for device '%s' must be an object", sensorName, device)
			}
			readFileVal, ok := sensorVal["readfile"]
			if !ok {
				return nil, fmt.Errorf("sensor '%s' for device '%s' must have 'readfile' property", sensorName, device)
			}
			readFile, ok := mconf_values.MconfUnwrapString(readFileVal)
			if !ok {
				return nil, fmt.Errorf("sensor '%s' for device '%s' 'readfile' property must be a string", sensorName, device)
			}
			readFile = fmt.Sprintf("%s/%s", devices[device].Path, readFile)
			if _, err := os.Stat(readFile); os.IsNotExist(err) {
				return nil, fmt.Errorf("sensor '%s' for device '%s' 'readfile' property must be a valid file", sensorName, device)
			}
			offset := 0.0
			if offsetVal, ok := sensorVal["offset"]; ok {
				offsetF, ok := mconf_values.MconfUnwrapFloat(offsetVal)
				if !ok {
					offsetInt, ok := mconf_values.MconfUnwrapInt(offsetVal)
					if !ok {
						return nil, fmt.Errorf("sensor '%s' for device '%s' 'offset' property must be a number", sensorName, device)
					}
					offset = float64(offsetInt.Int64())
				} else {
					offset, _ = offsetF.Float64()
				}
			}

			minValue := -10000.0
			if minValueVal, ok := sensorVal["min"]; ok {
				minValueF, ok := mconf_values.MconfUnwrapFloat(minValueVal)
				if !ok {
					minValueInt, ok := mconf_values.MconfUnwrapInt(minValueVal)
					if !ok {
						return nil, fmt.Errorf("sensor '%s' for device '%s' 'min' property must be a number", sensorName, device)
					}
					minValue = float64(minValueInt.Int64())
				} else {
					minValue, _ = minValueF.Float64()
				}
			}

			maxValue := 10000.0
			if maxValueVal, ok := sensorVal["max"]; ok {
				maxValueF, ok := mconf_values.MconfUnwrapFloat(maxValueVal)
				if !ok {
					maxValueInt, ok := mconf_values.MconfUnwrapInt(maxValueVal)
					if !ok {
						return nil, fmt.Errorf("sensor '%s' for device '%s' 'max' property must be a number", sensorName, device)
					}
					maxValue = float64(maxValueInt.Int64())
				} else {
					maxValue, _ = maxValueF.Float64()
				}
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
