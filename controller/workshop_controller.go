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
	cp                  *ControlPanel
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
	HouseholdToControlPanel(cp, hp, workshop.Household)
	wc := &WorkshopController{workshopPanel: wp, householdPanel: hp, workshop: workshop, cp: cp}

	hcy := HouseholdControllerGUIBottomY * ControlPanelSY
	tasks := economy.GetManufactureNames(workshop.Household.Building.Plan.GetExtensions())
	wc.manufactureDropDown = &gui.DropDown{
		X:        float64(24),
		Y:        hcy,
		SX:       IconS + gui.FontSize*16,
		SY:       IconS,
		Options:  tasks,
		Icons:    toTaskNames(tasks),
		Selected: -1,
	}
	if workshop.Manufacture != nil {
		wc.manufactureDropDown.SetSelectedValue(workshop.Manufacture.Name)
	}
	wp.AddDropDown(wc.manufactureDropDown)

	wp.AddLabel(&gui.DynamicImageLabel{
		X:  IconS + gui.FontSize*16 + 32,
		Y:  hcy,
		SX: IconS,
		SY: IconS,
		Icon: func() string {
			if workshop.IsManufactureProfitable() {
				return "profitable"
			} else {
				return "not_profitable"
			}
		},
	})

	wp.AddButton(gui.SimpleButton{
		ButtonGUI: gui.ButtonGUI{Icon: "tasks/calculate", X: 24, Y: hcy + LargeIconD, SX: LargeIconS, SY: LargeIconS},
		Highlight: func() bool { return workshop.AutoSwitch },
		ClickImpl: func() { workshop.AutoSwitch = !workshop.AutoSwitch }})
	wp.AddTextLabel("pick most profitable", 24+LargeIconD, hcy+LargeIconD*1.5)

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
	HouseholdToControlPanel(wc.cp, wc.householdPanel, wc.workshop.Household)
}
