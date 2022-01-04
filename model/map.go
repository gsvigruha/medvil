package model

import (
	"medvil/model/building"
)

type Map struct {
	SX        uint16
	SY        uint16
	Fields    [][]Field
	Buildings []building.Building
}
