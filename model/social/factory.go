package social

import (
	"medvil/model/artifacts"
	"medvil/model/economy"
	"medvil/model/navigation"
	"medvil/model/time"
	//"medvil/model/vehicles"
)

type VehicleOrder struct {
	T *economy.VehicleConstruction
}

type Factory struct {
	Household Household
	Orders    []*VehicleOrder
}

func (f *Factory) ElapseTime(Calendar *time.CalendarType, m navigation.IMap) {
	f.Household.ElapseTime(Calendar, m)

	//home := m.GetField(f.Household.Building.X, f.Household.Building.Y)
	mp := f.Household.Town.Marketplace

	var newOrders []*VehicleOrder
	for _, order := range f.Orders {
		needs := artifacts.Purchasable(order.T.Inputs)
		tag := "order_input"
		if needs != nil && f.Household.Money >= mp.Price(needs) {
			f.Household.AddTask(&economy.BuyTask{
				Exchange:       mp,
				HouseholdMoney: &f.Household.Money,
				Goods:          needs,
				MaxPrice:       mp.Price(needs) * 2,
				TaskTag:        tag,
			})
		} else {
			newOrders = append(newOrders, order)
		}
	}
	f.Orders = newOrders
}

func (f *Factory) NumOrders(vc *economy.VehicleConstruction) int {
	var n int
	for _, order := range f.Orders {
		if order.T == vc {
			n++
		}
	}
	return n
}
