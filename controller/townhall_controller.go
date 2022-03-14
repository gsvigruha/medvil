package controller

import (
	"medvil/model/social"
	"medvil/view/gui"
)

func TownhallToControlPanel(cp *ControlPanel, th *social.Townhall) {
	hp := &gui.Panel{X: 0, Y: ControlPanelDynamicPanelTop, SX: ControlPanelSX, SY: HouseholdControllerSY}
	HouseholdToControlPanel(hp, &th.Household)

	hp.AddPanel(gui.CreateNumberPanel(10, 600, 120, 20, 0, 100, 10, "farm tax rate %v", &th.Household.Town.Transfers.Farm.TaxRate).P)
	hp.AddPanel(gui.CreateNumberPanel(10, 625, 120, 20, 0, 1000, 10, "farm threshold %v", &th.Household.Town.Transfers.Farm.TaxThreshold).P)
	hp.AddPanel(gui.CreateNumberPanel(10, 650, 120, 20, 0, 1000, 10, "farm subsidy %v", &th.Household.Town.Transfers.Farm.Subsidy).P)

	hp.AddPanel(gui.CreateNumberPanel(10, 675, 120, 20, 0, 100, 10, "shop tax rate %v", &th.Household.Town.Transfers.Workshop.TaxRate).P)
	hp.AddPanel(gui.CreateNumberPanel(10, 700, 120, 20, 0, 1000, 10, "shop threshold %v", &th.Household.Town.Transfers.Workshop.TaxThreshold).P)
	hp.AddPanel(gui.CreateNumberPanel(10, 725, 120, 20, 0, 1000, 10, "shop subsidy %v", &th.Household.Town.Transfers.Workshop.Subsidy).P)

	hp.AddPanel(gui.CreateNumberPanel(150, 600, 120, 20, 0, 100, 10, "mine tax rate %v", &th.Household.Town.Transfers.Mine.TaxRate).P)
	hp.AddPanel(gui.CreateNumberPanel(150, 625, 120, 20, 0, 1000, 10, "mine threshold %v", &th.Household.Town.Transfers.Mine.TaxThreshold).P)
	hp.AddPanel(gui.CreateNumberPanel(150, 650, 120, 20, 0, 1000, 10, "mine subsidy %v", &th.Household.Town.Transfers.Mine.Subsidy).P)

	hp.AddPanel(gui.CreateNumberPanel(150, 675, 120, 20, 0, 100, 10, "market funding %v", &th.Household.Town.Transfers.MarketFundingRate).P)

	cp.SetDynamicPanel(hp)
}
