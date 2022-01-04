package model

type FieldObject interface {
	Walkable() bool
	LiftN() int8
	LiftE() int8
	LiftS() int8
	LiftW() int8
}
