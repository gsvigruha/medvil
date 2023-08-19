package social

import (
	"medvil/model/artifacts"
	"medvil/model/economy"
	"medvil/model/navigation"
	"medvil/model/time"
)

type Townhall struct {
	Household     *Household
	StorageTarget map[*artifacts.Artifact]int
	Traders       []*Trader
}

const StorageRefillBudgetPercentage = 0.5
const ConstructionStorageCapacity = 0.7
const PaperBudgetRatio = 0.1

const TownhallMaxDistance = 25

func (t *Townhall) ElapseTime(Calendar *time.CalendarType, m navigation.IMap) {
	t.Household.ElapseTime(Calendar, m)
	mp := t.Household.Town.Marketplace

	for _, a := range artifacts.All {
		tag := "storage_target#" + a.Name
		transportQuantity := ProductTransportQuantity(a)
		goods := []artifacts.Artifacts{artifacts.Artifacts{A: a, Quantity: transportQuantity}}
		var q uint16 = 0
		if storageQ, ok := t.Household.Resources.Artifacts[a]; ok {
			q = storageQ
		}
		if t.Household.NumTasks("exchange", tag) == 0 {
			targetQ := uint16(t.StorageTarget[a])
			if q > targetQ {
				qToSell := t.Household.ArtifactToSell(a, q, false, false)
				if qToSell > 0 {
					t.Household.AddTask(&economy.SellTask{
						Exchange: mp,
						Goods:    goods,
						TaskTag:  tag,
					})
				}
			} else if q < targetQ {
				maxPrice := uint32(float64(t.Household.Money) * StorageRefillBudgetPercentage / float64(len(t.Household.Resources.Artifacts)))
				if t.Household.Money >= mp.Price(goods) && mp.HasTraded(a) {
					t.Household.AddTask(&economy.BuyTask{
						Exchange:        mp,
						HouseholdWallet: t.Household,
						Goods:           goods,
						MaxPrice:        maxPrice,
						TaskTag:         tag,
					})
				}
			}
		}
	}

	if t.Household.Resources.Get(Paper) < ProductTransportQuantity(Paper) && t.Household.NumTasks("exchange", "paper_purchase") == 0 {
		needs := []artifacts.Artifacts{artifacts.Artifacts{A: Paper, Quantity: ProductTransportQuantity(Paper)}}
		if t.Household.Money >= mp.Price(needs) && mp.HasTraded(Paper) {
			t.Household.AddTask(&economy.BuyTask{
				Exchange:        mp,
				HouseholdWallet: t.Household,
				Goods:           needs,
				MaxPrice:        uint32(float64(t.Household.Money) * PaperBudgetRatio),
				TaskTag:         "paper_purchase",
			})
		}
	}

	for _, trader := range t.Traders {
		trader.ElapseTime(Calendar, m)
	}
}

func (t *Townhall) GetFields() []navigation.FieldWithContext {
	fields := make([]navigation.FieldWithContext, len(t.Household.Town.Roads)+len(t.Household.Town.Walls))
	for i := range t.Household.Town.Roads {
		fields[i] = t.Household.Town.Roads[i]
	}
	for i := range t.Household.Town.Walls {
		fields[i+len(t.Household.Town.Roads)] = t.Household.Town.Walls[i]
	}
	return fields
}

func (t *Townhall) FieldWithinDistance(field *navigation.Field) bool {
	if t.Household.Building == nil {
		return true
	}
	return WithinDistance(t.Household.Building, field, TownhallMaxDistance)
}

func (t *Townhall) getTraderDestField(trader *Trader, m navigation.IMap) *navigation.Field {
	hf := t.Household.RandomField(m, trader.Vehicle.T.BuildingCheckFn)
	dest := &navigation.BuildingDestination{B: t.Household.Town.Marketplace.Building, ET: trader.Vehicle.T.BuildingExtensionType}
	if hf != nil {
		path := m.ShortPath(hf.GetLocation(), dest, trader.Person.Traveller.PathType())
		if path != nil && len(path.P) > 2 {
			pe := path.P[len(path.P)-2]
			return m.GetField(pe.GetLocation().X, pe.GetLocation().Y)
		}
	}
	return nil
}

func (t *Townhall) CreateTrader(m navigation.IMap) {
	for i, v := range t.Household.Vehicles {
		if !v.InUse {
			for j, p := range t.Household.People {
				if p.Task == nil {
					var r artifacts.Resources
					r.Init(v.T.MaxVolume)
					trader := &Trader{
						Money:          0,
						Person:         p,
						Vehicle:        v,
						Resources:      &r,
						SourceExchange: t.Household.Town.Marketplace,
					}
					t.Traders = append(t.Traders, trader)
					p.Home = trader
					p.Traveller.UseVehicle(v)
					t.Household.Vehicles = append(t.Household.Vehicles[:i], t.Household.Vehicles[i+1:]...)
					t.Household.People = append(t.Household.People[:j], t.Household.People[j+1:]...)
					homeField := t.getTraderDestField(trader, m)
					p.Task = &economy.GoHomeTask{F: homeField, P: p}
					return
				}
			}
		}
	}
}

func (t *Townhall) Filter(Calendar *time.CalendarType, m navigation.IMap) {
	var newTraders = make([]*Trader, 0, len(t.Traders))
	for _, trader := range t.Traders {
		if trader.Person.Health == 0 {
			field := m.GetField(trader.Person.Traveller.FX, trader.Person.Traveller.FY)
			field.UnregisterTraveller(trader.Person.Traveller)
			field.UnregisterTraveller(trader.Vehicle.Traveller)
			t.Household.Money += trader.Money
		} else {
			newTraders = append(newTraders, trader)
		}
	}
	t.Traders = newTraders
}
