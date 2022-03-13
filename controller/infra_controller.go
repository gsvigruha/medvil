package controller

import (
	"github.com/tfriedel6/canvas"
	"image/color"
	"medvil/model/building"
	"medvil/model/terrain"
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
	} else if ic.it == InfraTypeCanal && c.Map.HasNeighborField(rf.F.X, rf.F.Y, terrain.Water) && rf.F.Buildable() {
		c.Map.AddInfraConstruction(c.Country, rf.F.X, rf.F.Y, building.CanalType)
	} else if ic.it == InfraTypeAqueduct {
	} else if ic.it == InfraTypeBridge {
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
		b:  gui.ButtonGUI{Icon: "infra/aqueduct", X: float64(170), Y: float64(InfraPanelTop), SX: 32, SY: 32},
		it: InfraTypeAqueduct,
		ic: ic,
	})

	p.AddButton(InfraBuildButton{
		b:  gui.ButtonGUI{Icon: "infra/bridge", X: float64(210), Y: float64(InfraPanelTop), SX: 32, SY: 32},
		it: InfraTypeBridge,
		ic: ic,
	})

	cp.SetDynamicPanel(p)
	cp.C.ClickHandler = ic
}
