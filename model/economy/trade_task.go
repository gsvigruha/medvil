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
	SourceMarketD     navigation.Destination
	TargetMarketD     navigation.Destination
	SourceExchange    Exchange
	TargetExchange    Exchange
	TraderR           *artifacts.Resources
	TraderWallet      Wallet
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
		return t.SourceMarketD
	case TradeTaskStateDropoffAtDest:
		return t.TargetMarketD
	case TradeTaskStatePickupAtDest:
		return t.TargetMarketD
	case TradeTaskStateDropoffAtSource:
		return t.SourceMarketD
	}
	return nil
}

func (t *TradeTask) Complete(Calendar *time.CalendarType, tool bool) bool {
	switch t.state {
	case TradeTaskStatePickupAtSource:
		t.Goods = []artifacts.Artifacts{}
		if t.SourceExchange.Price(t.GoodsSourceToDest) <= t.TraderWallet.GetMoney() {
			t.Goods = t.SourceExchange.BuyAsManyAsPossible(t.GoodsSourceToDest, t.TraderWallet)
		}
		t.state = TradeTaskStateDropoffAtDest
	case TradeTaskStateDropoffAtDest:
		if t.TargetExchange.CanSell(t.Goods) {
			t.TargetExchange.Sell(t.Goods, t.TraderWallet)
		}
		t.state = TradeTaskStatePickupAtDest
	case TradeTaskStatePickupAtDest:
		t.Goods = []artifacts.Artifacts{}
		if t.TargetExchange.Price(t.GoodsDestToSource) <= t.TraderWallet.GetMoney() {
			t.Goods = t.TargetExchange.BuyAsManyAsPossible(t.GoodsDestToSource, t.TraderWallet)
		}
		t.state = TradeTaskStateDropoffAtSource
	case TradeTaskStateDropoffAtSource:
		if t.SourceExchange.CanSell(t.Goods) {
			t.SourceExchange.Sell(t.Goods, t.TraderWallet)
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
