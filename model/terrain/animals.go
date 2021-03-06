package terrain

import (
	"medvil/model/artifacts"
	"medvil/model/time"
)

type AnimalType struct {
	Name             string
	MaturityAgeYears uint8
	EndOfYearYield   artifacts.Artifacts
	EndOfLifeYield   artifacts.Artifacts
}

type Animal struct {
	T             *AnimalType
	BirthDateDays uint32
	Fed           bool
	Corralled     bool
}

func (a *Animal) ElapseTime(Calendar *time.CalendarType) {
	a.Fed = !a.Corralled && Calendar.Month >= 9
}

func (a *Animal) AgeYears(Calendar *time.CalendarType) uint32 {
	return (Calendar.DaysElapsed() - a.BirthDateDays) / (30 * 12)
}

func (a *Animal) IsMature(Calendar *time.CalendarType) bool {
	return a.AgeYears(Calendar) >= uint32(a.T.MaturityAgeYears)
}

var Sheep = &AnimalType{
	Name:             "sheep",
	MaturityAgeYears: 3,
	EndOfYearYield:   artifacts.Artifacts{A: artifacts.GetArtifact("wool"), Quantity: 1},
	EndOfLifeYield:   artifacts.Artifacts{A: artifacts.GetArtifact("sheep"), Quantity: 1},
}
