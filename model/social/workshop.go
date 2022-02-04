package social

import (
	"medvil/model/economy"
	"medvil/model/navigation"
	"medvil/model/time"
)

type Workshop struct {
	Household   Household
	Manufacture *economy.Manufacture
}

func (w *Workshop) ElapseTime(Calendar *time.CalendarType, m navigation.IMap) {
	w.Household.ElapseTime(Calendar, m)
	needs := w.Household.Resources.Needs(w.Manufacture.Inputs)
	if needs != nil && w.Household.Town.Marketplace.Storage.Has(needs) {
	}
}
