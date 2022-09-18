package controller

import (
	"github.com/tfriedel6/canvas"
	"medvil/model/artifacts"
	"medvil/model/navigation"
	"medvil/model/social"
	"medvil/renderer"
	"medvil/view/gui"
)

type TraderController struct {
	traderPanel *gui.Panel
	trader      *social.Trader
}

func TraderToControlPanel(cp *ControlPanel, trader *social.Trader) {
	p := &gui.Panel{X: 0, Y: ControlPanelDynamicPanelTop, SX: ControlPanelSX, SY: HouseholdControllerSY}
	tc := &TraderController{traderPanel: p, trader: trader}
	TraderToPanel(p, trader)
	cp.SetDynamicPanel(tc)
}

func TraderToPanel(p *gui.Panel, trader *social.Trader) {
	MoneyToControlPanel(p, trader.SourceExchange.Town, &trader.Money, 100, 10, float64(IconH+50))
	PersonToPanel(p, 0, trader.Person, IconW, PersonGUIY*ControlPanelSY)
	p.AddScaleLabel("heating", 10, ArtifactsGUIY*ControlPanelSY, IconS, IconS, 4, trader.GetHeating(), false)
	p.AddScaleLabel("barrel", 10+float64(IconW), ArtifactsGUIY*ControlPanelSY, IconS, IconS, 4, trader.Resources.UsedVolumeCapacity(), false)
	var aI = 2
	for _, a := range artifacts.All {
		if q, ok := trader.Resources.Artifacts[a]; ok {
			ArtifactsToControlPanel(p, aI, a, q, ArtifactsGUIY*ControlPanelSY)
			aI++
		}
	}
	for i, task := range trader.Tasks {
		if i >= MaxNumTasks {
			break
		}
		TaskToControlPanel(p, i%IconRowMax, TaskGUIY*ControlPanelSY+float64(i/IconRowMax*IconH), task, IconW)
	}
}

func (tc *TraderController) CaptureClick(x, y float64) {
	tc.traderPanel.CaptureClick(x, y)
}

func (tc *TraderController) Render(cv *canvas.Canvas) {
	tc.traderPanel.Render(cv)
}

func (tc *TraderController) Clear() {}

func (tc *TraderController) Refresh() {
	tc.traderPanel.Clear()
	TraderToPanel(tc.traderPanel, tc.trader)
}

func (tc *TraderController) GetActiveFields(c *Controller, rf *renderer.RenderedField) []navigation.FieldWithContext {
	return nil
}

func (tc *TraderController) HandleClick(c *Controller, rf *renderer.RenderedField) bool {
	return false
}
