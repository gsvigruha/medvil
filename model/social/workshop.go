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
		if mp.Price(w.Manufacture.Inputs) < mp.Price(w.Manufacture.Outputs) {
			for _, a := range w.Manufacture.Inputs {
				needs := w.Household.Resources.Needs(a.Multiply(ProductTransportQuantity).Wrap())
				tag := "manufacture_input#" + a.A.Name
				if needs != nil && w.Household.NumTasks("exchange", tag) == 0 {
					if w.Household.Money >= mp.Price(needs) {
						if mp.Storage.Has(needs) {
							mx, my := mp.Building.GetRandomBuildingXY()
							w.Household.AddTask(&economy.ExchangeTask{
								HomeF:          home,
								MarketF:        m.GetField(mx, my),
								Exchange:       mp,
								HouseholdR:     &w.Household.Resources,
								HouseholdMoney: &w.Household.Money,
								GoodsToBuy:     needs,
								GoodsToSell:    nil,
								TaskTag:        tag,
							})
						} else if Calendar.Day == 30 && Calendar.Hour == 0 {
							mp.RegisterDemand(needs)
						}
					}
				}
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
				if qToSell > 0 {
					tag := "sell_artifacts#" + a.Name
					goods := []artifacts.Artifacts{artifacts.Artifacts{A: a, Quantity: ProductTransportQuantity}}
					if w.Household.Town.Marketplace.CanSell(goods) && NumBatchesSimple(int(qToSell), ProductTransportQuantity) > w.Household.NumTasks("exchange", tag) {
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
