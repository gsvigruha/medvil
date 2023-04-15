package controller

import (
	"github.com/tfriedel6/canvas"
	"image/color"
	"medvil/model/building"
	"medvil/model/materials"
	"medvil/model/navigation"
	"medvil/model/social"
	"medvil/renderer"
	"medvil/view/gui"
	"strconv"
)

var BuildingButtonPanelTop = 0.15
var BuildingBasePanelTop = 0.5

var DX = 24.0
var DY = 16.0
var DZ = 12.0

func ScaleBuildingControllerElements(scale float64) {
	DX = scale * 24
	DY = scale * 16
	DZ = scale * 12
}

type BuildingsController struct {
	PlanName    string
	Plan        *building.BuildingPlan
	RoofM       *materials.Material
	UnitM       *materials.Material
	ExtensionT  *building.BuildingExtensionType
	Direction   uint8
	Perspective *uint8
	bt          building.BuildingType
	p           *gui.Panel
	activeTown  *social.Town
	del         bool
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
			if b.bc.Plan.BaseShape[b.i][b.j].Extension == nil && len(b.bc.Plan.BaseShape[b.i][b.j].Floors) < building.MaxNumFloors(b.bc.bt) {
				b.bc.Plan.BaseShape[b.i][b.j].Floors = append(b.bc.Plan.BaseShape[b.i][b.j].Floors, building.Floor{M: b.bc.UnitM})
			}
		}
	} else if b.bc.RoofM != nil {
		if !b.bc.del && b.bc.Plan.BaseShape[b.i][b.j] != nil && len(b.bc.Plan.BaseShape[b.i][b.j].Floors) > 0 {
			b.bc.Plan.BaseShape[b.i][b.j].Roof = &building.Roof{RoofType: building.RoofTypeSplit, M: b.bc.RoofM}
		}
	} else if b.bc.ExtensionT != nil && b.bc.Plan.HasNeighborUnit(uint8(b.i), uint8(b.j), 0) && len(b.bc.Plan.GetExtensions()) == 0 {
		if !b.bc.del && b.bc.Plan.BaseShape[b.i][b.j] == nil {
			b.bc.Plan.BaseShape[b.i][b.j] = &building.PlanUnits{}
			b.bc.Plan.BaseShape[b.i][b.j].Extension = &building.BuildingExtension{T: b.bc.ExtensionT}
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
	if b.ET == nil || b.ET == building.Forge {
		if b.M != nil {
			cv.SetFillStyle("texture/building/" + b.M.Name + ".png")
		} else if b.ET == building.Forge {
			cv.SetFillStyle("texture/building/stone.png")
		}
		cv.SetStrokeStyle("#666")
		cv.SetLineWidth(2)
		cv.BeginPath()
		for _, p := range b.p.Points {
			cv.LineTo(p.X, p.Y)
		}
		cv.ClosePath()
		if b.M != nil || b.ET == building.Forge {
			cv.Fill()
		}
		cv.Stroke()
	} else {
		if b.ET == building.WaterMillWheel {
			img := "icon/gui/building/" + b.ET.Name + ".png"
			cv.DrawImage(img, b.p.Points[0].X-IconS/2, b.p.Points[0].Y+4, IconS, IconS)
		}
	}
}

func (b BuildingBaseButton) Contains(x float64, y float64) bool {
	return b.p.Contains(x, y)
}

func (b BuildingBaseButton) Enabled() bool {
	return true
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
	if FloorM == nil && RoofM == nil && ExtensionT == nil {
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
		cv.FillRect(b.b.X, b.b.Y, IconS, IconS)
	}
}

func (b FloorButton) Contains(x float64, y float64) bool {
	return b.b.Contains(x, y)
}

func (b FloorButton) Enabled() bool {
	return true
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
		cv.FillRect(b.b.X, b.b.Y, IconS, IconS)
	}
}

