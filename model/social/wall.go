package social

import (
	"math/rand"
	"medvil/model/building"
	"medvil/model/navigation"
	"medvil/model/time"
)

const WallBrokenRate = 1.0 / (24.0 * 365.0 * 10.0)

type Wall struct {
	Building *building.Building
	Town     *Town
	F        *navigation.Field
}

func (w *Wall) Field() *navigation.Field {
	return w.F
}

func (w *Wall) Context() string {
	return ""
}

func (w *Wall) ElapseTime(Calendar *time.CalendarType, m navigation.IMap) {
	if rand.Float64() < WallBrokenRate {
		w.Building.Broken = true
	}
}
