package terrain

import (
	"math"
	"medvil/model/artifacts"
	"medvil/model/time"
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
	Ripe          bool
}

func (p *Plant) IsTree() bool {
	return p.T.TreeT != nil
}

func (p *Plant) Maturity(Calendar *time.CalendarType) float64 {
	return math.Min(float64(p.AgeYears(Calendar)), float64(p.T.MaturityAgeYears)) / float64(p.T.MaturityAgeYears)
}

func (p *Plant) IsMature(Calendar *time.CalendarType) bool {
	return p.Maturity(Calendar) == 1.0
}

func (p *Plant) AgeYears(Calendar *time.CalendarType) uint32 {
	return (Calendar.DaysElapsed() - p.BirthDateDays) / (30 * 12)
}

func (p *Plant) ElapseTime(Calendar *time.CalendarType) {
	if p.T.IsAnnual() {
		p.Ripe = (Calendar.DaysElapsed() - p.BirthDateDays) > 90
	} else if p.T.TreeT != nil {
		p.Ripe = Calendar.Month >= 7
	}
}

func (p *PlantType) IsAnnual() bool {
	return p.TreeT == nil && p.MaturityAgeYears <= 1.0
}
