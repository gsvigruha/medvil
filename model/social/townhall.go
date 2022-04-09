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

func (t *Townhall) ElapseTime(Calendar *time.CalendarType, m navigation.IMap) {
	t.Household.ElapseTime(Calendar, m)
	mp := t.Household.Town.Marketplace
	if t.Household.Resources.UsedVolumeCapacity() < 0.8 {
		for _, a := range building.ConstructionInputs {
			tag := "construction_input#" + a.Name
			goods := []artifacts.Artifacts{artifacts.Artifacts{A: a, Quantity: ProductTransportQuantity}}
			if t.Household.NumTasks("exchange", tag) == 0 && mp.Storage.Has(goods) && t.Household.Money >= mp.Price(goods) {
				t.Household.AddTask(&economy.BuyTask{
					Exchange:       mp,
					HouseholdMoney: &t.Household.Money,
					Goods:          goods,
					MaxPrice:       mp.Price(goods),
					TaskTag:        tag,
				})
			}
		}
	}
}
