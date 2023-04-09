package economy

import (
	"medvil/model/artifacts"
	"medvil/model/navigation"
	"medvil/model/time"
	"medvil/model/vehicles"
)

const TradeTaskStatePickupAtSource uint8 = 0
const TradeTaskStateDropoffAtDest uint8 = 1
const TradeTaskStatePickupAtDest uint8 = 2
const TradeTaskStateDropoffAtSource uint8 = 3

type TradeTask struct {
	TaskBase
	SourceMarketF     *navigation.Field
	TargetMarketF     *navigation.Field
	SourceExchange    Exchange
	TargetExchange    Exchange
	TraderR           *artifacts.Resources
	TraderMoney       *uint32
	Vehicle           *vehicles.Vehicle
	GoodsSourceToDest []artifacts.Artifacts
	GoodsDestToSource []artifacts.Artifacts
	TaskTag           string
	Goods             []artifacts.Artifacts
	state             uint8
}

func (t *TradeTask) Destination() navigation.Destination {
	switch t.state {
	case TradeTaskStatePickupAtSource:
		return t.SourceMarketF
	case TradeTaskStateDropoffAtDest:
		return t.TargetMarketF
	case TradeTaskStatePickupAtDest:
		return t.TargetMarketF
	case TradeTaskStateDropoffAtSource:
		return t.SourceMarketF
	}
	return nil
}

func (t *TradeTask) Complete(Calendar *time.CalendarType, tool bool) bool {
	switch t.state {
	case TradeTaskStatePickupAtSource:
		t.Goods = []artifacts.Artifacts{}
		if t.SourceExchange.Price(t.GoodsSourceToDest) <= *t.TraderMoney {
			t.Goods = t.SourceExchange.BuyAsManyAsPossible(t.GoodsSourceToDest, t.TraderMoney)
		}
		t.state = TradeTaskStateDropoffAtDest
	case TradeTaskStateDropoffAtDest:
		if t.TargetExchange.CanSell(t.Goods) {
			t.TargetExchange.Sell(t.Goods, t.TraderMoney)
		}
		t.state = TradeTaskStatePickupAtDest
	case TradeTaskStatePickupAtDest:
		t.Goods = []artifacts.Artifacts{}
		if t.TargetExchange.Price(t.GoodsDestToSource) <= *t.TraderMoney {
			t.Goods = t.TargetExchange.BuyAsManyAsPossible(t.GoodsDestToSource, t.TraderMoney)
		}
		t.state = TradeTaskStateDropoffAtSource
	case TradeTaskStateDropoffAtSource:
		if t.SourceExchange.CanSell(t.Goods) {
			t.SourceExchange.Sell(t.Goods, t.TraderMoney)
		}
		return true
	}
	return false
}

func (t *TradeTask) Blocked() bool {
	switch t.state {
	case TradeTaskStatePickupAtSource:
		return !t.SourceExchange.HasAny(t.GoodsSourceToDest)
	case TradeTaskStatePickupAtDest:
		return !t.TargetExchange.HasAny(t.GoodsDestToSource)
	}
	return false
}

func (t *TradeTask) Name() string {
	return "trade"
}

func (t *TradeTask) Tag() string {
	return t.TaskTag
}

func (t *TradeTask) Expired(Calendar *time.CalendarType) bool {
	return false
}

func (t *TradeTask) Motion() uint8 {
	return navigation.MotionStand
}
