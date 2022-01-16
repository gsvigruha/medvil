package controller

import (
	//"fmt"
	"github.com/go-gl/glfw/v3.3/glfw"
	"medvil/model"
	"medvil/model/navigation"
	"medvil/model/social"
	"medvil/model/time"
	"medvil/renderer"
	"medvil/view/gui"
)

const ControlPanelSX = 300

const PerspectiveNE uint8 = 0
const PerspectiveSE uint8 = 1
const PerspectiveSW uint8 = 2
const PerspectiveNW uint8 = 3

type Controller struct {
	X                     float64
	Y                     float64
	W                     int
	H                     int
	ScrollX               int
	ScrollY               int
	Perspective           uint8
	Calendar              *time.CalendarType
	RenderedFields        []*renderer.RenderedField
	RenderedBuildingUnits []*renderer.RenderedBuildingUnit
	SelectedField         *navigation.Field
	SelectedHousehold     *social.Household
	ReverseReferences     *model.ReverseReferences
	ControlPanel          *gui.Panel
}

func (c *Controller) KeyboardCallback(wnd *glfw.Window, key glfw.Key, code int, action glfw.Action, mod glfw.ModifierKey) {
	if key == glfw.KeyEnter && action == glfw.Release {
		c.Perspective = (c.Perspective + 1) % 4
	}
	if action == glfw.Press {
		if key == glfw.KeyUp {
			c.ScrollY -= 256
		}
		if key == glfw.KeyDown {
			c.ScrollY += 256
		}
		if key == glfw.KeyLeft {
			c.ScrollX -= 256
		}
		if key == glfw.KeyRight {
			c.ScrollX += 256
		}
	}
}

func (c *Controller) MouseButtonCallback(wnd *glfw.Window, button glfw.MouseButton, action glfw.Action, mod glfw.ModifierKey) {
	if action == glfw.Press && button == glfw.MouseButton1 {
		for i := range c.RenderedBuildingUnits {
			rbu := c.RenderedBuildingUnits[i]
			if rbu.Contains(c.X, c.Y) {
				c.SelectedHousehold = c.ReverseReferences.BuildingToHousehold[rbu.Unit.B]
				c.SelectedField = nil
				HouseholdToControlPanel(c.ControlPanel, c.SelectedHousehold)
				return
			}
		}
		for i := range c.RenderedFields {
			rf := c.RenderedFields[i]
			if rf.Contains(c.X, c.Y) {
				c.SelectedField = rf.F
				c.SelectedHousehold = nil
				return
			}
		}
	}
}

func (c *Controller) MouseMoveCallback(wnd *glfw.Window, x float64, y float64) {
	c.X = x
	c.Y = y
}

func (c *Controller) MouseScrollCallback(wnd *glfw.Window, x float64, y float64) {
}

func Link(wnd *glfw.Window) *Controller {
	W, H := wnd.GetSize()
	Calendar := &time.CalendarType{
		Year:  1000,
		Month: 2,
		Day:   1,
		Hour:  0,
	}
	controlPanel := &gui.Panel{X: 0, Y: 0, SX: ControlPanelSX, SY: float64(H)}
	C := Controller{H: H, W: W, Calendar: Calendar, ControlPanel: controlPanel}
	wnd.SetKeyCallback(C.KeyboardCallback)
	wnd.SetMouseButtonCallback(C.MouseButtonCallback)
	wnd.SetCursorPosCallback(C.MouseMoveCallback)
	wnd.SetScrollCallback(C.MouseScrollCallback)
	return &C
}

func (c *Controller) UpdateReverseReferences(rr *model.ReverseReferences) {
	c.ReverseReferences = rr
}

func (c *Controller) ResetRenderedObjects() {
	c.RenderedFields = []*renderer.RenderedField{}
	c.RenderedBuildingUnits = []*renderer.RenderedBuildingUnit{}
}

func (c *Controller) AddRenderedField(rf *renderer.RenderedField) {
	c.RenderedFields = append(c.RenderedFields, rf)
}

func (c *Controller) AddRenderedBuildingUnit(rbu *renderer.RenderedBuildingUnit) {
	c.RenderedBuildingUnits = append(c.RenderedBuildingUnits, rbu)
}
