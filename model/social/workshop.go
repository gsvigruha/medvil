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
				qToSell := w.Household.ArtifactToSell(a, q, w.Manufacture.IsOutput(a))
				if qToSell > ProductTransportQuantity {
					tag := "sell_artifacts#" + a.Name
					goods := []artifacts.Artifacts{artifacts.Artifacts{A: a, Quantity: ProductTransportQuantity}}
					if w.Household.Town.Marketplace.CanSell(goods) && int(qToSell)/ProductTransportQuantity > w.Household.NumTasks("exchange", tag) {
						mx, my := mp.Building.GetRandomBuildingXY()
						w.Household.AddTask(&economy.ExchangeTask{
							HomeF:          home,
							MarketF:        m.GetField(mx, my),
							Exchange:       w.Household.Town.Marketplace,
							HouseholdR:     &w.Household.Resources,
							HouseholdMoney: &w.Household.Money,
							GoodsToBuy:     nil,
							GoodsToSell:    goods,
							TaskTag:        tag,
						})
					}
				}
			}
		}
	}
}
