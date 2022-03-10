package controller

import (
	"github.com/tfriedel6/canvas"
	"image/color"
	"medvil/model/building"
	"medvil/model/materials"
	"medvil/renderer"
	"medvil/view/gui"
)

const RoofPanelTop = 100
const FloorsPanelTop = 200
const BuildingBasePanelTop = 400

const DX = 24
const DY = 16

type BuildingsController struct {
	PlanName string
	Plan     *building.BuildingPlan
	RoofM    *materials.Material
	UnitM    *materials.Material
	bt       building.BuildingType
}

type BuildingBaseButton struct {
	i        uint8
	j        uint8
	k        uint8
	p        renderer.Polygon
	panel    *gui.Panel
	bc       *BuildingsController
	addFloor bool
	M        *materials.Material
}

func (b BuildingBaseButton) Click() {
	if b.bc.UnitM != nil {
		if !b.bc.Plan.HasUnit(b.i, b.j, b.k) {
			if b.bc.Plan.BaseShape[b.i][b.j] == nil {
				b.bc.Plan.BaseShape[b.i][b.j] = &building.PlanUnits{}
			}
			b.bc.Plan.BaseShape[b.i][b.j].Floors = append(b.bc.Plan.BaseShape[b.i][b.j].Floors, building.Floor{M: b.bc.UnitM})
			b.panel.AddButton(createBuildingBaseButton(b.bc, b.panel, b.i, b.j, b.k+1, b.p.Points[0].X, b.p.Points[0].Y-DY/2))
		}
	} else if b.bc.RoofM != nil {
		if b.bc.Plan.BaseShape[b.i][b.j] != nil {
			b.bc.Plan.BaseShape[b.i][b.j].Roof = &building.Roof{Flat: false, M: b.bc.RoofM}
		}
	}
}

func (b BuildingBaseButton) Render(cv *canvas.Canvas) {
	cv.SetStrokeStyle("#AAA")
	cv.BeginPath()
	for _, p := range b.p.Points {
		cv.LineTo(p.X, p.Y)
	}
	cv.ClosePath()
	cv.Stroke()
}

func (b BuildingBaseButton) Contains(x float64, y float64) bool {
	return b.p.Contains(x, y)
}

func createBuildingBaseButton(bc *BuildingsController, panel *gui.Panel, i, j, k uint8, x, y float64) *BuildingBaseButton {
	return &BuildingBaseButton{
		i:     i,
		j:     j,
		k:     k,
		bc:    bc,
		panel: panel,
		p: renderer.Polygon{Points: []renderer.Point{
			renderer.Point{x, y},
			renderer.Point{x - DX, y + DY},
			renderer.Point{x, y + DY*2},
			renderer.Point{x + DX, y + DY},
		}},
		addFloor: true,
	}
}

type FloorButton struct {
	b  gui.ButtonGUI
	m  *materials.Material
	bc *BuildingsController
}

func (b FloorButton) Click() {
	b.bc.UnitM = b.m
	b.bc.RoofM = nil
}

func (b FloorButton) Render(cv *canvas.Canvas) {
	b.b.Render(cv)
	if b.bc.UnitM != b.m {
		cv.SetFillStyle(color.RGBA{R: 0, G: 0, B: 0, A: 128})
		cv.FillRect(b.b.X, b.b.Y, 32, 32)
	}
}

func (b FloorButton) Contains(x float64, y float64) bool {
	return b.b.Contains(x, y)
}

type RoofButton struct {
	b  gui.ButtonGUI
	m  *materials.Material
	bc *BuildingsController
}

func (b RoofButton) Click() {
	b.bc.RoofM = b.m
	b.bc.UnitM = nil
}

func (b RoofButton) Render(cv *canvas.Canvas) {
	b.b.Render(cv)
	if b.bc.RoofM != b.m {
		cv.SetFillStyle(color.RGBA{R: 0, G: 0, B: 0, A: 128})
		cv.FillRect(b.b.X, b.b.Y, 32, 32)
	}
}

func (b RoofButton) Contains(x float64, y float64) bool {
	return b.b.Contains(x, y)
}

func (bc *BuildingsController) HandleClick(c *Controller, rf *renderer.RenderedField) bool {
	if c.ActiveBuildingPlan.IsComplete() {
		c.Map.AddBuildingConstruction(c.Country, rf.F.X, rf.F.Y, c.ActiveBuildingPlan, bc.bt)
		return true
	}
	return false
}

func BuildingsToControlPanel(cp *ControlPanel, bt building.BuildingType) {
	p := &gui.Panel{X: 0, Y: ControlPanelDynamicPanelTop, SX: ControlPanelSX, SY: ControlPanelDynamicPanelSY}
	bc := &BuildingsController{Plan: &building.BuildingPlan{}, bt: bt}

	for i, m := range building.RoofMaterials(bt) {
		p.AddButton(RoofButton{
			b:  gui.ButtonGUI{Texture: "building/" + m.Name, X: float64(i*40 + 10), Y: float64(RoofPanelTop), SX: 32, SY: 32},
			m:  m,
			bc: bc,
		})
	}

	for i, m := range building.FloorMaterials(bt) {
		p.AddButton(FloorButton{
			b:  gui.ButtonGUI{Texture: "building/" + m.Name, X: float64(i*40 + 10), Y: float64(FloorsPanelTop), SX: 32, SY: 32},
			m:  m,
			bc: bc,
		})
	}

	for i := 0; i < building.BuildingBaseMaxSize; i++ {
		for j := 0; j < building.BuildingBaseMaxSize; j++ {
			p.AddButton(createBuildingBaseButton(bc, p, uint8(i), uint8(j), 0, float64(120+i*DX-j*DX+10), float64(j*DY+i*DY+BuildingBasePanelTop)))
		}
	}
	cp.SetDynamicPanel(p)
	cp.C.ActiveBuildingPlan = bc.Plan
	cp.C.ClickHandler = bc
}
