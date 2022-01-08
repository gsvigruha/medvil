package terrain

import (
	"math"
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
}

type PlantType struct {
	Name        string
	MaturityAge uint8
	TreeT       *TreeType
}

var Oak = TreeType{
	BranchWidth0:        7.0,
	BranchLength0:       25.0,
	BranchWidthD:        []float64{0.7, 0.6},
	BranchLengthD:       []float64{0.9, 0.7},
	BranchAngles:        []float64{math.Pi / 8, math.Pi / 2.4},
	LeavesMinIterarion:  1,
	LeavesSize:          10.0,
	BranchingIterations: 6,
}

var AllPlantTypes = [...]PlantType{
	PlantType{Name: "oak tree", MaturityAge: 10, TreeT: &Oak},
}

type Plant struct {
	T     *PlantType
	X     uint16
	Y     uint16
	Age   uint8
	Shape uint8
}

func (p *Plant) IsTree() bool {
	return p.T.TreeT != nil
}

func (p *Plant) Maturity() float64 {
	return math.Min(float64(p.Age), float64(p.T.MaturityAge)) / float64(p.T.MaturityAge)
}
