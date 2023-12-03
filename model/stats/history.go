package stats

import (
	"medvil/model/artifacts"
	"medvil/model/economy"
)

var MaxHistory = 2400

type HistoryElement struct {
	Stats
}

func (he HistoryElement) GetDeaths() uint32 {
	return he.Deaths
}

func (he HistoryElement) GetDepartures() uint32 {
	return he.Departures
}

func (he HistoryElement) GetPoverty() uint32 {
	return he.Poverty
}

func (he HistoryElement) GetPeople() uint32 {
	return he.Global.People
}

func (he HistoryElement) GetArtifacts() uint32 {
	return he.Global.Artifacts
}

func (he HistoryElement) GetExchangedQuantity() uint32 {
	var quantity uint32 = 0
	for _, q := range he.TradeQ {
		quantity += q
	}
	return quantity
}

func (he HistoryElement) GetExchangedPrice() uint32 {
	var moneyAmount uint32 = 0
	for _, q := range he.TradeM {
		moneyAmount += q
	}
	return moneyAmount
}

func (he HistoryElement) computePrice(as []*artifacts.Artifact) uint32 {
	var quantity uint32 = 0
	var price uint32 = 0
	for _, a := range as {
		quantity += he.TradeQ[a]
		price += he.TradeM[a]
	}
	if quantity > 0 {
		return price / quantity
	}
	return 0
}

func (he HistoryElement) GetFoodPrice() uint32 {
	return he.computePrice(economy.FoodArtifacts)
}

func (he HistoryElement) GetHouseholdItemPrices() uint32 {
	return he.computePrice(economy.HouseholdItems)
}

func (he HistoryElement) GetBuildingMaterialsPrice() uint32 {
	return he.computePrice(economy.BuildingMaterials)
}

func (he HistoryElement) GetTransportTaskTime() uint32 {
	return he.CompletedT["TransportTask"]
}

func (he HistoryElement) GetExchangeTaskTime() uint32 {
	return he.CompletedT["ExchangeTask"]
}

func (he HistoryElement) GetAgricultureTaskTime() uint32 {
	return he.CompletedT["AgriculturalTask"]
}

func (he HistoryElement) GetManufactureTaskTime() uint32 {
	return he.CompletedT["ManufactureTask"]
}

type History struct {
	Elements []HistoryElement
}

func (h *History) Archive(stats *Stats) {
	stats.PendingT = nil
	h.Elements = append(h.Elements, HistoryElement{*stats})
	if len(h.Elements) > MaxHistory { // buffer
		h.Elements = h.Elements[len(h.Elements)-MaxHistory:]
	}
}

func (he HistoryElement) GetFarmMoney() uint32 {
	return he.Farm.Money
}

func (he HistoryElement) GetWorkshopMoney() uint32 {
	return he.Workshop.Money
}

func (he HistoryElement) GetMineMoney() uint32 {
	return he.Mine.Money
}

func (he HistoryElement) GetGovernmentMoney() uint32 {
	return he.Gov.Money
}

func (he HistoryElement) GetTraderMoney() uint32 {
	return he.Trader.Money
}

func (he HistoryElement) GetFarmPeople() uint32 {
	return he.Farm.People
}

func (he HistoryElement) GetWorkshopPeople() uint32 {
	return he.Workshop.People
}

func (he HistoryElement) GetMinePeople() uint32 {
	return he.Mine.People
}

func (he HistoryElement) GetGovernmentPeople() uint32 {
	return he.Gov.People
}

func (he HistoryElement) GetTraderPeople() uint32 {
	return he.Trader.People
}
