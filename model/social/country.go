package social

import (
	"math/rand"
	"medvil/model/artifacts"
	"medvil/model/building"
	"medvil/model/economy"
	"medvil/model/stats"
	"medvil/model/time"
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

func (c *Country) ArchiveHistory(Calendar *time.CalendarType) {
	c.History.Archive(c.Stats())
	for _, town := range c.Towns {
		town.ArchiveHistory(Calendar)
	}
}

func (c *Country) PickTownName() string {
	if c.T == CountryTypeOutlaw {
		return OutlawColonyNames[rand.Intn(len(OutlawColonyNames))]
	} else {
		return TownNames[rand.Intn(len(TownNames))]
	}
}

func (c *Country) CreateNewTown(Calendar *time.CalendarType, b *building.Building, supplier Supplier) {
	name := c.PickTownName()
	newTown := &Town{Country: c, Supplier: supplier, Settings: DefaultTownSettings, Name: name}
	newTown.Townhall = &Townhall{Household: &Household{Building: b, Town: newTown, Resources: &artifacts.Resources{}, BoatEnabled: true}}
	newTown.Townhall.Household.Resources.VolumeCapacity = uint32(b.Plan.Area()) * StoragePerArea
	newTown.Init(Calendar, len(c.History.Elements))
	newTown.Townhall.Household.TargetNumPeople = newTown.Townhall.Household.Building.Plan.Area()
	for a, q := range DefaultStorageTarget {
		newTown.Townhall.StorageTarget[artifacts.GetArtifact(a)] = q
	}
	c.Towns = append(c.Towns, newTown)
}
