package social

import (
	"medvil/model/artifacts"
	"medvil/model/economy"
	"medvil/model/navigation"
	"medvil/model/time"
)

type Townhall struct {
	Household     Household
	StorageTarget map[*artifacts.Artifact]*int
	Traders       []*Trader
}

const StorageRefillBudgetPercentage = 0.5
const ConstructionStorageCapacity = 0.7

const TownhallMaxDistance = 25

func (t *Townhall) ElapseTime(Calendar *time.CalendarType, m navigation.IMap) {
	t.Household.ElapseTime(Calendar, m)
	mp := t.Household.Town.Marketplace

	for _, a := range artifacts.All {
		tag := "storage_target#" + a.Name
		transportQuantity := ProductTransportQuantity(a)
		goods := []artifacts.Artifacts{artifacts.Artifacts{A: a, Quantity: transportQuantity}}
		if q, ok := t.Household.Resources.Artifacts[a]; ok {
			if t.Household.NumTasks("exchange", tag) == 0 {
				targetQ := uint16(*(t.StorageTarget[a]))
				if q > targetQ {
					qToSell := t.Household.ArtifactToSell(a, q, false)
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
							Exchange:       mp,
							HouseholdMoney: &t.Household.Money,
							Goods:          goods,
							MaxPrice:       maxPrice,
							TaskTag:        tag,
						})
					}
				}
			}
		}
	}

	for _, trader := range t.Traders {
		trader.ElapseTime(Calendar, m)
	}
}

func (t *Townhall) GetFields() []navigation.FieldWithContext {
	fields := make([]navigation.FieldWithContext, len(t.Household.Town.Roads))
	for i := range t.Household.Town.Roads {
		fields[i] = t.Household.Town.Roads[i]
	}
	return fields
}

func (t *Townhall) FieldWithinDistance(field *navigation.Field) bool {
	if t.Household.Building == nil {
		return true
	}
	return WithinDistance(t.Household.Building, field, TownhallMaxDistance)
}

func (t *Townhall) CreateTrader() {
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
						Resources:      r,
						SourceExchange: t.Household.Town.Marketplace,
					}
					t.Traders = append(t.Traders, trader)
					p.Home = trader
					p.Traveller.UseVehicle(v)
					t.Household.Vehicles = append(t.Household.Vehicles[:i], t.Household.Vehicles[i+1:]...)
					t.Household.People = append(t.Household.People[:j], t.Household.People[j+1:]...)
					return
				}
			}
		}
	}
}
