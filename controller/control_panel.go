package controller

import (
	"github.com/tfriedel6/canvas"
	"github.com/tfriedel6/canvas/backend/goglbackend"
	"medvil/model/building"
	"medvil/view/gui"
	"reflect"
	"strconv"
)

var ControlPanelSX = 300.0
var ControlPanelSY = 800.0
var ControlPanelDynamicPanelTop = 0.1
var ControlPanelDynamicPanelSY = 0.6

const CPButtonHighlightNone = 0
const CPButtonHighlightSmall = 1
const CPButtonHighlightLarge = 2

type ControlPanel struct {
	topPanel       *gui.Panel
	dynamicPanel   Panel
	helperPanel    *gui.Panel
	dateLabel      *gui.TextLabel
	moneyLabel     *gui.TextLabel
	peopleLabel    *gui.TextLabel
	artifactsLabel *gui.TextLabel
	buildingsLabel *gui.TextLabel
	timeButton     *ControlPanelButton
	C              *Controller
	buffer         *canvas.Canvas
}

type ControlPanelButton struct {
	b         gui.ButtonGUI
	c         *Controller
	action    func(*Controller)
	highlight uint8
}

func (b ControlPanelButton) Click() {
	b.action(b.c)
}

func (b ControlPanelButton) Render(cv *canvas.Canvas) {
	if b.highlight == CPButtonHighlightSmall {
		cv.SetFillStyle("#48C")
		cv.FillRect(b.b.X, b.b.Y, b.b.SX, b.b.SY)
	} else if b.highlight == CPButtonHighlightLarge {
		cv.SetFillStyle("#8AD")
		cv.FillRect(b.b.X, b.b.Y, b.b.SX, b.b.SY)
	}
	b.b.Render(cv)
}

func (b ControlPanelButton) Contains(x float64, y float64) bool {
	return b.b.Contains(x, y)
}

func (b ControlPanelButton) Enabled() bool {
	return b.b.Enabled()
}

func CPActionCancel(c *Controller) {
	c.Reset()
	c.ControlPanel.dynamicPanel = nil
}

func CPActionTimeScaleChange(c *Controller) {
	if c.TimeSpeed == 1 {
		c.TimeSpeed = 5
		c.ControlPanel.timeButton.highlight = CPButtonHighlightSmall
	} else if c.TimeSpeed == 5 {
		c.TimeSpeed = 20
		c.ControlPanel.timeButton.highlight = CPButtonHighlightLarge
	} else {
		c.TimeSpeed = 1
		c.ControlPanel.timeButton.highlight = CPButtonHighlightNone
	}
}

func (p *ControlPanel) Refresh() {
	p.dateLabel.Text = strconv.Itoa(
		int(p.C.Map.Calendar.Day)) + ", " +
		strconv.Itoa(int(p.C.Map.Calendar.Month)) + ", " +
		strconv.Itoa(int(p.C.Map.Calendar.Year))
	stats := p.C.Country.Stats()
	p.moneyLabel.Text = strconv.Itoa(int(stats.Money))
	p.peopleLabel.Text = strconv.Itoa(int(stats.People))
	p.artifactsLabel.Text = strconv.Itoa(int(stats.Artifacts))
	p.buildingsLabel.Text = strconv.Itoa(int(stats.Buildings))
	if p.dynamicPanel != nil {
		p.dynamicPanel.Refresh()
	}
}

func (p *ControlPanel) Clear() {
	if p.dynamicPanel != nil {
		p.dynamicPanel.Clear()
	}
}

func (p *ControlPanel) GetHelperPanel() *gui.Panel {
	p.helperPanel.Clear()
	return p.helperPanel
}

func (p *ControlPanel) Setup(c *Controller, ctx *goglbackend.GLContext) {
	p.C = c
	if c.W < 2000 {
		ControlPanelSX = 400.0
		ControlPanelSY = float64(c.H)
		IconS = 32.0
		IconW = 40
		IconH = 40
		LargeIconS = 48.0
		LargeIconD = 60.0
		gui.FontSize = 16.0
		ScaleBuildingControllerElements(1.0)
	} else {
		ControlPanelSX = 600.0
		ControlPanelSY = float64(c.H)
		IconS = 48.0
		IconW = 60
		IconH = 60
		LargeIconS = 72.0
		LargeIconD = 80.0
		gui.FontSize = 24.0
		ScaleBuildingControllerElements(1.5)
	}

	p.topPanel = &gui.Panel{X: 0, Y: 0, SX: ControlPanelSX, SY: ControlPanelSY}
	p.helperPanel = &gui.Panel{X: 0, Y: ControlPanelSY * 0.95, SX: ControlPanelSX, SY: ControlPanelSY * 0.05}

	offscreen, _ := goglbackend.NewOffscreen(int(ControlPanelSX), int(ControlPanelSY), false, ctx)
	p.buffer = canvas.New(offscreen)
}

