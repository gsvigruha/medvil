package main

import (
	"fmt"
	"github.com/pkg/profile"
	"log"
	"math/rand"
	"medvil/controller"
	"medvil/view"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"time"
)

const (
	sx uint16 = 25
	sy uint16 = 25
)

var FrameRenderTimeNs int64 = 25000000

func init() {
	if val, exists := os.LookupEnv("MEDVIL_FRAME_RENDER_TIME_MS"); exists {
		if time, err := strconv.Atoi(val); err == nil {
			FrameRenderTimeNs = int64(time) * 1000 * 1000
		}
	}
}

func main() {
	if os.Getenv("MEDVIL_PROFILE") == "1" {
		// This crashes the Mac app bundle for some reason
		defer profile.Start(profile.ProfilePath(".")).Stop()
	}

	rand.Seed(time.Now().UnixNano())
	wnd, cv, ctx, viewSettings, err := view.CreateWindow("Medville")
	if err != nil {
		panic(err)
	}
	ic := view.NewImageCache(ctx)

	c := controller.Link(wnd, ctx)
	c.ViewSettings = *viewSettings
	controller.LibraryToControlPanel(c.ControlPanel)

	fmt.Println("Init done")
	fmt.Println("CPUs: " + strconv.Itoa(runtime.NumCPU()))

	wnd.MainLoop(func() {
		start := time.Now()
		if c.Map != nil {
			c.MapLock.Lock()
			view.Render(ic, cv, *c.Map, c)
			c.MapLock.Unlock()
			elapsed := time.Since(start)

			if elapsed.Nanoseconds() < FrameRenderTimeNs {
				time.Sleep(time.Duration(FrameRenderTimeNs-elapsed.Nanoseconds()) * time.Nanosecond)
			}

			c.Refresh()
			ic.Clean()
			if os.Getenv("MEDVIL_VERBOSE") == "2" {
				log.Printf("Rendering took %s (fps %s)", elapsed, wnd.FPS())
				log.Printf("%s", c.Map.Calendar)
			}
			for i := 0; i < c.TimeSpeed; i++ {
				c.MapLock.Lock()
				c.Map.Calendar.Tick()
				c.Map.ElapseTime()
				c.MapLock.Unlock()
			}
		} else {
			c.ControlPanel.Render(cv, c)
			cv.DrawImage(filepath.FromSlash("icon/gui/background.png"), controller.ControlPanelSX, 0, float64(cv.Width())-controller.ControlPanelSX, float64(cv.Height()))
		}
		c.RenderTick()
	})
}
