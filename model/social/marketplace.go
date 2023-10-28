package social

import (
	"medvil/model/artifacts"
	"medvil/model/building"
	"medvil/model/economy"
	"medvil/model/navigation"
	"medvil/model/stats"
	"medvil/model/time"
)

const StorageToSoldRatio = 12

type Marketplace struct {
	Town      *Town `json:"-"`
	Building  *building.Building
	Money     uint32
	Storage   artifacts.Resources
	Prices    map[*artifacts.Artifact]uint32
	Sold      map[*artifacts.Artifact]uint32
	Bought    map[*artifacts.Artifact]uint32
	BuyTasks  map[*economy.BuyTask]bool
	SellTasks map[*economy.SellTask]bool
}

func (mp *Marketplace) prices(as []*artifacts.Artifact) []uint32 {
	prices := make([]uint32, len(as))
	for i, a := range as {
		prices[i] = mp.Prices[a]
	}
	return prices
}

func (mp *Marketplace) Init() {
	mp.Prices = make(map[*artifacts.Artifact]uint32)
	mp.Sold = make(map[*artifacts.Artifact]uint32)
	mp.Bought = make(map[*artifacts.Artifact]uint32)
	mp.BuyTasks = make(map[*economy.BuyTask]bool)
	mp.SellTasks = make(map[*economy.SellTask]bool)
	for _, a := range artifacts.All {
		mp.Prices[a] = 10
		mp.Reset(a)
		mp.Storage.Add(a, 0)
		for _, m := range economy.AllManufacture {
			for _, o := range m.Outputs {
				if a == o.A {
					mp.Prices[a] = 20
				}
			}
		}
	}
}

func (mp *Marketplace) Reset(a *artifacts.Artifact) {
	mp.Sold[a] = 0
	mp.Bought[a] = 0
}

type SupplyAndDemand struct {
	Supply uint32
	Demand uint32
}

var GoldCoin = artifacts.GetArtifact("gold_coin")

func (mp *Marketplace) pendingSupplyAndDemand() map[*artifacts.Artifact]*SupplyAndDemand {
	sd := make(map[*artifacts.Artifact]*SupplyAndDemand)
	for _, a := range artifacts.All {
		sd[a] = &SupplyAndDemand{Supply: 0, Demand: 0}
	}
	for t := range mp.BuyTasks {
		for _, a := range t.Goods {
			sd[a.A].Demand += uint32(a.Quantity)
		}
	}
	for t := range mp.SellTasks {
		for _, a := range t.Goods {
			sd[a.A].Supply += uint32(a.Quantity)
		}
	}
	return sd
}

func (mp *Marketplace) IncPrice(a *artifacts.Artifact) {
	var delta = float64(mp.Prices[a]) * 0.1
	if delta < 1.0 {
		delta = 1.0
	}
	mp.Prices[a] += uint32(delta)
	if mp.Prices[a] > 999 {
		mp.Prices[a] = 999
	}
}

func (mp *Marketplace) DecPrice(a *artifacts.Artifact) {
	var delta = float64(mp.Prices[a]) * 0.1
	if delta < 1.0 {
		delta = 1.0
	}
	if mp.Prices[a] > 1 {
		mp.Prices[a] -= uint32(delta)
	}
}

func (mp *Marketplace) ElapseTime(Calendar *time.CalendarType, m navigation.IMap) {
	if Calendar.Hour == 0 && Calendar.Day == 1 {
		allGold := []artifacts.Artifacts{artifacts.Artifacts{A: GoldCoin, Quantity: mp.Storage.Get(GoldCoin)}}
		price := mp.Price(allGold)
		wallet := &mp.Town.Townhall.Household.Money
		*wallet += price * 2
		mp.Buy(allGold, mp.Town.Townhall.Household)
		sd := mp.pendingSupplyAndDemand()

		for _, a := range artifacts.All {
			if a != GoldCoin {
				supply := sd[a].Supply + uint32(mp.Storage.Artifacts[a])/StorageToSoldRatio
				demand := sd[a].Demand
				if supply < demand {
					mp.IncPrice(a)
				} else if supply > demand {
					mp.DecPrice(a)
				} else if mp.Storage.Artifacts[a] >= ProductTransportQuantity(a) && demand == 0 {
					mp.DecPrice(a)
				}

				mp.Reset(a)
			}
		}
	}
}

