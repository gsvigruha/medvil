package stats

type SocietyStats struct {
	Deaths                 uint32
	Departures             uint32
	ProducedNum            uint32
	ExchangedNum           uint32
	ProducedPrice          uint32
	ExchangedPrice         uint32
	FoodPrice              uint32
	HouseholdItemPrice     uint32
	BuildingMaterialsPrice uint32
}

func (s *SocietyStats) RegisterTrade(price uint32, quantity uint16) {
	s.ExchangedNum += uint32(quantity)
	s.ExchangedPrice += uint32(price)
}

func (s *SocietyStats) RegisterDeath() {
	s.Deaths++
}

func (s *SocietyStats) RegisterDeparture() {
	s.Departures++
}

func (s *SocietyStats) RegisterFoodPrices(items []uint32) {
	for _, item := range items {
		s.FoodPrice += item
	}
	s.FoodPrice /= uint32(len(items))
}

func (s *SocietyStats) RegisterBuildingMaterialsPrices(items []uint32) {
	for _, item := range items {
		s.BuildingMaterialsPrice += item
	}
	s.BuildingMaterialsPrice /= uint32(len(items))
}

func (s *SocietyStats) RegisterHouseholdItemPrices(items []uint32) {
	for _, item := range items {
		s.HouseholdItemPrice += item
	}
	s.HouseholdItemPrice /= uint32(len(items))
}
