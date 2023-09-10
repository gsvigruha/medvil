package social

import (
	"math"
	"medvil/model/artifacts"
	"medvil/model/building"
	"medvil/model/economy"
	"medvil/model/navigation"
	"medvil/model/vehicles"
	"strings"
)

func NumBatchesSimple(totalQuantity, transportQuantity uint16) int {
	return NumBatches(totalQuantity, transportQuantity, transportQuantity)
}

func NumBatches(totalQuantity, minTransportQuantity, maxTransportQuantity uint16) int {
	if totalQuantity < minTransportQuantity {
		return 0
	}
	return int(math.Ceil(float64(totalQuantity) / float64(maxTransportQuantity)))
}

func CountTags(task economy.Task, name, tag string) int {
	if task.Name() != name {
		return 0
	}
	var i = 0
	taskTags := strings.Split(task.Tag(), ";")
	for _, taskTag := range taskTags {
		if strings.Contains(taskTag, tag) {
			i++
		}
	}
	return i
}

func IsExchangeBaseTask(t economy.Task) bool {
	_, sok := t.(*economy.SellTask)
	_, bok := t.(*economy.BuyTask)
	return sok || bok
}

func BuildingsConnectedWithWater(b1, b2 *building.Building, m navigation.IMap) bool {
	for _, e1 := range b1.GetExtensionsWithCoords(building.Deck) {
		for _, e2 := range b2.GetExtensionsWithCoords(building.Deck) {
			if m.GetField(e1.X, e1.Y).IslandLabel == m.GetField(e2.X, e2.Y).IslandLabel {
				return true
			}
		}
	}
	return false
}

func CombineExchangeTasks(h Home, mp *Marketplace, m navigation.IMap) {
	waterOk := h.GetBuilding() == nil || BuildingsConnectedWithWater(h.GetBuilding(), mp.Building, m)
	var vehicle *vehicles.Vehicle
	var buildingCheckFn = navigation.Field.BuildingNonExtension
	var buildingExtension = building.NonExtension
	var et *economy.ExchangeTask
	var maxVolume uint16 = 0
	var batchStart = true

	var tasks []economy.Task
	for _, ot := range h.GetTasks() {
		if batchStart {
			vehicle = h.AllocateVehicle(waterOk)
			if vehicle != nil {
				if vehicle.T.Water {
					buildingCheckFn = navigation.Field.BoatDestination
					buildingExtension = building.Deck
				} else {
					buildingCheckFn = navigation.Field.BuildingNonExtension
					buildingExtension = building.NonExtension
				}
				maxVolume = vehicle.T.MaxVolume
			} else {
				buildingCheckFn = navigation.Field.BuildingNonExtension
				buildingExtension = building.NonExtension
				maxVolume = ExchangeTaskMaxVolumePedestrian
			}
			var hf navigation.Destination = h.RandomField(m, buildingCheckFn)
			if h.IsHomeVehicle() {
				hf = &navigation.TravellerDestination{T: vehicle.Traveller}
			} else if vehicle != nil && vehicle.Parking != nil {
				hf = vehicle.Parking
			}
			et = &economy.ExchangeTask{
				HomeD:           hf,
				MarketD:         &navigation.BuildingDestination{B: mp.Building, ET: buildingExtension},
				Exchange:        mp,
				HouseholdR:      h.GetResources(),
				HouseholdWallet: h,
				Vehicle:         vehicle,
				GoodsToBuy:      nil,
				GoodsToSell:     nil,
				TaskTag:         "market",
			}
			batchStart = false
		}

		var combined = false
		bt, bok := ot.(*economy.BuyTask)
		if bok && !bt.Blocked() && !bt.IsPaused() && artifacts.GetVolume(et.GoodsToBuy) < maxVolume {
			et.AddBuyTask(bt)
			combined = true
		}
		st, sok := ot.(*economy.SellTask)
		if sok && !st.Blocked() && !st.IsPaused() && artifacts.GetVolume(et.GoodsToSell) < maxVolume {
			et.AddSellTask(st)
			combined = true
		}
		if !combined {
			tasks = append(tasks, ot)
		} else if artifacts.GetVolume(et.GoodsToBuy) >= maxVolume || artifacts.GetVolume(et.GoodsToSell) >= maxVolume {
			tasks = append([]economy.Task{et}, tasks...)
			et = nil
			vehicle = nil
			batchStart = true
		}
	}
	if et != nil {
		if artifacts.GetVolume(et.GoodsToSell) > 0 || artifacts.GetVolume(et.GoodsToBuy) > 0 {
			tasks = append([]economy.Task{et}, tasks...)
		} else if vehicle != nil {
			vehicle.SetInUse(false)
		}
	}
	h.SetTasks(tasks)
}

func FirstUnblockedTask(h Home, e *economy.EquipmentType) economy.Task {
	if len(h.GetTasks()) == 0 {
		return nil
	}
	var i = 0
	for i < len(h.GetTasks()) {
		t := h.GetTasks()[i]
		if !t.Blocked() && !t.IsPaused() && t.Equipped(e) {
			break
		}
		i++
	}
	if i == len(h.GetTasks()) {
		return nil
	}
	return h.GetTasks()[i]
}

func GetNextTask(h Home, e *economy.EquipmentType) economy.Task {
	if len(h.GetTasks()) == 0 {
		return nil
	}
	var i = 0
	for i < len(h.GetTasks()) {
		t := h.GetTasks()[i]
		_, sok := t.(*economy.SellTask)
		_, bok := t.(*economy.BuyTask)
		if !sok && !bok && !t.Blocked() && !t.IsPaused() && t.Equipped(e) {
			break
		}
		i++
	}
	if i == len(h.GetTasks()) {
		return nil
	}
	t := h.GetTasks()[i]
	h.SetTasks(append(h.GetTasks()[0:i], h.GetTasks()[i+1:]...))
	return t
}
