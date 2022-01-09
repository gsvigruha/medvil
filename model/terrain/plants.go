package terrain

import (
	"math"
	"medvil/controller"
	"medvil/model/artifacts"
)

type PlantType struct {
	Name             string
	MaturityAgeYears uint8
	TreeT            *TreeType
	Yield            artifacts.Artifacts
}

type Plant struct {
	T             *PlantType
	X             uint16
	Y             uint16
	BirthDateDays uint32
	Shape         uint8
}

func (p *Plant) IsTree() bool {
	return p.T.TreeT != nil
}

func (p *Plant) Maturity(Calendar *controller.CalendarType) float64 {
	return math.Min(float64(p.AgeYears(Calendar)), float64(p.T.MaturityAgeYears)) / float64(p.T.MaturityAgeYears)
}

func (p *Plant) AgeYears(Calendar *controller.CalendarType) uint32 {
	return (Calendar.DaysElapsed() - p.BirthDateDays) / (30 * 12)
}
