package curve

type Curve struct {
	IsStep          bool
	Sensor          string // sensor name
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

func (c Curve) GetPoint(temp float64) (float64, bool) {
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

func interpolate(p1, p2 CurvePoint, temp float64) float64 {
	if p1.Temp == p2.Temp {
		return p1.Percent
	}
	return p1.Percent + (p2.Percent-p1.Percent)*(temp-p1.Temp)/(p2.Temp-p1.Temp)
}
