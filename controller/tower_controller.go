package controller

import (
	"github.com/tfriedel6/canvas"
	"medvil/model/social"
	"medvil/view/gui"
)

type TowerController struct {
	householdPanel *gui.Panel
	towerPanel     *gui.Panel
	tower          *social.Tower
}

func TowerToControlPanel(cp *ControlPanel, tower *social.Tower) {
	hp := &gui.Panel{X: 0, Y: ControlPanelDynamicPanelTop, SX: ControlPanelSX, SY: HouseholdControllerSY}
	tp := &gui.Panel{X: 0, Y: ControlPanelDynamicPanelTop + HouseholdControllerSY, SX: ControlPanelSX, SY: ControlPanelDynamicPanelSY - HouseholdControllerSY}
	HouseholdToControlPanel(hp, &tower.Household)
	tc := &TowerController{towerPanel: tp, householdPanel: hp, tower: tower}

	cp.SetDynamicPanel(tc)
}

func (tc *TowerController) CaptureClick(x, y float64) {
	tc.householdPanel.CaptureClick(x, y)
	tc.towerPanel.CaptureClick(x, y)
}

func (tc *TowerController) Render(cv *canvas.Canvas) {
	tc.householdPanel.Render(cv)
	tc.towerPanel.Render(cv)
}

func (tc *TowerController) Clear() {}

func (tc *TowerController) Refresh() {
	tc.householdPanel.Clear()
	HouseholdToControlPanel(tc.householdPanel, &tc.tower.Household)
}
