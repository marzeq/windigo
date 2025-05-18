package curve

import (
	"fmt"

	"github.com/marzeq/windigo/sensor"
)

type Curve struct {
	IsStep          bool
	Sensors         []string // sensor names
	AggregateFunc   AggregateFunc
	Points          []CurvePoint
	Hysteresis      float64
	Period          int64 // in seconds
	PreviousPercent float64
	PreviousTemp    float64
}

type CurvePoint struct {
	Temp    float64
	Percent float64
}

type Curves = map[string]*Curve // mapping asigned name by user to curve

type AggregateFunc func(readings []float64) float64

var AggregateAvg AggregateFunc = func(readings []float64) float64 {
	if len(readings) == 0 {
		return 0
	}
	sum := 0.0
	for _, reading := range readings {
		sum += reading
	}
	return sum / float64(len(readings))
}

var AggregateMax AggregateFunc = func(readings []float64) float64 {
	if len(readings) == 0 {
		return 0
	}
	max := readings[0]
	for _, reading := range readings {
		if reading > max {
			max = reading
		}
	}
	return max
}

var AggregateMin AggregateFunc = func(readings []float64) float64 {
	if len(readings) == 0 {
		return 0
	}
	min := readings[0]
	for _, reading := range readings {
		if reading < min {
			min = reading
		}
	}
	return min
}

func (c Curve) GetPointFromTemp(temp float64) (float64, bool) {
	if c.IsStep {
		for i := range c.Points {
			if temp >= c.Points[i].Temp && (i == len(c.Points)-1 || temp <= c.Points[i+1].Temp) {
				return c.Points[i].Percent, true
			}
		}

		return 0, false
	}

	for i := range len(c.Points) - 1 {
		if c.Points[i].Temp <= temp && c.Points[i+1].Temp >= temp {
			return interpolate(c.Points[i], c.Points[i+1], temp), true
		}
	}

	return 0, false
}

func (c Curve) GetAggregateTemp(sensors sensor.Sensors) (float64, error) {
	temps := []float64{}
	for sensorName, sensor := range sensors {
		sensorTemp, err := sensor.ReadTemperature()
		if err != nil {
			return 0, fmt.Errorf("error reading temperature from sensor '%s': %w", sensorName, err)
		}
		temps = append(temps, sensorTemp)
	}
	if len(temps) == 0 {
		return 0, fmt.Errorf("no valid temperatures for curve '%s'", c.Sensors)
	}
	return c.AggregateFunc(temps), nil
}

func interpolate(p1, p2 CurvePoint, temp float64) float64 {
	if p1.Temp == p2.Temp {
		return p1.Percent
	}
	return p1.Percent + (p2.Percent-p1.Percent)*(temp-p1.Temp)/(p2.Temp-p1.Temp)
}
