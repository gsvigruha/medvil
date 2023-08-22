package social

import (
	"medvil/model/stats"
)

const CountryTypePlayer uint8 = 0
const CountryTypeOutlaw uint8 = 1
const CountryTypeOtherCivilization uint8 = 2

type Country struct {
	Towns        []*Town
	T            uint8
	History      *stats.History
	SocietyStats *stats.SocietyStats
}

func (c *Country) Stats() *stats.Stats {
	s := &stats.Stats{}
	for _, town := range c.Towns {
		s.Add(town.Stats)
	}
	return s
}

func (c *Country) AddTownIfDoesNotExist(town *Town) {
	for _, t := range c.Towns {
		if t == town {
			return
		}
	}
	c.Towns = append(c.Towns, town)
}

func (c *Country) ArchiveHistory() {
	c.History.Archive(c.Stats(), c.SocietyStats)
	c.SocietyStats = &stats.SocietyStats{}
}