func (b RoofButton) Contains(x float64, y float64) bool {
	return b.b.Contains(x, y)
}

func (b RoofButton) Enabled() bool {
	return true
}

type ExtensionButton struct {
	b   gui.ButtonGUI
	t   *building.BuildingExtensionType
	bc  *BuildingsController
	del bool
}

func (b ExtensionButton) Click() {
	b.bc.UnitM = nil
	b.bc.RoofM = nil
	b.bc.ExtensionT = b.t
	b.bc.del = b.del
}

func (b ExtensionButton) Render(cv *canvas.Canvas) {
	b.b.Render(cv)
	if b.bc.ExtensionT == nil || b.bc.ExtensionT != b.t || (!b.bc.del && b.del) {
		cv.SetFillStyle(color.RGBA{R: 0, G: 0, B: 0, A: 128})
		cv.FillRect(b.b.X, b.b.Y, IconS, IconS)
	}
}

func (b ExtensionButton) Contains(x float64, y float64) bool {
	return b.b.Contains(x, y)
}

func (b ExtensionButton) Enabled() bool {
	return true
}

type RotationButton struct {
	b  *gui.ButtonGUI
	bc *BuildingsController
}

func (b RotationButton) Click() {
	b.bc.Direction = (b.bc.Direction + 1) % 4
	b.b.Icon = "building/dir_" + strconv.Itoa(int(b.bc.Direction))
	b.bc.RotatePlan()
}

func (b RotationButton) Render(cv *canvas.Canvas) {
	b.b.Render(cv)
}

func (b RotationButton) Contains(x float64, y float64) bool {
	return b.b.Contains(x, y)
}

func (b RotationButton) Enabled() bool {
	return true
}

func (bc *BuildingsController) RotatePlan() {
	newPlan := &building.BuildingPlan{BuildingType: bc.Plan.BuildingType}
	for i := 0; i < building.BuildingBaseMaxSize-1; i++ {
		for j := 0; j < building.BuildingBaseMaxSize-1; j++ {
			newPlan.BaseShape[i][j] = bc.Plan.BaseShape[j][building.BuildingBaseMaxSize-1-i]
		}
	}
	bc.Plan = newPlan
	bc.GenerateButtons()
}

func (bc *BuildingsController) GetActiveFields(c *Controller, rf *renderer.RenderedField) []navigation.FieldWithContext {
	if bc.activeTown.Townhall != nil && !bc.activeTown.Townhall.FieldWithinDistance(rf.F) {
		return nil
	}
	return c.Map.GetBuildingBaseFields(rf.F.X, rf.F.Y, bc.Plan, building.DirectionNone)
}

func (bc *BuildingsController) HandleClick(c *Controller, rf *renderer.RenderedField) bool {
	if bc.activeTown == nil {
		return false
	}
	if bc.activeTown.Townhall != nil && !bc.activeTown.Townhall.FieldWithinDistance(rf.F) {
		return false
	}
	if bc.Plan.IsComplete() {
		c.Map.AddBuildingConstruction(bc.activeTown, rf.F.X, rf.F.Y, bc.Plan, bc.Direction)
		return true
	}
	return false
}

