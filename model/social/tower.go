package social

import (
	"medvil/model/navigation"
)

type PatrolLand struct {
	X uint16
	Y uint16
	F *navigation.Field
}

type Tower struct {
	Household Household
	Land      []PatrolLand
}
