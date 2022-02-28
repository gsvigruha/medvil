package building

import (
	"medvil/model/artifacts"
)

type BuildingConstruction struct {
	Building      *Building
	Progress      uint16
	RemainingCost []artifacts.Artifacts
	T             BuildingType
}
