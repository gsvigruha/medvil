package social

import (
	"medvil/model/artifacts"
	"medvil/model/building"
	"medvil/model/economy"
	"medvil/model/navigation"
	"medvil/model/stats"
	"medvil/model/time"
	"medvil/model/vehicles"
	"medvil/util"
)

const TradeProfitThreshold = 2.0
const TradingCapitalRatio = 0.5

type Trader struct {
	Money          uint32
	Person         *Person
	Vehicle        *vehicles.Vehicle
	Resources      *artifacts.Resources
	SourceExchange *Marketplace
	TargetExchange *Marketplace
	Town           *Town
	Tasks          []economy.Task
}

func (t *Trader) ElapseTime(Calendar *time.CalendarType, m navigation.IMap) {
	t.Person.ElapseTime(Calendar, m)
	t.Person.Traveller.UseVehicle(t.Vehicle)
	if Calendar.Hour == 0 && Calendar.Day == 15 {
		CombineExchangeTasks(t, t.SourceExchange, m)
	}
	if t.Person.IsHome {
		FindWaterTask(t, 1, m)
		GetFoodTasks(t, 1, t.SourceExchange)
		if t.NumTasks("trade", "trade") == 0 {
			task := t.GetTradeTask(m)
			if task != nil {
				t.AddTask(task)
			}
		}
	}
	if Calendar.Hour == 0 && Calendar.Day == 1 && t.Person.Task != nil {
		if tt, ok := t.Person.Task.(*economy.TradeTask); ok {
			tt.Pause(true)
			t.AddPriorityTask(tt)
			t.Person.Task = nil
		}
	}
	if !economy.HasPersonalTask(t.Tasks) {
		for i := 0; i < len(t.Tasks); i++ {
			if t.Tasks[i].IsPaused() {
				t.Tasks[i].Pause(false)
			}
		}
	}
}

func (t *Trader) GetArtifactToTrade(pickupMP, dropoffMP *Marketplace) *artifacts.Artifact {
	var weights []float64
	var tradableArtifacts []*artifacts.Artifact
	for _, a := range artifacts.All {
		if pickupMP.HasTraded(a) && dropoffMP.HasTraded(a) && pickupMP.Storage.Get(a) > ProductTransportQuantity(a) && dropoffMP.Prices[a] > 2 {
			if pickupMP.Prices[a]*TradeProfitThreshold <= dropoffMP.Prices[a] {
				profit := float64(dropoffMP.Prices[a]) / float64(pickupMP.Prices[a])
				weights = append(weights, profit)
				tradableArtifacts = append(tradableArtifacts, a)
			}
		}
	}
	if len(weights) > 0 {
		return tradableArtifacts[util.RandomIndexWeighted(weights)]
	}
	return nil
}

func (t *Trader) GetGoodsToTrade(a *artifacts.Artifact, src *Marketplace, dst *Marketplace) []artifacts.Artifacts {
	if a != nil {
		buyQuantity := uint16(float64(t.Money) * TradingCapitalRatio / float64(src.Prices[a]))
		sellQuantity := uint16(float64(dst.Money) / float64(dst.Prices[a]))
		var quantity = buyQuantity
		if sellQuantity < buyQuantity {
			quantity = sellQuantity
		}
		capacity := t.Vehicle.T.MaxVolume / a.V
		if capacity < quantity {
			quantity = capacity
		}
		return []artifacts.Artifacts{artifacts.Artifacts{A: a, Quantity: quantity}}
	}
	return []artifacts.Artifacts{}
}

