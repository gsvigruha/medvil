package main

import (
	"fmt"
	"github.com/pkg/profile"
	"log"
	"medvil/controller"
	"medvil/maps"
	"medvil/view"
	"os"
	"time"
)

const (
	sx uint16 = 25
	sy uint16 = 25
)

func main() {
	defer profile.Start(profile.ProfilePath(".")).Stop()

	wnd, cv, ctx, err := view.CreateWindow(1280, 720, "Medvil")
	if err != nil {
		panic(err)
	}
	ic := view.NewImageCache(ctx)

	m := maps.LoadMap("samples/map/coast_1")

	c := controller.Link(wnd.Window, &m)

	fmt.Println("Init done")

	wnd.MainLoop(func() {
		start := time.Now()
		view.Render(ic, cv, m, c)
		elapsed := time.Since(start)

		if elapsed.Nanoseconds() < 25000000 {
			time.Sleep(time.Duration(25000000-elapsed.Nanoseconds()) * time.Nanosecond)
		}

		c.Refresh()
		ic.Clean()
		if os.Getenv("MEDVIL_VERBOSE") == "1" {
			log.Printf("Rendering took %s (fps %s)", elapsed, wnd.FPS())
			log.Printf("%s", c.Calendar)
		}
		for i := 0; i < c.TimeSpeed; i++ {
			c.Calendar.Tick()
			m.ElapseTime(c.Calendar)
		}
	})
}
