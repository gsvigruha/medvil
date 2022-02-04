package economy

import (
	"medvil/model/artifacts"
	"medvil/model/navigation"
	"medvil/model/time"
)

type Exchange interface {
	Buy([]artifacts.Artifacts)
	Sell([]artifacts.Artifacts)
	CanBuy([]artifacts.Artifacts) bool
	CanSell([]artifacts.Artifacts) bool
	Price([]artifacts.Artifacts) uint16
}

type ExchangeTask struct {
	PickupF        *navigation.Field
	DropoffF       *navigation.Field
	Exchange       Exchange
	HouseholdR     *artifacts.Resources
	HouseholdMoney *uint16
	GoodsToBuy     []artifacts.Artifacts
	GoodsToSell    []artifacts.Artifacts
	dropoff        bool
}

func (t *ExchangeTask) Field() *navigation.Field {
	if t.dropoff {
		return t.DropoffF
	} else {
		return t.PickupF
	}
}

func (t *ExchangeTask) Complete(Calendar *time.CalendarType) bool {
	if t.dropoff {
		t.HouseholdR.AddAll(t.GoodsToBuy)
		return true
	} else {
		t.Exchange.Buy(t.GoodsToBuy)
		t.Exchange.Sell(t.GoodsToSell)
		t.dropoff = true
	}
	return false
}

func (t *ExchangeTask) Blocked() bool {
	if t.dropoff {
		return false
	} else {
		return t.Exchange.CanBuy(t.GoodsToBuy) && t.Exchange.CanSell(t.GoodsToSell) && t.Exchange.Price(t.GoodsToBuy) >= *t.HouseholdMoney
	}
}

func (t *ExchangeTask) Name() string {
	return "exchange"
}
