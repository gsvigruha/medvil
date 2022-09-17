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
	goods             []artifacts.Artifacts
	state             uint8
	waittime          uint16
}

func (t *TradeTask) Field() *navigation.Field {
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
		t.goods = t.SourceExchange.BuyAsManyAsPossible(t.GoodsSourceToDest, t.TraderMoney)
		t.state = TradeTaskStateDropoffAtDest
		t.waittime = 0
	case TradeTaskStateDropoffAtDest:
		if t.TargetExchange.CanSell(t.goods) {
			t.TargetExchange.Sell(t.goods, t.TraderMoney)
			t.state = TradeTaskStatePickupAtDest
		} else {
			t.waittime++
		}
		if t.waittime > MaxWaitTime {
			t.state = TradeTaskStatePickupAtDest
		}
	case TradeTaskStatePickupAtDest:
		t.goods = t.TargetExchange.BuyAsManyAsPossible(t.GoodsDestToSource, t.TraderMoney)
		t.state = TradeTaskStateDropoffAtSource
		t.waittime = 0
	case TradeTaskStateDropoffAtSource:
		if t.SourceExchange.CanSell(t.goods) {
			t.SourceExchange.Sell(t.goods, t.TraderMoney)
			return true
		} else {
			t.waittime++
		}
		if t.waittime > MaxWaitTime {
			return true
		}
	}
	return false
}

func (t *TradeTask) Blocked() bool {
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
