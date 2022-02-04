package social

import (
	"medvil/model/artifacts"
	"medvil/model/building"
	"medvil/model/navigation"
	"medvil/model/time"
)

type Marketplace struct {
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

func (mp *Marketplace) ElapseTime(Calendar *time.CalendarType, m navigation.IMap) {
	if Calendar.Hour == 0 {
		for _, a := range artifacts.All {
			r := float64(mp.Supply[a]) / float64(mp.Demand[a])
			if r >= 1.1 {
				mp.Prices[a]--
			} else if r <= 0.9 {
				mp.Prices[a]++
			}
		}
	}
}

func (mp *Marketplace) Buy(as []artifacts.Artifacts, wallet *uint32) {
	price := mp.Price(as)
	mp.Storage.RemoveAll(as)
	mp.Money += price
	*wallet -= price
}

func (mp *Marketplace) Sell(as []artifacts.Artifacts, wallet *uint32) {
	price := mp.Price(as)
	mp.Storage.AddAll(as)
	mp.Money -= price
	*wallet += price
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
