package social

import (
	"medvil/model/artifacts"
	"medvil/model/economy"
	"medvil/model/navigation"
	"medvil/model/time"
	"medvil/model/vehicles"
)

const OrderStateOrdered = 0
const OrderStateStarted = 1
const OrderStateBuilt = 2

type VehicleOrder struct {
	T     *economy.VehicleConstruction
	State uint8
	F     *Factory
}

func (o *VehicleOrder) CompleteBuild(f *navigation.Field) {
	o.State = OrderStateBuilt
	vehicle := &vehicles.Vehicle{T: o.T.Output, Traveller: &navigation.Traveller{
		FX:      f.X,
		FY:      f.Y,
		FZ:      0,
		PX:      50,
		PY:      50,
		Visible: true,
		T:       navigation.TravellerTypeBoat,
	}}
	o.F.Household.Vehicles = append(o.F.Household.Vehicles, vehicle)
	f.RegisterTraveller(vehicle.Traveller)
}

type Factory struct {
	Household Household
	Orders    []*VehicleOrder
}

func (f *Factory) ElapseTime(Calendar *time.CalendarType, m navigation.IMap) {
	f.Household.ElapseTime(Calendar, m)

	//home := m.GetField(f.Household.Building.X, f.Household.Building.Y)
	mp := f.Household.Town.Marketplace

	needs := artifacts.Resources{}
	needs.Init(0)
	for _, order := range f.Orders {
		if order.State == OrderStateOrdered {
			needs.AddAll(artifacts.Purchasable(order.T.Inputs))
		}
	}

	for a, q := range needs.Artifacts {
		tag := "order_input#" + a.Name
		transportQ := ProductTransportQuantity(a)
		e := f.Household.Resources.Get(a)
		if q > e && NumBatchesSimple(q-e, transportQ)+1 > f.Household.NumTasks("exchange", tag) {
			inputs := artifacts.Artifacts{A: a, Quantity: transportQ}.Wrap()
			if f.Household.Money >= mp.Price(inputs) {
				f.Household.AddTask(&economy.BuyTask{
					Exchange:       mp,
					HouseholdMoney: &f.Household.Money,
					Goods:          inputs,
					MaxPrice:       mp.Price(inputs) * 2,
					TaskTag:        tag,
				})
			}
		}
	}

	for _, order := range f.Orders {
		if order.State == OrderStateOrdered && f.Household.Resources.RemoveAll(order.T.Inputs) {
			_, x, y := f.Household.Building.GetExtensionWithCoords()
			deck := m.GetField(x, y)
			f.Household.AddTask(&economy.VehicleConstructionTask{
				T: order.T,
				O: order,
				F: deck,
				R: &f.Household.Resources,
			})
			order.State = OrderStateStarted
		}
	}
}

func (f *Factory) Price(vc *economy.VehicleConstruction) uint32 {
	return f.Household.Town.Marketplace.Price(artifacts.Purchasable(vc.Inputs)) * 2
}

func (f *Factory) CreateOrder(vc *economy.VehicleConstruction) {
	f.Orders = append(f.Orders, &VehicleOrder{T: vc, F: f, State: OrderStateOrdered})
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

func PickFactory(fs []*Factory) *Factory {
	var f *Factory = nil
	var orders = 0
	for _, fi := range fs {
		if f == nil || orders < len(fi.Orders) {
			f = fi
		}
	}
	return f
}
