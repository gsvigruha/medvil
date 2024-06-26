package view

import (
	"fmt"
	"github.com/go-gl/gl/v3.2-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/tfriedel6/canvas"
	"github.com/tfriedel6/canvas/backend/goglbackend"
	"image"
	_ "image/gif" // Imported here so that applications based on this package support these formats by default
	_ "image/jpeg"
	_ "image/png"
	"medvil/controller"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

// Window represents the opened window with GL context. The Mouse* and Key*
// functions can be set for callbacks
type Window struct {
	Window     *glfw.Window
	canvas     *canvas.Canvas
	c          *controller.Controller
	frameTimes [10]time.Time
	frameIndex int
	frameCount int
	fps        float32
	close      bool
	MouseDown  func(button, x, y int)
	MouseMove  func(x, y int)
	MouseUp    func(button, x, y int)
	MouseWheel func(x, y int)
	KeyDown    func(scancode int, rn rune, name string)
	KeyUp      func(scancode int, rn rune, name string)
	KeyChar    func(rn rune)
	SizeChange func(w, h int)
}

// CreateWindow creates a window using SDL and initializes the OpenGL context
func CreateWindow(title string) (*Window, *canvas.Canvas, *goglbackend.GLContext, *controller.ViewSettings, error) {
	runtime.LockOSThread()

	// init GLFW
	err := glfw.Init()
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("Error initializing GLFW: %v", err)
	}

	viewSettings := controller.LoadSettings()
	var w, h int
	if viewSettings.Resolution == controller.ResolutionHD {
		w, h = 1280, 720
	} else if viewSettings.Resolution == controller.ResolutionFHD {
		w, h = 1920, 1080
	} else if viewSettings.Resolution == controller.ResolutionQHD {
		w, h = 2560, 1440
	}

	// the stencil size setting is required for the canvas to work
	glfw.WindowHint(glfw.StencilBits, 8)
	glfw.WindowHint(glfw.DepthBits, 0)
	glfw.WindowHint(glfw.Samples, 4)
	glfw.WindowHint(glfw.CocoaRetinaFramebuffer, 1)

	// create window
	var monitor *glfw.Monitor
	if viewSettings.FullScreen {
		monitor = glfw.GetPrimaryMonitor()
	}
	window, err := glfw.CreateWindow(w, h, title, monitor, nil)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("Error creating window: %v", err)
	}
	window.MakeContextCurrent()

	// init GL
	err = gl.Init()
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("Error initializing GL: %v", err)
	}

	// set vsync on, enable multisample (if available)
	glfw.SwapInterval(1)
	gl.Enable(gl.MULTISAMPLE)
	//gl.Enable(gl.SAMPLE_ALPHA_TO_COVERAGE)

	// context
	ctx, err := goglbackend.NewGLContext()
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("Error initializing GLContext: %v", err)
	}

	// load canvas GL backend
	fbw, fbh := window.GetFramebufferSize()
	backend, err := goglbackend.New(0, 0, fbw, fbh, ctx)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("Error loading GoGL backend: %v", err)
	}
	fmt.Println("Frame buffer size: ", fbw, fbh)

	cv := canvas.New(backend)
	wnd := &Window{
		Window: window,
		canvas: cv,
	}

	window.SetSizeCallback(func(w *glfw.Window, width, height int) {
		if wnd.SizeChange != nil {
			wnd.SizeChange(width, height)
		} else {
			fbw, fbh := window.GetFramebufferSize()
			backend.SetBounds(0, 0, fbw, fbh)
		}
		fbw, fbh := window.GetFramebufferSize()
		wnd.c.ControlPanel.SetupDims(fbw, fbh)
	})
	window.SetCloseCallback(func(w *glfw.Window) {
		wnd.c.Save("latest_autosave.mdvl")
		wnd.Close()
	})

	icon, err := os.Open(filepath.FromSlash("icon/gui/house.png"))
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("Error loading icon: %v", err)
	}
	iconImage, _, err := image.Decode(icon)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("Error loading icon: %v", err)
	}
	window.SetIcon([]image.Image{iconImage})

	return wnd, cv, ctx, viewSettings, nil
}

func (wnd *Window) SetController(c *controller.Controller) {
	wnd.c = c
}

func (wnd *Window) GetGLFWWindow() *glfw.Window {
	return wnd.Window
}

// FPS returns the frames per second (averaged over 10 frames)
func (wnd *Window) FPS() float32 {
	return wnd.fps
}

// Close can be used to end a call to MainLoop
func (wnd *Window) Close() {
	wnd.close = true
}

// StartFrame handles events and gets the window ready for rendering
func (wnd *Window) StartFrame() {
	wnd.Window.MakeContextCurrent()
	glfw.PollEvents()
}

// FinishFrame updates the FPS count and displays the frame
func (wnd *Window) FinishFrame() {
	now := time.Now()
	wnd.frameTimes[wnd.frameIndex] = now
	wnd.frameIndex++
	wnd.frameIndex %= len(wnd.frameTimes)
	if wnd.frameCount < len(wnd.frameTimes) {
		wnd.frameCount++
	} else {
		diff := now.Sub(wnd.frameTimes[wnd.frameIndex]).Seconds()
		wnd.fps = float32(wnd.frameCount-1) / float32(diff)
	}

	wnd.Window.SwapBuffers()
}

// MainLoop runs a main loop and calls run on every frame
func (wnd *Window) MainLoop(run func()) {
	for !wnd.close {
		wnd.StartFrame()
		run()
		wnd.FinishFrame()
	}
}

// Size returns the current width and height of the window.
// Note that this size may not be the same as the size of the
// framebuffer, since some operating systems scale the window.
// Use the Width/Height/Size function on Canvas to determine
// the drawing size
func (wnd *Window) Size() (int, int) {
	return wnd.Window.GetSize()
}

// FramebufferSize returns the current width and height of
// the framebuffer, which is also the internal size of the
// canvas
func (wnd *Window) FramebufferSize() (int, int) {
	return wnd.Window.GetFramebufferSize()
}
