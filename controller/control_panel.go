package controller

import (
	"medvil/view/gui"
	"strconv"
)


func showBuildingControllerClick(i interface{}) {
	c := i.(*Controller)
	c.ShowBuildingController()
}

func SetupControlPanel(p *gui.Panel, c *Controller) {
	p.Clear()
	dateStr := strconv.Itoa(int(c.Calendar.Day)) + ", " + strconv.Itoa(int(c.Calendar.Month)) + ", " + strconv.Itoa(int(c.Calendar.Year))
	p.AddTextLabel(dateStr, 10, 20)
	p.AddButton("building", 30, 20, 32, 32, showBuildingControllerClick)
}
