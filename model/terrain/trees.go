package terrain

import (
	"math"
	"medvil/model/artifacts"
)

const TreeNumShapes = 8

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
	LogYield            uint16
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
	LogYield:            10,
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
	LogYield:            5,
}

var AllTreeTypes = [...]*PlantType{
	&PlantType{Name: "oak tree", MaturityAgeYears: 10, TreeT: &Oak, Habitat: Land},
	&PlantType{Name: "apple tree", MaturityAgeYears: 10, TreeT: &Apple, Yield: artifacts.Artifacts{A: artifacts.GetArtifact("fruit"), Quantity: 3}, Habitat: Land},
}
