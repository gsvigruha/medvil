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
	if ic.it == InfraTypeNone {
		rf.F.Road = building.Road{T: nil}
	} else if ic.it == InfraTypeDirtRoad && rf.F.Walkable() && rf.F.Buildable() {
		rf.F.Road = building.Road{T: building.DirtRoadType}
	} else if ic.it == InfraTypeCobbleRoad && rf.F.Walkable() && rf.F.Buildable() {
		rf.F.Road = building.Road{T: building.CobbleRoadType}
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

	cp.SetDynamicPanel(p)
	cp.C.ClickHandler = ic
}
