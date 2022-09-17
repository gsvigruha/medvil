package economy

import (
	"medvil/model/artifacts"
	"medvil/model/navigation"
	"medvil/model/time"
	"medvil/model/vehicles"
)

const TradeTaskStatePickup uint8 = 0
const TradeTaskStateDrop uint8 = 1

type TradeTask struct {
	TaskBase
	SourceMarketF  *navigation.Field
	TargetMarketF  *navigation.Field
	SourceExchange Exchange
	TargetExchange Exchange
	TraderR        *artifacts.Resources
	TraderMoney    *uint32
	Vehicle        *vehicles.Vehicle
	GoodsToTrade   []artifacts.Artifacts
	TaskTag        string
	goods          []artifacts.Artifacts
	state          uint8
	waittime       uint16
}

func (t *TradeTask) Field() *navigation.Field {
	switch t.state {
	case TradeTaskStatePickup:
		return t.SourceMarketF
	case TradeTaskStateDrop:
		return t.TargetMarketF
	}
	return nil
}

func (t *TradeTask) Complete(Calendar *time.CalendarType, tool bool) bool {
	switch t.state {
	case TradeTaskStatePickup:
		t.goods = t.SourceExchange.BuyAsManyAsPossible(t.GoodsToTrade, t.TraderMoney)
		t.state = TradeTaskStateDrop
	case TradeTaskStateDrop:
		if t.TargetExchange.CanSell(t.goods) {
			t.TargetExchange.Sell(t.goods, t.TraderMoney)
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
	switch t.state {
	case TradeTaskStatePickup:
		return !t.SourceExchange.CanSell(t.GoodsToTrade) || t.SourceExchange.Price(t.GoodsToTrade) > *t.TraderMoney
	case TradeTaskStateDrop:
		return false
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
