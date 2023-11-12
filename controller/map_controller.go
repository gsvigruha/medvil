package controller

import (
	"github.com/tfriedel6/canvas"
	"github.com/tfriedel6/canvas/backend/goglbackend"
	"image/color"
	"math"
	"medvil/model"
	"medvil/model/building"
	"medvil/model/terrain"
	"medvil/view/gui"
)

type MapLabel struct {
	cp  *ControlPanel
	m   *model.Map
	img *canvas.Canvas
	d   float64
}

func (l MapLabel) Render(cv *canvas.Canvas) {
	cv.DrawImage(l.img, 24, ControlPanelSY*0.15, float64(l.img.Width()), float64(l.img.Height()))
	cv.SetLineWidth(2)
	cv.SetStrokeStyle("#D00")
	x := float64(l.cp.C.CenterX)*l.d + 24
	y := float64(l.cp.C.CenterY)*l.d + ControlPanelSY*0.15
	cv.BeginPath()
	cv.Arc(x, y, l.d*12, 0, math.Pi*2.0, true)
	cv.ClosePath()
	cv.Stroke()
}

func (l MapLabel) CaptureClick(x float64, y float64) {
	if x >= 24 && x <= float64(24+l.img.Width()) && y >= ControlPanelSY*0.15 && y <= ControlPanelSY*0.15+float64(l.img.Height()) {
		l.cp.C.CenterX = int((x - 24) / l.d)
		l.cp.C.CenterY = int((y - ControlPanelSY*0.15) / l.d)
	}
}

type MapController struct {
	p *gui.Panel
}

func (mc *MapController) CaptureMove(x, y float64) {
	mc.p.CaptureMove(x, y)
}

func (mc *MapController) CaptureClick(x, y float64) {
	mc.p.CaptureClick(x, y)
}

func (mc *MapController) Render(cv *canvas.Canvas) {
	mc.p.Render(cv)
}

func (mc *MapController) Clear() {}

func (mc *MapController) Refresh() {
}

func (mc *MapController) GetHelperSuggestions() *gui.Suggestion {
	return nil
}

func MapToControlPanel(cp *ControlPanel) {
	p := &gui.Panel{X: 0, Y: ControlPanelDynamicPanelTop, SX: ControlPanelSX, SY: HouseholdControllerSY}
	mc := &MapController{p: p}
	var sx = ControlPanelSX - 48
	if sx > ControlPanelSY*0.3 {
		sx = ControlPanelSY * 0.3
	}
	d := sx / float64(cp.C.Map.SX)
	offscreen, _ := goglbackend.NewOffscreen(int(sx), int(sx)*int(cp.C.Map.SY)/int(cp.C.Map.SX), true, cp.C.ctx)
	cv := canvas.New(offscreen)
	for i, fields := range cp.C.Map.Fields {
		for j, field := range fields {
			if field.Deposit != nil {
				if field.Deposit.T == terrain.IronBog {
					cv.SetFillStyle("#A22")
				} else if field.Deposit.T == terrain.Gold {
					cv.SetFillStyle("#EA2")
				} else if field.Deposit.T == terrain.Mud {
					cv.SetFillStyle("#862")
				} else if field.Deposit.T == terrain.Rock {
					cv.SetFillStyle("#BBB")
				}
			} else if field.Terrain.T == terrain.Water {
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
	mb := MapLabel{cp: cp, m: cp.C.Map, img: cv, d: d}
	p.AddLabel(mb)

	DrawStats(cp, p)

	cp.SetDynamicPanel(mc)
}
