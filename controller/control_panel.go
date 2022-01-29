package controller

import (
	"github.com/tfriedel6/canvas"
	"medvil/view/gui"
	"strconv"
)

type ControlPanel struct {
	P         *gui.Panel
	dateLabel *gui.TextLabel
	c         *Controller
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
	p.dateLabel.Text = strconv.Itoa(int(p.c.Calendar.Day)) + ", " + strconv.Itoa(int(p.c.Calendar.Month)) + ", " + strconv.Itoa(int(p.c.Calendar.Year))
}

func (p *ControlPanel) Setup(c *Controller) {
	p.c = c
	p.P = &gui.Panel{X: 0, Y: 0, SX: ControlPanelSX, SY: ControlPanelSY}
	p.P.Clear()
	p.dateLabel = p.P.AddTextLabel("", 10, 20)
	p.P.AddButton(ControlPanelButton{b: gui.ButtonGUI{Icon: "building", X: 10, Y: 30, SX: 32, SY: 32}, c: c})
}
