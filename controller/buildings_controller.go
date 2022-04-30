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
const ExtensionPanelTop = 300
const BuildingBasePanelTop = 500

const DX = 24
const DY = 16
const DZ = 12

type BuildingsController struct {
	PlanName   string
	Plan       *building.BuildingPlan
	RoofM      *materials.Material
	UnitM      *materials.Material
	ExtensionT *building.BuildingExtensionType
	bt         building.BuildingType
	p          *gui.Panel
	del        bool
}

type BuildingBaseButton struct {
	i     int
	j     int
	k     int
	p     renderer.Polygon
	panel *gui.Panel
	bc    *BuildingsController
	M     *materials.Material
	ET    *building.BuildingExtensionType
}

func (b BuildingBaseButton) Click() {
	if b.bc.UnitM != nil {
		if !b.bc.del && !b.bc.Plan.HasUnit(uint8(b.i), uint8(b.j), uint8(b.k)) {
			if b.bc.Plan.BaseShape[b.i][b.j] == nil {
				b.bc.Plan.BaseShape[b.i][b.j] = &building.PlanUnits{}
				b.bc.Plan.BaseShape[b.i][b.j].Roof = &building.Roof{RoofType: building.RoofTypeFlat, M: b.bc.UnitM}
			}
			if b.bc.Plan.BaseShape[b.i][b.j].Extension == nil {
				b.bc.Plan.BaseShape[b.i][b.j].Floors = append(b.bc.Plan.BaseShape[b.i][b.j].Floors, building.Floor{M: b.bc.UnitM})
			}
		}
	} else if b.bc.RoofM != nil {
		if !b.bc.del && b.bc.Plan.BaseShape[b.i][b.j] != nil && len(b.bc.Plan.BaseShape[b.i][b.j].Floors) > 0 {
			b.bc.Plan.BaseShape[b.i][b.j].Roof = &building.Roof{RoofType: building.RoofTypeSplit, M: b.bc.RoofM}
		}
	} else if b.bc.ExtensionT != nil && b.bc.Plan.HasNeighborUnit(uint8(b.i), uint8(b.j), 0) && b.bc.Plan.GetExtension() == nil {
		if !b.bc.del && b.bc.Plan.BaseShape[b.i][b.j] == nil {
			b.bc.Plan.BaseShape[b.i][b.j] = &building.PlanUnits{}
			b.bc.Plan.BaseShape[b.i][b.j].Extension = &building.BuildingExtension{T: *b.bc.ExtensionT}
		}
	} else if b.bc.del {
		if b.bc.Plan.BaseShape[b.i][b.j] != nil {
			maxFloor := len(b.bc.Plan.BaseShape[b.i][b.j].Floors) - 1
			if maxFloor >= 0 {
				if b.bc.Plan.BaseShape[b.i][b.j].Roof.Flat() {
					b.bc.Plan.BaseShape[b.i][b.j].Floors = b.bc.Plan.BaseShape[b.i][b.j].Floors[0:maxFloor]
				} else {
					b.bc.Plan.BaseShape[b.i][b.j].Roof.RoofType = building.RoofTypeFlat
					b.bc.Plan.BaseShape[b.i][b.j].Roof.M = b.bc.Plan.BaseShape[b.i][b.j].Floors[maxFloor].M
				}
			}
			if b.bc.Plan.BaseShape[b.i][b.j].Extension != nil {
				b.bc.Plan.BaseShape[b.i][b.j].Extension = nil
			}
			if len(b.bc.Plan.BaseShape[b.i][b.j].Floors) == 0 {
				b.bc.Plan.BaseShape[b.i][b.j] = nil
			}
		}
	}
	b.bc.GenerateButtons()
}

func (b BuildingBaseButton) Render(cv *canvas.Canvas) {
	if b.M != nil {
		cv.SetFillStyle("texture/building/" + b.M.Name + ".png")
	}
	cv.SetStrokeStyle("#666")
	cv.SetLineWidth(2)
	cv.BeginPath()
	for _, p := range b.p.Points {
		cv.LineTo(p.X, p.Y)
	}
	cv.ClosePath()
	if b.M != nil {
		cv.Fill()
	}
	cv.Stroke()
}

