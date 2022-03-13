package stats

type Stats struct {
	Money     uint32
	Artifacts uint32
	People    uint32
	Buildings uint32
}

func (s *Stats) Combine(os *Stats) *Stats {
	return &Stats{
		Money:     s.Money + os.Money,
		Artifacts: s.Artifacts + os.Artifacts,
		People:    s.People + os.People,
		Buildings: s.Buildings + os.Buildings,
	}
}

func (s *Stats) Add(os *Stats) {
	s.Money += os.Money
	s.Artifacts += os.Artifacts
	s.People += os.People
	s.Buildings += os.Buildings
}
