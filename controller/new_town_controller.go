package controller

import (
	"github.com/tfriedel6/canvas"
	"medvil/model/artifacts"
	"medvil/model/building"
	"medvil/model/economy"
	"medvil/model/social"
	"medvil/renderer"
	"medvil/view/gui"
	"strconv"
)

const NewTownRowH = IconH + 32

const NewTownControllerStatePickBuildTownhall = 1
const NewTownControllerStatePickBuildMarket = 2
const NewTownControllerStatePickResources = 3
const NewTownControllerStateStart = 4

type NewTownControllerButton struct {
	c     *NewTownController
	b     gui.ButtonGUI
	state uint8
}

func (b *NewTownControllerButton) Click() {
	b.c.state = b.state
	if b.state == NewTownControllerStatePickResources {
		b.c.bc = nil
	} else if b.state == NewTownControllerStateStart {
		b.c.cp.C.ActiveBuildingPlan = nil
		b.c.cp.C.ClickHandler = nil
		srcH := &b.c.sourceTH.Household
		dstH := &b.c.newTown.Townhall.Household
		for a, q := range b.c.resources {
			if *q > 0 {
				srcH.AddTask(&economy.TransportTask{
					PickupF:  b.c.cp.C.Map.GetField(srcH.Building.X, srcH.Building.Y),
					DropoffF: b.c.cp.C.Map.GetField(dstH.Building.X, dstH.Building.Y),
					PickupR:  &srcH.Resources,
					DropoffR: &dstH.Resources,
					A:        a,
					Quantity: uint16(*q),
				})
			}
		}
		for i := 0; i < *b.c.numPeople; i++ {
			srcH.ReassignFirstPerson(dstH, b.c.cp.C.Map)
			dstH.TargetNumPeople++
		}
		targetMoney := uint32(*b.c.money)
		if srcH.Money > targetMoney {
			dstH.Money += targetMoney
			srcH.Money -= targetMoney
		} else {
			dstH.Money += srcH.Money
			srcH.Money = 0
		}
		srcH.Town.Country.AddTown(b.c.newTown)
	} else {
		if b.state == NewTownControllerStatePickBuildTownhall {
			b.c.bc = CreateBuildingsController(b.c.cp, building.BuildingTypeTownhall)
		} else if b.state == NewTownControllerStatePickBuildMarket {
			b.c.bc = CreateBuildingsController(b.c.cp, building.BuildingTypeMarket)
		}
		b.c.cp.C.ActiveBuildingPlan = b.c.bc.Plan
		b.c.cp.C.ClickHandler = b.c
	}
}

func (b *NewTownControllerButton) Render(cv *canvas.Canvas) {
	b.b.Render(cv)
}

func (b *NewTownControllerButton) Contains(x float64, y float64) bool {
	return b.b.Contains(x, y)
}

type NewTownController struct {
	p         *gui.Panel
	bc        *BuildingsController
	resources map[*artifacts.Artifact]*int
	numPeople *int
	money     *int
	sourceTH  *social.Townhall
	newTown   *social.Town
	cp        *ControlPanel
	state     uint8
}

func NewTownToControlPanel(cp *ControlPanel, th *social.Townhall) {
	if th == nil {
		return
	}
	p := &gui.Panel{X: 0, Y: ControlPanelDynamicPanelTop, SX: ControlPanelSX, SY: HouseholdControllerSY}
	resources := make(map[*artifacts.Artifact]*int)
	for _, a := range artifacts.All {
		var n int
		resources[a] = &n
	}
	newTown := &social.Town{Country: th.Household.Town.Country}
	newTown.Townhall = &social.Townhall{Household: social.Household{Town: newTown}}
	newTown.Marketplace = &social.Marketplace{Town: newTown}
	newTown.Init()
	newTown.Marketplace.Init()
	cp.C.ActiveTown = newTown
	var money int
	var numPeople int
	c := &NewTownController{
		p:         p,
		resources: resources,
		money:     &money,
		numPeople: &numPeople,
		sourceTH:  th,
		newTown:   newTown,
		cp:        cp,
		state:     NewTownControllerStatePickBuildTownhall,
	}

	SetupNewTownController(c)
	cp.SetDynamicPanel(c)
}