func (b BuildingBaseButton) Contains(x float64, y float64) bool {
	return b.p.Contains(x, y)
}

func createBuildingBaseButton(
	bc *BuildingsController,
	i, j, k int,
	x, y float64,
	FloorM *materials.Material,
	RoofM *materials.Material,
	ExtensionT *building.BuildingExtensionType) *BuildingBaseButton {

	var polygon renderer.Polygon
	var M *materials.Material
	var ET *building.BuildingExtensionType
	if FloorM == nil && RoofM == nil {
		polygon = renderer.Polygon{Points: []renderer.Point{
			renderer.Point{x, y},
			renderer.Point{x - DX, y + DY},
			renderer.Point{x, y + DY*2},
			renderer.Point{x + DX, y + DY},
		}}
	} else if RoofM != nil {
		polygon = renderer.Polygon{Points: []renderer.Point{
			renderer.Point{x, y + DY*2},
			renderer.Point{x - DX, y + DY},
			renderer.Point{x, y - DZ},
			renderer.Point{x + DX, y + DY},
		}}
		M = RoofM
	} else if FloorM != nil {
		polygon = renderer.Polygon{Points: []renderer.Point{
			renderer.Point{x, y},
			renderer.Point{x - DX, y + DY},
			renderer.Point{x - DX, y + DY + DZ},
			renderer.Point{x, y + DY*2 + DZ},
			renderer.Point{x + DX, y + DY + DZ},
			renderer.Point{x + DX, y + DY},
		}}
		M = FloorM
	} else if ExtensionT != nil {
		polygon = renderer.Polygon{Points: []renderer.Point{
			renderer.Point{x, y},
			renderer.Point{x - DX, y + DY},
			renderer.Point{x - DX, y + DY + DZ},
			renderer.Point{x, y + DY*2 + DZ},
			renderer.Point{x + DX, y + DY + DZ},
			renderer.Point{x + DX, y + DY},
		}}
		ET = ExtensionT
	}
	return &BuildingBaseButton{
		i:     i,
		j:     j,
		k:     k,
		bc:    bc,
		panel: bc.p,
		p:     polygon,
		M:     M,
		ET:    ET,
	}
}

type FloorButton struct {
	b   gui.ButtonGUI
	m   *materials.Material
	bc  *BuildingsController
	del bool
}

func (b FloorButton) Click() {
	b.bc.UnitM = b.m
	b.bc.RoofM = nil
	b.bc.ExtensionT = nil
	b.bc.del = b.del
}

func (b FloorButton) Render(cv *canvas.Canvas) {
	b.b.Render(cv)
	if b.bc.UnitM != b.m || (!b.bc.del && b.del) {
		cv.SetFillStyle(color.RGBA{R: 0, G: 0, B: 0, A: 128})
		cv.FillRect(b.b.X, b.b.Y, 32, 32)
	}
}

func (b FloorButton) Contains(x float64, y float64) bool {
	return b.b.Contains(x, y)
}

type RoofButton struct {
	b   gui.ButtonGUI
	m   *materials.Material
	bc  *BuildingsController
	del bool
}

func (b RoofButton) Click() {
	b.bc.RoofM = b.m
	b.bc.UnitM = nil
	b.bc.ExtensionT = nil
	b.bc.del = b.del
}

func (b RoofButton) Render(cv *canvas.Canvas) {
	b.b.Render(cv)
	if b.bc.RoofM != b.m || (!b.bc.del && b.del) {
		cv.SetFillStyle(color.RGBA{R: 0, G: 0, B: 0, A: 128})
		cv.FillRect(b.b.X, b.b.Y, 32, 32)
	}
}

func (b RoofButton) Contains(x float64, y float64) bool {
	return b.b.Contains(x, y)
}

type ExtensionButton struct {
	b   gui.ButtonGUI
	t   building.BuildingExtensionType
	bc  *BuildingsController
	del bool
}

func (b ExtensionButton) Click() {
	b.bc.UnitM = nil
	b.bc.RoofM = nil
	b.bc.ExtensionT = &b.t
	b.bc.del = b.del
}

