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
	workshopSubPanel    *gui.Panel
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
	wsp := &gui.Panel{X: 0, Y: ControlPanelDynamicPanelTop + HouseholdControllerSY, SX: ControlPanelSX, SY: ControlPanelDynamicPanelSY - HouseholdControllerSY}
	HouseholdToControlPanel(cp, hp, workshop.Household, "workshop")
	wc := &WorkshopController{workshopPanel: wp, workshopSubPanel: wsp, householdPanel: hp, workshop: workshop, cp: cp}

	hcy := HouseholdControllerGUIBottomY * ControlPanelSY
	wp.AddTextLabel("Select manufacture", 24, hcy-IconS/4.0)
	tasks := economy.GetManufactureNames(workshop.Household.Building.Plan.GetExtensions())
	wc.manufactureDropDown = &gui.DropDown{
		X:        float64(24),
		Y:        hcy,
		SX:       IconS + gui.FontSize*10,
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
		X:  IconS + gui.FontSize*10 + 32,
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

	wp.AddPanel(wsp)
	wc.UpdateSubPanel()

	wp.AddButton(&gui.SimpleButton{
		ButtonGUI: gui.ButtonGUI{Icon: "tasks/calculate", X: 24 + IconS + gui.FontSize*10 + LargeIconD, Y: hcy - gui.FontSize/2.0, SX: LargeIconS, SY: LargeIconS, OnHoover: func() {
			cp.HelperMessage("Optimize tasks based on profitability. Needs paper.")
		}},
		Highlight: func() bool { return workshop.AutoSwitch },
		ClickImpl: func() {
			workshop.AutoSwitch = !workshop.AutoSwitch
		}})

	cp.SetDynamicPanel(wc)
}

func (wc *WorkshopController) UpdateSubPanel() {
	hcy := HouseholdControllerGUIBottomY * ControlPanelSY
	if wc.workshop.Manufacture != nil {
		var aI = 0
		for _, a := range wc.workshop.Manufacture.Inputs {
			ArtifactsToControlPanel(wc.cp, wc.workshopSubPanel, aI, a.A, a.Quantity, hcy+LargeIconD)
			aI++
		}
		wc.workshopSubPanel.AddImageLabel("arrow", float64(24+aI*IconW), hcy+LargeIconD, IconS, IconS, gui.ImageLabelStyleRegular)
		aI++
		for _, a := range wc.workshop.Manufacture.Outputs {
			ArtifactsToControlPanel(wc.cp, wc.workshopSubPanel, aI, a.A, a.Quantity, hcy+LargeIconD)
			aI++
		}
	}
}

func (wc *WorkshopController) CaptureMove(x, y float64) {
	wc.householdPanel.CaptureMove(x, y)
	wc.workshopPanel.CaptureMove(x, y)
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
	HouseholdToControlPanel(wc.cp, wc.householdPanel, wc.workshop.Household, "workshop")
	wc.workshopSubPanel.Clear()
	wc.UpdateSubPanel()
	wc.CaptureMove(wc.cp.C.X, wc.cp.C.Y)
}

func (wc *WorkshopController) GetHelperSuggestions() *gui.Suggestion {
	suggestion := GetHouseholdHelperSuggestions(wc.workshop.Household)
	if suggestion != nil {
		return suggestion
	}
	hcy := HouseholdControllerGUIBottomY * ControlPanelSY
	if wc.workshop.Manufacture == nil {
		return &gui.Suggestion{
			Message: "Select goods to produce. Different building extensions,\nlike waterwheels or forges let you produce different\ntypes of goods like flour or metal. You can change\nwhat each workshop produces as the needs of your town evolve.",
			Icon:    "workshop_mixed", X: IconS + gui.FontSize*10 + 32 + float64(IconW), Y: hcy + float64(IconH)/2.0,
		}
	}
	return nil
}