func (p *ControlPanel) GenerateButtons() {
	p.topPanel.Clear()
	c := p.C
	ih := 4.0
	th := IconS/2 - gui.FontSize/2
	p.dateLabel = p.topPanel.AddTextLabel("", ControlPanelSX*0.03, th+gui.FontSize)
	p.topPanel.AddImageLabel("coin", ControlPanelSX*0.25, ih, IconS, IconS, gui.ImageLabelStyleRegular)
	p.moneyLabel = p.topPanel.AddTextLabel("", ControlPanelSX*0.25+IconS, th+gui.FontSize)
	p.topPanel.AddImageLabel("person", ControlPanelSX*0.5, ih, IconS, IconS, gui.ImageLabelStyleRegular)
	p.peopleLabel = p.topPanel.AddTextLabel("", ControlPanelSX*0.5+IconS, th+gui.FontSize)
	p.topPanel.AddImageLabel("barrel", ControlPanelSX*0.65, ih, IconS, IconS, gui.ImageLabelStyleRegular)
	p.artifactsLabel = p.topPanel.AddTextLabel("", ControlPanelSX*0.65+IconS, th+gui.FontSize)
	p.topPanel.AddImageLabel("house", ControlPanelSX*0.85, ih, IconS, IconS, gui.ImageLabelStyleRegular)
	p.buildingsLabel = p.topPanel.AddTextLabel("", ControlPanelSX*0.85+IconS, th+gui.FontSize)

	iconTop := 15 + IconS
	p.topPanel.AddButton(gui.SimpleButton{
		ButtonGUI: gui.ButtonGUI{Icon: "house", X: float64(24 + LargeIconD*0), Y: iconTop, SX: LargeIconS, SY: LargeIconS},
		Highlight: func() bool { return p.IsBuildingType() },
		ClickImpl: func() { c.ShowBuildingController() }})
	p.topPanel.AddButton(gui.SimpleButton{
		ButtonGUI: gui.ButtonGUI{Icon: "infra", X: float64(24 + LargeIconD*1), Y: iconTop, SX: LargeIconS, SY: LargeIconS},
		Highlight: func() bool { return p.IsInfraType() },
		ClickImpl: func() { c.ShowInfraController() }})
	p.topPanel.AddButton(gui.SimpleButton{
		ButtonGUI: gui.ButtonGUI{Icon: "town", X: float64(24 + LargeIconD*2), Y: iconTop, SX: LargeIconS, SY: LargeIconS},
		Highlight: func() bool { return p.IsDynamicPanelType("NewTownController") },
		ClickImpl: func() { c.ShowNewTownController() }})
	p.topPanel.AddButton(gui.SimpleButton{
		ButtonGUI: gui.ButtonGUI{Icon: "demolish", X: float64(24 + LargeIconD*3), Y: iconTop, SX: LargeIconS, SY: LargeIconS},
		Highlight: func() bool { return p.IsDynamicPanelType("DemolishController") },
		ClickImpl: func() { c.ShowDemolishController() }})
	p.topPanel.AddButton(gui.SimpleButton{
		ButtonGUI: gui.ButtonGUI{Icon: "library", X: float64(24 + LargeIconD*4), Y: iconTop, SX: LargeIconS, SY: LargeIconS},
		Highlight: func() bool { return p.IsDynamicPanelType("LibraryController") },
		ClickImpl: func() { c.ShowLibraryController() }})
	p.topPanel.AddButton(ControlPanelButton{
		b: gui.ButtonGUI{Icon: "cancel", X: float64(24 + LargeIconD*5), Y: iconTop, SX: LargeIconS, SY: LargeIconS},
		c: c, action: CPActionCancel})
	p.timeButton = &ControlPanelButton{
		b: gui.ButtonGUI{Icon: "time", X: float64(24 + LargeIconD*6), Y: iconTop, SX: LargeIconS, SY: LargeIconS},
		c: c, action: CPActionTimeScaleChange}
	p.topPanel.AddButton(p.timeButton)
}

func (p *ControlPanel) SetDynamicPanel(dp Panel) {
	p.Clear()
	p.dynamicPanel = dp
}

func (p *ControlPanel) CaptureClick(x, y float64) {
	p.topPanel.CaptureClick(x, y)
	if p.dynamicPanel != nil {
		p.dynamicPanel.CaptureClick(x, y)
	}
}

func (p *ControlPanel) Render(cv *canvas.Canvas, c *Controller) {
	if c.RenderCnt == 0 {
		p.topPanel.Render(p.buffer)
		if p.dynamicPanel != nil {
			p.dynamicPanel.Render(p.buffer)
		}
		p.helperPanel.Render(p.buffer)
	}
	cv.DrawImage(p.buffer, 0, 0, ControlPanelSX, ControlPanelSY)
}

func (p *ControlPanel) IsDynamicPanelType(typeName string) bool {
	if p.dynamicPanel == nil {
		return false
	}
	return reflect.TypeOf(p.dynamicPanel).String() == ("*controller." + typeName)
}

func (p *ControlPanel) IsBuildingTypeOf(bt building.BuildingType) bool {
	if p.C.ClickHandler == nil {
		return false
	}
	if bc, ok := p.C.ClickHandler.(*BuildingsController); ok {
		return bc.Plan.BuildingType == bt
	}
	return false
}

func (p *ControlPanel) IsBuildingType() bool {
	if p.C.ClickHandler == nil {
		return false
	}
	_, ok := p.C.ClickHandler.(*BuildingsController)
	return ok
}

func (p *ControlPanel) IsInfraType() bool {
	if p.C.ClickHandler == nil {
		return false
	}
	_, ok := p.C.ClickHandler.(*InfraController)
	return ok
}

func (p *ControlPanel) HelperMessage(msg string) {
	p.GetHelperPanel().AddTextLabel(msg, ControlPanelSX*0.03, ControlPanelSY*0.95)
}
