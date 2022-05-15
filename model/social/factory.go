package social

import (
	"medvil/model/artifacts"
	"medvil/model/economy"
	"medvil/model/navigation"
	"medvil/model/time"
	//"medvil/model/vehicles"
)

const OrderStateStart = 0
const OrderStateInputsPurchased = 1
const OrderStateBuilt = 2
const OrderStateDone = 3

type VehicleOrder struct {
	T     *economy.VehicleConstruction
	State uint8
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
		if order.State == OrderStateStart {
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
				order.State = OrderStateInputsPurchased
			}
		} else if order.State == OrderStateInputsPurchased {
		} else if order.State == OrderStateBuilt {
		} else {
			newOrders = append(newOrders, order)
		}
	}
	f.Orders = newOrders
}

func (f *Factory) CreateOrder(vc *economy.VehicleConstruction) {
	f.Orders = append(f.Orders, &VehicleOrder{T: vc, State: OrderStateStart})
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
