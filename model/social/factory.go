package social

import (
	"medvil/model/navigation"
	"medvil/model/time"
)

type Factory struct {
	Household Household
}

func (f *Factory) ElapseTime(Calendar *time.CalendarType, m navigation.IMap) {
	f.Household.ElapseTime(Calendar, m)
}
