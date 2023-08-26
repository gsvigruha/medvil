package social

import (
	"medvil/model/navigation"
)

type House interface {
	HomeProvider
	GetFields() []navigation.FieldWithContext
}
