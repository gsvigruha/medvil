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

func (w *Workshop) IsManufactureProfitable() bool {
	if w.Manufacture != nil {
		mp := w.Household.Town.Marketplace
		return float64(mp.Price(w.Manufacture.Outputs)) >= float64(mp.Price(artifacts.Purchasable(w.Manufacture.Inputs)))*ProfitCostRatio
	}
	return false
}

func (w *Workshop) ElapseTime(Calendar *time.CalendarType, m navigation.IMap) {
	w.Household.ElapseTime(Calendar, m)
	home := m.GetField(w.Household.Building.X, w.Household.Building.Y)
	if w.Manufacture != nil {
		mp := w.Household.Town.Marketplace
		purchasableInputs := artifacts.Purchasable(w.Manufacture.Inputs)
		maxUnitCost := float64(mp.Price(w.Manufacture.Outputs)) / ProfitCostRatio
		if float64(mp.Price(purchasableInputs)) < maxUnitCost {
			transportQ := MinProductTransportQuantity(purchasableInputs)
			batch := artifacts.Multiply(purchasableInputs, transportQ)
			tag := "manufacture_input"
			if w.Household.Resources.Needs(batch) != nil && w.Household.NumTasks("exchange", tag) == 0 {
				if w.Household.Money >= mp.Price(batch) {
					w.Household.AddTask(&economy.BuyTask{
						Exchange:       mp,
						HouseholdMoney: &w.Household.Money,
						Goods:          batch,
						MaxPrice:       uint32(maxUnitCost * float64(transportQ)),
						TaskTag:        tag,
					})
					numP := uint16(len(w.Household.People))
					water := artifacts.GetArtifact("water")
					if w.Manufacture.IsInput(water) &&
						w.Household.Resources.Get(water) < economy.MinFoodOrDrinkPerPerson*numP+WaterTransportQuantity &&
						w.Household.NumTasks("transport", "water") == 0 {
						hf := w.Household.RandomField(m, navigation.Field.BuildingNonExtension)
						pickup := m.FindDest(navigation.Location{X: hf.X, Y: hf.Y, Z: 0}, economy.WaterDestination{}, navigation.PathTypePedestrian)
						if pickup != nil {
							w.Household.AddPriorityTask(&economy.TransportTask{
								PickupD:  pickup,
								DropoffD: w.Household.Destination(navigation.Field.BuildingNonExtension),
								PickupR:  &pickup.Terrain.Resources,
								DropoffR: &w.Household.Resources,
								A:        water,
								Quantity: WaterTransportQuantity,
							})
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

		w.Household.SellArtifacts(w.Manufacture.IsInput, w.Manufacture.IsOutput)
	}
	w.Household.MaybeBuyBoat(Calendar, m)
	w.Household.MaybeBuyCart(Calendar, m)
}

func (w *Workshop) GetFields() []navigation.FieldWithContext {
	return []navigation.FieldWithContext{}
}

func (w *Workshop) GetHousehold() *Household {
	return &w.Household
}
