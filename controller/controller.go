package controller

import (
	"github.com/go-gl/glfw/v3.3/glfw"
)

var X int
var Y int
var W int
var H int

var ScrollX int = 0
var ScrollY int = 0

const PerspectiveNE uint8 = 0
const PerspectiveSE uint8 = 1
const PerspectiveSW uint8 = 2
const PerspectiveNW uint8 = 3

var Perspective uint8 = PerspectiveNE

// Interface for keyboard input events.
// Implement and register to receive keyboard input.
type KeyboardListener interface {
	OnKeyEvent(glfw.Key, int, glfw.Action, glfw.ModifierKey)
}

// Interface for mouse input events.
// Implement and register to receive mouse input.
type MouseListener interface {
	OnMouseButton(glfw.MouseButton, glfw.Action, glfw.ModifierKey)
	OnMouseMove(float64, float64)
	OnMouseScroll(float64, float64)
}

func KeyboardCallback(wnd *glfw.Window, key glfw.Key, code int, action glfw.Action, mod glfw.ModifierKey) {
	if key == glfw.KeyEnter && action == glfw.Release {
		Perspective = (Perspective + 1) % 4
	}
	if action == glfw.Press {
		if key == glfw.KeyUp {
			ScrollY -= 128
		}
		if key == glfw.KeyDown {
			ScrollY += 128
		}
		if key == glfw.KeyLeft {
			ScrollX -= 128
		}
		if key == glfw.KeyRight {
			ScrollX += 128
		}
	}
}

func MouseButtonCallback(wnd *glfw.Window, button glfw.MouseButton, action glfw.Action, mod glfw.ModifierKey) {
}

func MouseMoveCallback(wnd *glfw.Window, x float64, y float64) {
	X = int(x)
	Y = int(y)
}

func MouseScrollCallback(wnd *glfw.Window, x float64, y float64) {
}

func Link(wnd *glfw.Window) {
	W, H = wnd.GetSize()
	wnd.SetKeyCallback(KeyboardCallback)
	wnd.SetMouseButtonCallback(MouseButtonCallback)
	wnd.SetCursorPosCallback(MouseMoveCallback)
	wnd.SetScrollCallback(MouseScrollCallback)
}
