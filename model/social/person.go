package social

import (
	"math/rand"
	"medvil/model/artifacts"
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
	Tool      bool
}

func (p *Person) releaseTask() {
	p.Task = nil
	if p.Tool {
		p.Tool = false
		p.Household.Resources.Add(Tools, 1)
	}
}

func (p *Person) ElapseTime(Calendar *time.CalendarType, m navigation.IMap) {
	if p.Task != nil {
		if p.Traveller.FX == p.Task.Field().X && p.Traveller.FY == p.Task.Field().Y {
			// Work on task
			if p.Task.Motion() == navigation.MotionStand {
				p.Traveller.ResetPhase()
			} else {
				p.Traveller.Motion = p.Task.Motion()
				p.Traveller.IncPhase()
			}
			if p.Task.Complete(Calendar, p.Tool) {
				p.releaseTask()
			}
		} else {
			if p.Traveller.EnsurePath(p.Task.Field(), navigation.TravellerTypePedestrian, m) {
				if p.IsHome {
					// Start on path
					p.IsHome = false
					p.Traveller.ResetPhase()
				} else {
					// Move on path
					p.Traveller.Move(m)
				}
			} else {
				if !economy.IsPersonalTask(p.Task.Name()) {
					p.Task.Pause(true)
					p.Household.AddTask(p.Task)
				}
				p.releaseTask()
			}
		}
	} else {
		// no task
		home := m.GetField(p.Household.Building.X, p.Household.Building.Y)
		if p.Water < WaterThreshold && p.Household.HasDrink() {
			p.Task = &economy.DrinkTask{F: home, P: p}
		} else if p.Food < FoodThreshold && p.Household.HasFood() {
			p.Task = &economy.EatTask{F: home, P: p}
		} else if p.Household.HasTask() {
			p.Task = p.Household.getNextTaskCombineExchange(m)
			if !p.Tool && p.Household.Resources.Remove(Tools, 1) == 1 {
				p.Tool = true
			}
		} else if !p.IsHome {
			p.Task = &economy.GoHomeTask{F: home, P: p}
		}
	}
	p.Traveller.Visible = !(p.Household.Building == m.GetField(p.Traveller.FX, p.Traveller.FY).Building.GetBuilding())
	if Calendar.Hour == 0 {
		if p.Food > 0 {
			p.Food--
		}
		if p.Water > 0 {
			p.Water--
		}
		if p.Food == 0 && p.Health > 0 {
			p.Health--
		}
		if p.Food == 0 && p.Happiness > 0 {
			p.Happiness--
		}
		if p.Water == 0 && p.Health > 0 {
			p.Health--
		}
		if p.Household.Heating < rand.Float64() && p.Health > 0 {
			p.Health--
		}
	}
}

func (p *Person) changePersonState(a *artifacts.Artifact) {
	c := economy.ArtifactToPersonState[a]
	if int(p.Food)+int(c.Food) > MaxPersonState {
		p.Food = MaxPersonState
	} else {
		p.Food += c.Food
	}
	if int(p.Water)+int(c.Water) > MaxPersonState {
		p.Water = MaxPersonState
	} else {
		p.Water += c.Water
	}
	if int(p.Health)+int(c.Health) > MaxPersonState {
		p.Health = MaxPersonState
	} else {
		p.Health += c.Health
	}
	if int(p.Happiness)+int(c.Happiness) > MaxPersonState {
		p.Happiness = MaxPersonState
	} else {
		p.Happiness += c.Happiness
	}
}

func (p *Person) Eat() {
	if p.Household.HasFood() {
		available := economy.AvailableFood(p.Household.Resources)
		a := available[rand.Intn(len(available))]
		p.Household.Resources.Remove(a, 1)
		p.changePersonState(a)
	}
}

func (p *Person) Drink() {
	if p.Household.HasDrink() {
		available := economy.AvailableDrink(p.Household.Resources)
		a := available[rand.Intn(len(available))]
		p.Household.Resources.Remove(a, 1)
		p.changePersonState(a)
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
