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
	cp             *ControlPanel
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
	HouseholdToControlPanel(cp, hp, tower.Household, "tower")
	tc := &TowerController{towerPanel: tp, householdPanel: hp, tower: tower, cp: cp}
	tc.UseType = military.MilitaryLandUseTypeNone

	hcy := HouseholdControllerGUIBottomY * ControlPanelSY
	tp.AddButton(&LandUseButton{
		b:       gui.ButtonGUI{Texture: "terrain/grass", X: float64(10), Y: hcy, SX: IconS, SY: IconS},
		luc:     tc,
		useType: military.MilitaryLandUseTypeNone,
	})
	tp.AddButton(&LandUseButton{
		b:       gui.ButtonGUI{Icon: "artifacts/shield", X: float64(10 + IconW*1), Y: hcy, SX: IconS, SY: IconS},
		luc:     tc,
		useType: military.MilitaryLandUseTypePatrol,
	})

	cp.SetDynamicPanel(tc)
	cp.C.ClickHandler = tc
}

func (tc *TowerController) CaptureMove(x, y float64) {
	tc.householdPanel.CaptureMove(x, y)
	tc.towerPanel.CaptureMove(x, y)
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
	HouseholdToControlPanel(tc.cp, tc.householdPanel, tc.tower.Household, "tower")
	tc.CaptureMove(tc.cp.C.X, tc.cp.C.Y)
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
	if !rf.F.Allocated && tc.UseType == military.MilitaryLandUseTypePatrol && tc.tower.FieldWithinDistance(rf.F) {
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

func (tc *TowerController) GetHelperSuggestions() *gui.Suggestion {
	return nil
}