func (mp *Marketplace) RegisterSellTask(t *economy.SellTask, add bool) {
	if add {
		mp.SellTasks[t] = true
	} else {
		delete(mp.SellTasks, t)
	}
}

func (mp *Marketplace) RegisterBuyTask(t *economy.BuyTask, add bool) {
	if add {
		mp.BuyTasks[t] = true
	} else {
		delete(mp.BuyTasks, t)
	}
}

func (mp *Marketplace) Buy(as []artifacts.Artifacts, wallet economy.Wallet) {
	price := mp.Price(as)
	mp.Storage.RemoveAll(as)
	mp.Money += price
	wallet.Spend(price)
	for _, a := range as {
		mp.Bought[a.A] += uint32(a.Quantity)
		mp.Town.Stats.RegisterTrade(a.A, mp.Prices[a.A], a.Quantity)
	}
}

func (mp *Marketplace) BuyAsManyAsPossible(as []artifacts.Artifacts, wallet economy.Wallet) []artifacts.Artifacts {
	existingArtifacts := mp.Storage.GetAsManyAsPossible(as)
	price := mp.Price(existingArtifacts)
	mp.Money += price
	wallet.Spend(price)
	for _, a := range existingArtifacts {
		mp.Bought[a.A] += uint32(a.Quantity)
		mp.Town.Stats.RegisterTrade(a.A, mp.Prices[a.A], a.Quantity)
	}
	return existingArtifacts
}

func (mp *Marketplace) Sell(as []artifacts.Artifacts, wallet economy.Wallet) {
	price := mp.Price(as)
	mp.Storage.AddAll(as)
	mp.Money -= price
	wallet.Earn(price)
	for _, a := range as {
		mp.Sold[a.A] += uint32(a.Quantity)
		mp.Town.Stats.RegisterTrade(a.A, mp.Prices[a.A], a.Quantity)
	}
}

func (mp *Marketplace) SellAsManyAsPossible(as []artifacts.Artifacts, wallet economy.Wallet) []artifacts.Artifacts {
	var leftover []artifacts.Artifacts
	for _, a := range as {
		var price uint32
		var quantity uint16
		if mp.Prices[a.A]*uint32(a.Quantity) <= mp.Money {
			price = mp.Prices[a.A] * uint32(a.Quantity)
			quantity = a.Quantity
		} else {
			quantity = uint16(mp.Money / mp.Prices[a.A])
			price = mp.Prices[a.A] * uint32(quantity)
		}
		wallet.Earn(price)
		mp.Money -= price
		mp.Storage.Add(a.A, quantity)
		mp.Sold[a.A] += uint32(quantity)
		mp.Town.Stats.RegisterTrade(a.A, mp.Prices[a.A], quantity)
		leftover = append(leftover, artifacts.Artifacts{A: a.A, Quantity: a.Quantity - quantity})
	}
	return leftover
}

func (mp *Marketplace) CanBuy(as []artifacts.Artifacts) bool {
	return mp.Storage.HasAll(as)
}

func (mp *Marketplace) CanSell(as []artifacts.Artifacts) bool {
	return mp.Price(as) <= mp.Money
}

func (mp *Marketplace) HasAny(as []artifacts.Artifacts) bool {
	return mp.Storage.HasAny(as)
}

func (mp *Marketplace) Price(as []artifacts.Artifacts) uint32 {
	var price uint32 = 0
	for _, a := range as {
		if a.A != artifacts.Water {
			price += mp.Prices[a.A] * uint32(a.Quantity)
		}
	}
	return price
}

func (mp *Marketplace) Stats() *stats.HouseholdStats {
	return &stats.HouseholdStats{
		Money:     mp.Money,
		People:    0,
		Buildings: 1,
		Artifacts: mp.Storage.NumArtifacts(),
	}
}

func (mp *Marketplace) HasTraded(a *artifacts.Artifact) bool {
	_, ok := mp.Storage.Artifacts[a]
	return ok
}

func (mp *Marketplace) UnitPrice(a *artifacts.Artifact) uint32 {
	return mp.Prices[a]
}

func (mp *Marketplace) Spend(amount uint32) {
	mp.Money -= amount
}

func (mp *Marketplace) Earn(amount uint32) {
	mp.Money += amount
}

func (mp *Marketplace) GetMoney() uint32 {
	return mp.Money
}
