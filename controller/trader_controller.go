package controller

import (
	"github.com/tfriedel6/canvas"
	"medvil/model/artifacts"
	"medvil/model/economy"
	"medvil/model/navigation"
	"medvil/model/social"
	"medvil/renderer"
	"medvil/view/gui"
)

type TraderController struct {
	traderPanel *gui.Panel
	trader      *social.Trader
	cp          *ControlPanel
}

func TraderToControlPanel(cp *ControlPanel, trader *social.Trader) {
	p := &gui.Panel{X: 0, Y: ControlPanelDynamicPanelTop, SX: ControlPanelSX, SY: HouseholdControllerSY}
	tc := &TraderController{traderPanel: p, trader: trader, cp: cp}
	TraderToPanel(cp, p, trader)
	cp.SetDynamicPanel(tc)
	cp.C.ClickHandler = tc
}

func TraderToPanel(cp *ControlPanel, p *gui.Panel, trader *social.Trader) {
	MoneyToControlPanel(p, trader.SourceExchange.Town, &trader.Money, 100, 10, float64(IconH+50))
	PersonToPanel(cp, p, 0, trader.Person, IconW, PersonGUIY*ControlPanelSY)
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
		TaskToControlPanel(cp, p, i%IconRowMax, TaskGUIY*ControlPanelSY+float64(i/IconRowMax*IconH), task, IconW)
	}
	if trader.Person.Task != nil {
		if tradeTask, ok := trader.Person.Task.(*economy.TradeTask); ok {
			for i, as := range tradeTask.Goods {
				ArtifactsToControlPanel(p, i, as.A, as.Quantity, VehicleGUIY*ControlPanelSY)
			}
		}
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
	TraderToPanel(tc.cp, tc.traderPanel, tc.trader)
}

func (tc *TraderController) GetActiveFields(c *Controller, rf *renderer.RenderedField) []navigation.FieldWithContext {
	var fs []navigation.FieldWithContext
	for _, coords := range tc.trader.SourceExchange.Building.GetBuildingXYs(true) {
		fs = append(fs, c.Map.GetField(coords[0], coords[1]))
	}
	if tc.trader.TargetExchange != nil {
		for _, coords := range tc.trader.TargetExchange.Building.GetBuildingXYs(true) {
			fs = append(fs, c.Map.GetField(coords[0], coords[1]))
		}
	}
	for _, f := range tc.trader.Person.Traveller.GetPathFields(c.Map) {
		fs = append(fs, f)
	}
	return fs
}

func HandleClickForTrader(trader *social.Trader, c *Controller, rf *renderer.RenderedField) bool {
	th := c.ReverseReferences.BuildingToTownhall[rf.F.Building.GetBuilding()]
	if th != nil && th != trader.SourceExchange.Town.Townhall {
		trader.TargetExchange = th.Household.Town.Marketplace
		return true
	}
	mp := c.ReverseReferences.BuildingToMarketplace[rf.F.Building.GetBuilding()]
	if mp != nil && mp != trader.SourceExchange.Town.Marketplace {
		trader.TargetExchange = mp
		return true
	}
	return true
}

func (tc *TraderController) HandleClick(c *Controller, rf *renderer.RenderedField) bool {
	return HandleClickForTrader(tc.trader, c, rf)
}
