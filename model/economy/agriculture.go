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
const AgriculturalTaskReedCutting = 7
const AgriculturalTaskGrazing = 8
const AgriculturalTaskCorralling = 9

const AgriculturalTaskDurationPloughing = 24 * 30
const AgriculturalTaskDurationSowing = 24 * 15
const AgriculturalTaskDurationHarvesting = 24 * 30
const AgriculturalTaskDurationPlanting = 24 * 5
const AgriculturalTaskDurationTreeCutting = 24 * 10
const AgriculturalTaskDurationReedCutting = 24 * 10
const AgriculturalTaskDurationGrazing = 24 * 5
const AgriculturalTaskDurationCorralling = 24 * 10

const FarmFieldUseTypeBarren uint8 = 0
const FarmFieldUseTypeWheat uint8 = 1
const FarmFieldUseTypeOrchard uint8 = 2
const FarmFieldUseTypePasture uint8 = 3
const FarmFieldUseTypeVegetables uint8 = 4
const FarmFieldUseTypeForestry uint8 = 5
const FarmFieldUseTypeReed uint8 = 6

type AgriculturalTask struct {
	TaskBase
	T        uint8
	F        *navigation.Field
	Progress uint16
	UseType  uint8
	Start    time.CalendarType
}

func (t *AgriculturalTask) Field() *navigation.Field {
	return t.F
}

func (t *AgriculturalTask) Complete(Calendar *time.CalendarType, tool bool) bool {
	t.Progress++
	if tool {
		t.Progress++
	}
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
	case AgriculturalTaskReedCutting:
		if t.Progress >= AgriculturalTaskDurationReedCutting {
			t.F.Terrain.Resources.Add(t.F.Plant.T.Yield.A, t.F.Plant.T.Yield.Quantity)
			t.F.Plant = nil
			return true
		}
	case AgriculturalTaskGrazing:
		if t.Progress >= AgriculturalTaskDurationGrazing {
			if t.F.Animal == nil {
				t.F.Animal = &terrain.Animal{
					T:             terrain.Sheep,
					BirthDateDays: Calendar.DaysElapsed(),
					Fed:           false,
					Corralled:     false,
				}
			} else {
				t.F.Animal.Corralled = false
				t.F.Animal.Fed = false
			}
			return true
		}
	case AgriculturalTaskCorralling:
		if t.Progress >= AgriculturalTaskDurationCorralling {
			if t.F.Animal != nil {
				t.F.Terrain.Resources.Add(t.F.Animal.T.EndOfYearYield.A, t.F.Animal.T.EndOfYearYield.Quantity)
				if t.F.Animal.IsMature(Calendar) {
					t.F.Terrain.Resources.Add(t.F.Animal.T.EndOfLifeYield.A, t.F.Animal.T.EndOfLifeYield.Quantity)
					t.F.Animal = nil
				} else {
					t.F.Animal.Corralled = true
				}
			}
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
	case AgriculturalTaskReedCutting:
		return t.F.Plant == nil || t.F.Plant.T.Name != "reed"
	case AgriculturalTaskGrazing:
		return t.F.Plant != nil || t.F.Terrain.T != terrain.Grass
	case AgriculturalTaskCorralling:
		return t.F.Animal == nil || !t.F.Animal.Fed
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
	case AgriculturalTaskReedCutting:
		return "reedcutting"
	case AgriculturalTaskGrazing:
		return "grazing"
	case AgriculturalTaskCorralling:
		return "corralling"
	}
	return ""
}

func (t *AgriculturalTask) Tag() string {
	return ""
}

func (t *AgriculturalTask) Expired(Calendar *time.CalendarType) bool {
	return Calendar.DaysElapsed()-t.Start.DaysElapsed() >= 365
}

func (t *AgriculturalTask) Motion() uint8 {
	return navigation.MotionFieldWork
}
