package social

import (
	"medvil/model/artifacts"
	"medvil/model/building"
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
	Person         *Person
	Resources      artifacts.Resources
	SourceExchange *Marketplace
	TargetExchange *Marketplace
}

func (t *Trader) ElapseTime(Calendar *time.CalendarType, m navigation.IMap) {
	t.Person.ElapseTime(Calendar, m)
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

func (t *Trader) AddTask(economy.Task) {

}

func (t *Trader) HasFood() bool {
	return economy.HasFood(t.Resources)
}

func (t *Trader) HasDrink() bool {
	return economy.HasDrink(t.Resources)
}

func (t *Trader) HasMedicine() bool {
	return economy.HasMedicine(t.Resources)
}

func (t *Trader) HasBeer() bool {
	return economy.HasBeer(t.Resources)
}

func (t *Trader) Field(m navigation.IMap) *navigation.Field {
	return m.GetField(t.Person.Traveller.FX, t.Person.Traveller.FY)
}

func (t *Trader) NextTask(m navigation.IMap, e economy.Equipment) economy.Task {
	return t.GetTradeTask(m)
}

func (t *Trader) GetResources() *artifacts.Resources {
	return &t.Resources
}

func (t *Trader) GetBuilding() *building.Building {
	return nil
}

func (t *Trader) GetHeating() float64 {
	return 0
}

func (t *Trader) HasEnoughTextile() bool {
	return false
}

func (t *Trader) AddVehicle(v *vehicles.Vehicle) {
}

func (t *Trader) GetVehicle() *vehicles.Vehicle {
	return nil
}
