package economy

import (
	"medvil/model/navigation"
	"medvil/model/terrain"
	"medvil/model/time"
	"math/rand"
	//"fmt"
)

var ArgicultureCycleStartTime = time.TimeOfYear{Month: 3, Day: 1}

const AgriculturalTaskPloughing = 1
const AgriculturalTaskSowing = 2
const AgriculturalTaskHarvesting = 3

const AgriculturalTaskDurationPloughing = 24 * 10
const AgriculturalTaskDurationSowing = 24 * 5
const AgriculturalTaskDurationHarvesting = 24 * 10


const FarmFieldUseTypeBarren uint8 = 0
const FarmFieldUseTypeWheat uint8 = 1
const FarmFieldUseTypeOrchard uint8 = 2
const FarmFieldUseTypePasture uint8 = 3
const FarmFieldUseTypeVegetables uint8 = 4


type AgriculturalTask struct {
	T        uint8
	L        navigation.Location
	Progress uint16
	UseType uint8
}

func (t *AgriculturalTask) Location() navigation.Location {
	return t.L
}

func (t *AgriculturalTask) Complete(Calendar *time.CalendarType) bool {
	t.Progress++
	switch t.T {
	case AgriculturalTaskPloughing:
		if t.Progress >= AgriculturalTaskDurationPloughing {
			t.L.F.Terrain.T = terrain.Dirt
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
			t.L.F.Plant = &terrain.Plant{
				T:             cropType,
				X:             t.L.X,
				Y:             t.L.Y,
				BirthDateDays: Calendar.DaysElapsed(),
				Shape:         uint8(rand.Intn(10)),
			}
			return true
		}
	case AgriculturalTaskHarvesting:
		return t.Progress >= AgriculturalTaskDurationHarvesting
	}
	return false
}

func (t *AgriculturalTask) Blocked() bool {
	switch t.T {
	case AgriculturalTaskPloughing:
		return t.L.F.Plant != nil
	case AgriculturalTaskSowing:
		return t.L.F.Plant != nil || t.L.F.Terrain.T != terrain.Dirt
	case AgriculturalTaskHarvesting:
		return t.L.F.Plant == nil || !t.L.F.Plant.Ripe
	}
	return false
}
