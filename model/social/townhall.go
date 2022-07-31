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

const ConstructionBudgetPercentage = 0.3
const ConstructionStorageCapacity = 0.7

func (t *Townhall) ElapseTime(Calendar *time.CalendarType, m navigation.IMap) {
	t.Household.ElapseTime(Calendar, m)
	mp := t.Household.Town.Marketplace

	for _, a := range artifacts.All {
		tag := "exchange#" + a.Name
		purchaseQuantity := ProductTransportQuantity(a)
		goods := []artifacts.Artifacts{artifacts.Artifacts{A: a, Quantity: purchaseQuantity}}
		if q, ok := t.Household.Resources.Artifacts[a]; ok {
			if t.Household.NumTasks("exchange", tag) == 0 {
				targetQ := uint16(*(t.StorageTarget[a]))
				if q > targetQ {
					t.Household.AddTask(&economy.SellTask{
						Exchange: mp,
						Goods:    goods,
						TaskTag:  tag,
					})
				} else if q < targetQ {
					if t.Household.Money >= mp.Price(goods) && mp.HasTraded(a) {
						t.Household.AddTask(&economy.BuyTask{
							Exchange:       mp,
							HouseholdMoney: &t.Household.Money,
							Goods:          goods,
							MaxPrice:       uint32(mp.Price(goods) * 2),
							TaskTag:        tag,
						})
					}
				}
			}
		}
	}
}
