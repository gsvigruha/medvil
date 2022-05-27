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

const ProfitCostRatio = 2.0

func (w *Workshop) ElapseTime(Calendar *time.CalendarType, m navigation.IMap) {
	w.Household.ElapseTime(Calendar, m)
	home := m.GetField(w.Household.Building.X, w.Household.Building.Y)
	if w.Manufacture != nil {
		mp := w.Household.Town.Marketplace
		purchasableInputs := artifacts.Purchasable(w.Manufacture.Inputs)
		maxUnitCost := float64(mp.Price(w.Manufacture.Outputs)) / ProfitCostRatio
		if float64(mp.Price(purchasableInputs)) < maxUnitCost {
			transportQ := MinProductTransportQuantity(purchasableInputs)
			needs := w.Household.Resources.Needs(artifacts.Multiply(purchasableInputs, transportQ))
			tag := "manufacture_input"
			if needs != nil && w.Household.NumTasks("exchange", tag) == 0 {
				if w.Household.Money >= mp.Price(needs) {
					w.Household.AddTask(&economy.BuyTask{
						Exchange:       mp,
						HouseholdMoney: &w.Household.Money,
						Goods:          needs,
						MaxPrice:       uint32(maxUnitCost * float64(transportQ)),
						TaskTag:        tag,
					})
					numP := uint16(len(w.Household.People))
					water := artifacts.GetArtifact("water")
					if w.Manufacture.IsInput(water) && w.Household.Resources.Get(water) < economy.MinFoodOrDrinkPerPerson*numP+WaterTransportQuantity {
						hx, hy, ok := GetRandomBuildingXY(w.Household.Building, m, navigation.Field.Walkable)
						if ok {
							dest := m.FindDest(navigation.Location{X: hx, Y: hy, Z: 0}, economy.WaterDestination{}, navigation.TravellerTypePedestrian)
							if dest != nil {
								w.Household.AddTask(&economy.TransportTask{
									PickupF:  dest,
									DropoffF: m.GetField(hx, hy),
									PickupR:  &dest.Terrain.Resources,
									DropoffR: &w.Household.Resources,
									A:        water,
									Quantity: WaterTransportQuantity,
								})
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
		}

		for a, q := range w.Household.Resources.Artifacts {
			if !w.Manufacture.IsInput(a) {
				qToSell := w.Household.ArtifactToSell(a, q, w.Manufacture.IsOutput(a))
				if qToSell > 0 {
					tag := "sell_artifacts#" + a.Name
					goods := []artifacts.Artifacts{artifacts.Artifacts{A: a, Quantity: ProductTransportQuantity(a)}}
					if NumBatchesSimple(qToSell, ProductTransportQuantity(a)) > w.Household.NumTasks("exchange", tag) {
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
