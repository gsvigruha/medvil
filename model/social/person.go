package social

import (
	"math/rand"
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
		if p.Traveller.FX == p.Task.Field().X && p.Traveller.FY == p.Task.Field().Y {
			// Work on task
			p.Traveller.ResetPhase()
			if p.Task.Complete(Calendar) {
				p.Task = nil
			}
		} else {
			p.Traveller.EnsurePath(p.Task.Field(), navigation.TravellerTypePedestrian, m)
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
		// no task
		home := m.GetField(p.Household.Building.X, p.Household.Building.Y)
		if p.Water < WaterThreshold && p.Household.HasDrink() {
			p.Task = &economy.DrinkTask{F: home, P: p}
		} else if p.Food < FoodThreshold && p.Household.HasFood() {
			p.Task = &economy.EatTask{F: home, P: p}
		} else if !p.IsHome {
			p.Task = &economy.GoHomeTask{F: home, P: p}
		}
	}
	p.Traveller.Visible = !(p.Traveller.FX == p.Household.Building.X && p.Traveller.FY == p.Household.Building.Y)
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
	if p.Household.HasFood() {
		available := economy.AvailableFood(p.Household.Resources)
		p.Household.Resources.Remove(available[rand.Intn(len(available))], 1)
		p.Food = MaxPersonState
	}
}

func (p *Person) Drink() {
	if p.Household.HasDrink() {
		available := economy.AvailableDrink(p.Household.Resources)
		p.Household.Resources.Remove(available[rand.Intn(len(available))], 1)
		p.Water = MaxPersonState
	}
}

func (p *Person) SetHome() {
	p.IsHome = true
}

func (p *Person) HasFood() bool {
	return p.Household.HasFood()
}

func (p *Person) HasDrink() bool {
	return p.Household.HasDrink()
}
