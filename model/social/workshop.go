package social

import (
	"medvil/model/artifacts"
	"medvil/model/building"
	"medvil/model/economy"
	"medvil/model/navigation"
	"medvil/model/time"
)

type Workshop struct {
	Household   *Household
	Manufacture *economy.Manufacture
	AutoSwitch  bool
}

const ProfitCostRatio = 2.0

var Paper = artifacts.GetArtifact("paper")

func (w *Workshop) IsManufactureProfitable() bool {
	mp := w.Household.Town.Marketplace
	if w.Manufacture != nil && mp != nil {
		return float64(mp.Price(w.Manufacture.Outputs)) >= float64(mp.Price(artifacts.Purchasable(w.Manufacture.Inputs)))*ProfitCostRatio
	}
	return false
}

func (w *Workshop) ElapseTime(Calendar *time.CalendarType, m navigation.IMap) {
	w.Household.ElapseTime(Calendar, m)
	home := m.GetField(w.Household.Building.X, w.Household.Building.Y)
	mp := w.Household.Town.Marketplace
	numP := uint16(len(w.Household.People))

	if mp != nil {
		if w.AutoSwitch && Calendar.Day == 30 && Calendar.Hour == 0 && Calendar.Month%3 == 0 && w.Household.Resources.Remove(Paper, 1) > 0 {
			var maxProfit = 0.0
			for _, mName := range economy.GetManufactureNames(w.Household.Building.Plan.GetExtensions()) {
				manufacture := economy.GetManufacture(mName)
				profit := float64(mp.Price(manufacture.Outputs)) / float64(mp.Price(artifacts.Purchasable(manufacture.Inputs)))
				if profit > maxProfit && (w.Household.Town.Settings.Coinage || manufacture.Name != "goldsmith") {
					maxProfit = profit
					w.Manufacture = manufacture
				}
			}
		}
		if w.Manufacture != nil && (w.Household.Town.Settings.Coinage || w.Manufacture.Name != "goldsmith") {
			purchasableInputs := artifacts.Purchasable(w.Manufacture.Inputs)
			maxUnitCost := float64(mp.Price(w.Manufacture.Outputs)) / ProfitCostRatio
			if float64(mp.Price(purchasableInputs)) < maxUnitCost {
				transportQ := MinProductTransportQuantity(purchasableInputs)
				batch := artifacts.Multiply(purchasableInputs, transportQ)
				tag := economy.SingleTag(economy.TagManufactureInput)
				if w.Household.Resources.Needs(batch) != nil && w.Household.NumTasks("exchange", tag) < int(numP)/2 && len(w.Household.Tasks) < int(numP)*5 {
					if w.Household.Money >= mp.Price(batch) {
						w.Household.AddTask(&economy.BuyTask{
							Exchange:        mp,
							HouseholdWallet: w.Household,
							Goods:           batch,
							MaxPrice:        uint32(maxUnitCost * float64(transportQ)),
							TaskTag:         tag,
						})
						water := artifacts.GetArtifact("water")
						if w.Manufacture.IsInput(water) &&
							w.Household.Resources.Get(water) < economy.MinFoodOrDrinkPerPerson*numP+WaterTransportQuantity &&
							w.Household.NumTasks("transport", economy.SingleTag(water.Idx)) == 0 {
							hf := w.Household.RandomField(m, navigation.Field.BuildingNonExtension)
							pickup := m.FindDest(navigation.Location{X: hf.X, Y: hf.Y, Z: 0}, economy.WaterDestination{}, navigation.PathTypePedestrian)
							if pickup != nil {
								w.Household.AddPriorityTask(&economy.TransportTask{
									PickupD:        pickup,
									DropoffD:       w.Household.Destination(building.NonExtension),
									PickupR:        pickup.Terrain.Resources,
									DropoffR:       w.Household.Resources,
									A:              water,
									TargetQuantity: WaterTransportQuantity,
								})
							}
						}
					}
				}

				if w.Household.Resources.RemoveAll(w.Manufacture.Inputs) {
					w.Household.AddTask(&economy.ManufactureTask{
						M: w.Manufacture,
						F: home,
						R: w.Household.Resources,
					})
				}
			}

			w.Household.SellArtifacts(w.Manufacture.IsInput, w.Manufacture.IsOutput)
		}

		paperTag := economy.SingleTag(economy.TagPaperPurchase)
		if w.AutoSwitch && w.Household.Resources.Get(Paper) < ProductTransportQuantity(Paper) && w.Household.NumTasks("exchange", paperTag) == 0 {
			needs := []artifacts.Artifacts{artifacts.Artifacts{A: Paper, Quantity: ProductTransportQuantity(Paper)}}
			if w.Household.Money >= mp.Price(needs) && mp.HasTraded(Paper) {
				w.Household.AddTask(&economy.BuyTask{
					Exchange:        mp,
					HouseholdWallet: w.Household,
					Goods:           needs,
					MaxPrice:        uint32(float64(w.Household.Money) * ExtrasBudgetRatio),
					TaskTag:         paperTag,
				})
			}
		}

		w.Household.MaybeBuyBoat(Calendar, m)
		w.Household.MaybeBuyCart(Calendar, m)
	}
}

func (w *Workshop) GetFields() []navigation.FieldWithContext {
	return []navigation.FieldWithContext{}
}

func (w *Workshop) GetHome() Home {
	return w.Household
}
