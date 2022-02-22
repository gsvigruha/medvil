package social

import (
	"encoding/json"
	"medvil/model/artifacts"
	"medvil/model/economy"
	"medvil/model/navigation"
	"medvil/model/terrain"
	"medvil/model/time"
)

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
	}
	return ""
}

type Farm struct {
	Household Household
	Land      []FarmLand
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
	if l.F.Terrain.Resources.HasRealArtifacts() {
		for a, q := range l.F.Terrain.Resources.Artifacts {
			tag := economy.TransportTaskTag(l.F, a)
			if f.Household.NumTasks("transport", tag) == 0 {
				f.Household.AddTask(&economy.TransportTask{
					PickupF:  l.F,
					DropoffF: home,
					PickupR:  &l.F.Terrain.Resources,
					DropoffR: &f.Household.Resources,
					A:        a,
					Quantity: q,
				})
			}
		}
	}
}

func (f *Farm) ElapseTime(Calendar *time.CalendarType, m navigation.IMap) {
	f.Household.ElapseTime(Calendar, m)
	home := m.GetField(f.Household.Building.X, f.Household.Building.Y)
	if economy.ArgicultureCycleStartTime.Matches(Calendar) {
		for i := range f.Land {
			l := f.Land[i]
			if (l.UseType == economy.FarmFieldUseTypeWheat || l.UseType == economy.FarmFieldUseTypeVegetables) && l.F.Plant == nil {
				if l.F.Terrain.T == terrain.Dirt {
					f.Household.AddTask(&economy.AgriculturalTask{T: economy.AgriculturalTaskSowing, F: l.F, UseType: l.UseType})
					f.Household.AddTask(&economy.AgriculturalTask{T: economy.AgriculturalTaskHarvesting, F: l.F, UseType: l.UseType})
				} else if l.F.Terrain.T == terrain.Grass {
					f.Household.AddTask(&economy.AgriculturalTask{T: economy.AgriculturalTaskPloughing, F: l.F, UseType: l.UseType})
					f.Household.AddTask(&economy.AgriculturalTask{T: economy.AgriculturalTaskSowing, F: l.F, UseType: l.UseType})
					f.Household.AddTask(&economy.AgriculturalTask{T: economy.AgriculturalTaskHarvesting, F: l.F, UseType: l.UseType})
				}
			} else if l.UseType == economy.FarmFieldUseTypeOrchard {
				if l.F.Plant == nil {
					f.Household.AddTask(&economy.AgriculturalTask{T: economy.AgriculturalTaskPlantingAppleTree, F: l.F, UseType: l.UseType})
				} else if l.F.Plant.T.TreeT == &terrain.Apple && l.F.Plant.IsMature(Calendar) {
					f.Household.AddTask(&economy.AgriculturalTask{T: economy.AgriculturalTaskHarvesting, F: l.F, UseType: l.UseType})
				}
			} else if l.UseType == economy.FarmFieldUseTypeForestry {
				if l.F.Plant == nil {
					f.Household.AddTask(&economy.AgriculturalTask{T: economy.AgriculturalTaskPlantingOakTree, F: l.F, UseType: l.UseType})
				}
			}
			if l.F.Plant != nil && l.F.Plant.IsTree() && l.F.Plant.IsMature(Calendar) && l.UseType != economy.FarmFieldUseTypeOrchard {
				f.Household.AddTask(&economy.AgriculturalTask{T: economy.AgriculturalTaskTreeCutting, F: l.F, UseType: l.UseType})
			}
		}
	}
	if Calendar.Hour == 0 {
		for _, land := range f.Land {
			f.AddTransportTask(land, m)
		}
	}
	for a, q := range f.Household.Resources.Artifacts {
		qToSell := f.Household.ArtifactToSell(a, q)
		if qToSell > 0 {
			goods := []artifacts.Artifacts{artifacts.Artifacts{A: a, Quantity: qToSell}}
			if f.Household.Town.Marketplace.CanSell(goods) && f.Household.Resources.RemoveAll(goods) {
				mx, my := f.Household.Town.Marketplace.Building.GetRandomBuildingXY()
				f.Household.AddTask(&economy.ExchangeTask{
					HomeF:          home,
					MarketF:        m.GetField(mx, my),
					Exchange:       f.Household.Town.Marketplace,
					HouseholdR:     &f.Household.Resources,
					HouseholdMoney: &f.Household.Money,
					GoodsToBuy:     nil,
					GoodsToSell:    goods,
				})
			}
		}
	}
}

func (f *Farm) GetFields() []navigation.FieldWithContext {
	fields := make([]navigation.FieldWithContext, len(f.Land))
	for i := range f.Land {
		fields[i] = f.Land[i]
	}
	return fields
}
