package controller

import (
	"github.com/tfriedel6/canvas"
	"medvil/model/economy"
	"medvil/model/social"
	"medvil/view/gui"
)

type WorkshopController struct {
	householdPanel      *gui.Panel
	workshopPanel       *gui.Panel
	workshop            *social.Workshop
	manufactureDropDown *gui.DropDown
}

func toTaskNames(names []string) []string {
	var taskNames []string
	for _, name := range names {
		taskNames = append(taskNames, "tasks/"+name)
	}
	return taskNames
}

func WorkshopToControlPanel(cp *ControlPanel, workshop *social.Workshop) {
	hp := &gui.Panel{X: 0, Y: ControlPanelDynamicPanelTop, SX: ControlPanelSX, SY: HouseholdControllerSY}
	wp := &gui.Panel{X: 0, Y: ControlPanelDynamicPanelTop + HouseholdControllerSY, SX: ControlPanelSX, SY: ControlPanelDynamicPanelSY - HouseholdControllerSY}
	HouseholdToControlPanel(hp, &workshop.Household)
	wc := &WorkshopController{workshopPanel: wp, householdPanel: hp, workshop: workshop}

	wc.manufactureDropDown = &gui.DropDown{
		X:        float64(10),
		Y:        float64(HouseholdControllerGUIBottomY),
		SX:       128,
		SY:       20,
		Options:  economy.GetAllManufactureNames(),
		Icons:    toTaskNames(economy.GetAllManufactureNames()),
		Selected: -1,
	}
	if workshop.Manufacture != nil {
		wc.manufactureDropDown.SetSelectedValue(workshop.Manufacture.Name)
	}
	wp.AddDropDown(wc.manufactureDropDown)

	cp.SetDynamicPanel(wc)
}

func (wc *WorkshopController) CaptureClick(x, y float64) {
	wc.householdPanel.CaptureClick(x, y)
	wc.workshopPanel.CaptureClick(x, y)
	wc.workshop.Manufacture = economy.GetManufacture(wc.manufactureDropDown.GetSelectedValue())
}

func (wc *WorkshopController) Render(cv *canvas.Canvas) {
	wc.householdPanel.Render(cv)
	wc.workshopPanel.Render(cv)
}

func (wc *WorkshopController) Clear() {}

func (wc *WorkshopController) Refresh() {
	wc.householdPanel.Clear()
	HouseholdToControlPanel(wc.householdPanel, &wc.workshop.Household)
}
