package social

import (
	"medvil/model/navigation"
)

type IMap interface {
	navigation.IMap
	GetCountries(t uint8) []*Country
	GetNearbyGuard(t *navigation.Traveller) *Person
}
