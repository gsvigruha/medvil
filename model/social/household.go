package social

import (
	"medvil/model/building"
)

type HouseHold struct {
	People   []*Person
	Money    uint32
	Building *building.Building
	Town     *Town
}
