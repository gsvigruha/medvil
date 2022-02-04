package social

import (
	"medvil/model/artifacts"
	"medvil/model/building"
	"medvil/model/navigation"
	"medvil/model/time"
)

type Marketplace struct {
	Building *building.Building
	Money    uint32
	Storage  artifacts.Resources
	Price    map[*artifacts.Artifact]uint16
	Supply   map[*artifacts.Artifact]uint16
	Demand   map[*artifacts.Artifact]uint16
}

func (mp *Marketplace) Init() {
	for _, a := range artifacts.All {
		mp.Price[a] = 10
		mp.Supply[a] = 0
		mp.Demand[a] = 0
	}
}

func (mp *Marketplace) ElapseTime(Calendar *time.CalendarType, m navigation.IMap) {
	if Calendar.Hour == 0 {
		for _, a := range artifacts.All {
			r := float64(mp.Supply[a]) / float64(mp.Demand[a])
			if r >= 1.1 {
				mp.Price[a]--
			} else if r <= 0.9 {
				mp.Price[a]++
			}
		}
	}

}
