package controller

import (
	"github.com/tfriedel6/canvas"
	"medvil/model/artifacts"
	"medvil/model/building"
	"medvil/model/social"
	"medvil/view/gui"
	"strconv"
)

const NewTownRowH = IconH + 32

const NewTownControllerStatePickResources = 1
const NewTownControllerStatePickBuildTownhall = 2
const NewTownControllerStatePickBuildMarket = 3

type NewTownControllerButton struct {
	c     *NewTownController
	b     gui.ButtonGUI
	state uint8
}

func (b *NewTownControllerButton) Click() {
	b.c.state = b.state
	if b.state == NewTownControllerStatePickResources {
		b.c.bc = nil
	} else {
		if b.state == NewTownControllerStatePickBuildTownhall {
			b.c.bc = CreateBuildingsController(b.c.cp, building.BuildingTypeTownhall)
		} else if b.state == NewTownControllerStatePickBuildMarket {
			b.c.bc = CreateBuildingsController(b.c.cp, building.BuildingTypeMarket)
		}
		b.c.cp.C.ActiveBuildingPlan = b.c.bc.Plan
		b.c.cp.C.ClickHandler = b.c.bc
	}
}

func (b *NewTownControllerButton) Render(cv *canvas.Canvas) {
	b.b.Render(cv)
}

func (b *NewTownControllerButton) Contains(x float64, y float64) bool {
	return b.b.Contains(x, y)
}

type NewTownController struct {
	p        *gui.Panel
	bc       *BuildingsController
	r        map[*artifacts.Artifact]*int
	sourceTH *social.Townhall
	newTown  *social.Town
	cp       *ControlPanel
	state    uint8
}

func NewTownToControlPanel(cp *ControlPanel, th *social.Townhall) {
	if th == nil {
		return
	}
	p := &gui.Panel{X: 0, Y: ControlPanelDynamicPanelTop, SX: ControlPanelSX, SY: HouseholdControllerSY}
	r := make(map[*artifacts.Artifact]*int)
	for _, a := range artifacts.All {
		var n int
		r[a] = &n
	}
	newTown := &social.Town{Country: th.Household.Town.Country}
	cp.C.ActiveTown = newTown
	c := &NewTownController{p: p, r: r, sourceTH: th, newTown: newTown, cp: cp, state: NewTownControllerStatePickResources}

	SetupNewTownController(c)
	cp.SetDynamicPanel(c)
}

func SetupNewTownController(c *NewTownController) {
	if c.bc != nil {
		c.p.AddPanel(c.bc.p)
	}
	if c.state == NewTownControllerStatePickResources {
		var aI = 0
		for _, a := range artifacts.All {
			if q, ok := c.sourceTH.Household.Resources.Artifacts[a]; ok {
				ArtifactsPickerToControlPanel(c, aI, a, q, 140)
				aI++
			}
		}
	}
	c.p.AddButton(&NewTownControllerButton{
		c: c, state: NewTownControllerStatePickResources,
		b: gui.ButtonGUI{Icon: "barrel", X: float64(10), Y: float64(100), SX: 32, SY: 32},
	})
	c.p.AddButton(&NewTownControllerButton{
		c: c, state: NewTownControllerStatePickBuildTownhall,
		b: gui.ButtonGUI{Icon: "town", X: float64(50), Y: float64(100), SX: 32, SY: 32},
	})
	c.p.AddButton(&NewTownControllerButton{
		c: c, state: NewTownControllerStatePickBuildMarket,
		b: gui.ButtonGUI{Icon: "market", X: float64(90), Y: float64(100), SX: 32, SY: 32},
	})
}

func ArtifactsPickerToControlPanel(c *NewTownController, i int, a *artifacts.Artifact, q uint16, top float64) {
	xI := i % IconRowMax
	yI := i / IconRowMax
	c.p.AddImageLabel("artifacts/"+a.Name, float64(10+xI*IconW), top+float64(yI)*NewTownRowH, 32, 32, gui.ImageLabelStyleRegular)
	c.p.AddTextLabel(strconv.Itoa(int(q)), float64(10+xI*IconW), top+float64(yI)*NewTownRowH+IconH+4)
	c.p.AddPanel(gui.CreateNumberPanel(float64(10+xI*IconW), top+float64(yI)*NewTownRowH+IconH+8, 32, 20, 0, int(q), 5, "%v", c.r[a]).P)
}

func (c *NewTownController) CaptureClick(x, y float64) {
	c.p.CaptureClick(x, y)
}

func (c *NewTownController) Render(cv *canvas.Canvas) {
	c.p.Render(cv)
}

func (c *NewTownController) Clear() {}

func (c *NewTownController) Refresh() {
	c.p.Clear()
	SetupNewTownController(c)
}
