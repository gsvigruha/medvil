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
	SellAsManyAsPossible([]artifacts.Artifacts, Wallet)
	CanBuy([]artifacts.Artifacts) bool
	CanSell([]artifacts.Artifacts) bool
	Price([]artifacts.Artifacts) uint32
	RegisterSellTask(*SellTask, bool)
	RegisterBuyTask(*BuyTask, bool)
	HasAny([]artifacts.Artifacts) bool
}

const ExchangeTaskStatePickupAtHome uint8 = 0
const ExchangeTaskStateMarket uint8 = 1
const ExchangeTaskStateDropoffAtHome uint8 = 2
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
	TaskTag         string
	Goods           []artifacts.Artifacts
	State           uint8
	Waittime        uint16
}

func (t *ExchangeTask) Destination() navigation.Destination {
	switch t.State {
	case ExchangeTaskStatePickupAtHome:
		return t.HomeD
	case ExchangeTaskStateMarket:
		return t.MarketD
	case ExchangeTaskStateDropoffAtHome:
		return t.HomeD
	}
	return nil
}

func (t *ExchangeTask) Complete(Calendar *time.CalendarType, tool bool) bool {
	switch t.State {
	case ExchangeTaskStatePickupAtHome:
		t.Goods = t.HouseholdR.GetAsManyAsPossible(t.GoodsToSell)
		if t.Vehicle != nil {
			t.Traveller.UseVehicle(t.Vehicle)
		}
		t.State = ExchangeTaskStateMarket
	case ExchangeTaskStateMarket:
		if t.Exchange.CanSell(t.Goods) && t.Exchange.Price(t.GoodsToBuy) <= t.HouseholdWallet.GetMoney() {
			t.Exchange.Sell(t.Goods, t.HouseholdWallet)
			t.Goods = t.Exchange.BuyAsManyAsPossible(t.GoodsToBuy, t.HouseholdWallet)
			t.State = ExchangeTaskStateDropoffAtHome
		} else {
			t.Waittime++
		}
		if t.Waittime > MaxWaitTime {
			if t.Exchange.CanSell(t.Goods) {
				t.Exchange.Sell(t.Goods, t.HouseholdWallet)
				t.Goods = []artifacts.Artifacts{}
			}
			t.State = ExchangeTaskStateDropoffAtHome
		}
	case ExchangeTaskStateDropoffAtHome:
		t.HouseholdR.AddAll(t.Goods)
		t.Traveller.ExitVehicle()
		return true
	}
	return false
}

func (t *ExchangeTask) Blocked() bool {
	switch t.State {
	case ExchangeTaskStatePickupAtHome:
		return !t.Exchange.CanSell(t.Goods) || t.Exchange.Price(t.GoodsToBuy) > t.HouseholdWallet.GetMoney() || !t.Exchange.CanBuy(t.GoodsToBuy)
	case ExchangeTaskStateMarket:
		return !t.Exchange.CanSell(t.Goods) || t.Exchange.Price(t.GoodsToBuy) > t.HouseholdWallet.GetMoney() || !t.Exchange.CanBuy(t.GoodsToBuy)
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
