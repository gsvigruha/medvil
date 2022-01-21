package social

import (
	"encoding/json"
	"medvil/model/artifacts"
	"medvil/model/economy"
	"medvil/model/navigation"
	"medvil/model/terrain"
	"medvil/model/time"
	//"fmt"
)

type FarmLand struct {
	X       uint16
	Y       uint16
	UseType uint8
	F       *navigation.Field
}

func (f FarmLand) ToLocation() navigation.Location {
	return navigation.Location{X: f.X, Y: f.Y, F: f.F}
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
	home := navigation.Location{X: f.Household.Building.X, Y: f.Household.Building.Y, F: m.GetField(f.Household.Building.X, f.Household.Building.Y)}
	var a *artifacts.Artifact = nil
	if l.UseType == economy.FarmFieldUseTypeVegetables {
		a = artifacts.GetArtifact("vegetable")
	} else if l.UseType == economy.FarmFieldUseTypeWheat {
		a = artifacts.GetArtifact("grain")
	}
	if a != nil {
		f.Household.AddTask(&economy.TransportTask{
			PickupL:  l.ToLocation(),
			DropoffL: home,
			PickupR:  &l.F.Terrain.Resources,
			DropoffR: &f.Household.Resources,
			A:        a,
			Quantity: 10,
		})
	}
}

func (f *Farm) ElapseTime(Calendar *time.CalendarType, m navigation.IMap) {
	f.Household.ElapseTime(Calendar)
	if economy.ArgicultureCycleStartTime.Matches(Calendar) {
		for i := range f.Land {
			l := f.Land[i]
			location := navigation.Location{X: l.X, Y: l.Y, F: l.F}
			if (l.UseType == economy.FarmFieldUseTypeWheat || l.UseType == economy.FarmFieldUseTypeVegetables) && l.F.Plant == nil {
				if l.F.Terrain.T == terrain.Dirt {
					f.Household.AddTask(&economy.AgriculturalTask{T: economy.AgriculturalTaskSowing, L: location, UseType: l.UseType})
					f.Household.AddTask(&economy.AgriculturalTask{T: economy.AgriculturalTaskHarvesting, L: location, UseType: l.UseType})
					f.AddTransportTask(l, m)
				} else if l.F.Terrain.T == terrain.Grass {
					f.Household.AddTask(&economy.AgriculturalTask{T: economy.AgriculturalTaskPloughing, L: location, UseType: l.UseType})
					f.Household.AddTask(&economy.AgriculturalTask{T: economy.AgriculturalTaskSowing, L: location, UseType: l.UseType})
					f.Household.AddTask(&economy.AgriculturalTask{T: economy.AgriculturalTaskHarvesting, L: location, UseType: l.UseType})
					f.AddTransportTask(l, m)
				}
			} else if l.UseType == economy.FarmFieldUseTypeOrchard && l.F.Plant != nil && l.F.Plant.T.TreeT == &terrain.Apple {
				f.Household.AddTask(&economy.AgriculturalTask{T: economy.AgriculturalTaskHarvesting, L: location, UseType: l.UseType})
				f.AddTransportTask(l, m)
			}
		}
	}
}
