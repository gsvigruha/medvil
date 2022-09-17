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
const HealthThreshold = 100
const HappinessThreshold = 100

type Person struct {
	Food      uint8
	Water     uint8
	Happiness uint8
	Health    uint8
	Home      Home
	Task      economy.Task
	IsHome    bool
	Traveller *navigation.Traveller
	Equipment economy.Equipment
}

func (p *Person) releaseTask() {
	p.Task = nil
	if p.Equipment.Tool() {
		p.Equipment = &economy.NoEquipment{}
		p.Home.GetResources().Add(Tools, 1)
	}
	p.Traveller.ExitVehicle()
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
			if p.Task.Complete(Calendar, p.Equipment.Tool()) {
				p.releaseTask()
			}
		} else {
			if p.Traveller.EnsurePath(p.Task.Field(), m) {
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
					p.Home.AddTask(p.Task)
				}
				p.releaseTask()
			}
		}
	} else {
		// no task
		home := p.Home.Field(m)
		if p.Water < WaterThreshold && p.Home.HasDrink() {
			p.Task = &economy.DrinkTask{F: home, P: p}
		} else if p.Food < FoodThreshold && p.Home.HasFood() {
			p.Task = &economy.EatTask{F: home, P: p}
		} else if p.Health < HealthThreshold && p.Home.HasMedicine() {
			p.Task = &economy.HealTask{F: home, P: p}
		} else if p.Happiness < HappinessThreshold && p.Home.HasBeer() {
			p.Task = &economy.DrinkTask{F: home, P: p}
		} else if p.Task = p.Home.NextTask(m, p.Equipment); p.Task != nil {
			if !p.Equipment.Tool() && p.Home.GetResources().Remove(Tools, 1) == 1 {
				p.Equipment = &economy.Tool{}
			}
			p.Task.SetUp(p.Traveller, p.Home)
		} else if !p.IsHome {
			p.Task = &economy.GoHomeTask{F: home, P: p}
		}
	}
	p.Traveller.SetHome(p.Home.GetBuilding() != nil && p.Home.GetBuilding() == m.GetField(p.Traveller.FX, p.Traveller.FY).Building.GetBuilding())
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
		if p.Home.GetHeating() < rand.Float64() && p.Health > 0 {
			p.Health--
		}
		if !p.Home.HasEnoughTextile() && p.Happiness > 0 {
			p.Happiness--
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
	if p.Home.HasFood() {
		available := economy.AvailableFood(*p.Home.GetResources())
		a := available[rand.Intn(len(available))]
		p.Home.GetResources().Remove(a, 1)
		p.changePersonState(a)
	}
}

func (p *Person) Drink() {
	if p.Home.HasDrink() {
		available := economy.AvailableDrink(*p.Home.GetResources())
		a := available[rand.Intn(len(available))]
		p.Home.GetResources().Remove(a, 1)
		p.changePersonState(a)
	}
}

func (p *Person) Heal() {
	if p.Home.HasMedicine() {
		p.Home.GetResources().Remove(economy.Medicine, 1)
		p.changePersonState(economy.Medicine)
	}
}

func (p *Person) DrinkBeer() {
	if p.Home.HasBeer() {
		p.Home.GetResources().Remove(economy.Beer, 1)
		p.changePersonState(economy.Beer)
	}
}

func (p *Person) SetHome() {
	p.IsHome = true
}

func (p *Person) HasFood() bool {
	return p.Home.HasFood()
}

func (p *Person) HasDrink() bool {
	return p.Home.HasDrink()
}

func (p *Person) HasMedicine() bool {
	return p.Home.HasMedicine()
}

func (p *Person) HasBeer() bool {
	return p.Home.HasBeer()
}
