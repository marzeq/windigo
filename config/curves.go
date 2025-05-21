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

	curvesVal, ok := mconf_values.UnwrapObject(curvesGot)
	if !ok {
		return nil, fmt.Errorf("section 'curves' must be an object")
	}
	curves := make(curve.Curves)

	for name, cur := range curvesVal {
		curveVal, ok := mconf_values.UnwrapObject(cur)
		if !ok {
			return nil, fmt.Errorf("curve '%s' must be an object", name)
		}

		tpe, exists, okType := mconf_values.ObjGetString(curveVal, "type")
		if !exists {
			return nil, fmt.Errorf("curve '%s' must have 'type' attribute", name)
		} else if !okType {
			return nil, fmt.Errorf("curve '%s' 'type' must be a string", name)
		}

		var curveType curve.CurveType
		switch tpe {
		case "step":
			curveType = curve.CurveTypeStep
		case "linear":
			curveType = curve.CurveTypeLinear
		default:
			return nil, fmt.Errorf("curve '%s' 'type' must be 'step' or 'linear'", name)
		}

		curveSensors := []string{}

		processSensor := func(sensorGot mconf_values.MconfValue) error {
			sensorVal, ok := mconf_values.UnwrapString(sensorGot)
			if !ok {
				return fmt.Errorf("curve '%s' 'sensor' must be a string", name)
			}
			sensorKey := ""
			for k := range sensors {
				if k == sensorVal {
					sensorKey = k
					break
				}
			}
			if sensorKey == "" {
				return fmt.Errorf("curve '%s' 'sensor' must be a valid sensor", name)
			}
			curveSensors = append(curveSensors, sensorKey)
			return nil
		}

		sensorGot, ok := curveVal["sensor"]
		if ok {
			if err := processSensor(sensorGot); err != nil {
				return nil, err
			}
		} else {
			sensors, exists, typeOk := mconf_values.ObjGetList(curveVal, "sensors")
			if !exists {
				return nil, fmt.Errorf("curve '%s' must have 'sensor' or 'sensors' attribute", name)
			} else if !typeOk {
				return nil, fmt.Errorf("curve '%s' 'sensors' must be a list", name)
			}
			for _, sensorGot := range sensors {
				if err := processSensor(sensorGot); err != nil {
					return nil, err
				}
			}
		}

		if len(curveSensors) == 0 {
			return nil, fmt.Errorf("curve '%s' must have 'sensor' or 'sensors' attribute", name)
		}

		aggreggate := "avg"
		aggregateGot, ok := curveVal["aggregate"]
		if ok {
			aggregateVal, ok := mconf_values.UnwrapString(aggregateGot)
			if !ok {
				return nil, fmt.Errorf("curve '%s' 'aggregate' must be a string", name)
			}
			if aggregateVal != "avg" && aggregateVal != "max" && aggregateVal != "min" {
				return nil, fmt.Errorf("curve '%s' 'aggregate' must be 'avg', 'max' or 'min'", name)
			}
			aggreggate = aggregateVal
		}
		var aggregateFunc curve.AggregateFunc
		switch aggreggate {
		case "avg":
			aggregateFunc = curve.AggregateAvg
		case "max":
			aggregateFunc = curve.AggregateMax
		case "min":
			aggregateFunc = curve.AggregateMin
		}

		pointsGot, ok := curveVal["points"]
		if !ok {
			return nil, fmt.Errorf("curve '%s' must have 'points' attribute", name)
		}
		pointsVal, ok := mconf_values.UnwrapList(pointsGot)
		if !ok {
			return nil, fmt.Errorf("curve '%s' 'points' must be a list", name)
		}

		points := make([]curve.CurvePoint, len(pointsVal))

		for i, point := range pointsVal {
			pointVal, ok := mconf_values.UnwrapList(point)
			if !ok {
				return nil, fmt.Errorf("curve '%s' 'points' must be a list of lists", name)
			}

			if len(pointVal) != 2 {
				return nil, fmt.Errorf("curve '%s' 'points' must be a list of lists of length 2", name)
			}

			tempGot := pointVal[0]
			percentGot := pointVal[1]

			temp, ok := mconf_values.UnwrapFloat(tempGot)
			if !ok {
				tempInt, ok := mconf_values.UnwrapInt(tempGot)
				if !ok {
					return nil, fmt.Errorf("curve '%s' 'points' must be a list of lists of numbers", name)
				}
				temp = float64(tempInt)
			}

			if temp < 0 {
				return nil, fmt.Errorf("curve '%s' 'points' must be >= 0", name)
			}

			if i != 0 && temp < points[i-1].Temp {
				return nil, fmt.Errorf("curve '%s' 'points' must be in ascending order", name)
			}

			percent, ok := mconf_values.UnwrapFloat(percentGot)
			if !ok {
				percentInt, ok := mconf_values.UnwrapInt(percentGot)
				if !ok {
					return nil, fmt.Errorf("curve '%s' 'points' must be a list of lists of numbers", name)
				}
				percent = float64(percentInt)
			}

			points[i] = curve.CurvePoint{
				Temp:    temp,
				Percent: percent,
			}
		}

		hysteresis, exists, typeOk := mconf_values.ObjGetFloat(curveVal, "hysteresis")
		if !exists {
			hysteresis = 0.0
		} else if !typeOk {
			hysteresisInt, _, typeOk := mconf_values.ObjGetInt(curveVal, "hysteresis")
			if !typeOk {
				return nil, fmt.Errorf("curve '%s' 'hysteresis' must be a float or int", name)
			}
			hysteresis = float64(hysteresisInt)

			if hysteresis < 0 {
				return nil, fmt.Errorf("curve '%s' 'hysteresis' must be >= 0", name)
			}
		}

		period, exists, typeOk := mconf_values.ObjGetInt(curveVal, "period")
		if !exists {
			period = 1
		} else if !typeOk {
			if !ok {
				return nil, fmt.Errorf("curve '%s' 'readevery' must be an int", name)
			}
		}

		if period < 1 {
			return nil, fmt.Errorf("curve '%s' 'readevery' must be >= 1", name)
		}

		curves[name] = &curve.Curve{
			Type:          curveType,
			Sensors:       curveSensors,
			AggregateFunc: aggregateFunc,
			Points:        points,
			Hysteresis:    hysteresis,
			Period:        period,
		}
	}

	return curves, nil
}
