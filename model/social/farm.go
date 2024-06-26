package social

import (
	"encoding/json"
	"medvil/model/artifacts"
	"medvil/model/economy"
	"medvil/model/navigation"
	"medvil/model/terrain"
	"medvil/model/time"
	"medvil/util"
)

const FarmMaxDistance = 6
const FarmMaxDistanceClearing = 12
const PastureMaxSlope = 2

var fruit = artifacts.GetArtifact("fruit")
var vegetable = artifacts.GetArtifact("vegetable")
var grain = artifacts.GetArtifact("grain")
var herb = artifacts.GetArtifact("herb")

type FarmLand struct {
	X       uint16
	Y       uint16
	UseType uint8
	F       *navigation.Field
}

func (l FarmLand) Field() *navigation.Field {
	return l.F
}

func (l FarmLand) Context() string {
	switch l.UseType {
	case economy.FarmFieldUseTypeWheat:
		return "grain"
	case economy.FarmFieldUseTypeVegetables:
		return "vegetable"
	case economy.FarmFieldUseTypeOrchard:
		return "fruit"
	case economy.FarmFieldUseTypeForestry:
		return "log"
	case economy.FarmFieldUseTypeReed:
		return "reed"
	case economy.FarmFieldUseTypePasture:
		return "sheep"
	case economy.FarmFieldUseTypeHerb:
		return "herb"
	}
	return ""
}

type Farm struct {
	Household  *Household
	Land       []FarmLand
	AutoSwitch bool
}

func (f *Farm) UnmarshalJSON(data []byte) error {
	var j map[string]json.RawMessage
	if err := json.Unmarshal(data, &j); err != nil {
		return err
	}
	if err := json.Unmarshal(j["household"], &f.Household); err != nil {
		return err
	}
	var l [][]uint16
	if err := json.Unmarshal(j["land"], &l); err != nil {
		return err
	}
	f.Land = make([]FarmLand, len(l))
	for i := range l {
		f.Land[i].X = l[i][0]
		f.Land[i].Y = l[i][1]
		f.Land[i].UseType = uint8(l[i][2])
	}
	return nil
}

func (f *Farm) AddTransportTask(l FarmLand, m navigation.IMap) {
	home := m.GetField(f.Household.Building.X, f.Household.Building.Y)
	for a, q := range l.F.Terrain.Resources.Artifacts {
		if l.F.Terrain.Resources.IsRealArtifact(a) && q > 0 {
			tag := economy.TransportTaskTag(l.F, a)
			if f.Household.NumTasks("transport", tag) == 0 {
				f.Household.AddTask(&economy.TransportTask{
					PickupD:        l.F,
					DropoffD:       home,
					PickupR:        l.F.Terrain.Resources,
					DropoffR:       f.Household.Resources,
					A:              a,
					TargetQuantity: q,
				})
			}
		}
	}
}

