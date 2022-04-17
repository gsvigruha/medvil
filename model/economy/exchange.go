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
	RegisterSellTask(*SellTask, bool)
	RegisterBuyTask(*BuyTask, bool)
}

const ExchangeTaskStatePickupAtHome uint8 = 0
const ExchangeTaskStateMarket uint8 = 1
const ExchangeTaskStateDropoffAtHome uint8 = 2
const MaxWaitTime = 24 * 10

type ExchangeTask struct {
	TaskBase
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
	waittime       uint16
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

func (t *ExchangeTask) Complete(Calendar *time.CalendarType, tool bool) bool {
	switch t.state {
	case ExchangeTaskStatePickupAtHome:
		t.goods = t.HouseholdR.GetAsManyAsPossible(t.GoodsToSell)
		t.state = ExchangeTaskStateMarket
	case ExchangeTaskStateMarket:
		if t.Exchange.CanSell(t.goods) && t.Exchange.Price(t.GoodsToBuy) <= *t.HouseholdMoney {
			t.Exchange.Sell(t.goods, t.HouseholdMoney)
			t.goods = t.Exchange.BuyAsManyAsPossible(t.GoodsToBuy, t.HouseholdMoney)
			t.state = ExchangeTaskStateDropoffAtHome
		} else {
			t.waittime++
		}
		if t.waittime > MaxWaitTime {
			if t.Exchange.CanSell(t.goods) {
				t.Exchange.Sell(t.goods, t.HouseholdMoney)
				t.goods = []artifacts.Artifacts{}
			}
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
		return !t.Exchange.CanSell(t.goods) || t.Exchange.Price(t.GoodsToBuy) > *t.HouseholdMoney || !t.Exchange.CanBuy(t.GoodsToBuy)
	case ExchangeTaskStateMarket:
		return !t.Exchange.CanSell(t.goods) || t.Exchange.Price(t.GoodsToBuy) > *t.HouseholdMoney || !t.Exchange.CanBuy(t.GoodsToBuy)
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

func (t *ExchangeTask) Expired(Calendar *time.CalendarType) bool {
	return false
}

func (t *ExchangeTask) AddBuyTask(bt *BuyTask) {
	t.GoodsToBuy = append(t.GoodsToBuy, bt.Goods...)
	t.TaskTag = t.TaskTag + ";" + bt.TaskTag
	t.Exchange.RegisterBuyTask(bt, false)
}

func (t *ExchangeTask) AddSellTask(st *SellTask) {
	t.GoodsToSell = append(t.GoodsToSell, st.Goods...)
	t.TaskTag = t.TaskTag + ";" + st.TaskTag
	t.Exchange.RegisterSellTask(st, false)
}

func (t *ExchangeTask) Motion() uint8 {
	return navigation.MotionStand
}
