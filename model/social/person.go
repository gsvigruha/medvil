package social

import (
	"medvil/model/economy"
	"medvil/model/navigation"
	"medvil/model/time"
)

const MaxPersonState = 100

type Person struct {
	Hunger    uint8
	Thirst    uint8
	Happiness uint8
	Health    uint8
	Household *Household
	Task      economy.Task
	IsHome    bool
	Traveller *navigation.Traveller
}

func (p *Person) ElapseTime(Calendar *time.CalendarType, Map navigation.IMap) {
	if p.Task != nil {
		if p.Traveller.FX == p.Task.Location().X && p.Traveller.FY == p.Task.Location().Y {
			if p.Task.Complete(Calendar) {
				p.Task = nil
			}
		} else {
			if p.IsHome {
				p.IsHome = false
				b := p.Household.Building
				Map.GetField(b.X, b.Y).RegisterTraveller(p.Traveller)
			} else {
				p.Traveller.Move(p.Task.Location(), Map)
			}
		}
	}
}
