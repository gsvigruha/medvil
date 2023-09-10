package social

import (
	"medvil/model/artifacts"
	"medvil/model/building"
	"medvil/model/economy"
	"medvil/model/stats"
)

const CountryTypePlayer uint8 = 0
const CountryTypeOutlaw uint8 = 1
const CountryTypeOtherCivilization uint8 = 2

type Country struct {
	Towns   []*Town
	T       uint8
	History *stats.History
}

func (c *Country) Stats() *stats.Stats {
	s := &stats.Stats{}
	s.Init(make(map[economy.Task]uint32))
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
	c.History.Archive(c.Stats())
	for _, town := range c.Towns {
		town.ArchiveHistory()
	}
}

func (c *Country) CreateNewTown(b *building.Building, supplier *Town) {
	newTown := &Town{Country: c, Supplier: supplier, Settings: DefaultTownSettings}
	newTown.Townhall = &Townhall{Household: &Household{Building: b, Town: newTown, Resources: &artifacts.Resources{}}}
	newTown.Townhall.Household.Resources.VolumeCapacity = b.Plan.Area() * StoragePerArea
	newTown.Init()
	newTown.Townhall.Household.TargetNumPeople = newTown.Townhall.Household.Building.Plan.Area()
	for a, q := range DefaultStorageTarget {
		newTown.Townhall.StorageTarget[artifacts.GetArtifact(a)] = q
	}
	c.Towns = append(c.Towns, newTown)
}
