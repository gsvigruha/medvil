package controller

import (
	"fmt"
	"github.com/tfriedel6/canvas"
	"medvil/model/artifacts"
	"medvil/model/economy"
	"medvil/model/navigation"
	"medvil/model/social"
	"medvil/renderer"
	"medvil/view/gui"
	"strconv"
)

type TownhallControllerButton struct {
	tc       *TownhallController
	b        gui.ButtonGUI
	subPanel *gui.Panel
}

func (b *TownhallControllerButton) Click() {
	b.tc.subPanel = b.subPanel
}

func (b *TownhallControllerButton) Render(cv *canvas.Canvas) {
	b.b.Render(cv)
}

func (b *TownhallControllerButton) Contains(x float64, y float64) bool {
	return b.b.Contains(x, y)
}

type TownhallController struct {
	householdPanel *gui.Panel
	buttons        []*TownhallControllerButton
	subPanel       *gui.Panel
	th             *social.Townhall
}

func TownhallToControlPanel(cp *ControlPanel, th *social.Townhall) {
	hp := &gui.Panel{X: 0, Y: ControlPanelDynamicPanelTop, SX: ControlPanelSX, SY: HouseholdControllerSY}
	tp := &gui.Panel{X: 0, Y: ControlPanelDynamicPanelTop + HouseholdControllerSY, SX: ControlPanelSX, SY: ControlPanelDynamicPanelTop}
	sp := &gui.Panel{X: 0, Y: ControlPanelDynamicPanelTop + HouseholdControllerSY, SX: ControlPanelSX, SY: ControlPanelDynamicPanelTop}
	fp := &gui.Panel{X: 0, Y: ControlPanelDynamicPanelTop + HouseholdControllerSY, SX: ControlPanelSX, SY: ControlPanelDynamicPanelTop}

	tc := &TownhallController{householdPanel: hp, th: th}
	tc.buttons = []*TownhallControllerButton{
		&TownhallControllerButton{tc: tc, subPanel: tp, b: gui.ButtonGUI{Icon: "taxes", X: 10, Y: 560, SX: 32, SY: 32}},
		&TownhallControllerButton{tc: tc, subPanel: sp, b: gui.ButtonGUI{Icon: "barrel", X: 50, Y: 560, SX: 32, SY: 32}},
		&TownhallControllerButton{tc: tc, subPanel: fp, b: gui.ButtonGUI{Icon: "factory", X: 90, Y: 560, SX: 32, SY: 32}},
	}

	HouseholdToControlPanel(hp, &th.Household)

	tp.AddPanel(gui.CreateNumberPanel(10, 600, 120, 20, 0, 100, 10, "farm tax rate %v", &th.Household.Town.Transfers.Farm.TaxRate).P)
	tp.AddPanel(gui.CreateNumberPanel(10, 625, 120, 20, 0, 1000, 10, "farm threshold %v", &th.Household.Town.Transfers.Farm.TaxThreshold).P)
	tp.AddPanel(gui.CreateNumberPanel(10, 650, 120, 20, 0, 1000, 10, "farm subsidy %v", &th.Household.Town.Transfers.Farm.Subsidy).P)

	tp.AddPanel(gui.CreateNumberPanel(10, 675, 120, 20, 0, 100, 10, "shop tax rate %v", &th.Household.Town.Transfers.Workshop.TaxRate).P)
	tp.AddPanel(gui.CreateNumberPanel(10, 700, 120, 20, 0, 1000, 10, "shop threshold %v", &th.Household.Town.Transfers.Workshop.TaxThreshold).P)
	tp.AddPanel(gui.CreateNumberPanel(10, 725, 120, 20, 0, 1000, 10, "shop subsidy %v", &th.Household.Town.Transfers.Workshop.Subsidy).P)

	tp.AddPanel(gui.CreateNumberPanel(150, 600, 120, 20, 0, 100, 10, "mine tax rate %v", &th.Household.Town.Transfers.Mine.TaxRate).P)
	tp.AddPanel(gui.CreateNumberPanel(150, 625, 120, 20, 0, 1000, 10, "mine threshold %v", &th.Household.Town.Transfers.Mine.TaxThreshold).P)
	tp.AddPanel(gui.CreateNumberPanel(150, 650, 120, 20, 0, 1000, 10, "mine subsidy %v", &th.Household.Town.Transfers.Mine.Subsidy).P)

	tp.AddPanel(gui.CreateNumberPanel(150, 675, 120, 20, 0, 100, 10, "factory tax rate %v", &th.Household.Town.Transfers.Factory.TaxRate).P)
	tp.AddPanel(gui.CreateNumberPanel(150, 700, 120, 20, 0, 1000, 10, "factory threshold %v", &th.Household.Town.Transfers.Factory.TaxThreshold).P)
	tp.AddPanel(gui.CreateNumberPanel(150, 725, 120, 20, 0, 1000, 10, "factory subsidy %v", &th.Household.Town.Transfers.Factory.Subsidy).P)

	tp.AddPanel(gui.CreateNumberPanel(10, 750, 120, 20, 0, 100, 10, "military funding %v", &th.Household.Town.Transfers.Tower.Subsidy).P)
	tp.AddPanel(gui.CreateNumberPanel(150, 750, 120, 20, 0, 100, 10, "market funding %v", &th.Household.Town.Transfers.MarketFundingRate).P)

	var aI = 0
	for _, a := range artifacts.All {
		if q, ok := th.Household.Resources.Artifacts[a]; ok {
			ArtifactStorageToControlPanel(sp, th, aI, a, q, 600)
			aI++
		}
	}

	for i, vc := range social.GetVehicleConstructions(th.Household.Town.Factories) {
		fp.AddPanel(CreateOrderPanelForTownhall(10, float64(i*40+600), 60, 20, th, vc))
	}

	cp.SetDynamicPanel(tc)
	cp.C.ClickHandler = tc
}

