package model

type Map struct {
	SX        uint16
	SY        uint16
	Fields    [][]Field
	Buildings []Building
}
