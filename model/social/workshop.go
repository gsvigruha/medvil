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
		if mp.Price(w.Manufacture.Inputs)*2 < mp.Price(w.Manufacture.Outputs) {
			needs := w.Household.Resources.Needs(artifacts.Multiply(w.Manufacture.Inputs, ProductTransportQuantity))
			tag := "manufacture_input"
			if needs != nil && w.Household.NumTasks("exchange", tag) == 0 {
				if w.Household.Money >= mp.Price(needs) {
					w.Household.AddTask(&economy.BuyTask{
						Exchange:       mp,
						HouseholdMoney: &w.Household.Money,
						Goods:          needs,
						MaxPrice:       mp.Price(w.Manufacture.Outputs) * ProductTransportQuantity / 2,
						TaskTag:        tag,
					})
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
					if NumBatchesSimple(int(qToSell), ProductTransportQuantity) > w.Household.NumTasks("exchange", tag) {
						w.Household.AddTask(&economy.SellTask{
							Exchange: mp,
							Goods:    goods,
							TaskTag:  tag,
						})
					}
				}
			}
		}
	}
}