func ArtifactStorageToControlPanel(p *gui.Panel, th *social.Townhall, i int, a *artifacts.Artifact, q uint16, top float64) {
	xI := i % IconRowMax
	yI := i / IconRowMax
	p.AddImageLabel("artifacts/"+a.Name, float64(10+xI*IconW), top+float64(yI)*NewTownRowH, 32, 32, gui.ImageLabelStyleRegular)
	p.AddTextLabel(strconv.Itoa(int(q)), float64(10+xI*IconW), top+float64(yI)*NewTownRowH+IconH+4)
	p.AddPanel(gui.CreateNumberPanel(float64(10+xI*IconW), top+float64(yI)*NewTownRowH+IconH+8, 32, 20, 0, 100, 5, "%v", th.StorageTarget[a]).P)
}

func (tc *TownhallController) CaptureClick(x, y float64) {
	tc.householdPanel.CaptureClick(x, y)
}

func (tc *TownhallController) Render(cv *canvas.Canvas) {
	tc.householdPanel.Render(cv)
}

func (tc *TownhallController) Clear() {}

func (tc *TownhallController) Refresh() {
	tc.householdPanel.Clear()
	HouseholdToControlPanel(tc.householdPanel, &tc.th.Household)
	for _, button := range tc.buttons {
		tc.householdPanel.AddButton(button)
	}
	if tc.subPanel != nil {
		tc.householdPanel.AddPanel(tc.subPanel)
	}
}

func (tc *TownhallController) GetActiveFields(c *Controller, rf *renderer.RenderedField) []navigation.FieldWithContext {
	return tc.th.GetFields()
}

func (tc *TownhallController) HandleClick(c *Controller, rf *renderer.RenderedField) bool {
	for i := range tc.th.Household.Town.Roads {
		r := tc.th.Household.Town.Roads[i]
		if r.X == rf.F.X && r.Y == rf.F.Y {
			r.Allocated = false
			tc.th.Household.Town.Roads = append(tc.th.Household.Town.Roads[:i], tc.th.Household.Town.Roads[i+1:]...)
			return true
		}
	}
	if !rf.F.Allocated && rf.F.Road != nil {
		tc.th.Household.Town.Roads = append(tc.th.Household.Town.Roads, rf.F)
		rf.F.Allocated = true
		return true
	}
	return false
}

func CreateOrderPanelForTownhall(x, y, sx, sy float64, th *social.Townhall, vc *economy.VehicleConstruction) *gui.Panel {
	p := &gui.Panel{}
	l := p.AddTextLabel("", x, y+sy*2/3)
	var factories []*social.Factory
	for _, factory := range th.Household.Town.Factories {
		if economy.ConstructionCompatible(vc, factory.Household.Building.Plan.GetExtension()) {
			factories = append(factories, factory)
		}
	}
	p.AddButton(OrderButton{
		b:         gui.ButtonGUI{Icon: "plus", X: x + sx, Y: y, SX: sy, SY: sy},
		factories: factories,
		vc:        vc,
		l:         l,
	})
	p.AddTextLabel(fmt.Sprintf("$%v", factories[0].Price(vc)), x+sx+sy*2, y+sy*2/3)
	return p
}