func (b ExtensionButton) Render(cv *canvas.Canvas) {
	b.b.Render(cv)
	if b.bc.ExtensionT == nil || *b.bc.ExtensionT != b.t || (!b.bc.del && b.del) {
		cv.SetFillStyle(color.RGBA{R: 0, G: 0, B: 0, A: 128})
		cv.FillRect(b.b.X, b.b.Y, 32, 32)
	}
}

func (b ExtensionButton) Contains(x float64, y float64) bool {
	return b.b.Contains(x, y)
}

func (bc *BuildingsController) HandleClick(c *Controller, rf *renderer.RenderedField) bool {
	if c.ActiveBuildingPlan.IsComplete() {
		c.Map.AddBuildingConstruction(c.Country, rf.F.X, rf.F.Y, c.ActiveBuildingPlan)
		return true
	}
	return false
}

func (bc *BuildingsController) GenerateButtons() {
	bc.p.Buttons = nil
	bc.p.AddButton(RoofButton{
		b:   gui.ButtonGUI{Icon: "cancel", X: 10, Y: float64(RoofPanelTop), SX: 32, SY: 32},
		del: true,
		bc:  bc,
	})
	for i, m := range building.RoofMaterials(bc.bt) {
		bc.p.AddButton(RoofButton{
			b:  gui.ButtonGUI{Texture: "building/" + m.Name, X: float64(i*40 + 50), Y: float64(RoofPanelTop), SX: 32, SY: 32},
			m:  m,
			bc: bc,
		})
	}

	for i, m := range building.FloorMaterials(bc.bt) {
		bc.p.AddButton(FloorButton{
			b:  gui.ButtonGUI{Texture: "building/" + m.Name, X: float64(i*40 + 10), Y: float64(FloorsPanelTop), SX: 32, SY: 32},
			m:  m,
			bc: bc,
		})
	}

	if bc.bt == building.BuildingTypeWorkshop {
		bc.p.AddButton(ExtensionButton{
			b:  gui.ButtonGUI{Icon: "building/watermill_wheel", X: 10, Y: float64(ExtensionPanelTop), SX: 32, SY: 32},
			t:  building.WaterMillWheel,
			bc: bc,
		})

		bc.p.AddButton(ExtensionButton{
			b:  gui.ButtonGUI{Icon: "building/forge", X: 50, Y: float64(ExtensionPanelTop), SX: 32, SY: 32},
			t:  building.Forge,
			bc: bc,
		})
	}

	for i := 0; i < building.BuildingBaseMaxSize; i++ {
		for j := 0; j < building.BuildingBaseMaxSize; j++ {
			x := float64(120 + i*DX - j*DX + 10)
			y := float64(j*DY + i*DY + BuildingBasePanelTop)
			bc.p.AddButton(createBuildingBaseButton(bc, i, j, 0, x, y, nil, nil, nil))
			if bc.Plan.BaseShape[i][j] != nil {
				var k int
				for k = range bc.Plan.BaseShape[i][j].Floors {
					bc.p.AddButton(createBuildingBaseButton(bc, i, j, k+1, x, y-DZ*float64(k+1), bc.Plan.BaseShape[i][j].Floors[k].M, nil, nil))
				}
				if bc.Plan.BaseShape[i][j].Roof != nil && !bc.Plan.BaseShape[i][j].Roof.Flat() {
					bc.p.AddButton(createBuildingBaseButton(bc, i, j, k+1, x, y-DZ*float64(k+1), nil, bc.Plan.BaseShape[i][j].Roof.M, nil))
				}
				if bc.Plan.BaseShape[i][j].Extension != nil {
					bc.p.AddButton(createBuildingBaseButton(bc, i, j, k+1, x, y-DZ*float64(k+1), nil, nil, &bc.Plan.BaseShape[i][j].Extension.T))
				}
			}
		}
	}
}

func BuildingsToControlPanel(cp *ControlPanel, bt building.BuildingType) {
	p := &gui.Panel{X: 0, Y: ControlPanelDynamicPanelTop, SX: ControlPanelSX, SY: ControlPanelDynamicPanelSY, SingleClick: true}
	bc := &BuildingsController{Plan: &building.BuildingPlan{BuildingType: bt}, bt: bt, p: p}

	bc.GenerateButtons()

	cp.SetDynamicPanel(p)
	cp.C.ActiveBuildingPlan = bc.Plan
	cp.C.ClickHandler = bc
}
