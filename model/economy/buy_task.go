package economy

import (
	"medvil/model/artifacts"
	"medvil/model/navigation"
	"medvil/model/time"
)

type BuyTask struct {
	Exchange       Exchange
	HouseholdMoney *uint32
	Artifact       *artifacts.Artifact
	Quantity       uint16
	MaxPrice       uint16
	TaskTag        string
}

func (t *BuyTask) GoodsToBuy() []artifacts.Artifacts {
	return []artifacts.Artifacts{artifacts.Artifacts{A: t.Artifact, Quantity: t.Quantity}}
}

func (t *BuyTask) Field() *navigation.Field {
	return nil
}

func (t *BuyTask) Complete(Calendar *time.CalendarType) bool {
	return false
}

func (t *BuyTask) Blocked() bool {
	return t.Exchange.Price(t.GoodsToBuy()) > *t.HouseholdMoney || !t.Exchange.CanBuy(t.GoodsToBuy())
}

func (t *BuyTask) Name() string {
	return "buy"
}

func (t *BuyTask) Tag() string {
	return t.TaskTag
}

func (t *BuyTask) Expired(Calendar *time.CalendarType) bool {
	return false
}
