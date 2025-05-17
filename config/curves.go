package config

import (
	"fmt"

	"github.com/marzeq/mconf/mconf_values"
	"github.com/marzeq/windigo/curve"
	"github.com/marzeq/windigo/sensor"
)

func GetCurves(config ConfigFile, sensors sensor.Sensors) (curve.Curves, error) {
	curvesGot, ok := config["curves"]
	if !ok {
		return nil, fmt.Errorf("section 'curves' not found")
	}

	curvesVal, ok := mconf_values.MconfUnwrapObject(curvesGot)
	if !ok {
		return nil, fmt.Errorf("section 'curves' must be an object")
	}
	curves := make(curve.Curves)

	for name, cur := range curvesVal {
		curveVal, ok := mconf_values.MconfUnwrapObject(cur)
		if !ok {
			return nil, fmt.Errorf("curve '%s' must be an object", name)
		}
		typeGot, ok := curveVal["type"]
		if !ok {
			return nil, fmt.Errorf("curve '%s' must have 'type' attribute", name)
		}

		typeVal, ok := mconf_values.MconfUnwrapString(typeGot)
		if !ok {
			return nil, fmt.Errorf("curve '%s' 'type' must be a string", name)
		}
		if typeVal != "step" && typeVal != "linear" {
			return nil, fmt.Errorf("curve '%s' 'type' must be 'step' or 'linear'", name)
		}
		isStep := typeVal == "step"

		sensorGot, ok := curveVal["sensor"]
		if !ok {
			return nil, fmt.Errorf("curve '%s' must have 'sensor' attribute", name)
		}
		sensorVal, ok := mconf_values.MconfUnwrapString(sensorGot)
		if !ok {
			return nil, fmt.Errorf("curve '%s' 'sensor' must be a string", name)
		}
		sensorKey := ""
		for k := range sensors {
			if k == sensorVal {
				sensorKey = k
				break
			}
		}
		if sensorKey == "" {
			return nil, fmt.Errorf("curve '%s' 'sensor' must be a valid sensor", name)
		}

		pointsGot, ok := curveVal["points"]
		if !ok {
			return nil, fmt.Errorf("curve '%s' must have 'points' attribute", name)
		}
		pointsVal, ok := mconf_values.MconfUnwrapList(pointsGot)
		if !ok {
			return nil, fmt.Errorf("curve '%s' 'points' must be a list", name)
		}

		points := make([]curve.CurvePoint, len(pointsVal))

		for i, point := range pointsVal {
			pointVal, ok := mconf_values.MconfUnwrapList(point)
			if !ok {
				return nil, fmt.Errorf("curve '%s' 'points' must be a list of lists", name)
			}

			if len(pointVal) != 2 {
				return nil, fmt.Errorf("curve '%s' 'points' must be a list of lists of length 2", name)
			}

			tempGot := pointVal[0]
			percentGot := pointVal[1]

			tempVal, ok := mconf_values.MconfUnwrapFloat(tempGot)
			temp := 0.0
			if !ok {
				tempVal, ok := mconf_values.MconfUnwrapInt(tempGot)
				if !ok {
					return nil, fmt.Errorf("curve '%s' 'points' must be a list of lists of numbers", name)
				}
				tempInt := tempVal.Int64()
				temp = float64(tempInt)
			} else {
				temp, _ = tempVal.Float64()
			}

			if temp < 0 {
				return nil, fmt.Errorf("curve '%s' 'points' must be >= 0", name)
			}

			percentVal, ok := mconf_values.MconfUnwrapFloat(percentGot)
			percent := 0.0
			if !ok {
				percentVal, ok := mconf_values.MconfUnwrapInt(percentGot)
				if !ok {
					return nil, fmt.Errorf("curve '%s' 'points' must be a list of lists of numbers", name)
				}
				percentInt := percentVal.Int64()
				percent = float64(percentInt)
			} else {
				percent, _ = percentVal.Float64()
			}

			points[i] = curve.CurvePoint{
				Temp:    temp,
				Percent: percent,
			}
		}

		hysteresisGot, ok := curveVal["hysteresis"]
		hysteresis := 0.0
		if ok {
			hysteresisVal, ok := mconf_values.MconfUnwrapFloat(hysteresisGot)
			hysteresis = 0.0
			if !ok {
				hysteresisVal, ok := mconf_values.MconfUnwrapInt(hysteresisGot)
				if !ok {
					return nil, fmt.Errorf("curve '%s' 'hysteresis' must be a float", name)
				}
				hysteresisInt := hysteresisVal.Int64()
				hysteresis = float64(hysteresisInt)
			} else {
				hysteresis, _ = hysteresisVal.Float64()
			}

			if hysteresis < 0 {
				return nil, fmt.Errorf("curve '%s' 'hysteresis' must be >= 0", name)
			}
		}

		periodGot, ok := curveVal["period"]
		period := int64(1)

		if ok {
			readEveryVal, ok := mconf_values.MconfUnwrapInt(periodGot)
			period = 1
			if !ok {
				return nil, fmt.Errorf("curve '%s' 'readevery' must be an int", name)
			}
			period = readEveryVal.Int64()
			if period < 1 {
				return nil, fmt.Errorf("curve '%s' 'readevery' must be >= 1", name)
			}
		}

		curves[name] = &curve.Curve{
			IsStep:     isStep,
			Sensor:     sensorKey,
			Points:     points,
			Hysteresis: hysteresis,
			Period:     period,
		}
	}

	return curves, nil
}
