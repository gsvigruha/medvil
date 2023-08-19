package stats

type SocietyStats struct {
	Deaths         uint32
	Departures     uint32
	ProducedNum    uint32
	ExchangedNum   uint32
	ProducedPrice  uint32
	ExchangedPrice uint32
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
