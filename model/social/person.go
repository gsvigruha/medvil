package social

import (
	//"fmt"
	"medvil/model/economy"
	"medvil/model/navigation"
	"medvil/model/time"
)

const MaxPersonState = 250
const WaterThreshold = 100
const FoodThreshold = 100

type Person struct {
	Food      uint8
	Water     uint8
	Happiness uint8
	Health    uint8
	Household *Household
	Task      economy.Task
	IsHome    bool
	Traveller *navigation.Traveller
}

func (p *Person) ElapseTime(Calendar *time.CalendarType, m navigation.IMap) {
	if p.Task != nil {
		if p.Traveller.FX == p.Task.Location().X && p.Traveller.FY == p.Task.Location().Y {
			// Work on task
			p.Traveller.ResetPhase()
			if p.Task.Complete(Calendar) {
				p.Task = nil
			}
		} else {
			if p.Traveller.Path == nil || p.Traveller.Path.LastLocation() != p.Task.Location() {
				p.Traveller.Path = m.ShortPath(p.Traveller.FX, p.Traveller.FY, p.Task.Location().X, p.Task.Location().Y, navigation.TravellerTypePedestrian)
			}
			if p.IsHome {
				// Start on path
				p.IsHome = false
				p.Traveller.ResetPhase()
			} else {
				// Move on path
				p.Traveller.Move(m)
			}
		}
	} else {
		home := navigation.Location{X: p.Household.Building.X, Y: p.Household.Building.Y, F: m.GetField(p.Household.Building.X, p.Household.Building.Y)}
		if p.Water < WaterThreshold {
			p.Task = &economy.DrinkTask{L: home, P: p}
		} else if p.Food < FoodThreshold {
			p.Task = &economy.EatTask{L: home, P: p}
		} else if !p.IsHome {
			p.Task = &economy.GoHomeTask{L: home, P: p}
		}
	}
	if Calendar.Hour == 0 {
		if p.Food > 0 {
			p.Food--
		}
		if p.Water > 0 {
			p.Water--
		}
	}
}

func (p *Person) Eat() {
	p.Food = MaxPersonState
}

func (p *Person) Drink() {
	p.Water = MaxPersonState
}

func (p *Person) SetHome() {
	p.IsHome = true
}
