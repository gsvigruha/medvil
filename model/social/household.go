package social

import (
	"medvil/model/building"
)

type Household struct {
	People   []*Person
	Money    uint32
	Building *building.Building
	Town     *Town
}
