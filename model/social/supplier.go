package social

import (
	"medvil/model/navigation"
)

type Supplier interface {
	GetHome() Home
	ReassignFirstPerson(dstH Home, assingTask bool, m navigation.IMap)
}
