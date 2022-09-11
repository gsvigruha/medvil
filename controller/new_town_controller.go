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

const NewTownControllerStateBuildTownhall = 1
const NewTownControllerStateBuildMarket = 2
const NewTownControllerStatePickTown = 3
const NewTownControllerStatePickResources = 4
const NewTownControllerStateStart = 5

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
	} else if b.state == NewTownControllerStatePickTown {
		b.c.cp.C.ClickHandler = b.c
		b.c.bc = nil
	} else {
		if b.state == NewTownControllerStateBuildTownhall {
			b.c.bc = CreateBuildingsController(b.c.cp, building.BuildingTypeTownhall, b.c.newTown)
		} else if b.state == NewTownControllerStateBuildMarket {
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
		state:     NewTownControllerStateBuildTownhall,
	}

	SetupNewTownController(c)
	cp.SetDynamicPanel(c)
}

func SetupNewTownController(c *NewTownController) {
	if c.bc != nil {
		c.p.AddPanel(c.bc.p)
	}
	top := 0.15 * ControlPanelSY
	if c.state == NewTownControllerStatePickResources {
		c.p.AddImageLabel("person", 10, top, IconS, IconS, gui.ImageLabelStyleRegular)
		c.p.AddTextLabel(strconv.Itoa(len(c.sourceTH.Household.People)), 10, top+float64(IconH+4))
		c.p.AddPanel(gui.CreateNumberPanel(10, top+float64(IconH+8), IconS, 20, 0, len(c.sourceTH.Household.People), 1, "%v", c.numPeople).P)

		c.p.AddImageLabel("coin", float64(10+IconW), top, IconS, IconS, gui.ImageLabelStyleRegular)
		c.p.AddTextLabel(strconv.Itoa(int(c.sourceTH.Household.Money)), float64(10+IconW), top+float64(IconH+4))
		c.p.AddPanel(gui.CreateNumberPanel(float64(10+IconW), top+float64(IconH+8), IconS, 20, 0, int(c.sourceTH.Household.Money), 100, "%v", c.money).P)

		var aI = 2
		for _, a := range artifacts.All {
			if q, ok := c.sourceTH.Household.Resources.Artifacts[a]; ok {
				ArtifactsPickerToControlPanel(c, aI, a, q, top)
				aI++
			}
		}
	}

	c.p.AddButton(&NewTownControllerButton{
		c: c, state: NewTownControllerStateBuildTownhall,
		b: gui.ButtonGUI{Icon: "town", X: float64(10 + IconW*0), Y: float64(100), SX: IconS, SY: IconS},
	})
	c.p.AddButton(&NewTownControllerButton{
		c: c, state: NewTownControllerStateBuildMarket,
		b: gui.ButtonGUI{Icon: "market", X: float64(10 + IconW*1), Y: float64(100), SX: IconS, SY: IconS},
	})
	c.p.AddButton(&NewTownControllerButton{
		c: c, state: NewTownControllerStatePickTown,
		b: gui.ButtonGUI{Icon: "town", X: float64(10 + IconW*2), Y: float64(100), SX: IconS, SY: IconS},
	})
	c.p.AddButton(&NewTownControllerButton{
		c: c, state: NewTownControllerStatePickResources,
		b: gui.ButtonGUI{Icon: "barrel", X: float64(10 + IconW*3), Y: float64(100), SX: IconS, SY: IconS},
	})
	c.p.AddButton(&NewTownControllerButton{
		c: c, state: NewTownControllerStateStart,
		b: gui.ButtonGUI{Icon: "start", X: float64(10 + IconW*4), Y: float64(100), SX: IconS, SY: IconS},
	})
}

func ArtifactsPickerToControlPanel(c *NewTownController, i int, a *artifacts.Artifact, q uint16, top float64) {
	rowH := IconH + int(IconS)
	xI := i % IconRowMax
	yI := i / IconRowMax
	c.p.AddImageLabel("artifacts/"+a.Name, float64(10+xI*IconW), top+float64(yI*rowH), IconS, IconS, gui.ImageLabelStyleRegular)
	c.p.AddTextLabel(strconv.Itoa(int(q)), float64(10+xI*IconW), top+float64(yI*rowH+IconH+4))
	c.p.AddPanel(gui.CreateNumberPanel(float64(10+xI*IconW), top+float64(yI*rowH+IconH+8), IconS, 20, 0, int(q), 5, "%v", c.resources[a]).P)
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
	if c.ActiveTown != nil && ntc.newTown != nil && ntc.newTown.Townhall.Household.Building != nil {
		var fs []navigation.FieldWithContext
		for _, coords := range c.ActiveTown.Townhall.Household.Building.GetBuildingXYs(true) {
			fs = append(fs, c.Map.GetField(coords[0], coords[1]))
		}
		for _, coords := range ntc.newTown.Townhall.Household.Building.GetBuildingXYs(true) {
			fs = append(fs, c.Map.GetField(coords[0], coords[1]))
		}
		return fs
	}
	return []navigation.FieldWithContext{}
}

func (ntc *NewTownController) HandleClick(c *Controller, rf *renderer.RenderedField) bool {
	if c.ActiveTown == nil {
		return false
	}
	if ntc.state == NewTownControllerStatePickTown {
		if rf.F.Building.GetBuilding() != nil && rf.F.Building.GetBuilding().Plan.BuildingType == building.BuildingTypeTownhall {
			th := c.ReverseReferences.BuildingToTownhall[rf.F.Building.GetBuilding()]
			if th != nil {
				ntc.newTown = th.Household.Town
				return true
			}
		}
	} else if ntc.bc.Plan.IsComplete() {
		b := c.Map.AddBuilding(rf.F.X, rf.F.Y, ntc.bc.Plan, true, building.DirectionNone)
		if b != nil {
			if ntc.state == NewTownControllerStateBuildTownhall {
				ntc.newTown.Townhall.Household.Building = b
				ntc.newTown.Townhall.Household.Resources.VolumeCapacity = b.Plan.Area() * social.StoragePerArea
			} else if ntc.state == NewTownControllerStateBuildMarket {
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
