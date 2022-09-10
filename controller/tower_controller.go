package controller

import (
	"github.com/tfriedel6/canvas"
	"medvil/model/military"
	"medvil/model/navigation"
	"medvil/model/social"
	"medvil/renderer"
	"medvil/view/gui"
)

type TowerController struct {
	householdPanel *gui.Panel
	towerPanel     *gui.Panel
	tower          *social.Tower
	UseType        uint8
}

func (tc *TowerController) GetUseType() uint8 {
	return tc.UseType
}

func (tc *TowerController) SetUseType(ut uint8) {
	tc.UseType = ut
}

func TowerToControlPanel(cp *ControlPanel, tower *social.Tower) {
	hp := &gui.Panel{X: 0, Y: ControlPanelDynamicPanelTop, SX: ControlPanelSX, SY: HouseholdControllerSY}
	tp := &gui.Panel{X: 0, Y: ControlPanelDynamicPanelTop + HouseholdControllerSY, SX: ControlPanelSX, SY: ControlPanelDynamicPanelSY - HouseholdControllerSY}
	HouseholdToControlPanel(hp, &tower.Household)
	tc := &TowerController{towerPanel: tp, householdPanel: hp, tower: tower}
	tc.UseType = military.MilitaryLandUseTypeNone

	hcy := HouseholdControllerGUIBottomY * ControlPanelSY
	tp.AddButton(LandUseButton{
		b:       gui.ButtonGUI{Texture: "terrain/grass", X: float64(10), Y: hcy, SX: 32, SY: 32},
		luc:     tc,
		useType: military.MilitaryLandUseTypeNone,
	})
	tp.AddButton(LandUseButton{
		b:       gui.ButtonGUI{Icon: "artifacts/shield", X: float64(50), Y: hcy, SX: 32, SY: 32},
		luc:     tc,
		useType: military.MilitaryLandUseTypePatrol,
	})

	cp.SetDynamicPanel(tc)
	cp.C.ClickHandler = tc
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

func (tc *TowerController) GetActiveFields(c *Controller, rf *renderer.RenderedField) []navigation.FieldWithContext {
	return tc.tower.GetFields()
}

func (tc *TowerController) HandleClick(c *Controller, rf *renderer.RenderedField) bool {
	for i := range tc.tower.Land {
		l := &tc.tower.Land[i]
		if l.F.X == rf.F.X && l.F.Y == rf.F.Y {
			if tc.UseType == military.MilitaryLandUseTypeNone {
				// Disallocate land
				tc.tower.Land = append(tc.tower.Land[:i], tc.tower.Land[i+1:]...)
				rf.F.Allocated = false
			}
			return true
		}
	}
	if !rf.F.Allocated && tc.UseType == military.MilitaryLandUseTypePatrol {
		tc.tower.Land = append(tc.tower.Land, social.PatrolLand{
			X: rf.F.X,
			Y: rf.F.Y,
			F: rf.F,
		})
		rf.F.Allocated = true
		return true
	}
	return false
}
