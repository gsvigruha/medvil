package social

import (
	"medvil/model/artifacts"
)

const TradeProfitThreshold = 2.0

func (t *Townhall) GetProfitableTrades(otherMP *Marketplace) {
	mp := t.Household.Town.Marketplace
	for _, a := range artifacts.All {
		if mp.HasTraded(a) && otherMP.HasTraded(a) {
			if mp.Prices[a] >= otherMP.Prices[a]*TradeProfitThreshold {

			} else if mp.Prices[a]*TradeProfitThreshold <= otherMP.Prices[a] {

			}
		}
	}
}
