package controller

import (
	"github.com/tfriedel6/canvas"
	"medvil/model/social"
	"medvil/view/gui"
)

type TownhallController struct {
	householdPanel *gui.Panel
	taxPanel       *gui.Panel
	th             *social.Townhall
}

func TownhallToControlPanel(cp *ControlPanel, th *social.Townhall) {
	hp := &gui.Panel{X: 0, Y: ControlPanelDynamicPanelTop, SX: ControlPanelSX, SY: HouseholdControllerSY}
	tp := &gui.Panel{X: 0, Y: 600, SX: 300, SY: 100}
	tc := &TownhallController{householdPanel: hp, th: th, taxPanel: tp}

	HouseholdToControlPanel(hp, &th.Household)
	hp.AddPanel(tp)

	tp.AddPanel(gui.CreateNumberPanel(10, 600, 120, 20, 0, 100, 10, "farm tax rate %v", &th.Household.Town.Transfers.Farm.TaxRate).P)
	tp.AddPanel(gui.CreateNumberPanel(10, 625, 120, 20, 0, 1000, 10, "farm threshold %v", &th.Household.Town.Transfers.Farm.TaxThreshold).P)
	tp.AddPanel(gui.CreateNumberPanel(10, 650, 120, 20, 0, 1000, 10, "farm subsidy %v", &th.Household.Town.Transfers.Farm.Subsidy).P)

	tp.AddPanel(gui.CreateNumberPanel(10, 675, 120, 20, 0, 100, 10, "shop tax rate %v", &th.Household.Town.Transfers.Workshop.TaxRate).P)
	tp.AddPanel(gui.CreateNumberPanel(10, 700, 120, 20, 0, 1000, 10, "shop threshold %v", &th.Household.Town.Transfers.Workshop.TaxThreshold).P)
	tp.AddPanel(gui.CreateNumberPanel(10, 725, 120, 20, 0, 1000, 10, "shop subsidy %v", &th.Household.Town.Transfers.Workshop.Subsidy).P)

	tp.AddPanel(gui.CreateNumberPanel(150, 600, 120, 20, 0, 100, 10, "mine tax rate %v", &th.Household.Town.Transfers.Mine.TaxRate).P)
	tp.AddPanel(gui.CreateNumberPanel(150, 625, 120, 20, 0, 1000, 10, "mine threshold %v", &th.Household.Town.Transfers.Mine.TaxThreshold).P)
	tp.AddPanel(gui.CreateNumberPanel(150, 650, 120, 20, 0, 1000, 10, "mine subsidy %v", &th.Household.Town.Transfers.Mine.Subsidy).P)

	tp.AddPanel(gui.CreateNumberPanel(150, 675, 120, 20, 0, 100, 10, "market funding %v", &th.Household.Town.Transfers.MarketFundingRate).P)

	cp.SetDynamicPanel(tc)
}

func (tc *TownhallController) CaptureClick(x, y float64) {
	tc.householdPanel.CaptureClick(x, y)
}

func (tc *TownhallController) Render(cv *canvas.Canvas) {
	tc.householdPanel.Render(cv)
}

func (tc *TownhallController) Clear() {}

func (tc *TownhallController) Refresh() {
	tc.householdPanel.Clear()
	HouseholdToControlPanel(tc.householdPanel, &tc.th.Household)
	tc.householdPanel.AddPanel(tc.taxPanel)
}
