package controller

import (
	"medvil/model/social"
	"medvil/view/gui"
)

func TownhallToControlPanel(cp *ControlPanel, th *social.Townhall) {
	hp := &gui.Panel{X: 0, Y: ControlPanelDynamicPanelTop, SX: ControlPanelSX, SY: HouseholdControllerSY}
	HouseholdToControlPanel(hp, &th.Household)
	cp.SetDynamicPanel(hp)
}
