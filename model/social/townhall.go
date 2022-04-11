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

const PercentageSpentOnConstruction = 0.3

func (t *Townhall) ElapseTime(Calendar *time.CalendarType, m navigation.IMap) {
	t.Household.ElapseTime(Calendar, m)
	mp := t.Household.Town.Marketplace
	if t.Household.Resources.UsedVolumeCapacity() < 0.8 {
		maxPrice := uint32(float64(t.Household.Money) * PercentageSpentOnConstruction / float64(len(building.ConstructionInputs)))
		for _, a := range building.ConstructionInputs {
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
