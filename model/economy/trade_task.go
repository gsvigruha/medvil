package economy

import (
	"medvil/model/artifacts"
	"medvil/model/navigation"
	"medvil/model/time"
	"medvil/model/vehicles"
)

type TradeTask struct {
	TaskBase
	MarketF        *navigation.Field
	Exchange       Exchange
	TraderR        *artifacts.Resources
	TraderMoney    *uint32
	Vehicle        *vehicles.Vehicle
	GoodsToBuy     []artifacts.Artifacts
	GoodsToSell    []artifacts.Artifacts
	TaskTag        string
	goods          []artifacts.Artifacts
	state          uint8
	waittime       uint16
}

func (t *TradeTask) Field() *navigation.Field {
	return t.MarketF
}

func (t *TradeTask) Complete(Calendar *time.CalendarType, tool bool) bool {
	return false
}

func (t *TradeTask) Blocked() bool {
	return !t.Exchange.CanSell(t.goods) || t.Exchange.Price(t.GoodsToBuy) > *t.TraderMoney || !t.Exchange.CanBuy(t.GoodsToBuy)
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
