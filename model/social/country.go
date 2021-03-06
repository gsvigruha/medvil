package social

import (
	"medvil/model/stats"
)

type Country struct {
	Towns []*Town
}

func (c *Country) Stats() *stats.Stats {
	s := &stats.Stats{}
	for _, town := range c.Towns {
		s.Add(town.Stats)
	}
	return s
}

func (c *Country) AddTown(town *Town) {
	c.Towns = append(c.Towns, town)
}
