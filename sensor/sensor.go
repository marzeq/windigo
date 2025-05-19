package sensor

import (
	"fmt"
	"os"
)

type Sensor struct {
	Name     string
	Path     string
	Offset   float64
	MinValue float64
	MaxValue float64
}

type Sensors = map[string]Sensor // mapping asigned name by user to sensor path

func (s Sensor) ReadRealTemperature() (float64, error) {
	data, err := os.ReadFile(s.Path)
	if err != nil {
		return 0, err
	}

	var temperature float64
	_, err = fmt.Sscanf(string(data), "%f", &temperature)
	if err != nil {
		return 0, err
	}
	temperature /= 1000 // convert from millidegree Celsius to degree Celsius

	return temperature, nil
}

func (s Sensor) ReadTemperature() (float64, error) {
	temperature, err := s.ReadRealTemperature()
	if err != nil {
		return 0, err
	}
	temperature += s.Offset

	if temperature < s.MinValue {
		return s.MinValue, nil
	}
	if temperature > s.MaxValue {
		return s.MaxValue, nil
	}

	return temperature + s.Offset, nil
}
