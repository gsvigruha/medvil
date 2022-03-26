package controller

import (
	"github.com/tfriedel6/canvas"
	"image/color"
	"medvil/model/building"
	"medvil/renderer"
	"medvil/view/gui"
)

type InfraType uint8

const InfraTypeNone = 0
const InfraTypeDirtRoad = 1
const InfraTypeCobbleRoad = 2
const InfraTypeCanal = 3
const InfraTypeAqueduct = 4
const InfraTypeBridge = 5
const InfraTypeStoneWall1 = 11
const InfraTypeStoneWall2 = 12
const InfraTypeStoneWall3 = 13
const InfraTypeStoneWallRamp = 14

const InfraPanelTop = 100

type InfraController struct {
	it InfraType
}

type InfraBuildButton struct {
	b  gui.ButtonGUI
	it InfraType
	ic *InfraController
}

func (b InfraBuildButton) Click() {
	b.ic.it = b.it
}

func (b InfraBuildButton) Render(cv *canvas.Canvas) {
	b.b.Render(cv)
	if b.ic.it != b.it {
		cv.SetFillStyle(color.RGBA{R: 64, G: 0, B: 0, A: 128})
		cv.FillRect(b.b.X, b.b.Y, 32, 32)
	}
}

func (b InfraBuildButton) Contains(x float64, y float64) bool {
	return b.b.Contains(x, y)
}

func (ic *InfraController) HandleClick(c *Controller, rf *renderer.RenderedField) bool {
	if ic.it == InfraTypeDirtRoad && rf.F.Walkable() && rf.F.Buildable() {
		c.Map.AddRoadConstruction(c.Country, rf.F.X, rf.F.Y, building.DirtRoadType)
	} else if ic.it == InfraTypeCobbleRoad && rf.F.Walkable() && rf.F.Buildable() {
		c.Map.AddRoadConstruction(c.Country, rf.F.X, rf.F.Y, building.CobbleRoadType)
	} else if ic.it == InfraTypeCanal && rf.F.Buildable() {
		c.Map.AddInfraConstruction(c.Country, rf.F.X, rf.F.Y, building.CanalType)
	} else if ic.it == InfraTypeBridge && c.Map.Shore(rf.F.X, rf.F.Y) {
		c.Map.AddRoadConstruction(c.Country, rf.F.X, rf.F.Y, building.BridgeRoadType)
	} else if ic.it == InfraTypeStoneWall1 && rf.F.Buildable() {
		c.Map.AddBuildingConstruction(c.Country, rf.F.X, rf.F.Y, building.StoneWall1Type)
	} else if ic.it == InfraTypeStoneWall2 && rf.F.Buildable() {
		c.Map.AddBuildingConstruction(c.Country, rf.F.X, rf.F.Y, building.StoneWall2Type)
	} else if ic.it == InfraTypeStoneWall3 && rf.F.Buildable() {
		c.Map.AddBuildingConstruction(c.Country, rf.F.X, rf.F.Y, building.StoneWall3Type)
	} else if ic.it == InfraTypeStoneWallRamp {
		c.Map.AddWallRampConstruction(c.Country, rf.F.X, rf.F.Y)
	}
	return true
}

func InfraToControlPanel(cp *ControlPanel) {
	p := &gui.Panel{X: 0, Y: ControlPanelDynamicPanelTop, SX: ControlPanelSX, SY: ControlPanelDynamicPanelSY}
	ic := &InfraController{it: InfraTypeNone}

	p.AddButton(InfraBuildButton{
		b:  gui.ButtonGUI{Texture: "terrain/grass", X: float64(10), Y: float64(InfraPanelTop), SX: 32, SY: 32},
		it: InfraTypeNone,
		ic: ic,
	})

	p.AddButton(InfraBuildButton{
		b:  gui.ButtonGUI{Texture: "infra/dirt_road", X: float64(50), Y: float64(InfraPanelTop), SX: 32, SY: 32},
		it: InfraTypeDirtRoad,
		ic: ic,
	})

	p.AddButton(InfraBuildButton{
		b:  gui.ButtonGUI{Texture: "infra/cobble_road", X: float64(90), Y: float64(InfraPanelTop), SX: 32, SY: 32},
		it: InfraTypeCobbleRoad,
		ic: ic,
	})

	p.AddButton(InfraBuildButton{
		b:  gui.ButtonGUI{Texture: "infra/canal", X: float64(130), Y: float64(InfraPanelTop), SX: 32, SY: 32},
		it: InfraTypeCanal,
		ic: ic,
	})

	p.AddButton(InfraBuildButton{
		b:  gui.ButtonGUI{Icon: "infra/bridge", X: float64(210), Y: float64(InfraPanelTop), SX: 32, SY: 32},
		it: InfraTypeBridge,
		ic: ic,
	})

	p.AddButton(InfraBuildButton{
		b:  gui.ButtonGUI{Icon: "infra/stone_wall", X: float64(10), Y: float64(InfraPanelTop + 50), SX: 32, SY: 32},
		it: InfraTypeStoneWall1,
		ic: ic,
	})

	p.AddButton(InfraBuildButton{
		b:  gui.ButtonGUI{Icon: "infra/stone_wall", X: float64(50), Y: float64(InfraPanelTop + 50), SX: 32, SY: 32},
		it: InfraTypeStoneWall2,
		ic: ic,
	})

	p.AddButton(InfraBuildButton{
		b:  gui.ButtonGUI{Icon: "infra/stone_wall", X: float64(90), Y: float64(InfraPanelTop + 50), SX: 32, SY: 32},
		it: InfraTypeStoneWall3,
		ic: ic,
	})

	p.AddButton(InfraBuildButton{
		b:  gui.ButtonGUI{Icon: "infra/stone_wall", X: float64(130), Y: float64(InfraPanelTop + 50), SX: 32, SY: 32},
		it: InfraTypeStoneWallRamp,
		ic: ic,
	})

	cp.SetDynamicPanel(p)
	cp.C.ClickHandler = ic
}
