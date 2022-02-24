package economy

import (
	"medvil/model/artifacts"
	"medvil/model/navigation"
	"medvil/model/time"
)

type Exchange interface {
	Buy([]artifacts.Artifacts, *uint32)
	BuyAsManyAsPossible([]artifacts.Artifacts, *uint32) []artifacts.Artifacts
	Sell([]artifacts.Artifacts, *uint32)
	CanBuy([]artifacts.Artifacts) bool
	CanSell([]artifacts.Artifacts) bool
	Price([]artifacts.Artifacts) uint32
}

const ExchangeTaskStatePickupAtHome uint8 = 0
const ExchangeTaskStateMarket uint8 = 1
const ExchangeTaskStateDropoffAtHome uint8 = 2

type ExchangeTask struct {
	HomeF          *navigation.Field
	MarketF        *navigation.Field
	Exchange       Exchange
	HouseholdR     *artifacts.Resources
	HouseholdMoney *uint32
	GoodsToBuy     []artifacts.Artifacts
	GoodsToSell    []artifacts.Artifacts
	TaskTag        string
	goods          []artifacts.Artifacts
	state          uint8
}

func (t *ExchangeTask) Field() *navigation.Field {
	switch t.state {
	case ExchangeTaskStatePickupAtHome:
		return t.HomeF
	case ExchangeTaskStateMarket:
		return t.MarketF
	case ExchangeTaskStateDropoffAtHome:
		return t.HomeF
	}
	return nil
}

func (t *ExchangeTask) Complete(Calendar *time.CalendarType) bool {
	switch t.state {
	case ExchangeTaskStatePickupAtHome:
		t.goods = t.HouseholdR.GetAsManyAsPossible(t.GoodsToSell)
		t.state = ExchangeTaskStateMarket
	case ExchangeTaskStateMarket:
		if t.Exchange.CanSell(t.goods) && t.Exchange.Price(t.GoodsToBuy) <= *t.HouseholdMoney {
			t.Exchange.Sell(t.goods, t.HouseholdMoney)
			t.goods = t.Exchange.BuyAsManyAsPossible(t.GoodsToBuy, t.HouseholdMoney)
			t.state = ExchangeTaskStateDropoffAtHome
		}
	case ExchangeTaskStateDropoffAtHome:
		t.HouseholdR.AddAll(t.goods)
		return true
	}
	return false
}

func (t *ExchangeTask) Blocked() bool {
	switch t.state {
	case ExchangeTaskStatePickupAtHome:
		return false
	case ExchangeTaskStateMarket:
		return !t.Exchange.CanSell(t.goods) || t.Exchange.Price(t.GoodsToBuy) > *t.HouseholdMoney
	case ExchangeTaskStateDropoffAtHome:
		return false
	}
	return false
}

func (t *ExchangeTask) Name() string {
	return "exchange"
}

func (t *ExchangeTask) Tag() string {
	return t.TaskTag
}
