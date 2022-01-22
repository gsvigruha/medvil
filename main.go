package main

import (
	"fmt"
	"log"
	"medvil/controller"
	"medvil/maps"
	"medvil/view"
	"time"
)

const (
	sx uint16 = 25
	sy uint16 = 25
)

func main() {
	wnd, cv, ctx, _ := view.CreateWindow(1280, 720, "Medvil")
	ic := view.NewImageCache(ctx)

	c := controller.Link(wnd.Window)

	m := maps.LoadMap("samples/map/coast_1")

	fmt.Println("Init done")

	wnd.MainLoop(func() {
		w, h := float64(cv.Width()), float64(cv.Height())
		cv.SetFillStyle("#000")
		cv.FillRect(0, 0, w, h)
		start := time.Now()
		c.UpdateReverseReferences(m.ReverseReferences())
		view.Render(ic, cv, m, c)
		elapsed := time.Since(start)
		/*
			if elapsed.Nanoseconds() < 50000000 {
			    time.Sleep(30000000 * time.Nanosecond)
			}
		*/
		ic.Clean()
		if 0 == 1 {
			log.Printf("Rendering took %s (fps %s)", elapsed, wnd.FPS())
			log.Printf("%s", c.Calendar)
		}
		for i := 1; i < 2; i++ {
			c.Calendar.Tick()
			m.ElapseTime(c.Calendar)
		}
	})
}
