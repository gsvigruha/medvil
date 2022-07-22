package social

import (
	"medvil/model/navigation"
)

type PatrolLand struct {
	X uint16
	Y uint16
	F *navigation.Field
}

func (l PatrolLand) Field() *navigation.Field {
	return l.F
}

func (l PatrolLand) Context() string {
	return "shield"
}

type Tower struct {
	Household Household
	Land      []PatrolLand
}

func (t *Tower) GetFields() []navigation.FieldWithContext {
	fields := make([]navigation.FieldWithContext, len(t.Land))
	for i := range t.Land {
		fields[i] = t.Land[i]
	}
	return fields
}
