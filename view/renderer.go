package view

import (
	//"fmt"
	"github.com/tfriedel6/canvas"
	"image/color"
	"medvil/controller"
	"medvil/model"
	"medvil/renderer"
)

const (
	DX      float64 = 60.0
	DY      float64 = 40.0
	DZ      float64 = 15.0
	ViewSX  uint8   = 12
	ViewSY  uint8   = 10
	RadiusI int     = 20
)

func Render(ic *ImageCache, cv *canvas.Canvas, m model.Map, c *controller.Controller) {
	w := float64(cv.Width())
	h := float64(cv.Height())
	li := int(c.CenterX) - RadiusI
	hi := int(c.CenterX) + RadiusI
	lj := int(c.CenterY) - RadiusI
	hj := int(c.CenterY) + RadiusI
	for i := 0; i < hi-li; i++ {
		for j := 0; j < hj-lj; j++ {
			var pi, pj int
			switch c.Perspective {
			case controller.PerspectiveNE:
				pi, pj = i+li, hj-1-j
			case controller.PerspectiveSE:
				pi, pj = j+li, i+lj
			case controller.PerspectiveSW:
				pi, pj = hi-1-i, j+lj
			case controller.PerspectiveNW:
				pi, pj = hi-1-j, hj-1-i
			}
			if pi < 0 || pj < 0 || pi >= int(m.SX) || pj >= int(m.SY) {
				continue
			}
			var f = &m.Fields[pi][pj]
			var t = uint8(0)
			var r = uint8(0)
			var b = uint8(0)
			var l = uint8(0)
			switch c.Perspective {
			case controller.PerspectiveNE:
				t = f.SW
				r = f.NW
				b = f.NE
				l = f.SE
			case controller.PerspectiveSE:
				t = f.NW
				r = f.NE
				b = f.SE
				l = f.SW
			case controller.PerspectiveSW:
				t = f.NE
				r = f.SE
				b = f.SW
				l = f.NW
			case controller.PerspectiveNW:
				t = f.SE
				r = f.SW
				b = f.NW
				l = f.NE
			}
			x := w/2 - float64(i)*DX + float64(j)*DX
			y := float64(i)*DY + float64(j)*DY - float64(RadiusI)*DY*2 + h/2
			if x < controller.ControlPanelSX-DX || x > w+DX || y < -DY*2 || y > h+DY {
				continue
			}

			cv.SetFillStyle("texture/terrain/" + f.Terrain.T.Name + ".png")
			cv.SetStrokeStyle(color.RGBA{R: 192, G: 192, B: 192, A: 24})
			cv.SetLineWidth(2)

			rf := renderer.RenderedField{
				X: [4]float64{float64(x), float64(x - DX), float64(x), float64(x + DX)},
				Y: [4]float64{float64(y), float64(y + DY), float64(y + DY*2.0), float64(y + DY)},
				Z: [4]float64{DZ * float64(t), DZ * float64(l), DZ * float64(b), DZ * float64(r)},
				F: f,
			}
			RenderField(ic, cv, rf, t, l, b, r, m, f, c)
			if f.Travellers != nil {
				RenderTravellers(cv, f.Travellers, rf, c)
			}
			c.AddRenderedField(&rf)
		}
	}
	c.SwapRenderedObjects()
	RenderActiveBuildingPlanBase(cv, c)
	c.ControlPanel.Render(cv)
}
