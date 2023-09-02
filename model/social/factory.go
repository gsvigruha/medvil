package social

import (
	"medvil/model/artifacts"
	"medvil/model/building"
	"medvil/model/economy"
	"medvil/model/navigation"
	"medvil/model/time"
	"medvil/model/vehicles"
)

const OrderStateOrdered = 0
const OrderStateStarted = 1
const OrderStateBuilt = 2
const OrderStateComplete = 3
const OrderStateExpired = 4

type VehicleOrder struct {
	T       *economy.VehicleConstruction
	State   uint8
	F       *Factory
	Vehicle *vehicles.Vehicle
}

func (o *VehicleOrder) IsBuilt() bool {
	return o.State == OrderStateBuilt
}

func (o *VehicleOrder) IsExpired() bool {
	return o.State == OrderStateExpired
}

func (o *VehicleOrder) PickupVehicle() *vehicles.Vehicle {
	o.State = OrderStateComplete
	return o.Vehicle
}

func (o *VehicleOrder) Name() string {
	return o.T.Name
}

func (o *VehicleOrder) CompleteBuild(f *navigation.Field) {
	o.State = OrderStateBuilt
	var travellerType uint8
	if o.T.Output == vehicles.Boat {
		travellerType = navigation.TravellerTypeBoat
	} else if o.T.Output == vehicles.Cart {
		travellerType = navigation.TravellerTypeCart
	} else if o.T.Output == vehicles.TradingBoat {
		travellerType = navigation.TravellerTypeTradingBoat
	} else if o.T.Output == vehicles.TradingCart {
		travellerType = navigation.TravellerTypeTradingCart
	}
	vehicle := &vehicles.Vehicle{T: o.T.Output, Traveller: &navigation.Traveller{
		FX:      f.X,
		FY:      f.Y,
		FZ:      0,
		PX:      50,
		PY:      50,
		Visible: true,
		T:       travellerType,
	}}
	o.Vehicle = vehicle
	f.RegisterTraveller(vehicle.Traveller)
}

type Factory struct {
	Household *Household
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
					Exchange:        mp,
					HouseholdWallet: f.Household,
					Goods:           inputs,
					MaxPrice:        mp.Price(inputs) * 2,
					TaskTag:         tag,
				})
			}
		}
	}

	var newOrders []*VehicleOrder
	for _, order := range f.Orders {
		if order.State == OrderStateOrdered && f.Household.Resources.RemoveAll(order.T.Inputs) {
			var field *navigation.Field
			if order.T.Output.Water {
				ext, x, y := f.Household.Building.GetExtensionWithCoords(building.Deck)
				if ext == nil {
					order.State = OrderStateExpired
					continue
				}
				field = m.GetField(x, y)
			} else {
				field = f.Household.RandomField(m, navigation.Field.BuildingNonExtension)
			}
			f.Household.AddTask(&economy.VehicleConstructionTask{
				T: order.T,
				O: order,
				F: field,
				R: f.Household.Resources,
			})
			order.State = OrderStateStarted
		}
		if order.State != OrderStateComplete {
			newOrders = append(newOrders, order)
		}
	}
	f.Orders = newOrders
}

func VehiclePrice(mp *Marketplace, vc *economy.VehicleConstruction) uint32 {
	return mp.Price(artifacts.Purchasable(vc.Inputs)) * 2
}

func (f *Factory) Price(vc *economy.VehicleConstruction) uint32 {
	return VehiclePrice(f.Household.Town.Marketplace, vc)
}

func (f *Factory) CreateOrder(vc *economy.VehicleConstruction, h *Household) *VehicleOrder {
	price := f.Price(vc)
	if h.Money >= price {
		order := &VehicleOrder{T: vc, F: f, State: OrderStateOrdered}
		f.Orders = append(f.Orders, order)
		h.Money -= price
		f.Household.Money += price
		return order
	}
	return nil
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

func PickFactory(fs []*Factory, et *building.BuildingExtensionType) *Factory {
	var f *Factory = nil
	var orders = 0
	for _, fi := range fs {
		if et == building.NonExtension || fi.Household.Building.HasExtension(et) {
			if f == nil || orders < len(fi.Orders) {
				f = fi
				orders = len(fi.Orders)
			}
		}
	}
	return f
}

func GetVehicleConstructions(factories []*Factory, filter func(*economy.VehicleConstruction) bool) []*economy.VehicleConstruction {
	result := make([]*economy.VehicleConstruction, 0, len(economy.AllVehicleConstruction))
	for _, vc := range economy.AllVehicleConstruction {
		if filter(vc) {
			for _, factory := range factories {
				extensions := factory.Household.Building.Plan.GetExtensions()
				if economy.ConstructionCompatible(vc, extensions) {
					result = append(result, vc)
					break
				}
			}
		}
	}
	return result
}

func (f *Factory) GetFields() []navigation.FieldWithContext {
	return []navigation.FieldWithContext{}
}

func (f *Factory) GetHome() Home {
	return f.Household
}
