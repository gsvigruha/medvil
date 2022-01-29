package controller

import (
	"github.com/tfriedel6/canvas"
	"image/color"
	"medvil/model/building"
	"medvil/model/materials"
	"medvil/view/gui"
)

const RoofPanelTop = 100
const FloorsPanelTop = 200
const BuildingBasePanelTop = 400

type BuildingsController struct {
	PlanName string
	Plan     *building.BuildingPlan
}

type BuildingBaseButton struct {
	b  gui.ButtonGUI
	i  int
	j  int
	bc *BuildingsController
}

func (b BuildingBaseButton) Click() {
	b.bc.Plan.BaseShape[b.i][b.j] = !b.bc.Plan.BaseShape[b.i][b.j]
}

func (b BuildingBaseButton) Render(cv *canvas.Canvas) {
	if b.bc.Plan.BaseShape[b.i][b.j] {
		cv.SetFillStyle("#AAA")
		cv.FillRect(b.b.X, b.b.Y, 32, 32)
	}
	b.b.Render(cv)
	if !b.bc.Plan.BaseShape[b.i][b.j] {
		cv.SetFillStyle(color.RGBA{R: 64, G: 0, B: 0, A: 128})
		cv.FillRect(b.b.X, b.b.Y, 32, 32)
	}
}

func (b BuildingBaseButton) Contains(x float64, y float64) bool {
	return b.b.Contains(x, y)
}

type FloorButton struct {
	b  gui.ButtonGUI
	f  int
	m  *materials.Material
	bc *BuildingsController
}

func (b FloorButton) Click() {
	nf := len(b.bc.Plan.Floors)
	if nf > b.f {
		if b.bc.Plan.Floors[b.f].M == b.m {
			if b.f == nf-1 {
				b.bc.Plan.Floors = b.bc.Plan.Floors[:nf-1]
			}
		} else {
			b.bc.Plan.Floors[b.f].M = b.m
		}
	} else if nf == b.f {
		b.bc.Plan.Floors = append(b.bc.Plan.Floors, building.Floor{M: b.m})
	}
}

func (b FloorButton) Render(cv *canvas.Canvas) {
	b.b.Render(cv)
	floors := b.bc.Plan.Floors
	if len(floors) < b.f {
		cv.SetFillStyle(color.RGBA{R: 0, G: 0, B: 0, A: 192})
		cv.FillRect(b.b.X, b.b.Y, 32, 32)
	} else if len(floors) == b.f {
		cv.SetFillStyle(color.RGBA{R: 0, G: 0, B: 0, A: 128})
		cv.FillRect(b.b.X, b.b.Y, 32, 32)
	} else {
		if floors[b.f].M != b.m {
			cv.SetFillStyle(color.RGBA{R: 0, G: 0, B: 0, A: 128})
			cv.FillRect(b.b.X, b.b.Y, 32, 32)
		}
	}
}

func (b FloorButton) Contains(x float64, y float64) bool {
	return b.b.Contains(x, y)
}

type RoofButton struct {
	b    gui.ButtonGUI
	m    *materials.Material
	flat bool
	bc   *BuildingsController
}

func (b RoofButton) Click() {
	b.bc.Plan.Roof.M = b.m
	b.bc.Plan.Roof.Flat = b.flat
}

func (b RoofButton) Render(cv *canvas.Canvas) {
	b.b.Render(cv)
	if b.bc.Plan.Roof.M != b.m {
		cv.SetFillStyle(color.RGBA{R: 0, G: 0, B: 0, A: 128})
		cv.FillRect(b.b.X, b.b.Y, 32, 32)
	}
}

func (b RoofButton) Contains(x float64, y float64) bool {
	return b.b.Contains(x, y)
}

func BuildingsToControlPanel(cp *ControlPanel) {
	p := &gui.Panel{X: 0, Y: ControlPanelDynamicPanelTop, SX: ControlPanelSX, SY: ControlPanelDynamicPanelSY}
	bc := BuildingsController{Plan: &building.BuildingPlan{}}

	for i, m := range building.RoofMaterials {
		p.AddButton(RoofButton{
			b:    gui.ButtonGUI{Texture: "building/" + m.Name, X: float64(i*40 + 10), Y: float64(RoofPanelTop), SX: 32, SY: 32},
			m:    m,
			flat: false,
			bc:   &bc,
		})
	}

	for i, m := range building.FlatRoofMaterials {
		p.AddButton(RoofButton{
			b:    gui.ButtonGUI{Texture: "building/" + m.Name, X: float64((i+len(building.RoofMaterials))*40 + 10), Y: float64(RoofPanelTop), SX: 32, SY: 32},
			m:    m,
			flat: true,
			bc:   &bc,
		})
	}

	for j := 0; j < building.MaxFloors; j++ {
		for i, m := range building.FloorMaterials {
			p.AddButton(FloorButton{
				b:  gui.ButtonGUI{Texture: "building/" + m.Name, X: float64(i*40 + 10), Y: float64(j*40 + FloorsPanelTop), SX: 32, SY: 32},
				f:  building.MaxFloors - j - 1,
				m:  m,
				bc: &bc,
			})
		}
	}

	for i := 0; i < building.BuildingBaseMaxSize; i++ {
		for j := 0; j < building.BuildingBaseMaxSize; j++ {
			p.AddButton(BuildingBaseButton{
				b:  gui.ButtonGUI{Icon: "parcel", X: float64(i*40 + 10), Y: float64(j*40 + BuildingBasePanelTop), SX: 32, SY: 32},
				i:  i,
				j:  j,
				bc: &bc,
			})
		}
	}
	cp.SetDynamicPanel(p)
	cp.C.ActiveBuildingPlan = bc.Plan
}
