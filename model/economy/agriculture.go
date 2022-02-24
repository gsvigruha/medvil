package economy

import (
	"math/rand"
	"medvil/model/artifacts"
	"medvil/model/navigation"
	"medvil/model/terrain"
	"medvil/model/time"
)

var ArgicultureCycleStartTime = time.TimeOfYear{Month: 3, Day: 1}

const AgriculturalTaskPloughing = 1
const AgriculturalTaskSowing = 2
const AgriculturalTaskHarvesting = 3
const AgriculturalTaskPlantingAppleTree = 4
const AgriculturalTaskPlantingOakTree = 5
const AgriculturalTaskTreeCutting = 6

const AgriculturalTaskDurationPloughing = 24 * 30
const AgriculturalTaskDurationSowing = 24 * 15
const AgriculturalTaskDurationHarvesting = 24 * 30
const AgriculturalTaskDurationPlanting = 24 * 5
const AgriculturalTaskDurationTreeCutting = 24 * 10

const FarmFieldUseTypeBarren uint8 = 0
const FarmFieldUseTypeWheat uint8 = 1
const FarmFieldUseTypeOrchard uint8 = 2
const FarmFieldUseTypePasture uint8 = 3
const FarmFieldUseTypeVegetables uint8 = 4
const FarmFieldUseTypeForestry uint8 = 5

type AgriculturalTask struct {
	T        uint8
	F        *navigation.Field
	Progress uint16
	UseType  uint8
	Start    time.CalendarType
}

func (t *AgriculturalTask) Field() *navigation.Field {
	return t.F
}

func (t *AgriculturalTask) Complete(Calendar *time.CalendarType) bool {
	t.Progress++
	switch t.T {
	case AgriculturalTaskPloughing:
		if t.Progress >= AgriculturalTaskDurationPloughing {
			t.F.Terrain.T = terrain.Dirt
			return true
		}
	case AgriculturalTaskSowing:
		if t.Progress >= AgriculturalTaskDurationSowing {
			var cropType *terrain.PlantType
			if t.UseType == FarmFieldUseTypeWheat {
				cropType = &terrain.AllCropTypes[0]
			} else if t.UseType == FarmFieldUseTypeVegetables {
				cropType = &terrain.AllCropTypes[1]
			}
			t.F.Plant = &terrain.Plant{
				T:             cropType,
				X:             t.F.X,
				Y:             t.F.Y,
				BirthDateDays: Calendar.DaysElapsed(),
				Shape:         uint8(rand.Intn(10)),
			}
			return true
		}
	case AgriculturalTaskHarvesting:
		if t.Progress >= AgriculturalTaskDurationHarvesting {
			if t.F.Plant != nil {
				t.F.Terrain.Resources.Add(t.F.Plant.T.Yield.A, t.F.Plant.T.Yield.Quantity)
				if t.F.Plant.T.IsAnnual() {
					t.F.Plant = nil
				}
			}
			return true
		}
	case AgriculturalTaskPlantingAppleTree:
		if t.Progress >= AgriculturalTaskDurationPlanting {
			t.F.Plant = &terrain.Plant{
				T:             &terrain.AllTreeTypes[1],
				X:             t.F.X,
				Y:             t.F.Y,
				BirthDateDays: Calendar.DaysElapsed(),
				Shape:         uint8(rand.Intn(10)),
			}
			return true
		}
	case AgriculturalTaskPlantingOakTree:
		if t.Progress >= AgriculturalTaskDurationPlanting {
			t.F.Plant = &terrain.Plant{
				T:             &terrain.AllTreeTypes[0],
				X:             t.F.X,
				Y:             t.F.Y,
				BirthDateDays: Calendar.DaysElapsed(),
				Shape:         uint8(rand.Intn(10)),
			}
			return true
		}
	case AgriculturalTaskTreeCutting:
		if t.Progress >= AgriculturalTaskDurationTreeCutting {
			t.F.Terrain.Resources.Add(artifacts.GetArtifact("log"), t.F.Plant.T.TreeT.LogYield)
			t.F.Plant = nil
			return true
		}
	}
	return false
}

func (t *AgriculturalTask) Blocked() bool {
	switch t.T {
	case AgriculturalTaskPloughing:
		return t.F.Plant != nil
	case AgriculturalTaskSowing:
		return t.F.Plant != nil || t.F.Terrain.T != terrain.Dirt
	case AgriculturalTaskHarvesting:
		return t.F.Plant == nil || !t.F.Plant.Ripe
	case AgriculturalTaskPlantingAppleTree:
		return t.F.Plant != nil
	case AgriculturalTaskPlantingOakTree:
		return t.F.Plant != nil
	case AgriculturalTaskTreeCutting:
		return t.F.Plant == nil || !t.F.Plant.IsTree()
	}
	return false
}

func (t *AgriculturalTask) Name() string {
	switch t.T {
	case AgriculturalTaskPloughing:
		return "ploughing"
	case AgriculturalTaskSowing:
		return "sowing"
	case AgriculturalTaskHarvesting:
		return "harvesting"
	case AgriculturalTaskPlantingAppleTree:
		return "planting"
	case AgriculturalTaskPlantingOakTree:
		return "planting"
	case AgriculturalTaskTreeCutting:
		return "treecutting"
	}
	return ""
}

func (t *AgriculturalTask) Tag() string {
	return ""
}

func (t *AgriculturalTask) Expired(Calendar *time.CalendarType) bool {
	return Calendar.DaysElapsed()-t.Start.DaysElapsed() >= 365
}
