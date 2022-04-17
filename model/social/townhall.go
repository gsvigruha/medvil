package social

import (
	"medvil/model/artifacts"
	"medvil/model/building"
	"medvil/model/economy"
	"medvil/model/navigation"
	"medvil/model/time"
)

type Townhall struct {
	Household Household
}

const ConstructionBudgetPercentage = 0.3
const ConstructionStorageCapacity = 0.7

func (t *Townhall) ElapseTime(Calendar *time.CalendarType, m navigation.IMap) {
	t.Household.ElapseTime(Calendar, m)
	mp := t.Household.Town.Marketplace
	if t.Household.Resources.UsedVolumeCapacity() < ConstructionStorageCapacity {
		maxPrice := uint32(float64(t.Household.Money) * ConstructionBudgetPercentage / float64(len(building.ConstructionInputs)))
		maxVolumePerArtifact := t.Household.Resources.VolumeCapacity / uint16(len(building.ConstructionInputs))
		for _, a := range building.ConstructionInputs {
			if maxVolumePerArtifact > t.Household.Resources.Get(a)*a.V {
				tag := "construction_input#" + a.Name
				goods := []artifacts.Artifacts{artifacts.Artifacts{A: a, Quantity: ProductTransportQuantity(a)}}
				if t.Household.NumTasks("exchange", tag) == 0 && t.Household.Money >= mp.Price(goods) {
					t.Household.AddTask(&economy.BuyTask{
						Exchange:       mp,
						HouseholdMoney: &t.Household.Money,
						Goods:          goods,
						MaxPrice:       maxPrice,
						TaskTag:        tag,
					})
				}
			}
		}
	}
}
