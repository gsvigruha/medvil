package social

import (
	"medvil/model/artifacts"
	"medvil/model/economy"
	"medvil/model/navigation"
	"medvil/model/vehicles"
	"medvil/util"
)

const TradeProfitThreshold = 2.0

type Trader struct {
	Money          uint32
	Vehicle        *vehicles.Vehicle
	Resources      artifacts.Resources
	SourceExchange *Marketplace
	TargetExchange *Marketplace
	Task           *economy.TradeTask
}

func (t *Trader) ElapseTime(Calendar *time.CalendarType, m navigation.IMap) {
	if t.Task == nil {
		trader.Task = trader.GetTradeTask(m)
	}
}

func (t *Trader) GetTradeTask(m navigation.IMap) *economy.TradeTask {
	var weights []float64
	var as []*artifacts.Artifact
	for _, a := range artifacts.All {
		if t.SourceExchange.HasTraded(a) && t.TargetExchange.HasTraded(a) {
			if t.SourceExchange.Prices[a]*TradeProfitThreshold <= t.TargetExchange.Prices[a] {
				profit := float64(t.TargetExchange.Prices[a]) / float64(t.SourceExchange.Prices[a])
				weights = append(weights, profit)
				as = append(as, a)
			}
		}
	}
	if len(weights) > 0 {
		artifactToTrade := as[util.RandomIndexWeighted(weights)]
		smx, smy, smok := GetRandomBuildingXY(t.SourceExchange.Building, m, navigation.Field.BuildingNonExtension)
		tmx, tmy, tmok := GetRandomBuildingXY(t.TargetExchange.Building, m, navigation.Field.BuildingNonExtension)
		if smok && tmok {
			quantity := uint16(t.Money / t.SourceExchange.Prices[artifactToTrade])
			return &economy.TradeTask{
				SourceMarketF:  m.GetField(smx, smy),
				TargetMarketF:  m.GetField(tmx, tmy),
				SourceExchange: t.SourceExchange,
				TargetExchange: t.TargetExchange,
				TraderR:        &t.Resources,
				TraderMoney:    &t.Money,
				Vehicle:        nil,
				GoodsToTrade:   []artifacts.Artifacts{artifacts.Artifacts{A: artifactToTrade, Quantity: quantity}},
				TaskTag:        "",
			}
		}
	}
	return nil
}
