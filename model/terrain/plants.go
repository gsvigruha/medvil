package terrain

import (
	"math"
	"medvil/controller"
)

type TreeType struct {
	BranchWidth0        float64
	BranchLength0       float64
	BranchWidthD        []float64
	BranchLengthD       []float64
	BranchAngles        []float64
	LeavesMinIterarion  uint8
	LeavesSize          float64
	BranchingIterations uint8
	Blooms              bool
}

type PlantType struct {
	Name             string
	MaturityAgeYears uint8
	TreeT            *TreeType
}

var Oak = TreeType{
	BranchWidth0:        7.0,
	BranchLength0:       30.0,
	BranchWidthD:        []float64{0.7, 0.6},
	BranchLengthD:       []float64{0.9, 0.7},
	BranchAngles:        []float64{math.Pi / 8, math.Pi / 2.4},
	LeavesMinIterarion:  1,
	LeavesSize:          10.0,
	BranchingIterations: 6,
	Blooms:              false,
}

var Apple = TreeType{
	BranchWidth0:        5.0,
	BranchLength0:       25.0,
	BranchWidthD:        []float64{0.7, 0.7},
	BranchLengthD:       []float64{0.75, 0.7},
	BranchAngles:        []float64{math.Pi / 2.4, math.Pi / 2.4},
	LeavesMinIterarion:  1,
	LeavesSize:          8.0,
	BranchingIterations: 6,
	Blooms:              true,
}

var AllPlantTypes = [...]PlantType{
	PlantType{Name: "oak tree", MaturityAgeYears: 10, TreeT: &Oak},
	PlantType{Name: "apple tree", MaturityAgeYears: 10, TreeT: &Apple},
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
