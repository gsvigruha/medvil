package view

import (
	"github.com/tfriedel6/canvas"
	"image/color"
	"medvil/controller"
	"medvil/model"
	//"fmt"
)

const (
	DX float64 = 60.0
	DY float64 = 40.0
	DZ float64 = 15.0
)

func Render(cv *canvas.Canvas, m model.Map) {
	w := float64(cv.Width())
	for i := uint16(0); i < m.SX; i++ {
		for j := uint16(0); j < m.SY; j++ {
			var pi = i
			var pj = j
			var f = m.Fields[pi][pj]
			var t = uint8(0)
			var r = uint8(0)
			var b = uint8(0)
			var l = uint8(0)
			switch controller.Perspective {
			case controller.PerspectiveNE:
				pi = i
				pj = m.SY - 1 - j
				f = m.Fields[pi][pj]
				t = f.SW
				r = f.NW
				b = f.NE
				l = f.SE
			case controller.PerspectiveSE:
				pi = j
				pj = i
				f = m.Fields[pi][pj]
				t = f.NW
				r = f.NE
				b = f.SE
				l = f.SW
			case controller.PerspectiveSW:
				pi = m.SX - 1 - i
				pj = j
				f = m.Fields[pi][pj]
				t = f.NE
				r = f.SE
				b = f.SW
				l = f.NW
			case controller.PerspectiveNW:
				pi = m.SY - 1 - j
				pj = m.SX - 1 - i
				f = m.Fields[pi][pj]
				t = f.SE
				r = f.SW
				b = f.NW
				l = f.NE
			}
			cv.SetFillStyle("texture/terrain/" + f.Terrain.T.Name + ".png")
			cv.SetStrokeStyle(color.RGBA{R: 192, G: 192, B: 192, A: 24})
			cv.SetLineWidth(2)
			x := w/2 - float64(i)*DX + float64(j)*DX + float64(controller.ScrollX)
			y := float64(i)*DY + float64(j)*DY + float64(controller.ScrollY)

			rf := RenderedField{
				X: [4]float64{float64(x), float64(x - DX), float64(x), float64(x + DX)},
				Y: [4]float64{float64(y), float64(y + DY), float64(y + DY*2.0), float64(y + DY)},
				Z: [4]float64{DZ * float64(t), DZ * float64(l), DZ * float64(b), DZ * float64(r)}}
			rf.Draw(cv)
			cv.Fill()
			cv.Stroke()

			units := m.Fields[pi][pj].Building.BuildingUnits
			for k := 0; k < len(units); k++ {
				RenderBuildingUnit(cv, units[k], rf, k)
			}
			roof := m.Fields[pi][pj].Building.RoofUnit
			RenderBuildingRoof(cv, roof, rf, len(units))
		}
	}
}
