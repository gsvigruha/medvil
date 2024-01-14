package economy

import (
	"medvil/model/artifacts"
	"medvil/model/navigation"
	"medvil/model/time"
)

type BuyTask struct {
	TaskBase
	Exchange        Exchange
	HouseholdWallet Wallet
	Goods           []artifacts.Artifacts
	MaxPrice        uint32
	TaskTag         Tag
}

func (t *BuyTask) Destination() navigation.Destination {
	return nil
}

func (t *BuyTask) Complete(m navigation.IMap, tool bool) bool {
	return false
}

func (t *BuyTask) Blocked() bool {
	return t.Exchange.Price(t.Goods) > t.HouseholdWallet.GetMoney() || !t.Exchange.CanBuy(t.Goods)
}

func (t *BuyTask) Name() string {
	return "exchange"
}

func (t *BuyTask) Tags() Tags {
	return MakeTags(t.TaskTag)
}

func (t *BuyTask) Expired(Calendar *time.CalendarType) bool {
	price := t.Exchange.Price(t.Goods)
	expired := price > t.MaxPrice || price > t.HouseholdWallet.GetMoney()
	t.Exchange.RegisterBuyTask(t, !expired)
	return expired
}

func (t *BuyTask) Motion() uint8 {
	return navigation.MotionStand
}

func (t *BuyTask) IconName() string {
	return "buy"
}

func (t *BuyTask) Description() string {
	return "Buy goods at the market"
}
