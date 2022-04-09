package economy

import (
	"medvil/model/artifacts"
	"medvil/model/navigation"
	"medvil/model/time"
)

type SellTask struct {
	Exchange Exchange
	Artifact *artifacts.Artifact
	Quantity uint16
	TaskTag  string
}

func (t *SellTask) GoodsToSell() []artifacts.Artifacts {
	return []artifacts.Artifacts{artifacts.Artifacts{A: t.Artifact, Quantity: t.Quantity}}
}

func (t *SellTask) Field() *navigation.Field {
	return nil
}

func (t *SellTask) Complete(Calendar *time.CalendarType) bool {
	return false
}

func (t *SellTask) Blocked() bool {
	return !t.Exchange.CanSell(t.GoodsToSell())
}

func (t *SellTask) Name() string {
	return "buy"
}

func (t *SellTask) Tag() string {
	return t.TaskTag
}

func (t *SellTask) Expired(Calendar *time.CalendarType) bool {
	return false
}
