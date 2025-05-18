package daemon

import (
	"log"
	"time"

	"github.com/marzeq/windigo/config"
)

func Main(config config.Config) {
	log.Println("Starting windigo daemon")

	fatal := false
	for _, fan := range config.Fans {
		if err := fan.SetManualMode(); err != nil {
			log.Printf("%s", err.Error())
			fatal = true
			break
		}
	}

	if fatal {
		log.Println("Fatal error, exiting")
		return
	}

	t0 := time.Now().Unix()
	for {
		dt := time.Now().Unix() - t0
		for curveName, curve := range config.Curves {
			if dt%curve.Period != 0 {
				continue
			}

			temp, err := curve.GetAggregateTemp(config.Sensors)
			if err != nil {
				log.Printf("%v", err)
				continue
			}

			if curve.PreviousTemp != 0 && abs(temp-curve.PreviousTemp) < curve.Hysteresis {
				continue
			}
			percent, ok := curve.GetPointFromTemp(temp)
			if !ok {
				log.Printf("Error getting point from curve '%s'", curveName)
				continue
			}
			curve.PreviousTemp = temp

			setFor := ""
			for fanName, fan := range config.Fans {
				if fan.Curve != curveName {
					continue
				}

				if err := fan.SetSpeed(percent); err != nil {
					log.Printf("Error setting speed for fan '%s': %v", fanName, err)
					continue
				}
				setFor += "'" + fanName + "', "
			}

			if setFor != "" && percent != curve.PreviousPercent {
				setFor = setFor[:len(setFor)-2]
				log.Printf("Set %s to %.0f%%", setFor, percent)
			}
			curve.PreviousPercent = percent
			curve.PreviousTemp = temp
		}
		time.Sleep(1 * time.Second)
	}
}

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}
