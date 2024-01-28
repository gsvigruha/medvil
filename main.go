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
	"runtime/debug"
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
	defer func() {
		if r := recover(); r != nil {
			crashMsg := fmt.Sprintf("%v", r)
			crashLog := string(debug.Stack())
			fmt.Println(crashMsg)
			fmt.Println(crashLog)
			f, _ := os.Create("crash.log")
			f.Write([]byte(crashMsg))
			f.Write([]byte(crashLog))
			f.Close()
		}
	}()

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
		if c.Map != nil {
			start := time.Now()
			c.MapLock.Lock()
			view.Render(ic, cv, *c.Map, c)
			c.MapLock.Unlock()

			c.Refresh()
			ic.Clean()

			for i := 0; i < c.TimeSpeed; i++ {
				c.MapLock.Lock()
				c.Map.Calendar.Tick()
				c.Map.ElapseTime()
				c.MapLock.Unlock()
			}

			elapsed := time.Since(start)
			if os.Getenv("MEDVIL_VERBOSE") == "2" {
				log.Printf("Cycle took %s (fps %s)", elapsed, wnd.FPS())
				log.Printf("%s", c.Map.Calendar)
			}
			if elapsed.Nanoseconds() < FrameRenderTimeNs {
				time.Sleep(time.Duration(FrameRenderTimeNs-elapsed.Nanoseconds()) * time.Nanosecond)
			}
		} else {
			cv.DrawImage(filepath.FromSlash("icon/gui/background.png"), controller.ControlPanelSX, 0, float64(cv.Width())-controller.ControlPanelSX, float64(cv.Height()))
			c.ControlPanel.Render(cv, c)
		}
		c.RenderTick()
	})
}