func SetupNewTownController(c *NewTownController) {
	if c.bc != nil {
		c.p.AddPanel(c.bc.p)
	}
	if c.state == NewTownControllerStatePickResources {
		c.p.AddImageLabel("person", 10, 140, 32, 32, gui.ImageLabelStyleRegular)
		c.p.AddTextLabel(strconv.Itoa(len(c.sourceTH.Household.People)), 10, 140+IconH+4)
		c.p.AddPanel(gui.CreateNumberPanel(10, 140+IconH+8, 32, 20, 0, len(c.sourceTH.Household.People), 1, "%v", c.numPeople).P)

		c.p.AddImageLabel("artifacts/gold_coin", 50, 140, 32, 32, gui.ImageLabelStyleRegular)
		c.p.AddTextLabel(strconv.Itoa(int(c.sourceTH.Household.Money)), 50, 140+IconH+4)
		c.p.AddPanel(gui.CreateNumberPanel(50, 140+IconH+8, 32, 20, 0, int(c.sourceTH.Household.Money), 100, "%v", c.money).P)

		var aI = 2
		for _, a := range artifacts.All {
			if q, ok := c.sourceTH.Household.Resources.Artifacts[a]; ok {
				ArtifactsPickerToControlPanel(c, aI, a, q, 140)
				aI++
			}
		}
	}

	c.p.AddButton(&NewTownControllerButton{
		c: c, state: NewTownControllerStatePickBuildTownhall,
		b: gui.ButtonGUI{Icon: "town", X: float64(10), Y: float64(100), SX: 32, SY: 32},
	})
	c.p.AddButton(&NewTownControllerButton{
		c: c, state: NewTownControllerStatePickBuildMarket,
		b: gui.ButtonGUI{Icon: "market", X: float64(50), Y: float64(100), SX: 32, SY: 32},
	})
	c.p.AddButton(&NewTownControllerButton{
		c: c, state: NewTownControllerStatePickResources,
		b: gui.ButtonGUI{Icon: "barrel", X: float64(90), Y: float64(100), SX: 32, SY: 32},
	})
	c.p.AddButton(&NewTownControllerButton{
		c: c, state: NewTownControllerStateStart,
		b: gui.ButtonGUI{Icon: "start", X: float64(130), Y: float64(100), SX: 32, SY: 32},
	})
}

func ArtifactsPickerToControlPanel(c *NewTownController, i int, a *artifacts.Artifact, q uint16, top float64) {
	xI := i % IconRowMax
	yI := i / IconRowMax
	c.p.AddImageLabel("artifacts/"+a.Name, float64(10+xI*IconW), top+float64(yI)*NewTownRowH, 32, 32, gui.ImageLabelStyleRegular)
	c.p.AddTextLabel(strconv.Itoa(int(q)), float64(10+xI*IconW), top+float64(yI)*NewTownRowH+IconH+4)
	c.p.AddPanel(gui.CreateNumberPanel(float64(10+xI*IconW), top+float64(yI)*NewTownRowH+IconH+8, 32, 20, 0, int(q), 5, "%v", c.resources[a]).P)
}

func (ntc *NewTownController) CaptureClick(x, y float64) {
	ntc.p.CaptureClick(x, y)
}

func (ntc *NewTownController) Render(cv *canvas.Canvas) {
	ntc.p.Render(cv)
}

func (ntc *NewTownController) Clear() {}

func (ntc *NewTownController) Refresh() {
	ntc.p.Clear()
	SetupNewTownController(ntc)
}

func (ntc *NewTownController) HandleClick(c *Controller, rf *renderer.RenderedField) bool {
	if c.ActiveTown == nil {
		return false
	}
	if c.ActiveBuildingPlan.IsComplete() {
		b := c.Map.AddBuilding(rf.F.X, rf.F.Y, c.ActiveBuildingPlan, true)
		if b != nil {
			if ntc.state == NewTownControllerStatePickBuildTownhall {
				ntc.newTown.Townhall.Household.Building = b
				ntc.newTown.Townhall.Household.Resources.VolumeCapacity = b.Plan.Area() * social.StoragePerArea
			} else if ntc.state == NewTownControllerStatePickBuildMarket {
				ntc.newTown.Marketplace.Building = b
				ntc.newTown.Marketplace.Storage.VolumeCapacity = b.Plan.Area() * social.StoragePerArea
			}
			ntc.newTown.CreateBuildingConstruction(b, c.Map)
			return true
		} else {
			return false
		}
		return true
	}
	return false
}
