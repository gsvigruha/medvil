package social

import (
	"medvil/model/artifacts"
	"medvil/model/economy"
	"medvil/model/navigation"
	"medvil/model/vehicles"
)

const TradeProfitThreshold = 2.0

type Trader struct {
	Money          uint32
	Vehicle        *vehicles.Vehicle
	Resources      artifacts.Resources
	SourceExchange *Marketplace
	TargetExchange *Marketplace
}

func (t *Trader) GetProfitableTrades(m navigation.IMap) *economy.TradeTask {
	var maxProfit = 0.0
	var bestArtifact *artifacts.Artifact
	for _, a := range artifacts.All {
		if t.SourceExchange.HasTraded(a) && t.TargetExchange.HasTraded(a) {
			if t.SourceExchange.Prices[a]*TradeProfitThreshold <= t.TargetExchange.Prices[a] {
				profit := float64(t.TargetExchange.Prices[a]) / float64(t.SourceExchange.Prices[a])
				if profit > maxProfit {
					bestArtifact = a
					maxProfit = profit
				}
			}
		}
	}
	if bestArtifact != nil {
		smx, smy, smok := GetRandomBuildingXY(t.SourceExchange.Building, m, navigation.Field.BuildingNonExtension)
		tmx, tmy, tmok := GetRandomBuildingXY(t.TargetExchange.Building, m, navigation.Field.BuildingNonExtension)
		if smok && tmok {
			quantity := uint16(t.Money / t.SourceExchange.Prices[bestArtifact])
			return &economy.TradeTask{
				SourceMarketF:  m.GetField(smx, smy),
				TargetMarketF:  m.GetField(tmx, tmy),
				SourceExchange: t.SourceExchange,
				TargetExchange: t.TargetExchange,
				TraderR:        &t.Resources,
				TraderMoney:    &t.Money,
				Vehicle:        nil,
				GoodsToTrade:   []artifacts.Artifacts{artifacts.Artifacts{A: bestArtifact, Quantity: quantity}},
				TaskTag:        "",
			}
		}
	}
	return nil
}
