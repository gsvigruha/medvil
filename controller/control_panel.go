package controller

import (
	"github.com/tfriedel6/canvas"
	"medvil/view/gui"
	"strconv"
)

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

func SetupControlPanel(p *gui.Panel, c *Controller) {
	p.Clear()
	dateStr := strconv.Itoa(int(c.Calendar.Day)) + ", " + strconv.Itoa(int(c.Calendar.Month)) + ", " + strconv.Itoa(int(c.Calendar.Year))
	p.AddTextLabel(dateStr, 10, 20)
	p.AddButton(ControlPanelButton{b: gui.ButtonGUI{Icon: "building", X: 10, Y: 30, SX: 32, SY: 32}, c: c})
}
