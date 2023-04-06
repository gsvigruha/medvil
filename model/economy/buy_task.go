package economy

import (
	"medvil/model/artifacts"
	"medvil/model/navigation"
	"medvil/model/time"
)

type BuyTask struct {
	TaskBase
	Exchange       Exchange
	HouseholdMoney *uint32
	Goods          []artifacts.Artifacts
	MaxPrice       uint32
	TaskTag        string
}

func (t *BuyTask) Field() *navigation.Field {
	return nil
}

func (t *BuyTask) Complete(Calendar *time.CalendarType, tool bool) bool {
	return false
}

func (t *BuyTask) Blocked() bool {
	return t.Exchange.Price(t.Goods) > *t.HouseholdMoney || !t.Exchange.CanBuy(t.Goods)
}

func (t *BuyTask) Name() string {
	return "exchange"
}

func (t *BuyTask) Tag() string {
	return t.TaskTag
}

func (t *BuyTask) Expired(Calendar *time.CalendarType) bool {
	price := t.Exchange.Price(t.Goods)
	expired := price > t.MaxPrice || price > *t.HouseholdMoney
	t.Exchange.RegisterBuyTask(t, !expired)
	return expired
}

func (t *BuyTask) Motion() uint8 {
	return navigation.MotionStand
}

func (t *BuyTask) IconName() string {
	return "buy"
}
