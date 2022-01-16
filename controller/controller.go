package controller

import (
	"github.com/go-gl/glfw/v3.3/glfw"
	"medvil/model/time"
	"medvil/model/navigation"
	"medvil/renderer"
	"fmt"
)

const PerspectiveNE uint8 = 0
const PerspectiveSE uint8 = 1
const PerspectiveSW uint8 = 2
const PerspectiveNW uint8 = 3

type Controller struct {
	X              float64
	Y              float64
	W              int
	H              int
	ScrollX        int
	ScrollY        int
	Perspective    uint8
	Calendar       *time.CalendarType
	RenderedFields []*renderer.RenderedField
	SelectedField *navigation.Field
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
		for i := range c.RenderedFields {
			f := c.RenderedFields[i]
			if f.Contains(c.X, c.Y) {
				c.SelectedField = f.F
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
	C := Controller{H: H, W: W, Calendar: Calendar}
	wnd.SetKeyCallback(C.KeyboardCallback)
	wnd.SetMouseButtonCallback(C.MouseButtonCallback)
	wnd.SetCursorPosCallback(C.MouseMoveCallback)
	wnd.SetScrollCallback(C.MouseScrollCallback)
	return &C
}
