package controller

import (
	"github.com/tfriedel6/canvas"
	"github.com/tfriedel6/canvas/backend/goglbackend"
	"medvil/model/building"
	"medvil/view/gui"
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

func CPActionShowWorkshopController(c *Controller) {
	c.ShowBuildingController(building.BuildingTypeWorkshop)
}

func CPActionShowFarmController(c *Controller) {
	c.ShowBuildingController(building.BuildingTypeFarm)
}

func CPActionShowMineController(c *Controller) {
	c.ShowBuildingController(building.BuildingTypeMine)
}

func CPActionShowFactoryController(c *Controller) {
	c.ShowBuildingController(building.BuildingTypeFactory)
}

func CPActionShowInfraController(c *Controller) {
	c.ShowInfraController()
}

func CPActionShowNewTownController(c *Controller) {
	c.ShowNewTownController()
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
	p.dateLabel.Text = strconv.Itoa(int(p.C.Calendar.Day)) + ", " + strconv.Itoa(int(p.C.Calendar.Month)) + ", " + strconv.Itoa(int(p.C.Calendar.Year))
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

func (p *ControlPanel) Setup(c *Controller, ctx *goglbackend.GLContext) {
	p.C = c
	p.topPanel = &gui.Panel{X: 0, Y: 0, SX: ControlPanelSX, SY: ControlPanelSY}
	p.dateLabel = p.topPanel.AddTextLabel("", 10, 20)
	p.topPanel.AddImageLabel("coin", 80, 8, 16, 16, gui.ImageLabelStyleRegular)
	p.moneyLabel = p.topPanel.AddTextLabel("", 100, 20)
	p.topPanel.AddImageLabel("person", 150, 8, 16, 16, gui.ImageLabelStyleRegular)
	p.peopleLabel = p.topPanel.AddTextLabel("", 170, 20)
	p.topPanel.AddImageLabel("barrel", 200, 8, 16, 16, gui.ImageLabelStyleRegular)
	p.artifactsLabel = p.topPanel.AddTextLabel("", 220, 20)
	p.topPanel.AddImageLabel("workshop", 260, 8, 16, 16, gui.ImageLabelStyleRegular)
	p.buildingsLabel = p.topPanel.AddTextLabel("", 280, 20)
	p.topPanel.AddButton(ControlPanelButton{b: gui.ButtonGUI{Icon: "farm", X: 10, Y: 30, SX: 32, SY: 32}, c: c, action: CPActionShowFarmController})
	p.topPanel.AddButton(ControlPanelButton{b: gui.ButtonGUI{Icon: "mine", X: 50, Y: 30, SX: 32, SY: 32}, c: c, action: CPActionShowMineController})
	p.topPanel.AddButton(ControlPanelButton{b: gui.ButtonGUI{Icon: "workshop", X: 90, Y: 30, SX: 32, SY: 32}, c: c, action: CPActionShowWorkshopController})
	p.topPanel.AddButton(ControlPanelButton{b: gui.ButtonGUI{Icon: "factory", X: 130, Y: 30, SX: 32, SY: 32}, c: c, action: CPActionShowFactoryController})
	p.topPanel.AddButton(ControlPanelButton{b: gui.ButtonGUI{Icon: "infra", X: 170, Y: 30, SX: 32, SY: 32}, c: c, action: CPActionShowInfraController})
	p.topPanel.AddButton(ControlPanelButton{b: gui.ButtonGUI{Icon: "town", X: 210, Y: 30, SX: 32, SY: 32}, c: c, action: CPActionShowNewTownController})
	p.topPanel.AddButton(ControlPanelButton{b: gui.ButtonGUI{Icon: "cancel", X: 250, Y: 30, SX: 32, SY: 32}, c: c, action: CPActionCancel})
	p.timeButton = &ControlPanelButton{b: gui.ButtonGUI{Icon: "time", X: 250, Y: 70, SX: 32, SY: 32}, c: c, action: CPActionTimeScaleChange}
	p.topPanel.AddButton(p.timeButton)

	offscreen, _ := goglbackend.NewOffscreen(int(ControlPanelSX), int(ControlPanelSY), false, ctx)
	p.buffer = canvas.New(offscreen)
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
	}
	cv.DrawImage(p.buffer, 0, 0, ControlPanelSX, ControlPanelSY)
}
