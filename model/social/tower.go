package social

import (
	"medvil/model/artifacts"
	"medvil/model/economy"
	"medvil/model/military"
	"medvil/model/navigation"
	"medvil/model/time"
)

const WeaponBudgetRatio = 0.5

type PatrolLand struct {
	X uint16
	Y uint16
	F *navigation.Field
}

func (l PatrolLand) Field() *navigation.Field {
	return l.F
}

func (l PatrolLand) Context() string {
	return "shield"
}

type Tower struct {
	Household Household
	Land      []PatrolLand
}

func (t *Tower) GetFields() []navigation.FieldWithContext {
	fields := make([]navigation.FieldWithContext, len(t.Land))
	for i := range t.Land {
		fields[i] = t.Land[i]
	}
	return fields
}

var Sword = artifacts.GetArtifact("sword")
var Shield = artifacts.GetArtifact("shield")

func (t *Tower) ElapseTime(Calendar *time.CalendarType, m navigation.IMap) {
	h := &t.Household
	h.ElapseTime(Calendar, m)
	mp := h.Town.Marketplace

	unarmedPeople := t.numUnarmedPeople()
	if unarmedPeople > 0 {
		tag := "weapon_buying"
		var weapons = h.Resources.Get(Sword)
		if weapons > h.Resources.Get(Shield) {
			weapons = h.Resources.Get(Shield)
		}

		if h.NumTasks("exchange", tag) == 0 && weapons == 0 {
			var quantity = (ProductTransportQuantity(Sword) + ProductTransportQuantity(Shield)) / 2
			if quantity > unarmedPeople {
				quantity = unarmedPeople
			}
			needs := []artifacts.Artifacts{
				artifacts.Artifacts{A: Sword, Quantity: quantity},
				artifacts.Artifacts{A: Shield, Quantity: quantity}}
			if h.Money >= mp.Price(needs) && mp.HasTraded(Sword) && mp.HasTraded(Shield) {
				h.AddTask(&economy.BuyTask{
					Exchange:       mp,
					HouseholdMoney: &h.Money,
					Goods:          needs,
					MaxPrice:       uint32(float64(h.Money) * WeaponBudgetRatio),
					TaskTag:        tag,
				})
			}
		}

		if weapons > 0 {
			for _, p := range h.People {
				if !p.Equipment.Weapon() && h.Resources.Remove(Sword, 1) > 0 && h.Resources.Remove(Shield, 1) > 0 {
					p.Equipment = &economy.Weapon{}
					weapons--
					if weapons == 0 {
						break
					}
				}
			}
		}
	}

	if Calendar.Hour == 0 && Calendar.Day == 1 {
		patrolFields := t.getPatrolFields()
		if h.NumTasks("patrol", "") == 0 && len(patrolFields) > 0 {
			h.AddTask(&military.PatrolTask{
				Fields: patrolFields,
			})
		}
	}
}

func (t *Tower) getPatrolFields() []*navigation.Field {
	var f []*navigation.Field
	for _, l := range t.Land {
		f = append(f, l.F)
	}
	return f
}

func (t *Tower) numUnarmedPeople() uint16 {
	var i = uint16(0)
	for _, p := range t.Household.People {
		if !p.Equipment.Weapon() {
			i++
		}
	}
	return i
}
