package social

import (
	"medvil/model/navigation"
	"medvil/model/time"
)

type Townhall struct {
	Household Household
}

func (t *Townhall) ElapseTime(Calendar *time.CalendarType, m navigation.IMap) {
	t.Household.ElapseTime(Calendar)
}
