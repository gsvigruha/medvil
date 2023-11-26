package social

import (
	"log"
	"math/rand"
	"medvil/model/artifacts"
	"medvil/model/building"
	"medvil/model/economy"
	"medvil/model/navigation"
	"medvil/model/time"
	"medvil/model/vehicles"
	"os"
)

type Townhall struct {
	Household     *Household
	StorageTarget map[*artifacts.Artifact]int
	Traders       []*Trader
	Expeditions   []*Expedition
}

const StorageRefillBudgetPercentage = 0.5
const ConstructionStorageCapacity = 0.7

const TownhallMaxDistance = 25
const OutlawMaxDistance = 30

func (t *Townhall) ElapseTime(Calendar *time.CalendarType, m navigation.IMap) {
	t.Household.ElapseTime(Calendar, m)
	mp := t.Household.Town.Marketplace

	if mp != nil {
		resourcesFull := t.Household.Resources.Full()
		for _, a := range artifacts.All {
			tag := economy.SingleTag(economy.TagStorageTarget, a.Idx)
			transportQuantity := ProductTransportQuantity(a)
			goods := []artifacts.Artifacts{artifacts.Artifacts{A: a, Quantity: transportQuantity}}
			var q uint16 = 0
			if storageQ, ok := t.Household.Resources.Artifacts[a]; ok {
				q = storageQ
			}

			targetQ := uint16(t.StorageTarget[a])
			if q > targetQ {
				if t.Household.NumTasks("exchange", tag) == 0 {
					qToSell := t.Household.ArtifactToSell(a, q, false, false, resourcesFull)
					if qToSell > 0 {
						t.Household.AddTask(&economy.SellTask{
							Exchange: mp,
							Goods:    goods,
							TaskTag:  tag,
						})
					}
				}
			} else if q < targetQ {
				if uint16(t.Household.NumTasks("exchange", tag)) < (targetQ-q)/ProductTransportQuantity(a) {
					maxPrice := uint32(float64(t.Household.Money) * StorageRefillBudgetPercentage / float64(len(t.Household.Resources.Artifacts)))
					if t.Household.Money >= mp.Price(goods) && mp.HasTraded(a) {
						t.Household.AddTask(&economy.BuyTask{
							Exchange:        mp,
							HouseholdWallet: t.Household,
							Goods:           goods,
							MaxPrice:        maxPrice,
							TaskTag:         tag,
						})
					}
				}
			}
		}

		t.Household.MaybeBuyBoat(Calendar, m)
		t.Household.MaybeBuyCart(Calendar, m)
	}

	if t.Household.Town.Supplier != nil && t.Household.Town.Settings.UseSupplier && t.Household.Town.Supplier.CloseToTown(t.Household.Town, m) {
		src := t.Household.Town.Supplier
		dstH := t.Household
		if dstH.HasRoomForPeople() {
			src.ReassignFirstPerson(dstH, len(dstH.Tasks) == 0, m)
		}
		for _, a := range artifacts.All {
			var q uint16 = 0
			if storageQ, ok := t.Household.Resources.Artifacts[a]; ok {
				q = storageQ
			}
			pickupD := src.GetHome().Destination(building.NonExtension)
			if t.Household.NumTasks("transport", economy.TransportTaskTag(pickupD, a)) == 0 {
				targetQ := uint16(t.StorageTarget[a])
				if q < targetQ {
					if os.Getenv("MEDVIL_VERBOSE") == "2" {
						log.Printf("Created Transport Task")
					}
					t.Household.AddTask(&economy.TransportTask{
						PickupD:        pickupD,
						DropoffD:       m.GetField(dstH.Building.X, dstH.Building.Y),
						PickupR:        src.GetHome().GetResources(),
						DropoffR:       dstH.Resources,
						A:              a,
						TargetQuantity: ProductTransportQuantity(a),
					})
				}
			}
		}
	}

	for _, trader := range t.Traders {
		trader.ElapseTime(Calendar, m)
	}
	for _, expedition := range t.Expeditions {
		expedition.ElapseTime(Calendar, m)
	}
}

