package controller

import (
	"github.com/tfriedel6/canvas"
	"medvil/model/artifacts"
	"medvil/model/building"
	"medvil/model/economy"
	"medvil/model/navigation"
	"medvil/model/social"
	"medvil/renderer"
	"medvil/view/gui"
	"strconv"
)

var NewTownRowH = IconH + int(IconS)

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
		for a, q := range b.c.resources {
			b.c.newTown.Townhall.StorageTarget[a] = q
		}
		srcH.Town.Country.AddTown(b.c.newTown)
	} else {
		if b.state == NewTownControllerStatePickBuildTownhall {
			b.c.bc = CreateBuildingsController(b.c.cp, building.BuildingTypeTownhall, b.c.newTown)
		} else if b.state == NewTownControllerStatePickBuildMarket {
			b.c.bc = CreateBuildingsController(b.c.cp, building.BuildingTypeMarket, b.c.newTown)
		}
		b.c.cp.C.ClickHandler = b.c
	}
}

func (b *NewTownControllerButton) Render(cv *canvas.Canvas) {
	if b.state == NewTownControllerStatePickResources && b.c.newTown.Townhall.Household.Building != nil {
		if b.c.GetResourceVolume() > b.c.newTown.Townhall.Household.Building.Plan.Area()*social.StoragePerArea {
			cv.SetFillStyle("#822")
			cv.FillRect(b.b.X, b.b.Y, b.b.SX, b.b.SY)
		}
	}
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
		var n int = 0
		resources[a] = &n
	}
	newTown := &social.Town{Country: th.Household.Town.Country}
	newTown.Townhall = &social.Townhall{Household: social.Household{Town: newTown}}
	newTown.Marketplace = &social.Marketplace{Town: newTown}
	newTown.Init()
	newTown.Marketplace.Init()
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
		c.p.AddImageLabel("person", 10, 140, IconS, IconS, gui.ImageLabelStyleRegular)
		c.p.AddTextLabel(strconv.Itoa(len(c.sourceTH.Household.People)), 10, float64(140+IconH+4))
		c.p.AddPanel(gui.CreateNumberPanel(10, float64(140+IconH+8), IconS, 20, 0, len(c.sourceTH.Household.People), 1, "%v", c.numPeople).P)

		c.p.AddImageLabel("coin", 50, 140, IconS, IconS, gui.ImageLabelStyleRegular)
		c.p.AddTextLabel(strconv.Itoa(int(c.sourceTH.Household.Money)), 50, float64(140+IconH+4))
		c.p.AddPanel(gui.CreateNumberPanel(50, float64(140+IconH+8), IconS, 20, 0, int(c.sourceTH.Household.Money), 100, "%v", c.money).P)

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
		b: gui.ButtonGUI{Icon: "town", X: float64(10), Y: float64(100), SX: IconS, SY: IconS},
	})
	c.p.AddButton(&NewTownControllerButton{
		c: c, state: NewTownControllerStatePickBuildMarket,
		b: gui.ButtonGUI{Icon: "market", X: float64(50), Y: float64(100), SX: IconS, SY: IconS},
	})
	c.p.AddButton(&NewTownControllerButton{
		c: c, state: NewTownControllerStatePickResources,
		b: gui.ButtonGUI{Icon: "barrel", X: float64(90), Y: float64(100), SX: IconS, SY: IconS},
	})
	c.p.AddButton(&NewTownControllerButton{
		c: c, state: NewTownControllerStateStart,
		b: gui.ButtonGUI{Icon: "start", X: float64(130), Y: float64(100), SX: IconS, SY: IconS},
	})
}

func ArtifactsPickerToControlPanel(c *NewTownController, i int, a *artifacts.Artifact, q uint16, top float64) {
	xI := i % IconRowMax
	yI := i / IconRowMax
	c.p.AddImageLabel("artifacts/"+a.Name, float64(10+xI*IconW), top+float64(yI*NewTownRowH), IconS, IconS, gui.ImageLabelStyleRegular)
	c.p.AddTextLabel(strconv.Itoa(int(q)), float64(10+xI*IconW), top+float64(yI*NewTownRowH+IconH+4))
	c.p.AddPanel(gui.CreateNumberPanel(float64(10+xI*IconW), top+float64(yI*NewTownRowH+IconH+8), IconS, 20, 0, int(q), 5, "%v", c.resources[a]).P)
}

func (ntc *NewTownController) GetResourceVolume() uint16 {
	var v uint16 = 0
	for a, q := range ntc.resources {
		v += a.V * uint16(*q)
	}
	return v
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

func (ntc *NewTownController) GetActiveFields(c *Controller, rf *renderer.RenderedField) []navigation.FieldWithContext {
	if ntc.bc != nil {
		return ntc.bc.GetActiveFields(c, rf)
	}
	return []navigation.FieldWithContext{}
}

func (ntc *NewTownController) HandleClick(c *Controller, rf *renderer.RenderedField) bool {
	if c.ActiveTown == nil {
		return false
	}
	if ntc.bc.Plan.IsComplete() {
		b := c.Map.AddBuilding(rf.F.X, rf.F.Y, ntc.bc.Plan, true, building.DirectionNone)
		if b != nil {
			if ntc.state == NewTownControllerStatePickBuildTownhall {
				ntc.newTown.Townhall.Household.Building = b
				ntc.newTown.Townhall.Household.Resources.VolumeCapacity = b.Plan.Area() * social.StoragePerArea
			} else if ntc.state == NewTownControllerStatePickBuildMarket {
				ntc.newTown.Marketplace.Building = b
				ntc.newTown.Marketplace.Storage.VolumeCapacity = b.Plan.Area() * social.StoragePerArea
			}
			ntc.newTown.CreateBuildingConstruction(b, c.Map)
			for _, as := range ntc.bc.Plan.ConstructionCost() {
				*ntc.resources[as.A] = *ntc.resources[as.A] + int(as.Quantity)
			}
			return true
		} else {
			return false
		}
		return true
	}
	return false
}
