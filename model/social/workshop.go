package social

import (
	"fmt"
	"medvil/model/economy"
	"medvil/model/navigation"
	"medvil/model/time"
)

type Workshop struct {
	Household   Household
	Manufacture *economy.Manufacture
}

func (w *Workshop) ElapseTime(Calendar *time.CalendarType, m navigation.IMap) {
	w.Household.ElapseTime(Calendar, m)
	if w.Manufacture != nil {
		home := m.GetField(w.Household.Building.X, w.Household.Building.Y)
		mp := w.Household.Town.Marketplace
		market := m.GetField(mp.Building.X, mp.Building.Y)
		needs := w.Household.Resources.Needs(w.Manufacture.Inputs)
		fmt.Println(needs)
		if needs != nil && mp.Storage.Has(needs) && w.Household.Money >= mp.Price(needs) {
			w.Household.AddTask(&economy.ExchangeTask{
				PickupF:        home,
				DropoffF:       market,
				Exchange:       mp,
				HouseholdR:     &w.Household.Resources,
				HouseholdMoney: &w.Household.Money,
				GoodsToBuy:     needs,
				GoodsToSell:    nil,
			})
		}
	}
}
