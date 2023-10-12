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

var BuildingButtonPanelTop = 0.2
var BuildingBasePanelTop = 0.5
var BuildingCostTop = 0.75

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
	cp          *ControlPanel
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
			b.bc.Plan.BaseShape[b.i][b.j].Roof = &building.Roof{RoofType: building.GetRoofType(b.bc.RoofM), M: b.bc.RoofM}
		}
	} else if b.bc.ExtensionT != nil && len(b.bc.Plan.GetExtensions()) == 0 && !b.bc.del {
		if b.bc.ExtensionT.InUnit && b.bc.Plan.BaseShape[b.i][b.j] != nil && len(b.bc.Plan.BaseShape[b.i][b.j].Floors) > 0 {
			b.bc.Plan.BaseShape[b.i][b.j].Extension = &building.BuildingExtension{T: b.bc.ExtensionT}
		} else if !b.bc.ExtensionT.InUnit && b.bc.Plan.HasNeighborUnit(uint8(b.i), uint8(b.j), 0) && b.bc.Plan.BaseShape[b.i][b.j] == nil {
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
	b.bc.GenerateBuildingTypebuttons()
}

func (b BuildingBaseButton) Render(cv *canvas.Canvas) {
	if b.ET == nil || b.ET == building.Forge || b.ET == building.Kiln {
		if b.M != nil {
			cv.SetFillStyle("texture/building/" + b.M.Name + ".png")
		} else if b.ET == building.Forge {
			cv.SetFillStyle("texture/building/stone.png")
		} else if b.ET == building.Kiln {
			cv.SetFillStyle("texture/building/brick.png")
		}
		cv.SetStrokeStyle("#666")
		cv.SetLineWidth(2)
		cv.BeginPath()
		for _, p := range b.p.Points {
			cv.LineTo(p.X, p.Y)
		}
		cv.ClosePath()
		if b.M != nil || b.ET == building.Forge || b.ET == building.Kiln {
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

func (b BuildingBaseButton) SetHoover(h bool) {}

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
	RoofT *building.RoofType,
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
		var h = 0.0
		if *RoofT == building.RoofTypeSplit {
			h = DZ
		}
		polygon = renderer.Polygon{Points: []renderer.Point{
			renderer.Point{x, y + DY*2},
			renderer.Point{x - DX, y + DY},
			renderer.Point{x, y - h},
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
	b   *gui.ButtonGUI
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
		cv.FillRect(b.b.X, b.b.Y, LargeIconS, LargeIconS)
	}
}

func (b *FloorButton) SetHoover(h bool) {
	b.b.SetHoover(h)
}

func (b FloorButton) Contains(x float64, y float64) bool {
	return b.b.Contains(x, y)
}

func (b FloorButton) Enabled() bool {
	return true
}

type RoofButton struct {
	b   *gui.ButtonGUI
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
		cv.FillRect(b.b.X, b.b.Y, LargeIconS, LargeIconS)
	}
}

func (b *RoofButton) SetHoover(h bool) {
	b.b.SetHoover(h)
}

func (b RoofButton) Contains(x float64, y float64) bool {
	return b.b.Contains(x, y)
}

func (b RoofButton) Enabled() bool {
	return true
}

type ExtensionButton struct {
	b   *gui.ButtonGUI
	t   *building.BuildingExtensionType
	bc  *BuildingsController
	msg string
	del bool
}

func (b ExtensionButton) Click() {
	b.bc.UnitM = nil
	b.bc.RoofM = nil
	b.bc.ExtensionT = b.t
	b.bc.del = b.del
	b.bc.cp.HelperMessage(b.msg)
}

func (b ExtensionButton) Render(cv *canvas.Canvas) {
	b.b.Render(cv)
	if b.bc.ExtensionT == nil || b.bc.ExtensionT != b.t || (!b.bc.del && b.del) {
		cv.SetFillStyle(color.RGBA{R: 0, G: 0, B: 0, A: 128})
		cv.FillRect(b.b.X, b.b.Y, LargeIconS, LargeIconS)
	}
}

func (b *ExtensionButton) SetHoover(h bool) {
	b.b.SetHoover(h)
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

func (b *RotationButton) SetHoover(h bool) {
	b.b.SetHoover(h)
}

func (b RotationButton) Contains(x float64, y float64) bool {
	return b.b.Contains(x, y)
}

func (b RotationButton) Enabled() bool {
	return true
}

func (bc *BuildingsController) GetHelperSuggestions() *gui.Suggestion {
	if !bc.Plan.IsComplete() && bc.RoofM == nil && bc.UnitM == nil {
		return &gui.Suggestion{
			Message: "Design your house. First pick wall material\nand click to build units. Afterwards pick\na roof material and add a roof.",
			Icon:    "house", X: LargeIconD*2 + 24, Y: BuildingButtonPanelTop*ControlPanelSY + float64(LargeIconD),
		}
	}
	if bc.Plan.IsComplete() && bc.Plan.BuildingType == building.BuildingTypeWorkshop && len(bc.Plan.GetExtensions()) == 0 {
		if bc.ExtensionT == nil {
			return &gui.Suggestion{
				Message: ("Pick building extensions like workshop, water wheel or a forge,\nthen click on the building plan to add them.\n" +
					"Each extension lets your workshop to perform a few different tasks."),
				Icon: "building/workshop", X: LargeIconD*5 + 24, Y: BuildingButtonPanelTop*ControlPanelSY + float64(LargeIconD)*2.5,
			}
		} else if bc.ExtensionT == building.WaterMillWheel {
			return &gui.Suggestion{
				Message: "Waterwheels are needed to mill grain or paper, or saw wood.\nThe building needs to be adjacent to water.",
				Icon:    "building/water_mill_wheel", X: LargeIconD*5 + 24, Y: BuildingButtonPanelTop*ControlPanelSY + float64(LargeIconD)*2.5,
			}
		} else if bc.ExtensionT == building.Forge {
			return &gui.Suggestion{
				Message: "Forges let you work metals such as iron or gold, or create tools and weapons.",
				Icon:    "building/forge", X: LargeIconD*5 + 24, Y: BuildingButtonPanelTop*ControlPanelSY + float64(LargeIconD)*2.5,
			}
		} else if bc.ExtensionT == building.Kiln {
			return &gui.Suggestion{
				Message: "Kilns are needed to burn clay to produce bricks, tiles or pots.",
				Icon:    "building/kiln", X: LargeIconD*5 + 24, Y: BuildingButtonPanelTop*ControlPanelSY + float64(LargeIconD)*2.5,
			}
		} else if bc.ExtensionT == building.Cooker {
			return &gui.Suggestion{
				Message: "Cookers are used to bake bread, brew beer or make medicine.",
				Icon:    "building/cooker", X: LargeIconD*5 + 24, Y: BuildingButtonPanelTop*ControlPanelSY + float64(LargeIconD)*2.5,
			}
		} else if bc.ExtensionT == building.Workshop {
			return &gui.Suggestion{
				Message: "Workshop tools let you transform raw materials like logs or stones to\nbuilding materials, or to make textiles from wool.\nThey are also needed for butcher shops.",
				Icon:    "building/workshop", X: LargeIconD*5 + 24, Y: BuildingButtonPanelTop*ControlPanelSY + float64(LargeIconD)*2.5,
			}
		}
	}
	return nil
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
	bc.GenerateBuildingTypebuttons()
}

func (bc *BuildingsController) GetActiveFields(c *Controller, rf *renderer.RenderedField) []navigation.FieldWithContext {
	if bc.Plan.BuildingType != building.BuildingTypeTownhall && !bc.activeSupplier().FieldWithinDistance(rf.F) {
		return nil
	}
	if bc.Plan.BuildingType == building.BuildingTypeWorkshop && len(bc.Plan.GetExtensions()) == 0 {
		return nil
	}
	return c.Map.GetBuildingBaseFields(rf.F.X, rf.F.Y, bc.Plan, building.DirectionNone)
}

func (bc *BuildingsController) activeSupplier() social.Supplier {
	return bc.cp.C.ActiveSupplier
}

func (bc *BuildingsController) CaptureClick(x, y float64) {
	bc.p.CaptureClick(x, y)
}

func (bc *BuildingsController) CaptureMove(x, y float64) {
	bc.p.CaptureMove(x, y)
}

func (bc *BuildingsController) Clear() {
	bc.p.Clear()
}

func (bc *BuildingsController) Refresh() {
	bc.p.Refresh()
}

func (bc *BuildingsController) Render(cv *canvas.Canvas) {
	bc.p.Render(cv)
}

func (bc *BuildingsController) HandleClick(c *Controller, rf *renderer.RenderedField) bool {
	if bc.activeSupplier() == nil {
		return false
	}
	if bc.Plan.BuildingType != building.BuildingTypeTownhall && !bc.activeSupplier().FieldWithinDistance(rf.F) {
		return false
	}
	if bc.Plan.BuildingType == building.BuildingTypeWorkshop && len(bc.Plan.GetExtensions()) == 0 {
		return false
	}
	if bc.Plan.IsComplete() {
		c.Map.AddBuildingConstruction(bc.activeSupplier(), rf.F.X, rf.F.Y, bc.Plan, bc.Direction)
		return true
	}
	return false
}

func (bc *BuildingsController) GenerateButtons() {
	bc.p.Buttons = nil
	bc.p.Labels = nil

	roofPanelTop := BuildingButtonPanelTop * ControlPanelSY
	bc.p.AddButton(&RoofButton{
		b:   &gui.ButtonGUI{Icon: "cancel", X: float64(LargeIconD)*4 + 24, Y: roofPanelTop, SX: LargeIconS, SY: LargeIconS},
		del: true,
		bc:  bc,
	})
	for i, m := range building.RoofMaterials(bc.bt) {
		bc.p.AddButton(&RoofButton{
			b:  &gui.ButtonGUI{Texture: "building/" + m.Name, X: float64(i)*LargeIconD + 24, Y: roofPanelTop, SX: LargeIconS, SY: LargeIconS},
			m:  m,
			bc: bc,
		})
	}

	floorsPanelTop := BuildingButtonPanelTop*ControlPanelSY + float64(LargeIconD)
	for i, m := range building.FloorMaterials(bc.bt) {
		bc.p.AddButton(&FloorButton{
			b:  &gui.ButtonGUI{Texture: "building/" + m.Name, X: float64(i)*LargeIconD + 24, Y: floorsPanelTop, SX: LargeIconS, SY: LargeIconS},
			m:  m,
			bc: bc,
		})
	}

	extensionPanelTop := BuildingButtonPanelTop*ControlPanelSY + float64(LargeIconD*2)
	for i, e := range building.ExtensionTypes(bc.bt) {
		bc.p.AddButton(&ExtensionButton{
			b:   &gui.ButtonGUI{Icon: "building/" + e.Name, X: float64(i)*LargeIconD + 24, Y: extensionPanelTop, SX: LargeIconS, SY: LargeIconS},
			t:   e,
			bc:  bc,
			msg: e.Description,
		})
	}

	bc.p.AddButton(&RotationButton{
		b:  &gui.ButtonGUI{Icon: "building/rotate_" + strconv.Itoa(int(bc.Direction)), X: LargeIconD*5 + 24, Y: roofPanelTop, SX: LargeIconS, SY: LargeIconS},
		bc: bc,
	})

	for i, ext := range bc.Plan.GetExtensions() {
		bc.p.AddImageLabel("building/"+ext.T.Name, float64(i)*LargeIconD+24, extensionPanelTop+float64(LargeIconD*2), LargeIconS, LargeIconS, gui.ImageLabelStyleRegular)
	}

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
			bc.p.AddButton(createBuildingBaseButton(bc, pi, pj, 0, x, y, nil, nil, nil, nil))
			if bc.Plan.BaseShape[pi][pj] != nil {
				var k int
				for k = range bc.Plan.BaseShape[pi][pj].Floors {
					bc.p.AddButton(createBuildingBaseButton(bc, pi, pj, k+1, x, y-DZ*float64(k+1), bc.Plan.BaseShape[pi][pj].Floors[k].M, nil, nil, nil))
				}
				if bc.Plan.BaseShape[pi][pj].Roof != nil {
					bc.p.AddButton(createBuildingBaseButton(bc, pi, pj, k+1, x, y-DZ*float64(k+1), nil, bc.Plan.BaseShape[pi][pj].Roof.M, &bc.Plan.BaseShape[pi][pj].Roof.RoofType, nil))
				}
				if bc.Plan.BaseShape[pi][pj].Extension != nil {
					bc.p.AddButton(createBuildingBaseButton(bc, pi, pj, k+1, x, y-DZ*float64(k+1), nil, nil, nil, bc.Plan.BaseShape[pi][pj].Extension.T))
				}
			}
		}
	}

	for i, a := range bc.Plan.ConstructionCost() {
		ArtifactsToControlPanel(bc.p, i, a.A, a.Quantity, BuildingCostTop*ControlPanelSY)
	}
}

func CreateBuildingsController(cp *ControlPanel, bt building.BuildingType) *BuildingsController {
	p := &gui.Panel{
		X:           0,
		Y:           ControlPanelDynamicPanelTop * ControlPanelSY,
		SX:          ControlPanelSX - LargeIconD*2,
		SY:          ControlPanelDynamicPanelSY * ControlPanelSY,
		SingleClick: true,
	}
	bc := &BuildingsController{
		Plan:        &building.BuildingPlan{BuildingType: bt},
		bt:          bt,
		p:           p,
		cp:          cp,
		Direction:   building.DirectionN,
		Perspective: &cp.C.Perspective}

	bc.GenerateButtons()

	var helperMsg string
	switch bt {
	case building.BuildingTypeFarm:
		helperMsg = "Build farms to grow food or wood."
	case building.BuildingTypeWorkshop:
		helperMsg = "Build workshops to transform materials to products."
	case building.BuildingTypeMine:
		helperMsg = "Build mines to extract minerals, stone and clay."
	case building.BuildingTypeFactory:
		helperMsg = "Build factories to build vehicles."
	case building.BuildingTypeTownhall:
		helperMsg = "Establish a new town."
	}
	cp.HelperMessage(helperMsg)

	return bc
}

func (bc *BuildingsController) GenerateBuildingTypebuttons() {
	iconTop := 15 + IconS + LargeIconD
	if bc.cp.C.ActiveSupplier != nil && bc.cp.C.ActiveSupplier.BuildHousesEnabled() {
		bc.p.AddButton(&gui.SimpleButton{
			ButtonGUI: gui.ButtonGUI{Icon: "farm", X: float64(24 + LargeIconD*0), Y: iconTop, SX: LargeIconS, SY: LargeIconS},
			Highlight: func() bool { return bc.cp.IsBuildingTypeOf(building.BuildingTypeFarm) },
			ClickImpl: func() { SetupBuildingsController(bc.cp, building.BuildingTypeFarm) }})
		bc.p.AddButton(&gui.SimpleButton{
			ButtonGUI: gui.ButtonGUI{Icon: "mine", X: float64(24 + LargeIconD*1), Y: iconTop, SX: LargeIconS, SY: LargeIconS},
			Highlight: func() bool { return bc.cp.IsBuildingTypeOf(building.BuildingTypeMine) },
			ClickImpl: func() { SetupBuildingsController(bc.cp, building.BuildingTypeMine) }})
		bc.p.AddButton(&gui.SimpleButton{
			ButtonGUI: gui.ButtonGUI{Icon: "workshop", X: float64(24 + LargeIconD*2), Y: iconTop, SX: LargeIconS, SY: LargeIconS},
			Highlight: func() bool { return bc.cp.IsBuildingTypeOf(building.BuildingTypeWorkshop) },
			ClickImpl: func() { SetupBuildingsController(bc.cp, building.BuildingTypeWorkshop) }})
		bc.p.AddButton(&gui.SimpleButton{
			ButtonGUI: gui.ButtonGUI{Icon: "factory", X: float64(24 + LargeIconD*3), Y: iconTop, SX: LargeIconS, SY: LargeIconS},
			Highlight: func() bool { return bc.cp.IsBuildingTypeOf(building.BuildingTypeFactory) },
			ClickImpl: func() { SetupBuildingsController(bc.cp, building.BuildingTypeFactory) }})
	}
	if bc.cp.C.ActiveSupplier != nil {
		bc.p.AddButton(&gui.SimpleButton{
			ButtonGUI: gui.ButtonGUI{Icon: "town", X: float64(24 + LargeIconD*4), Y: iconTop, SX: LargeIconS, SY: LargeIconS},
			Highlight: func() bool { return bc.cp.IsBuildingTypeOf(building.BuildingTypeTownhall) },
			ClickImpl: func() { SetupBuildingsController(bc.cp, building.BuildingTypeTownhall) }})
	}
	if bc.cp.C.ActiveSupplier != nil && bc.cp.C.ActiveSupplier.BuildMarketplaceEnabled() {
		bc.p.AddButton(&gui.SimpleButton{
			ButtonGUI: gui.ButtonGUI{Icon: "market", X: float64(24 + LargeIconD*5), Y: iconTop, SX: LargeIconS, SY: LargeIconS},
			Highlight: func() bool { return bc.cp.IsBuildingTypeOf(building.BuildingTypeMarket) },
			ClickImpl: func() { SetupBuildingsController(bc.cp, building.BuildingTypeMarket) }})
	}
}

func SetupBuildingsController(cp *ControlPanel, bt building.BuildingType) *BuildingsController {
	bc := CreateBuildingsController(cp, bt)
	bc.GenerateBuildingTypebuttons()

	cp.SetDynamicPanel(bc)
	cp.C.ClickHandler = bc
	return bc
}

func BuildingsToControlPanel(cp *ControlPanel) {
	if cp.C.ActiveSupplier != nil && cp.C.ActiveSupplier.BuildHousesEnabled() {
		SetupBuildingsController(cp, building.BuildingTypeFarm)
	} else {
		SetupBuildingsController(cp, building.BuildingTypeTownhall)
	}
}
