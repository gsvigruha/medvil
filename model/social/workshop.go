package social

import (
	"medvil/model/artifacts"
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
	home := m.GetField(w.Household.Building.X, w.Household.Building.Y)
	if w.Manufacture != nil {
		mp := w.Household.Town.Marketplace
		market := m.GetField(mp.Building.X, mp.Building.Y)
		needs := w.Household.Resources.Needs(w.Manufacture.Inputs)
		if needs != nil && len(w.Household.Tasks) < 3 {
			if mp.Storage.Has(needs) && w.Household.Money >= mp.Price(needs) {
				w.Household.AddTask(&economy.ExchangeTask{
					HomeF:          home,
					MarketF:        market,
					Exchange:       mp,
					HouseholdR:     &w.Household.Resources,
					HouseholdMoney: &w.Household.Money,
					GoodsToBuy:     needs,
					GoodsToSell:    nil,
				})
			}
		}
		if w.Household.Resources.RemoveAll(w.Manufacture.Inputs) {
			w.Household.AddTask(&economy.ManufactureTask{
				M: w.Manufacture,
				F: home,
				R: &w.Household.Resources,
			})
		}

		for a, q := range w.Household.Resources.Artifacts {
			if !w.Manufacture.IsInput(a) {
				qToSell := w.Household.ArtifactToSell(a, q)
				if qToSell > 0 {
					goods := []artifacts.Artifacts{artifacts.Artifacts{A: a, Quantity: qToSell}}
					if w.Household.Town.Marketplace.CanSell(goods) && w.Household.Resources.RemoveAll(goods) {
						w.Household.AddTask(&economy.ExchangeTask{
							HomeF:          home,
							MarketF:        market,
							Exchange:       w.Household.Town.Marketplace,
							HouseholdR:     &w.Household.Resources,
							HouseholdMoney: &w.Household.Money,
							GoodsToBuy:     nil,
							GoodsToSell:    goods,
						})
					}
				}
			}
		}
	}
}