func (t *Trader) GetTradeTask(m navigation.IMap) economy.Task {
	if t.TargetExchange == nil || !t.SourceExchange.Town.Settings.Trading || !t.TargetExchange.Town.Settings.Trading {
		return nil
	}
	artifactSourceToDest := t.GetArtifactToTrade(t.SourceExchange, t.TargetExchange)
	artifactDestToSource := t.GetArtifactToTrade(t.TargetExchange, t.SourceExchange)
	if artifactSourceToDest != nil || artifactDestToSource != nil {
		goodsSourceToDest := t.GetGoodsToTrade(artifactSourceToDest, t.SourceExchange, t.TargetExchange)
		goodsDestToSource := t.GetGoodsToTrade(artifactDestToSource, t.TargetExchange, t.SourceExchange)
		if len(goodsSourceToDest) > 0 || len(goodsDestToSource) > 0 {
			return &economy.TradeTask{
				SourceMarketD:     &navigation.BuildingDestination{B: t.SourceExchange.Building, ET: t.Vehicle.T.BuildingExtensionType},
				TargetMarketD:     &navigation.BuildingDestination{B: t.TargetExchange.Building, ET: t.Vehicle.T.BuildingExtensionType},
				SourceExchange:    t.SourceExchange,
				TargetExchange:    t.TargetExchange,
				TraderR:           t.Resources,
				TraderWallet:      t,
				Vehicle:           nil,
				GoodsSourceToDest: goodsSourceToDest,
				GoodsDestToSource: goodsDestToSource,
				TaskTag:           "trade",
			}
		}
	}
	return nil
}

func (t *Trader) AddTask(task economy.Task) {
	t.Tasks = append(t.Tasks, task)
}

func (t *Trader) AddPriorityTask(task economy.Task) {
	t.Tasks = append([]economy.Task{task}, t.Tasks...)
}

func (t *Trader) GetTasks() []economy.Task {
	return t.Tasks
}

func (t *Trader) SetTasks(tasks []economy.Task) {
	t.Tasks = tasks
}

func (t *Trader) HasFood() bool {
	return economy.HasFood(*t.Resources)
}

func (t *Trader) HasDrink() bool {
	return economy.HasDrink(*t.Resources)
}

func (t *Trader) HasMedicine() bool {
	return economy.HasMedicine(*t.Resources)
}

func (t *Trader) HasBeer() bool {
	return economy.HasBeer(*t.Resources)
}

func (t *Trader) Field(m navigation.IMap) *navigation.Field {
	return m.GetField(t.Vehicle.Traveller.FX, t.Vehicle.Traveller.FY)
}

func (t *Trader) RandomField(m navigation.IMap, check func(navigation.Field) bool) *navigation.Field {
	return t.Field(m)
}

func (t *Trader) NextTask(m navigation.IMap, e *economy.EquipmentType) economy.Task {
	return GetNextTask(t, e)
}

func (t *Trader) GetResources() *artifacts.Resources {
	return t.Resources
}

func (t *Trader) GetBuilding() *building.Building {
	return nil
}

func (t *Trader) GetHeating() uint8 {
	return 100
}

func (t *Trader) HasEnoughClothes() bool {
	return true
}

func (t *Trader) AddVehicle(v *vehicles.Vehicle) {
}

func (t *Trader) AllocateVehicle(waterOk bool) *vehicles.Vehicle {
	return t.Vehicle
}

func (t *Trader) NumTasks(name string, tag string) int {
	var i = 0
	for _, t := range t.Tasks {
		i += CountTags(t, name, tag)
	}
	if t.Person.Task != nil {
		i += CountTags(t.Person.Task, name, tag)
	}
	return i
}

func (t *Trader) Spend(amount uint32) {
	t.Money -= amount
}

func (t *Trader) Earn(amount uint32) {
	t.Money += amount
}

func (t *Trader) GetMoney() uint32 {
	return t.Money
}

func (t *Trader) Destination(extensionType *building.BuildingExtensionType) navigation.Destination {
	return &navigation.TravellerDestination{T: t.Person.Traveller}
}

func (t *Trader) Stats() *stats.Stats {
	return &stats.Stats{
		Money:     t.Money,
		People:    1,
		Buildings: 0,
		Artifacts: t.Resources.NumArtifacts(),
	}
}

func (t *Trader) PendingCosts() uint32 {
	return PendingCosts(t.Tasks)
}

func (t *Trader) Broken() bool {
	return false
}

func (t *Trader) GetTown() *Town {
	return t.Town
}

func (t *Trader) GetPeople() []*Person {
	return []*Person{t.Person}
}

func (t *Trader) GetHome() Home {
	return t
}

func (t *Trader) GetExchange() economy.Exchange {
	return t.SourceExchange
}

func (t *Trader) IsHomeVehicle() bool {
	return true
}