func (bc *BuildingsController) GenerateButtons() {
	bc.p.Buttons = nil
	roofPanelTop := BuildingButtonPanelTop * ControlPanelSY
	bc.p.AddButton(RoofButton{
		b:   gui.ButtonGUI{Icon: "cancel", X: float64(IconW*4 + 10), Y: roofPanelTop, SX: IconS, SY: IconS},
		del: true,
		bc:  bc,
	})
	for i, m := range building.RoofMaterials(bc.bt) {
		bc.p.AddButton(RoofButton{
			b:  gui.ButtonGUI{Texture: "building/" + m.Name, X: float64(i*IconW + 10), Y: roofPanelTop, SX: IconS, SY: IconS},
			m:  m,
			bc: bc,
		})
	}

	floorsPanelTop := BuildingButtonPanelTop*ControlPanelSY + float64(IconH)
	for i, m := range building.FloorMaterials(bc.bt) {
		bc.p.AddButton(FloorButton{
			b:  gui.ButtonGUI{Texture: "building/" + m.Name, X: float64(i*IconW + 10), Y: floorsPanelTop, SX: IconS, SY: IconS},
			m:  m,
			bc: bc,
		})
	}

	extensionPanelTop := BuildingButtonPanelTop*ControlPanelSY + float64(IconH*2)
	for i, e := range building.ExtensionTypes(bc.bt) {
		bc.p.AddButton(ExtensionButton{
			b:  gui.ButtonGUI{Icon: "building/" + e.Name, X: float64(i*IconW + 10), Y: extensionPanelTop, SX: IconS, SY: IconS},
			t:  e,
			bc: bc,
		})
	}

	bc.p.AddButton(RotationButton{
		b:  &gui.ButtonGUI{Icon: "building/dir_" + strconv.Itoa(int(bc.Direction)), X: float64(IconW*5 + 10), Y: roofPanelTop, SX: IconS, SY: IconS},
		bc: bc,
	})

	m := building.BuildingBaseMaxSize
	for i := 0; i < m; i++ {
		for j := 0; j < m; j++ {
			var pi, pj int
			switch *bc.Perspective {
			case PerspectiveNE:
				pi, pj = i, m-1-j
			case PerspectiveSE:
				pi, pj = j, i
			case PerspectiveSW:
				pi, pj = m-1-i, j
			case PerspectiveNW:
				pi, pj = m-1-j, m-1-i
			}

			x := (ControlPanelSX-20)/2 - float64(i)*DX + float64(j)*DX + 10
			y := float64(j)*DY + float64(i)*DY + BuildingBasePanelTop*ControlPanelSY
			bc.p.AddButton(createBuildingBaseButton(bc, pi, pj, 0, x, y, nil, nil, nil))
			if bc.Plan.BaseShape[pi][pj] != nil {
				var k int
				for k = range bc.Plan.BaseShape[pi][pj].Floors {
					bc.p.AddButton(createBuildingBaseButton(bc, pi, pj, k+1, x, y-DZ*float64(k+1), bc.Plan.BaseShape[pi][pj].Floors[k].M, nil, nil))
				}
				if bc.Plan.BaseShape[pi][pj].Roof != nil && !bc.Plan.BaseShape[pi][pj].Roof.Flat() {
					bc.p.AddButton(createBuildingBaseButton(bc, pi, pj, k+1, x, y-DZ*float64(k+1), nil, bc.Plan.BaseShape[pi][pj].Roof.M, nil))
				}
				if bc.Plan.BaseShape[pi][pj].Extension != nil {
					bc.p.AddButton(createBuildingBaseButton(bc, pi, pj, k+1, x, y-DZ*float64(k+1), nil, nil, bc.Plan.BaseShape[pi][pj].Extension.T))
				}
			}
		}
	}
}

func CreateBuildingsController(cp *ControlPanel, bt building.BuildingType, activeTown *social.Town) *BuildingsController {
	p := &gui.Panel{X: 0, Y: ControlPanelDynamicPanelTop * ControlPanelSY, SX: ControlPanelSX, SY: ControlPanelDynamicPanelSY * ControlPanelSY, SingleClick: true}
	bc := &BuildingsController{
		Plan:        &building.BuildingPlan{BuildingType: bt},
		bt:          bt,
		p:           p,
		activeTown:  activeTown,
		Direction:   building.DirectionN,
		Perspective: &cp.C.Perspective}
	bc.GenerateButtons()
	return bc
}

func BuildingsToControlPanel(cp *ControlPanel, bt building.BuildingType) {
	bc := CreateBuildingsController(cp, bt, cp.C.ActiveTown)

	cp.SetDynamicPanel(bc.p)
	cp.C.ClickHandler = bc
}
