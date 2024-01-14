package economy

import (
	"medvil/model/artifacts"
	"medvil/model/navigation"
	"medvil/model/time"
	"medvil/model/vehicles"
)

type Wallet interface {
	Spend(uint32)
	Earn(uint32)
	GetMoney() uint32
}

type Exchange interface {
	Buy([]artifacts.Artifacts, Wallet)
	BuyAsManyAsPossible([]artifacts.Artifacts, Wallet) []artifacts.Artifacts
	Sell([]artifacts.Artifacts, Wallet)
	SellAsManyAsPossible([]artifacts.Artifacts, Wallet) []artifacts.Artifacts
	CanBuy([]artifacts.Artifacts) bool
	CanSell([]artifacts.Artifacts) bool
	Price([]artifacts.Artifacts) uint32
	RegisterSellTask(*SellTask, bool)
	RegisterBuyTask(*BuyTask, bool)
	HasAny([]artifacts.Artifacts) bool
	UnitPrice(*artifacts.Artifact) uint32
}

const ExchangeTaskStatePickupAtHome uint8 = 0
const ExchangeTaskStateMarketSell uint8 = 1
const ExchangeTaskStateMarketBuy uint8 = 2
const ExchangeTaskStateDropoffAtHome uint8 = 3
const MaxWaitTime = 24 * 10

type ExchangeTask struct {
	TaskBase
	HomeD           navigation.Destination
	MarketD         navigation.Destination
	Exchange        Exchange
	HouseholdR      *artifacts.Resources
	HouseholdWallet Wallet
	Vehicle         *vehicles.Vehicle
	GoodsToBuy      []artifacts.Artifacts
	GoodsToSell     []artifacts.Artifacts
	TaskTags        Tags
	Goods           []artifacts.Artifacts
	State           uint8
}

func (t *ExchangeTask) Destination() navigation.Destination {
	switch t.State {
	case ExchangeTaskStatePickupAtHome:
		return t.HomeD
	case ExchangeTaskStateMarketSell:
		return t.MarketD
	case ExchangeTaskStateMarketBuy:
		return t.MarketD
	case ExchangeTaskStateDropoffAtHome:
		return t.HomeD
	}
	return nil
}

func (t *ExchangeTask) Complete(m navigation.IMap, tool bool) bool {
	switch t.State {
	case ExchangeTaskStatePickupAtHome:
		t.Goods = t.HouseholdR.GetAsManyAsPossible(t.GoodsToSell)
		if t.Vehicle != nil {
			t.Traveller.UseVehicle(t.Vehicle)
		}
		t.State = ExchangeTaskStateMarketSell
	case ExchangeTaskStateMarketSell:
		t.Goods = artifacts.Filter(t.Exchange.SellAsManyAsPossible(t.Goods, t.HouseholdWallet))
		t.State = ExchangeTaskStateMarketBuy
	case ExchangeTaskStateMarketBuy:
		if t.Exchange.Price(t.GoodsToBuy) <= t.HouseholdWallet.GetMoney() {
			bought := artifacts.Filter(t.Exchange.BuyAsManyAsPossible(t.GoodsToBuy, t.HouseholdWallet))
			t.Goods = append(t.Goods, bought...)
		}
		t.State = ExchangeTaskStateDropoffAtHome
	case ExchangeTaskStateDropoffAtHome:
		t.HouseholdR.AddAll(t.Goods)
		return true
	}
	return false
}

func (t *ExchangeTask) Blocked() bool {
	return false
}

func (t *ExchangeTask) Name() string {
	return "exchange"
}

func (t *ExchangeTask) Tags() Tags {
	return t.TaskTags
}

func (t *ExchangeTask) Expired(Calendar *time.CalendarType) bool {
	if t.Vehicle != nil && t.Vehicle.Broken {
		return true
	}
	// Expire exchange task if the marketplace got recreated
	return t.Household != nil && t.Household.GetExchange() != t.Exchange
}

func (t *ExchangeTask) AddBuyTask(bt *BuyTask) {
	t.GoodsToBuy = append(t.GoodsToBuy, bt.Goods...)
	t.TaskTags = AppendTags(t.TaskTags, bt.TaskTag)
	t.Exchange.RegisterBuyTask(bt, false)
}

func (t *ExchangeTask) AddSellTask(st *SellTask) {
	t.GoodsToSell = append(t.GoodsToSell, st.Goods...)
	t.TaskTags = AppendTags(t.TaskTags, st.TaskTag)
	t.Exchange.RegisterSellTask(st, false)
}

func (t *ExchangeTask) Motion() uint8 {
	return navigation.MotionStand
}

func (t *ExchangeTask) Description() string {
	return "Buy and sell goods at the market"
}
