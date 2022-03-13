package social

import (
	"medvil/model/artifacts"
	"medvil/model/building"
	"medvil/model/navigation"
	"medvil/model/stats"
	"medvil/model/time"
)

type Marketplace struct {
	Town     *Town
	Building *building.Building
	Money    uint32
	Storage  artifacts.Resources
	Prices   map[*artifacts.Artifact]uint32
	Supply   map[*artifacts.Artifact]uint32
	Demand   map[*artifacts.Artifact]uint32
}

func (mp *Marketplace) Init() {
	mp.Prices = make(map[*artifacts.Artifact]uint32)
	mp.Supply = make(map[*artifacts.Artifact]uint32)
	mp.Demand = make(map[*artifacts.Artifact]uint32)
	for _, a := range artifacts.All {
		mp.Prices[a] = 10
		mp.Supply[a] = 0
		mp.Demand[a] = 0
	}
}

var gold = artifacts.GetArtifact("gold_coin")

func (mp *Marketplace) ElapseTime(Calendar *time.CalendarType, m navigation.IMap) {
	if Calendar.Hour == 0 && Calendar.Day == 1 {
		allGold := []artifacts.Artifacts{artifacts.Artifacts{A: gold, Quantity: mp.Storage.Get(gold)}}
		price := mp.Price(allGold)
		wallet := &mp.Town.Townhall.Household.Money
		*wallet += price * 2
		mp.Buy(allGold, wallet)

		for _, a := range artifacts.All {
			if mp.Supply[a] == 0 && mp.Demand[a] > 0 && mp.Demand[a] > uint32(mp.Storage.Get(a)) {
				mp.Prices[a]++
				mp.Supply[a] = 0
				mp.Demand[a] = 0
			} else if mp.Demand[a] == 0 && mp.Supply[a] > 0 && mp.Storage.Get(a) > 0 {
				if mp.Prices[a] > 1 {
					mp.Prices[a]--
					mp.Supply[a] = 0
					mp.Demand[a] = 0
				}
			} else if mp.Demand[a] > 0 && mp.Supply[a] > 0 {
				r := float64(mp.Supply[a]) / float64(mp.Demand[a])
				if r >= 1.1 && mp.Prices[a] > 1 {
					mp.Prices[a]--
					mp.Supply[a] = 0
					mp.Demand[a] = 0
				} else if r <= 0.9 {
					mp.Prices[a]++
					mp.Supply[a] = 0
					mp.Demand[a] = 0
				}
			}
		}
	}
}

func (mp *Marketplace) Buy(as []artifacts.Artifacts, wallet *uint32) {
	price := mp.Price(as)
	mp.Storage.RemoveAll(as)
	mp.Money += price
	*wallet -= price
	for _, a := range as {
		mp.Demand[a.A] += uint32(a.Quantity)
	}
}

func (mp *Marketplace) BuyAsManyAsPossible(as []artifacts.Artifacts, wallet *uint32) []artifacts.Artifacts {
	existingArtifacts := mp.Storage.GetAsManyAsPossible(as)
	price := mp.Price(existingArtifacts)
	mp.Money += price
	*wallet -= price
	for _, a := range as {
		mp.Demand[a.A] += uint32(a.Quantity)
	}
	return existingArtifacts
}

func (mp *Marketplace) Sell(as []artifacts.Artifacts, wallet *uint32) {
	price := mp.Price(as)
	mp.Storage.AddAll(as)
	mp.Money -= price
	*wallet += price
	for _, a := range as {
		mp.Supply[a.A] += uint32(a.Quantity)
	}
}

func (mp *Marketplace) CanBuy(as []artifacts.Artifacts) bool {
	return mp.Storage.Has(as)
}

func (mp *Marketplace) CanSell(as []artifacts.Artifacts) bool {
	return mp.Price(as) <= mp.Money
}

func (mp *Marketplace) Price(as []artifacts.Artifacts) uint32 {
	var price uint32 = 0
	for _, a := range as {
		price += mp.Prices[a.A] * uint32(a.Quantity)
	}
	return price
}

func (mp *Marketplace) Stats() *stats.Stats {
	return &stats.Stats{
		Money:     mp.Money,
		People:    0,
		Buildings: 1,
		Artifacts: mp.Storage.NumArtifacts(),
	}
}
