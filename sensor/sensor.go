package sensor

import (
	"fmt"
	"os"
)

type Sensor struct {
	Name string
	Path string
}

type Sensors = map[string]Sensor // mapping asigned name by user to sensor path

func (s Sensor) ReadTemperature() (float64, error) {
	data, err := os.ReadFile(s.Path)
	if err != nil {
		return 0, err
	}

	var temperature float64
	_, err = fmt.Sscanf(string(data), "%f", &temperature)
	if err != nil {
		return 0, err
	}
	temperature /= 1000 // Convert from millidegree Celsius to degree Celsius

	return temperature, nil
}
