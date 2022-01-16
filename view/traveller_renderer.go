package view

import (
	"github.com/tfriedel6/canvas"
	"medvil/controller"
	"medvil/model/navigation"
	"medvil/renderer"
)

const MaxPX = navigation.MaxPX
const MaxPY = navigation.MaxPY

func RenderTravellers(cv *canvas.Canvas, travellers []*navigation.Traveller, rf renderer.RenderedField, c *controller.Controller) {
	for i := range travellers {
		t := travellers[i]
		px := float64(t.PX)
		py := float64(t.PY)
		NEPX := rf.X[(2+c.Perspective)%4]
		NEPY := rf.Y[(2+c.Perspective)%4]
		SEPX := rf.X[(1+c.Perspective)%4]
		SEPY := rf.Y[(1+c.Perspective)%4]
		SWPX := rf.X[(0+c.Perspective)%4]
		SWPY := rf.Y[(0+c.Perspective)%4]
		NWPX := rf.X[(3+c.Perspective)%4]
		NWPY := rf.Y[(3+c.Perspective)%4]
		x := (NWPX*(MaxPX-px)*(MaxPY-py) +
			SWPX*(MaxPX-px)*py +
			NEPX*px*(MaxPY-py) +
			SEPX*px*py) / (MaxPX * MaxPY)
		y := (NWPY*(MaxPX-px)*(MaxPY-py) +
			SWPY*(MaxPX-px)*py +
			NEPY*px*(MaxPY-py) +
			SEPY*px*py) / (MaxPX * MaxPY)
		cv.DrawImage("texture/social/person.png", x-16, y-32, 32, 32)
	}
}
