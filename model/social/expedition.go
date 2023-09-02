package social

import (
	"medvil/model/artifacts"
	"medvil/model/economy"
	"medvil/model/vehicles"
)

type Expedition struct {
	Money     uint32
	People    []*Person
	Vehicle   *vehicles.Vehicle
	Resources *artifacts.Resources
	Tasks     []economy.Task
}
