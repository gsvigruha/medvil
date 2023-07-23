package main

import (
	"fmt"
	"github.com/pkg/profile"
	"log"
	"medvil/controller"
	"medvil/maps"
	"medvil/view"
	"os"
	"runtime"
	"strconv"
	"time"
)

const (
	sx uint16 = 25
	sy uint16 = 25
)

var PlantFrameRenderTimeNs int64 = 25000000

func init() {
	if val, exists := os.LookupEnv("MEDVIL_FRAME_RENDER_TIME_MS"); exists {
		if time, err := strconv.Atoi(val); err == nil {
			PlantFrameRenderTimeNs = int64(time) * 1000 * 1000
		}
	}
}

func main() {
	defer profile.Start(profile.ProfilePath(".")).Stop()

	wnd, cv, ctx, err := view.CreateWindow(1920, 1080, "Medvil")
	if err != nil {
		panic(err)
	}
	ic := view.NewImageCache(ctx)

	c := controller.Link(wnd.Window, ctx)
	c.Map = maps.NewMap(50, 50)
	c.LinkMap()

	fmt.Println("Init done")
	fmt.Println("CPUs: " + strconv.Itoa(runtime.NumCPU()))

	wnd.MainLoop(func() {
		start := time.Now()
		view.Render(ic, cv, *c.Map, c)
		elapsed := time.Since(start)

		if elapsed.Nanoseconds() < PlantFrameRenderTimeNs {
			time.Sleep(time.Duration(PlantFrameRenderTimeNs-elapsed.Nanoseconds()) * time.Nanosecond)
		}

		c.Refresh()
		ic.Clean()
		if os.Getenv("MEDVIL_VERBOSE") == "1" {
			log.Printf("Rendering took %s (fps %s)", elapsed, wnd.FPS())
			log.Printf("%s", c.Map.Calendar)
		}
		for i := 0; i < c.TimeSpeed; i++ {
			c.Map.Calendar.Tick()
			c.Map.ElapseTime()
		}
		c.RenderTick()
	})
}
