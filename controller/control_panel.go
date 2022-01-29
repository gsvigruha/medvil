package controller

import (
	"github.com/tfriedel6/canvas"
	"medvil/view/gui"
	"strconv"
)

const ControlPanelSX = 300
const ControlPanelSY = 700
const ControlPanelDynamicPanelTop = 100
const ControlPanelDynamicPanelSY = 600

type ControlPanel struct {
	topPanel     *gui.Panel
	dynamicPanel Panel
	dateLabel    *gui.TextLabel
	C            *Controller
}

type ControlPanelButton struct {
	b gui.ButtonGUI
	c *Controller
}

func (b ControlPanelButton) Click() {
	b.c.ShowBuildingController()
}

func (b ControlPanelButton) Render(cv *canvas.Canvas) {
	b.b.Render(cv)
}

func (b ControlPanelButton) Contains(x float64, y float64) bool {
	return b.b.Contains(x, y)
}

func (p *ControlPanel) Refresh() {
	p.dateLabel.Text = strconv.Itoa(int(p.C.Calendar.Day)) + ", " + strconv.Itoa(int(p.C.Calendar.Month)) + ", " + strconv.Itoa(int(p.C.Calendar.Year))
}

func (p *ControlPanel) Clear() {
	if p.dynamicPanel != nil {
		p.dynamicPanel.Clear()
	}
}

func (p *ControlPanel) Setup(c *Controller) {
	p.C = c
	p.topPanel = &gui.Panel{X: 0, Y: 0, SX: ControlPanelSX, SY: ControlPanelSY}
	p.dateLabel = p.topPanel.AddTextLabel("", 10, 20)
	p.topPanel.AddButton(ControlPanelButton{b: gui.ButtonGUI{Icon: "building", X: 10, Y: 30, SX: 32, SY: 32}, c: c})
}

func (p *ControlPanel) SetDynamicPanel(dp Panel) {
	p.Clear()
	p.dynamicPanel = dp
}

func (p *ControlPanel) CaptureClick(x, y float64) {
	p.topPanel.CaptureClick(x, y)
	p.dynamicPanel.CaptureClick(x, y)
}

func (p *ControlPanel) Render(cv *canvas.Canvas) {
	p.topPanel.Render(cv)
	if p.dynamicPanel != nil {
		p.dynamicPanel.Render(cv)
	}
}
