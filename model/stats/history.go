package stats

type HistoryElement struct {
	Stats
	SocietyStats
}

func (he HistoryElement) GetDeaths() uint32 {
	return he.Deaths
}

func (he HistoryElement) GetDepartures() uint32 {
	return he.Departures
}

func (he HistoryElement) GetPeople() uint32 {
	return he.People
}

func (he HistoryElement) GetArtifacts() uint32 {
	return he.Artifacts
}

func (he HistoryElement) GetExchangedNum() uint32 {
	return he.ExchangedNum
}

func (he HistoryElement) GetExchangedPrice() uint32 {
	return he.ExchangedPrice
}

func (he HistoryElement) GetFoodPrice() uint32 {
	return he.FoodPrice
}

func (he HistoryElement) GetHouseholdItemPrices() uint32 {
	return he.HouseholdItemPrice
}

func (he HistoryElement) GetBuildingMaterialsPrice() uint32 {
	return he.BuildingMaterialsPrice
}

type History struct {
	Elements []HistoryElement
}

func (h *History) Archive(stats *Stats, societyStats *SocietyStats) {
	h.Elements = append(h.Elements, HistoryElement{*stats, *societyStats})
}
