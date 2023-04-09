package economy

import (
	"medvil/model/artifacts"
	"medvil/model/navigation"
	"medvil/model/time"
)

type SellTask struct {
	TaskBase
	Exchange Exchange
	Goods    []artifacts.Artifacts
	TaskTag  string
}

func (t *SellTask) Destination() navigation.Destination {
	return nil
}

func (t *SellTask) Complete(Calendar *time.CalendarType, tool bool) bool {
	return false
}

func (t *SellTask) Blocked() bool {
	return !t.Exchange.CanSell(t.Goods)
}

func (t *SellTask) Name() string {
	return "exchange"
}

func (t *SellTask) Tag() string {
	return t.TaskTag
}

func (t *SellTask) Expired(Calendar *time.CalendarType) bool {
	t.Exchange.RegisterSellTask(t, true)
	return false
}

func (t *SellTask) Motion() uint8 {
	return navigation.MotionStand
}

func (t *SellTask) IconName() string {
	return "sell"
}
