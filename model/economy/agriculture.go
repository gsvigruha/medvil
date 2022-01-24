package economy

import (
	"math/rand"
	"medvil/model/navigation"
	"medvil/model/terrain"
	"medvil/model/time"
)

var ArgicultureCycleStartTime = time.TimeOfYear{Month: 3, Day: 1}

const AgriculturalTaskPloughing = 1
const AgriculturalTaskSowing = 2
const AgriculturalTaskHarvesting = 3

const AgriculturalTaskDurationPloughing = 24 * 30
const AgriculturalTaskDurationSowing = 24 * 15
const AgriculturalTaskDurationHarvesting = 24 * 30

const FarmFieldUseTypeBarren uint8 = 0
const FarmFieldUseTypeWheat uint8 = 1
const FarmFieldUseTypeOrchard uint8 = 2
const FarmFieldUseTypePasture uint8 = 3
const FarmFieldUseTypeVegetables uint8 = 4

type AgriculturalTask struct {
	T        uint8
	F        *navigation.Field
	Progress uint16
	UseType  uint8
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
			t.F.Terrain.Resources.Add(t.F.Plant.T.Yield.A, t.F.Plant.T.Yield.Quantity)
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
	}
	return ""
}
