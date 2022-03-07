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
	home := m.GetField(t.Household.Building.X, t.Household.Building.Y)
	for _, a := range building.ConstructionInputs {
		tag := "construction_input#" + a.Name
		goods := []artifacts.Artifacts{artifacts.Artifacts{A: a, Quantity: ProductTransportQuantity}}
		if t.Household.NumTasks("exchange", tag) == 0 && mp.Storage.Has(goods) && t.Household.Money >= mp.Price(goods) {
			mx, my := mp.Building.GetRandomBuildingXY()
			t.Household.AddTask(&economy.ExchangeTask{
				HomeF:          home,
				MarketF:        m.GetField(mx, my),
				Exchange:       mp,
				HouseholdR:     &t.Household.Resources,
				HouseholdMoney: &t.Household.Money,
				GoodsToBuy:     goods,
				GoodsToSell:    nil,
				TaskTag:        tag,
			})
		}
	}
}
