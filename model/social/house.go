package social

import (
	"medvil/model/navigation"
)

type House interface {
	GetHousehold() *Household
	GetFields() []navigation.FieldWithContext
}
