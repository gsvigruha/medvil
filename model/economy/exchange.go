package economy

import (
	"medvil/model/artifacts"
	"medvil/model/navigation"
	"medvil/model/time"
)

type Exchange interface {
	Buy([]artifacts.Artifacts, *uint32)
	Sell([]artifacts.Artifacts, *uint32)
	CanBuy([]artifacts.Artifacts) bool
	CanSell([]artifacts.Artifacts) bool
	Price([]artifacts.Artifacts) uint32
}

type ExchangeTask struct {
	HomeF          *navigation.Field
	MarketF        *navigation.Field
	Exchange       Exchange
	HouseholdR     *artifacts.Resources
	HouseholdMoney *uint32
	GoodsToBuy     []artifacts.Artifacts
	GoodsToSell    []artifacts.Artifacts
	backtrip       bool
}

func (t *ExchangeTask) Field() *navigation.Field {
	if t.backtrip {
		return t.HomeF
	} else {
		return t.MarketF
	}
}

func (t *ExchangeTask) Complete(Calendar *time.CalendarType) bool {
	if t.backtrip {
		t.HouseholdR.AddAll(t.GoodsToBuy)
		return true
	} else {
		if t.Exchange.CanBuy(t.GoodsToBuy) && t.Exchange.CanSell(t.GoodsToSell) {
			t.Exchange.Buy(t.GoodsToBuy, t.HouseholdMoney)
			t.Exchange.Sell(t.GoodsToSell, t.HouseholdMoney)
			t.backtrip = true
		}
	}
	return false
}

func (t *ExchangeTask) Blocked() bool {
	if t.backtrip {
		return false
	} else {
		return !t.Exchange.CanBuy(t.GoodsToBuy) || !t.Exchange.CanSell(t.GoodsToSell) || t.Exchange.Price(t.GoodsToBuy) > *t.HouseholdMoney
	}
}

func (t *ExchangeTask) Name() string {
	return "exchange"
}

func (t *ExchangeTask) Tag() string {
	return ""
}
