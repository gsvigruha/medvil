package controller

import (
	"medvil/view/gui"
	"strconv"
)

func SetupControlPanel(p *gui.Panel, c *Controller) {
	p.Clear()
	dateStr := strconv.Itoa(int(c.Calendar.Day)) + ", " + strconv.Itoa(int(c.Calendar.Month)) + ", " + strconv.Itoa(int(c.Calendar.Year))
	p.AddTextLabel(dateStr, 10, 20)
}
