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
}

const StorageRefillBudgetPercentage = 0.5
const ConstructionStorageCapacity = 0.7

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
}
