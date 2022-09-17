package social

import (
	"medvil/model/artifacts"
	"medvil/model/economy"
	"medvil/model/navigation"
	"medvil/model/time"
	"medvil/model/vehicles"
	"medvil/util"
)

const TradeProfitThreshold = 2.0
const TradingCapitalRatio = 0.5

type Trader struct {
	Money          uint32
	Vehicle        *vehicles.Vehicle
	Person         *Person
	Resources      artifacts.Resources
	SourceExchange *Marketplace
	TargetExchange *Marketplace
	Task           *economy.TradeTask
}

func (t *Trader) ElapseTime(Calendar *time.CalendarType, m navigation.IMap) {
	if t.Task == nil {
		t.Task = t.GetTradeTask(m)
	} else if t.Vehicle.Traveller.FX == t.Task.Field().X && t.Vehicle.Traveller.FY == t.Task.Field().Y {
		if t.Task.Complete(Calendar, false) {
			t.Task = nil
		}
	} else {
		if t.Vehicle.Traveller.EnsurePath(t.Task.Field(), m) {
			t.Vehicle.Traveller.Move(m)
		}
	}
}

func (t *Trader) GetArtifactToTrade(pickupMP, dropoffMP *Marketplace) *artifacts.Artifact {
	var weights []float64
	var tradableArtifacts []*artifacts.Artifact
	for _, a := range artifacts.All {
		if pickupMP.HasTraded(a) && dropoffMP.HasTraded(a) {
			if pickupMP.Prices[a]*TradeProfitThreshold <= dropoffMP.Prices[a] {
				profit := float64(dropoffMP.Prices[a]) / float64(pickupMP.Prices[a])
				weights = append(weights, profit)
				tradableArtifacts = append(tradableArtifacts, a)
			}
		}
	}
	if len(weights) > 0 {
		return tradableArtifacts[util.RandomIndexWeighted(weights)]
	}
	return nil
}

func (t *Trader) GetGoodsToTrade(a *artifacts.Artifact, mp *Marketplace) []artifacts.Artifacts {
	if a != nil {
		quantity := uint16(float64(t.Money) * TradingCapitalRatio / float64(mp.Prices[a]))
		return []artifacts.Artifacts{artifacts.Artifacts{A: a, Quantity: quantity}}
	}
	return []artifacts.Artifacts{}
}

func (t *Trader) GetTradeTask(m navigation.IMap) *economy.TradeTask {
	if t.TargetExchange == nil {
		return nil
	}
	artifactSourceToDest := t.GetArtifactToTrade(t.SourceExchange, t.TargetExchange)
	artifactDestToSource := t.GetArtifactToTrade(t.TargetExchange, t.SourceExchange)
	if artifactSourceToDest != nil || artifactDestToSource != nil {
		smx, smy, smok := GetRandomBuildingXY(t.SourceExchange.Building, m, navigation.Field.BuildingNonExtension)
		tmx, tmy, tmok := GetRandomBuildingXY(t.TargetExchange.Building, m, navigation.Field.BuildingNonExtension)
		if smok && tmok {
			return &economy.TradeTask{
				SourceMarketF:     m.GetField(smx, smy),
				TargetMarketF:     m.GetField(tmx, tmy),
				SourceExchange:    t.SourceExchange,
				TargetExchange:    t.TargetExchange,
				TraderR:           &t.Resources,
				TraderMoney:       &t.Money,
				Vehicle:           nil,
				GoodsSourceToDest: t.GetGoodsToTrade(artifactSourceToDest, t.SourceExchange),
				GoodsDestToSource: t.GetGoodsToTrade(artifactDestToSource, t.TargetExchange),
				TaskTag:           "",
			}
		}
	}
	return nil
}
