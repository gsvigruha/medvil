package controller

import (
	"medvil/model/social"
	"medvil/view/gui"
)

func PersonToControlPanel(cp *ControlPanel, person *social.Person) {
	p := &gui.Panel{X: 0, Y: ControlPanelDynamicPanelTop, SX: ControlPanelSX, SY: HouseholdControllerSY}
	PersonToPanel(cp, p, 0, person, IconW, PersonGUIY*ControlPanelSY)
	cp.SetDynamicPanel(p)
}