func (t *Townhall) GetFields() []navigation.FieldWithContext {
	fields := make([]navigation.FieldWithContext, len(t.Household.Town.Roads)+len(t.Household.Town.Walls))
	for i := range t.Household.Town.Roads {
		fields[i] = t.Household.Town.Roads[i]
	}
	for i := range t.Household.Town.Walls {
		fields[i+len(t.Household.Town.Roads)] = t.Household.Town.Walls[i]
	}
	return fields
}

func (t *Townhall) FieldWithinDistance(field *navigation.Field) bool {
	if t.Household.Building == nil {
		return true
	}
	return WithinDistance(t.Household.Building, field, TownhallMaxDistance)
}

func (t *Townhall) getTraderDestField(trader *Trader, m navigation.IMap) *navigation.Field {
	hf := t.Household.RandomField(m, trader.Vehicle.T.BuildingCheckFn)
	dest := &navigation.BuildingDestination{B: t.Household.Town.Marketplace.Building, ET: trader.Vehicle.T.BuildingExtensionType}
	if hf != nil {
		path := m.ShortPath(hf.GetLocation(), dest, trader.Person.Traveller.PathType())
		if path != nil && len(path.P) > 2 {
			pe := path.P[len(path.P)-2]
			return m.GetField(pe.GetLocation().X, pe.GetLocation().Y)
		}
	}
	return nil
}

func (t *Townhall) CreateExpedition(v *vehicles.Vehicle, p economy.Person) {
	for _, pI := range t.Household.People {
		person := p.(*Person)
		if pI == person {
			var r artifacts.Resources
			r.Init(uint32(v.T.MaxVolume))
			expedition := &Expedition{
				Name:            ExpeditionNames[rand.Intn(len(ExpeditionNames))],
				Money:           0,
				People:          []*Person{person},
				TargetNumPeople: 1,
				Vehicle:         v,
				Resources:       &r,
				Town:            t.Household.Town,
				StorageTarget:   make(map[*artifacts.Artifact]int),
			}
			for _, a := range artifacts.All {
				expedition.StorageTarget[a] = 0
			}
			t.Expeditions = append(t.Expeditions, expedition)
			person.Home = expedition
			person.SetHome()
			return
		}
	}
}

func (t *Townhall) CreateTrader(v *vehicles.Vehicle, p economy.Person) {
	for _, pI := range t.Household.People {
		person := p.(*Person)
		if pI == person {
			var r artifacts.Resources
			r.Init(uint32(v.T.MaxVolume))
			trader := &Trader{
				Money:          0,
				Person:         person,
				Vehicle:        v,
				Resources:      &r,
				SourceExchange: t.Household.Town.Marketplace,
				Town:           t.Household.Town,
			}
			t.Traders = append(t.Traders, trader)
			person.Home = trader
			person.Traveller.UseVehicle(v)
			person.SetHome()
			return
		}
	}
}

func (t *Townhall) Filter(Calendar *time.CalendarType, m navigation.IMap) {
	var newTraders = make([]*Trader, 0, len(t.Traders))
	for _, trader := range t.Traders {
		if trader.Person.Health == 0 {
			field := m.GetField(trader.Person.Traveller.FX, trader.Person.Traveller.FY)
			field.UnregisterTraveller(trader.Person.Traveller)
			field.UnregisterTraveller(trader.Vehicle.Traveller)
			t.Household.Money += trader.Money
		} else {
			newTraders = append(newTraders, trader)
		}
	}
	t.Traders = newTraders

	var newExpeditions = make([]*Expedition, 0, len(t.Expeditions))
	for _, expedition := range t.Expeditions {
		expedition.Filter(Calendar, m)
		if len(expedition.People) == 0 {
			f := m.GetField(expedition.Vehicle.Traveller.FX, expedition.Vehicle.Traveller.FY)
			f.UnregisterTraveller(expedition.Vehicle.Traveller)
			expedition.Resources.Deleted = true
			for _, otherTown := range t.Household.Town.Country.Towns {
				if otherTown.Supplier == expedition {
					otherTown.Supplier = nil
				}
			}
		} else {
			newExpeditions = append(newExpeditions, expedition)
		}
	}
	t.Expeditions = newExpeditions
}
