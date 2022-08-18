package controller

import (
	"medvil/model/social"
	"medvil/view/gui"
)

func PersonToControlPanel(cp *ControlPanel, person *social.Person) {
	p := &gui.Panel{X: 0, Y: ControlPanelDynamicPanelTop, SX: ControlPanelSX, SY: HouseholdControllerSY}
	PersonToPanel(p, 0, person, IconW)
	cp.SetDynamicPanel(p)
}