func (f *Farm) ElapseTime(Calendar *time.CalendarType, m navigation.IMap) {
	f.Household.ElapseTime(Calendar, m)

	if economy.ArgicultureCycleStartTime.Matches(Calendar) {
		for i := range f.Land {
			l := f.Land[i]
			if (l.UseType == economy.FarmFieldUseTypeWheat || l.UseType == economy.FarmFieldUseTypeVegetables || l.UseType == economy.FarmFieldUseTypeHerb) && l.F.Plant == nil {
				if l.F.Terrain.T == terrain.Dirt {
					f.Household.AddTask(&economy.AgriculturalTask{T: economy.AgriculturalTaskSowing, F: l.F, UseType: l.UseType, Start: *Calendar})
					f.Household.AddTask(&economy.AgriculturalTask{T: economy.AgriculturalTaskHarvesting, F: l.F, UseType: l.UseType, Start: *Calendar})
				} else if l.F.Terrain.T == terrain.Grass {
					f.Household.AddTask(&economy.AgriculturalTask{T: economy.AgriculturalTaskPloughing, F: l.F, UseType: l.UseType, Start: *Calendar})
					f.Household.AddTask(&economy.AgriculturalTask{T: economy.AgriculturalTaskSowing, F: l.F, UseType: l.UseType, Start: *Calendar})
					f.Household.AddTask(&economy.AgriculturalTask{T: economy.AgriculturalTaskHarvesting, F: l.F, UseType: l.UseType, Start: *Calendar})
				}
			} else if l.UseType == economy.FarmFieldUseTypeOrchard {
				if l.F.Plant == nil {
					f.Household.AddTask(&economy.AgriculturalTask{T: economy.AgriculturalTaskPlantingAppleTree, F: l.F, UseType: l.UseType, Start: *Calendar})
				} else if l.F.Plant.T.TreeT == &terrain.Apple {
					f.Household.AddTask(&economy.AgriculturalTask{T: economy.AgriculturalTaskHarvesting, F: l.F, UseType: l.UseType, Start: *Calendar})
				}
			} else if l.UseType == economy.FarmFieldUseTypeForestry {
				if l.F.Plant == nil {
					f.Household.AddTask(&economy.AgriculturalTask{T: economy.AgriculturalTaskPlantingOakTree, F: l.F, UseType: l.UseType, Start: *Calendar})
				}
			} else if l.UseType == economy.FarmFieldUseTypeReed {
				if l.F.Plant == nil {
					f.Household.AddTask(&economy.AgriculturalTask{T: economy.AgriculturalTaskPlantingReed, F: l.F, UseType: l.UseType, Start: *Calendar})
				} else if l.F.Plant != nil && l.F.Plant.T.Name == "reed" {
					f.Household.AddTask(&economy.AgriculturalTask{T: economy.AgriculturalTaskReedCutting, F: l.F, UseType: l.UseType, Start: *Calendar})
				}
			} else if l.UseType == economy.FarmFieldUseTypePasture && l.F.Terrain.T == terrain.Grass {
				f.Household.AddTask(&economy.AgriculturalTask{T: economy.AgriculturalTaskGrazing, F: l.F, UseType: l.UseType, Start: *Calendar})
				f.Household.AddTask(&economy.AgriculturalTask{T: economy.AgriculturalTaskCorralling, F: l.F, UseType: l.UseType, Start: *Calendar})
			}
			if l.F.Plant != nil && l.F.Plant.IsTree() && l.F.Plant.IsMature(Calendar) && l.UseType != economy.FarmFieldUseTypeOrchard {
				f.Household.AddTask(&economy.AgriculturalTask{T: economy.AgriculturalTaskTreeCutting, F: l.F, UseType: l.UseType, Start: *Calendar})
			} else if l.F.Plant != nil && l.F.Plant.IsTree() && l.UseType != economy.FarmFieldUseTypeForestry && l.UseType != economy.FarmFieldUseTypeOrchard {
				f.Household.AddTask(&economy.AgriculturalTask{T: economy.AgriculturalTaskTreeCutting, F: l.F, UseType: l.UseType, Start: *Calendar})
			}
		}
	}

	if Calendar.Hour == 0 {
		for _, land := range f.Land {
			f.AddTransportTask(land, m)
		}
		f.ReleaseClearedLand()
	}

	mp := f.Household.Town.Marketplace
	if mp != nil {
		f.Household.MaybeBuyPaper(f.AutoSwitch)
		f.Household.MaybeBuyBoat(Calendar, m)
		f.Household.MaybeBuyCart(Calendar, m)

		f.Household.SellArtifacts(NotInputOrProduct, f.IsOutput)

		if f.AutoSwitch && Calendar.Day == 30 && Calendar.Hour == 0 && Calendar.Month == 12 && f.Household.Resources.Remove(Paper, 1) > 0 {
			profits := []float64{
				float64(mp.Prices[vegetable]) * float64(terrain.Vegetables.Yield.Quantity),
				float64(mp.Prices[grain]) * float64(terrain.Grain.Yield.Quantity),
				float64(mp.Prices[herb]) * float64(terrain.Herb.Yield.Quantity),
			}
			for i, land := range f.Land {
				if land.UseType == economy.FarmFieldUseTypeVegetables || land.UseType == economy.FarmFieldUseTypeWheat || land.UseType == economy.FarmFieldUseTypeHerb {
					switch util.RandomIndexWeighted(profits) {
					case 0:
						f.Land[i].UseType = economy.FarmFieldUseTypeVegetables
					case 1:
						f.Land[i].UseType = economy.FarmFieldUseTypeWheat
					case 2:
						f.Land[i].UseType = economy.FarmFieldUseTypeHerb
					}
				}
			}
		}
	}
}

func (f *Farm) IsOutput(a *artifacts.Artifact) bool {
	return a == fruit || a == vegetable
}

func (f *Farm) GetFields() []navigation.FieldWithContext {
	fields := make([]navigation.FieldWithContext, len(f.Land))
	for i := range f.Land {
		fields[i] = f.Land[i]
	}
	return fields
}

func (f *Farm) FieldUsableFor(m navigation.IMap, field *navigation.Field, useType uint8) bool {
	if useType == economy.FarmFieldUseTypeReed {
		return m.Shore(field.X, field.Y)
	}
	if useType == economy.FarmFieldUseTypeOrchard || useType == economy.FarmFieldUseTypeForestry || useType == economy.FarmFieldUseTypeBarren {
		return field.Plantable(false)
	}
	if useType == economy.FarmFieldUseTypePasture {
		return field.Plantable(false) && field.MaxSlope() <= PastureMaxSlope
	}
	return field.Arable()
}

func (f *Farm) FieldWithinDistance(field *navigation.Field) bool {
	return WithinDistance(f.Household.Building, field, FarmMaxDistance)
}

func (f *Farm) FieldWithinDistanceClearing(field *navigation.Field) bool {
	return WithinDistance(f.Household.Building, field, FarmMaxDistanceClearing)
}

func (f *Farm) GetHome() Home {
	return f.Household
}

func (f *Farm) GetLandDistribution() map[uint8]int {
	result := make(map[uint8]int)
	for _, land := range f.Land {
		if cnt, ok := result[land.UseType]; ok {
			result[land.UseType] = cnt + 1
		} else {
			result[land.UseType] = 1
		}
	}
	return result
}

func (f *Farm) ReleaseClearedLand() {
	var newLand []FarmLand = make([]FarmLand, 0, len(f.Land))
	for _, land := range f.Land {
		if land.UseType != economy.FarmFieldUseTypeBarren || land.F.Plant != nil || !land.F.Terrain.Resources.IsEmpty() {
			newLand = append(newLand, land)
		} else {
			land.F.Allocated = false
		}
	}
	f.Land = newLand
}
