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

type NewTownController struct {
	p     *gui.Panel
	r     map[*artifacts.Artifact]*int
	th    *social.Townhall
	cp    *ControlPanel
	state uint8
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
	c := &NewTownController{p: p, r: r, th: th, cp: cp, state: NewTownControllerStatePickResources}
	SetupNewTownController(c)
	cp.SetDynamicPanel(c)
}

func SetupNewTownController(c *NewTownController) {
	if c.state == NewTownControllerStatePickResources {
		var aI = 0
		for _, a := range artifacts.All {
			if q, ok := c.th.Household.Resources.Artifacts[a]; ok {
				ArtifactsPickerToControlPanel(c, aI, a, q, 100)
				aI++
			}
		}
	} else if c.state == NewTownControllerStatePickBuildTownhall {
		BuildingsToControlPanel(c.cp, building.BuildingTypeTownhall)
	} else if c.state == NewTownControllerStatePickBuildMarket {
		BuildingsToControlPanel(c.cp, building.BuildingTypeMarket)
	}
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
