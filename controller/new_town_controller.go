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

const NewTownControllerStateNone = 0
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
	b.c.SetToState()
}

func (b *NewTownControllerButton) Enabled() bool {
	return b.b.Enabled()
}

func (ntc *NewTownController) SetToState() {
	if ntc.state == NewTownControllerStatePickResources {
		ntc.bc = nil
	} else if ntc.state == NewTownControllerStateStart {
		ntc.cp.C.ClickHandler = nil
		srcH := ntc.sourceTH.Household
		dstH := ntc.newTown.Townhall.Household
		for a, q := range ntc.resources {
			if q > 0 {
				srcH.AddTask(&economy.TransportTask{
					PickupD:  ntc.cp.C.Map.GetField(srcH.Building.X, srcH.Building.Y),
					DropoffD: ntc.cp.C.Map.GetField(dstH.Building.X, dstH.Building.Y),
					PickupR:  &srcH.Resources,
					DropoffR: &dstH.Resources,
					A:        a,
					Quantity: uint16(q),
				})
			}
		}
		for i := 0; i < *ntc.numPeople; i++ {
			srcH.ReassignFirstPerson(dstH, ntc.cp.C.Map)
			if len(dstH.People) > int(dstH.TargetNumPeople) {
				dstH.TargetNumPeople++
			}
		}
		targetMoney := uint32(*ntc.money)
		if srcH.Money > targetMoney {
			dstH.Money += targetMoney
			srcH.Money -= targetMoney
		} else {
			dstH.Money += srcH.Money
			srcH.Money = 0
		}
		for a, q := range ntc.resources {
			if q2, ok := ntc.newTown.Townhall.StorageTarget[a]; ok && q2 < q {
				ntc.newTown.Townhall.StorageTarget[a] = q
			}
		}
		srcH.Town.Country.AddTownIfDoesNotExist(ntc.newTown)
	} else if ntc.state == NewTownControllerStatePickTown {
		ntc.cp.C.ClickHandler = ntc
		ntc.bc = nil
	} else {
		if ntc.state == NewTownControllerStateBuildTownhall {
			ntc.bc = CreateBuildingsController(ntc.cp, ntc.p, building.BuildingTypeTownhall, ntc.newTown)
		} else if ntc.state == NewTownControllerStateBuildMarket {
			ntc.bc = CreateBuildingsController(ntc.cp, ntc.p, building.BuildingTypeMarket, ntc.newTown)
		}
		ntc.cp.C.ClickHandler = ntc
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
	resources map[*artifacts.Artifact]int
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
	resources := make(map[*artifacts.Artifact]int)
	for _, a := range artifacts.All {
		var n int = 0
		resources[a] = n
	}
	newTown := &social.Town{Country: th.Household.Town.Country}
	newTown.Townhall = &social.Townhall{Household: &social.Household{Town: newTown}}
	newTown.Marketplace = &social.Marketplace{Town: newTown}
	newTown.Init()
	newTown.Marketplace.Init()
	var money int = 100
	var numPeople int = 2
	c := &NewTownController{
		p:         p,
		resources: resources,
		money:     &money,
		numPeople: &numPeople,
		sourceTH:  th,
		newTown:   newTown,
		cp:        cp,
		state:     NewTownControllerStateNone,
	}

	SetupNewTownController(c)
	cp.SetDynamicPanel(c)
}

func SetupNewTownController(c *NewTownController) {
	if c.bc != nil {
		c.p.AddPanel(c.bc.p)
	}
	resTop := 0.15 * ControlPanelSY
	if c.state == NewTownControllerStatePickResources {
		c.p.AddImageLabel("person", 10, resTop, IconS, IconS, gui.ImageLabelStyleRegular)
		c.p.AddTextLabel(strconv.Itoa(len(c.sourceTH.Household.People)), 10, resTop+float64(IconH+4))
		c.p.AddPanel(gui.CreateNumberPaneFromVal(10, resTop+float64(IconH+8), IconS, 20, 0, len(c.sourceTH.Household.People), 1, "%v", c.numPeople).P)

		c.p.AddImageLabel("coin", float64(10+IconW), resTop, IconS, IconS, gui.ImageLabelStyleRegular)
		c.p.AddTextLabel(strconv.Itoa(int(c.sourceTH.Household.Money)), float64(10+IconW), resTop+float64(IconH+4))
		c.p.AddPanel(gui.CreateNumberPaneFromVal(float64(10+IconW), resTop+float64(IconH+8), IconS, 20, 0, int(c.sourceTH.Household.Money), 100, "%v", c.money).P)

		var aI = 2
		for _, a := range artifacts.All {
			if q, ok := c.sourceTH.Household.Resources.Artifacts[a]; ok {
				ArtifactsPickerToControlPanel(c, aI, a, q, resTop)
				aI++
			}
		}
	}

	top := float64(LargeIconD) + 50
	c.p.AddButton(&NewTownControllerButton{
		c: c, state: NewTownControllerStateBuildTownhall,
		b: gui.ButtonGUI{Icon: "new_town", X: float64(10 + LargeIconD*0), Y: top, SX: LargeIconS, SY: LargeIconS,
			Disabled: func() bool {
				return c.state != NewTownControllerStateNone && c.state != NewTownControllerStateBuildTownhall
			}},
	})
	c.p.AddButton(&NewTownControllerButton{
		c: c, state: NewTownControllerStateBuildMarket,
		b: gui.ButtonGUI{Icon: "new_market", X: float64(10 + LargeIconD*1), Y: top, SX: LargeIconS, SY: LargeIconS,
			Disabled: func() bool {
				return c.state != NewTownControllerStateBuildMarket
			}},
	})
	c.p.AddButton(&NewTownControllerButton{
		c: c, state: NewTownControllerStatePickTown,
		b: gui.ButtonGUI{Icon: "town", X: float64(10 + LargeIconD*0), Y: top + float64(LargeIconD), SX: LargeIconS, SY: LargeIconS,
			Disabled: func() bool {
				return c.state != NewTownControllerStateNone && c.state != NewTownControllerStatePickTown
			}},
	})
	c.p.AddButton(&NewTownControllerButton{
		c: c, state: NewTownControllerStatePickResources,
		b: gui.ButtonGUI{Icon: "barrel", X: float64(10 + LargeIconD*2), Y: top + float64(LargeIconD/2), SX: LargeIconS, SY: LargeIconS,
			Disabled: func() bool {
				return c.state != NewTownControllerStatePickResources
			}},
	})
	c.p.AddButton(&NewTownControllerButton{
		c: c, state: NewTownControllerStateStart,
		b: gui.ButtonGUI{Icon: "start", X: float64(10 + LargeIconD*3), Y: top + float64(LargeIconD/2), SX: LargeIconS, SY: LargeIconS,
			Disabled: func() bool {
				return c.state != NewTownControllerStatePickResources
			}},
	})
}

func ArtifactsPickerToControlPanel(c *NewTownController, i int, a *artifacts.Artifact, q uint16, top float64) {
	rowH := IconH + int(IconS)
	xI := i % IconRowMax
	yI := i / IconRowMax
	c.p.AddImageLabel("artifacts/"+a.Name, float64(10+xI*IconW), top+float64(yI*rowH), IconS, IconS, gui.ImageLabelStyleRegular)
	c.p.AddTextLabel(strconv.Itoa(int(q)), float64(10+xI*IconW), top+float64(yI*rowH+IconH+4))
	c.p.AddPanel(gui.CreateNumberPanel(float64(10+xI*IconW), top+float64(yI*rowH+IconH+8), IconS, 20, 0, int(q), 5, "%v",
		func() int { return c.resources[a] },
		func(v int) { c.resources[a] = v }).P)
}

func (ntc *NewTownController) GetResourceVolume() uint16 {
	var v uint16 = 0
	for a, q := range ntc.resources {
		v += a.V * uint16(q)
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
				ntc.state = NewTownControllerStatePickResources
				ntc.SetToState()
				return true
			}
		}
	} else if ntc.bc.Plan.IsComplete() {
		b := c.Map.AddBuilding(rf.F.X, rf.F.Y, ntc.bc.Plan, true, building.DirectionNone)
		if b != nil {
			if ntc.state == NewTownControllerStateBuildTownhall {
				ntc.newTown.Townhall.Household.Building = b
				ntc.newTown.Townhall.Household.Resources.VolumeCapacity = b.Plan.Area() * social.StoragePerArea
				ntc.state = NewTownControllerStateBuildMarket
			} else if ntc.state == NewTownControllerStateBuildMarket {
				ntc.newTown.Marketplace.Building = b
				ntc.newTown.Marketplace.Storage.VolumeCapacity = b.Plan.Area() * social.StoragePerArea
				ntc.state = NewTownControllerStatePickResources
			}
			ntc.newTown.CreateBuildingConstruction(b, c.Map)
			for _, as := range ntc.bc.Plan.ConstructionCost() {
				ntc.resources[as.A] = ntc.resources[as.A] + int(as.Quantity)
			}
			ntc.SetToState()
			return true
		} else {
			return false
		}
		return true
	}
	return false
}
