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
		needs := w.Household.Resources.Needs(w.Manufacture.Inputs)
		if needs != nil && w.Household.NumTasks("exchange", "manufacture_input") == 0 {
			if mp.Storage.Has(needs) && w.Household.Money >= mp.Price(needs) {
				mx, my := mp.Building.GetRandomBuildingXY()
				w.Household.AddTask(&economy.ExchangeTask{
					HomeF:          home,
					MarketF:        m.GetField(mx, my),
					Exchange:       mp,
					HouseholdR:     &w.Household.Resources,
					HouseholdMoney: &w.Household.Money,
					GoodsToBuy:     needs,
					GoodsToSell:    nil,
					TaskTag:        "manufacture_input",
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
					if w.Household.Town.Marketplace.CanSell(goods) && w.Household.NumTasks("exchange", "sell_artifacts") == 0 {
						mx, my := mp.Building.GetRandomBuildingXY()
						w.Household.AddTask(&economy.ExchangeTask{
							HomeF:          home,
							MarketF:        m.GetField(mx, my),
							Exchange:       w.Household.Town.Marketplace,
							HouseholdR:     &w.Household.Resources,
							HouseholdMoney: &w.Household.Money,
							GoodsToBuy:     nil,
							GoodsToSell:    goods,
							TaskTag:        "sell_artifacts",
						})
					}
				}
			}
		}
	}
}
