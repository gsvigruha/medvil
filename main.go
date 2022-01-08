package main

import (
	"fmt"
	"log"
	"time"

	"medvil/controller"
	"medvil/maps"
	"medvil/view"

	//"github.com/pkg/profile"
	"github.com/tfriedel6/canvas/glfwcanvas"
)

const (
	sx uint16 = 25
	sy uint16 = 25
)

func main() {
	wnd, cv, err := glfwcanvas.CreateWindow(1280, 720, "Medvil")
	if err != nil {
		panic(err)
	}

	controller.Link(wnd.Window)

	m := maps.NewMap(sx, sy)

	fmt.Println("Init done")

	wnd.MainLoop(func() {
		w, h := float64(cv.Width()), float64(cv.Height())
		cv.SetFillStyle("#000")
		cv.FillRect(0, 0, w, h)
		start := time.Now()
		view.Render(cv, m)
		elapsed := time.Since(start)
		/*
			if elapsed.Nanoseconds() < 50000000 {
			    time.Sleep(30000000 * time.Nanosecond)
			}
		*/
		if 0 == 1 {
			log.Printf("Rendering took %s", elapsed)
			log.Printf("%s", controller.Calendar)
		}
		for i := 1; i < 10; i++ {
			controller.Calendar.Tick()
		}

	})
}
