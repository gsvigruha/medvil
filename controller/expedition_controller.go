package controller

import (
	"github.com/tfriedel6/canvas"
	"medvil/model/artifacts"
	"medvil/model/navigation"
	"medvil/model/social"
	"medvil/renderer"
	"medvil/view/gui"
)

type ExpeditionController struct {
	expeditionPanel *gui.Panel
	expedition      *social.Expedition
	cp              *ControlPanel
}

const ExpeditionTaskGUIY = 0.8

func ExpeditionToControlPanel(cp *ControlPanel, expedition *social.Expedition) {
	p := &gui.Panel{X: 0, Y: ControlPanelDynamicPanelTop, SX: ControlPanelSX, SY: HouseholdControllerSY}
	tc := &ExpeditionController{expeditionPanel: p, expedition: expedition, cp: cp}
	ExpeditionToPanel(cp, p, expedition)
	cp.SetDynamicPanel(tc)
	cp.C.ClickHandler = tc
}

func ExpeditionToPanel(cp *ControlPanel, p *gui.Panel, expedition *social.Expedition) {
	MoneyToControlPanel(p, expedition.Town.Townhall.Household, expedition, 100, 10, LargeIconD+float64(IconH)+24)
	for i, person := range expedition.People {
		PersonToPanel(cp, p, i, person, IconW, PersonGUIY*ControlPanelSY)
	}
	for i := len(expedition.People); i < int(expedition.TargetNumPeople); i++ {
		p.AddImageLabel("person", float64(24+i*IconW), PersonGUIY*ControlPanelSY, IconS, IconS, gui.ImageLabelStyleDisabled)
	}
	p.AddButton(&gui.SimpleButton{
		ButtonGUI: gui.ButtonGUI{Icon: "plus", X: ControlPanelSX - 40, Y: PersonGUIY * ControlPanelSY, SX: 16, SY: 16},
		Highlight: func() bool { return false },
		ClickImpl: func() { expedition.IncTargetNumPeople() }})
	p.AddButton(&gui.SimpleButton{
		ButtonGUI: gui.ButtonGUI{Icon: "minus", X: ControlPanelSX - 40, Y: PersonGUIY*ControlPanelSY + 16, SX: 16, SY: 16},
		Highlight: func() bool { return false },
		ClickImpl: func() { expedition.DecTargetNumPeople() }})

	p.AddScaleLabel("barrel", 24, ArtifactsGUIY*ControlPanelSY, IconS, IconS, 4, expedition.Resources.UsedVolumeCapacity(), false)
	var aI = 1
	for _, a := range artifacts.All {
		var q uint16 = 0
		if storageQ, ok := expedition.Resources.Artifacts[a]; ok {
			q = storageQ
		}
		ArtifactStorageToControlPanel(p, expedition.StorageTarget, aI, a, q, ArtifactsGUIY*ControlPanelSY)
		aI++
	}

	for i, task := range expedition.Tasks {
		if i >= MaxNumTasks {
			break
		}
		TaskToControlPanel(cp, p, i%IconRowMax, ExpeditionTaskGUIY*ControlPanelSY+float64(i/IconRowMax*IconH), task, IconW)
	}
	if expedition.DestinationField != nil {
		if expedition.IsEveryoneBoarded() {
			p.AddImageLabel("tasks/goto", 24, ExpeditionTaskGUIY*ControlPanelSY+float64(IconH*2), IconS, IconS, gui.ImageLabelStyleRegular)
		} else {
			p.AddImageLabel("tasks/goto", 24, ExpeditionTaskGUIY*ControlPanelSY+float64(IconH*2), IconS, IconS, gui.ImageLabelStyleDisabled)
		}
	}
}

func (ec *ExpeditionController) CaptureClick(x, y float64) {
	ec.expeditionPanel.CaptureClick(x, y)
}

func (ec *ExpeditionController) Render(cv *canvas.Canvas) {
	ec.expeditionPanel.Render(cv)
}

func (ec *ExpeditionController) Clear() {}

func (ec *ExpeditionController) Refresh() {
	ec.expeditionPanel.Clear()
	ExpeditionToPanel(ec.cp, ec.expeditionPanel, ec.expedition)
}

func (ec *ExpeditionController) GetActiveFields(c *Controller, rf *renderer.RenderedField) []navigation.FieldWithContext {
	var fs []navigation.FieldWithContext
	if ec.expedition.CheckDestinationField(rf.F) {
		fs = append(fs, rf.F)
	} else if rf.F.Building.GetBuilding() != nil && c.ReverseReferences.BuildingToTownhall[rf.F.Building.GetBuilding()] != nil {
		fs = append(fs, rf.F)
	}
	return fs
}

func HandleClickForExpedition(expedition *social.Expedition, c *Controller, rf *renderer.RenderedField) bool {
	if expedition.CheckDestinationField(rf.F) {
		expedition.DestinationField = rf.F
		return true
	} else if rf.F.Building.GetBuilding() != nil && c.ReverseReferences.BuildingToTownhall[rf.F.Building.GetBuilding()] != nil {
		c.ReverseReferences.BuildingToTownhall[rf.F.Building.GetBuilding()].Household.Town.Supplier = expedition
	}
	return false
}

func (ec *ExpeditionController) HandleClick(c *Controller, rf *renderer.RenderedField) bool {
	return HandleClickForExpedition(ec.expedition, c, rf)
}
