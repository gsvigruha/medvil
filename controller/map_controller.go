package controller

import (
	"github.com/tfriedel6/canvas"
	"github.com/tfriedel6/canvas/backend/goglbackend"
	"image/color"
	"medvil/model"
	"medvil/model/terrain"
	"medvil/view/gui"
)

type MapButton struct {
	m   *model.Map
	img *canvas.Canvas
}

func (b MapButton) Click() {

}

func (b MapButton) Render(cv *canvas.Canvas) {
	cv.DrawImage(b.img, 24, ControlPanelSY*0.15, float64(b.img.Width()), float64(b.img.Height()))
}

func (b MapButton) Contains(x float64, y float64) bool {
	return x >= 24 && x <= float64(24+b.m.SX*2) && y >= ControlPanelSY*0.15 && y <= ControlPanelSY*0.15+float64(b.m.SY*2)
}

func (b MapButton) Enabled() bool {
	return true
}

func MapToControlPanel(cp *ControlPanel) {
	p := &gui.Panel{X: 0, Y: ControlPanelDynamicPanelTop, SX: ControlPanelSX, SY: HouseholdControllerSY}
	sx := ControlPanelSX - 48
	d := sx / float64(cp.C.Map.SX)
	offscreen, _ := goglbackend.NewOffscreen(int(sx), int(sx)*int(cp.C.Map.SY)/int(cp.C.Map.SX), true, cp.C.ctx)
	cv := canvas.New(offscreen)
	for i, fields := range cp.C.Map.Fields {
		for j, field := range fields {
			if field.Terrain.T == terrain.Water {
				cv.SetFillStyle("#48D")
			} else if field.Terrain.T == terrain.Grass {
				cv.SetFillStyle(color.RGBA{R: 0, G: 128 - field.NW, B: 0, A: 255})
			}
			if field.Building.GetBuilding() != nil {
				cv.SetFillStyle("#D00")
			}
			cv.FillRect(float64(i)*d, float64(j)*d, d, d)
		}
	}
	mb := MapButton{m: cp.C.Map, img: cv}
	p.AddButton(mb)
	cp.SetDynamicPanel(p)
}
