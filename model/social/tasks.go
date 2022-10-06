package social

import (
	"math"
	"medvil/model/artifacts"
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
	var i = 0
	taskTags := strings.Split(task.Tag(), ";")
	for _, taskTag := range taskTags {
		if task.Name() == name && strings.Contains(taskTag, tag) {
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

func GetExchangeTask(h Home, mp *Marketplace, m navigation.IMap, vehicle *vehicles.Vehicle) *economy.ExchangeTask {
	var maxVolume uint16 = ExchangeTaskMaxVolumePedestrian
	var buildingCheckFn = navigation.Field.BuildingNonExtension
	_, _, sailableMP := GetRandomBuildingXY(mp.Building, m, navigation.Field.BoatDestination)
	sailableH := h.RandomField(m, navigation.Field.BoatDestination) != nil
	if vehicle != nil {
		if vehicle.T.Water && sailableMP && sailableH {
			buildingCheckFn = navigation.Field.BoatDestination
		}
		maxVolume = vehicle.T.MaxVolume
	}

	mx, my, mok := GetRandomBuildingXY(mp.Building, m, buildingCheckFn)
	hf := h.RandomField(m, buildingCheckFn)
	if hf == nil || !mok {
		return nil
	}
	et := &economy.ExchangeTask{
		HomeF:          hf,
		MarketF:        m.GetField(mx, my),
		Exchange:       mp,
		HouseholdR:     h.GetResources(),
		HouseholdMoney: h.GetMoney(),
		Vehicle:        vehicle,
		GoodsToBuy:     nil,
		GoodsToSell:    nil,
		TaskTag:        "",
	}
	var empty = true
	var tasks []economy.Task
	for _, ot := range h.GetTasks() {
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
		} else {
			empty = false
		}
	}
	if !empty {
		h.SetTasks(tasks)
		return et
	}
	return nil
}

func GetNextTask(h Home, e economy.Equipment) economy.Task {
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
