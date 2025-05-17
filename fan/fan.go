package fan

import (
	"os"
	"strconv"
	"strings"
)

type Fan struct {
	Name       string // name of the fan
	Curve      string // curve name
	InputFile  string
	EnableFile string
	ReadFile   string
}

func NewFan(name, curve, inputFile, enableFile, readFile string) Fan {
	return Fan{
		Name:       name,
		Curve:      curve,
		InputFile:  inputFile,
		EnableFile: enableFile,
		ReadFile:   readFile,
	}
}

type Fans = map[string]Fan // mapping asigned name by user to fan

func (f Fan) ReadSpeed() (int, error) {
	data, err := os.ReadFile(f.ReadFile)
	if err != nil {
		return 0, err
	}
	s := strings.TrimSpace(string(data))
	speed, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}
	return speed, nil
}

func (f *Fan) SetManualMode() error {
	data := []byte("1")
	err := os.WriteFile(f.EnableFile, data, 0644)
	if err != nil {
		return err
	}
	return nil
}

func (f *Fan) SetSpeed(speed float64) error {
	if speed < 0 {
		speed = 0
	} else if speed > 100 {
		speed = 100
	}
	speedInt := int(speed * 255 / 100)

	speedStr := strconv.Itoa(speedInt)
	err := os.WriteFile(f.InputFile, []byte(speedStr), 0644)
	if err != nil {
		return err
	}

	return nil
}
