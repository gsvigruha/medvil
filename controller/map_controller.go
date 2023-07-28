package controller

import (
	"github.com/tfriedel6/canvas"
	"github.com/tfriedel6/canvas/backend/goglbackend"
	"image/color"
	"medvil/model"
	"medvil/model/building"
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
				bt := field.Building.GetBuilding().Plan.BuildingType
				if bt == building.BuildingTypeWall || bt == building.BuildingTypeGate || bt == building.BuildingTypeTower {
					cv.SetFillStyle("#888")
				} else {
					cv.SetFillStyle("#800")
				}
			}
			if field.Road != nil {
				cv.SetFillStyle("#445")
			}
			if field.Plant != nil && !field.Plant.IsTree() {
				cv.SetFillStyle("#A80")
			}
			cv.FillRect(float64(i)*d, float64(j)*d, d, d)
			if field.Animal != nil {
				cv.SetFillStyle("#BBB")
				cv.FillRect(float64(i)*d+d/3, float64(j)*d+d/3, d*2/3, d*2/3)
			}
			if field.Plant != nil && field.Plant.IsTree() {
				cv.SetFillStyle(color.RGBA{R: 0, G: 64 - field.NW/2, B: 0, A: 255})
				cv.FillRect(float64(i)*d+d/3, float64(j)*d+d/3, d*2/3, d*2/3)
			}
		}
	}
	mb := MapButton{m: cp.C.Map, img: cv}
	p.AddButton(mb)
	cp.SetDynamicPanel(p)
}
