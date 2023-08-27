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

var TownhallControllerGUIBottomY = 0.75

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

func (b *TownhallControllerButton) Enabled() bool {
	return b.b.Enabled()
}

type TownhallController struct {
	cp             *ControlPanel
	topPanel       *gui.Panel
	householdPanel *gui.Panel
	taxPanel       *gui.Panel
	storagePanel   *gui.Panel
	traderPanel    *gui.Panel
	buttons        []*TownhallControllerButton
	subPanel       *gui.Panel
	th             *social.Townhall
	activeTrader   *social.Trader
}

func TownhallToControlPanel(cp *ControlPanel, th *social.Townhall) {
	top := 15 + IconS + LargeIconD
	topPanel := &gui.Panel{X: 0, Y: ControlPanelDynamicPanelTop, SX: ControlPanelSX, SY: HouseholdControllerSY}
	hp := &gui.Panel{X: 0, Y: ControlPanelDynamicPanelTop, SX: ControlPanelSX, SY: HouseholdControllerSY}
	mp := &gui.Panel{X: 0, Y: ControlPanelDynamicPanelTop, SX: ControlPanelSX, SY: HouseholdControllerSY}
	sp := &gui.Panel{X: 0, Y: ControlPanelDynamicPanelTop, SX: ControlPanelSX, SY: HouseholdControllerSY}
	tp := &gui.Panel{X: 0, Y: ControlPanelDynamicPanelTop, SX: ControlPanelSX, SY: HouseholdControllerSY}

	tc := &TownhallController{cp: cp, th: th, topPanel: topPanel, householdPanel: hp, taxPanel: mp, storagePanel: sp, traderPanel: tp}
	tc.buttons = []*TownhallControllerButton{
		&TownhallControllerButton{tc: tc, subPanel: hp, b: gui.ButtonGUI{Icon: "house", X: float64(24 + LargeIconD*0), Y: top, SX: LargeIconS, SY: LargeIconS}},
		&TownhallControllerButton{tc: tc, subPanel: mp, b: gui.ButtonGUI{Icon: "taxes", X: float64(24 + LargeIconD*1), Y: top, SX: LargeIconS, SY: LargeIconS}},
		&TownhallControllerButton{tc: tc, subPanel: sp, b: gui.ButtonGUI{Icon: "barrel", X: float64(24 + LargeIconD*2), Y: top, SX: LargeIconS, SY: LargeIconS}},
		&TownhallControllerButton{tc: tc, subPanel: tp, b: gui.ButtonGUI{Icon: "trader", X: float64(24 + LargeIconD*3), Y: top, SX: LargeIconS, SY: LargeIconS}},
	}

	tc.subPanel = tc.householdPanel
	RefreshSubPanels(tc)

	cp.SetDynamicPanel(tc)
	cp.C.ClickHandler = tc
}

func RefreshSubPanels(tc *TownhallController) {
	th := tc.th
	town := th.Household.Town
	tp := tc.taxPanel
	sp := tc.storagePanel
	top := 15 + IconS + LargeIconD*2

	HouseholdToControlPanel(tc.cp, tc.householdPanel, th.Household)

	tpw := (ControlPanelSX - 30) / 2
	s := IconS / 2
	h := float64(LargeIconD / 3)
	tw := 24 + LargeIconD
	tp.AddImageLabel("farm", 24, top+h*2, LargeIconS, LargeIconS, gui.ImageLabelStyleRegular)
	tp.AddPanel(gui.CreateNumberPaneFromVal(tw, top+h*2, tpw-tw, s, 0, 100, 10, "tax rate %v", &th.Household.Town.Transfers.Farm.Rate).P)
	tp.AddPanel(gui.CreateNumberPaneFromVal(tw, top+h*3, tpw-tw, s, 0, 1000, 50, "threshold %v", &th.Household.Town.Transfers.Farm.Threshold).P)

	tp.AddImageLabel("workshop", 24, top+h*5, LargeIconS, LargeIconS, gui.ImageLabelStyleRegular)
	tp.AddPanel(gui.CreateNumberPaneFromVal(tw, top+h*5, tpw-tw, s, 0, 100, 10, "tax rate %v", &th.Household.Town.Transfers.Workshop.Rate).P)
	tp.AddPanel(gui.CreateNumberPaneFromVal(tw, top+h*6, tpw-tw, s, 0, 1000, 50, "threshold %v", &th.Household.Town.Transfers.Workshop.Threshold).P)

	tp.AddImageLabel("mine", 24+tpw, top+h*2, LargeIconS, LargeIconS, gui.ImageLabelStyleRegular)
	tp.AddPanel(gui.CreateNumberPaneFromVal(tw+tpw, top+h*2, tpw-tw, s, 0, 100, 10, "tax rate %v", &th.Household.Town.Transfers.Mine.Rate).P)
	tp.AddPanel(gui.CreateNumberPaneFromVal(tw+tpw, top+h*3, tpw-tw, s, 0, 1000, 50, "threshold %v", &th.Household.Town.Transfers.Mine.Threshold).P)

	tp.AddImageLabel("factory", 24+tpw, top+h*5, LargeIconS, LargeIconS, gui.ImageLabelStyleRegular)
	tp.AddPanel(gui.CreateNumberPaneFromVal(tw+tpw, top+h*5, tpw-tw, s, 0, 100, 10, "tax rate %v", &th.Household.Town.Transfers.Factory.Rate).P)
	tp.AddPanel(gui.CreateNumberPaneFromVal(tw+tpw, top+h*6, tpw-tw, s, 0, 1000, 50, "threshold %v", &th.Household.Town.Transfers.Factory.Threshold).P)

	tp.AddImageLabel("trader", 24, top+h*8, LargeIconS, LargeIconS, gui.ImageLabelStyleRegular)
	tp.AddPanel(gui.CreateNumberPaneFromVal(tw, top+h*8, tpw-tw, s, 0, 100, 10, "tax rate %v", &th.Household.Town.Transfers.Trader.Rate).P)
	tp.AddPanel(gui.CreateNumberPaneFromVal(tw, top+h*9, tpw-tw, s, 0, 1000, 50, "threshold %v", &th.Household.Town.Transfers.Trader.Threshold).P)

	tp.AddPanel(gui.CreateNumberPaneFromVal(tw+tpw, top+h*8, tpw-s, s, 0, 100, 50, "military %v", &th.Household.Town.Transfers.Tower.Threshold).P)
	tp.AddPanel(gui.CreateNumberPaneFromVal(tw+tpw, top+h*9, tpw-s, s, 0, 100, 10, "market %v", &th.Household.Town.Transfers.MarketFundingRate).P)

	tp.AddLargeTextLabel("Activities", 24, top+LargeIconD*4)
	tp.AddImageLabel("infra/cobble_road", 24, top+LargeIconD*5, LargeIconS, LargeIconS, gui.ImageLabelStyleRegular)
	tp.AddButton(gui.SimpleButton{
		ButtonGUI: gui.ButtonGUI{Icon: "tasks/repair", X: 24 + LargeIconD, Y: top + LargeIconD*5, SX: LargeIconS, SY: LargeIconS},
		Highlight: func() bool { return town.Settings.RoadRepairs },
		ClickImpl: func() {
			town.Settings.RoadRepairs = !town.Settings.RoadRepairs
			tc.cp.HelperMessage("Start or stop repairing roads")
		}})

	tp.AddImageLabel("infra/wall_small", 24, top+LargeIconD*6, LargeIconS, LargeIconS, gui.ImageLabelStyleRegular)
	tp.AddButton(gui.SimpleButton{
		ButtonGUI: gui.ButtonGUI{Icon: "tasks/repair", X: 24 + LargeIconD, Y: top + LargeIconD*6, SX: LargeIconS, SY: LargeIconS},
		Highlight: func() bool { return town.Settings.WallRepairs },
		ClickImpl: func() {
			town.Settings.WallRepairs = !town.Settings.WallRepairs
			tc.cp.HelperMessage("Start or stop repairing walls")
		}})

	tp.AddButton(gui.SimpleButton{
		ButtonGUI: gui.ButtonGUI{Icon: "trader", X: 24 + LargeIconD, Y: top + LargeIconD*7, SX: LargeIconS, SY: LargeIconS},
		Highlight: func() bool { return town.Settings.Trading },
		ClickImpl: func() {
			town.Settings.Trading = !town.Settings.Trading
			tc.cp.HelperMessage("Enable or disable trading with this city")
		}})

	tp.AddButton(gui.SimpleButton{
		ButtonGUI: gui.ButtonGUI{Icon: "barrel", X: 24 + LargeIconD, Y: top + LargeIconD*8, SX: LargeIconS, SY: LargeIconS},
		Highlight: func() bool { return town.Settings.ArtifactCollection },
		ClickImpl: func() {
			town.Settings.ArtifactCollection = !town.Settings.ArtifactCollection
			tc.cp.HelperMessage("Start or stop collecting nearby abandoned items")
		}})

	var aI = 0
	for _, a := range artifacts.All {
		var q uint16 = 0
		if storageQ, ok := th.Household.Resources.Artifacts[a]; ok {
			q = storageQ
		}
		ArtifactStorageToControlPanel(sp, th, aI, a, q, top+50)
		aI++
	}

	for i, vc := range social.GetVehicleConstructions(th.Household.Town.Factories) {
		tc.traderPanel.AddPanel(CreateOrderPanelForTownhall(24, float64(i+2)*IconS+top, gui.FontSize*8, s, th, vc, tc.cp.C.Map))
		if vc.Output.Trader {
			tc.traderPanel.AddButton(CreateTraderButtonForTownhall(24+tpw, float64(i+2)*IconS+top, float64(IconH), s, th, tc.cp.C.Map))
		}
	}
	for i, vehicle := range th.Household.Vehicles {
		VehicleToControlPanel(tc.traderPanel, i, 6*IconS+top, vehicle, IconW)
	}

	traderTop := top + ControlPanelSY*0.25
	for i, t := range th.Traders {
		tc.traderPanel.AddButton(CreateTraderButton(float64(24+i*IconW), traderTop, tc, t))
	}
	if tc.activeTrader != nil {
		MoneyToControlPanel(tc.traderPanel, th.Household.Town, &tc.activeTrader.Money, 24, 10, traderTop+float64(IconH)+IconS)
		for i, task := range tc.activeTrader.Tasks {
			TaskToControlPanel(tc.cp, tc.traderPanel, i, traderTop+float64(IconH*3)+IconS, task, IconW)
		}
	}
}

func ArtifactStorageToControlPanel(p *gui.Panel, th *social.Townhall, i int, a *artifacts.Artifact, q uint16, top float64) {
	rowH := int(IconS * 2)
	xI := i % IconRowMaxButtons
	yI := i / IconRowMaxButtons
	w := int(float64(IconW) * float64(IconRowMax) / float64(IconRowMaxButtons))
	p.AddImageLabel("artifacts/"+a.Name, float64(24+xI*w), top+float64(yI*rowH), IconS, IconS, gui.ImageLabelStyleRegular)
	p.AddTextLabel(strconv.Itoa(int(q)), float64(24+xI*w), top+float64(yI*rowH+IconH+4))
	p.AddPanel(gui.CreateNumberPanel(float64(24+xI*w), top+float64(yI*rowH+IconH+4), float64(IconW+8), gui.FontSize*1.5, 0, 250, 5, "%v",
		func() int { return th.StorageTarget[a] },
		func(v int) { th.StorageTarget[a] = v }).P)
}

func (tc *TownhallController) CaptureClick(x, y float64) {
	tc.topPanel.CaptureClick(x, y)
}

func (tc *TownhallController) Render(cv *canvas.Canvas) {
	tc.topPanel.Render(cv)
}

func (tc *TownhallController) Clear() {}

func (tc *TownhallController) Refresh() {
	tc.topPanel.Clear()
	tc.householdPanel.Clear()
	tc.taxPanel.Clear()
	tc.storagePanel.Clear()
	tc.traderPanel.Clear()
	RefreshSubPanels(tc)
	for _, button := range tc.buttons {
		tc.topPanel.AddButton(button)
	}
	if tc.subPanel != nil {
		tc.topPanel.AddPanel(tc.subPanel)
	}
}

func (tc *TownhallController) GetActiveFields(c *Controller, rf *renderer.RenderedField) []navigation.FieldWithContext {
	if tc.activeTrader != nil && tc.activeTrader.TargetExchange != nil {
		var fs []navigation.FieldWithContext
		for _, coords := range tc.activeTrader.TargetExchange.Building.GetBuildingXYs(true) {
			fs = append(fs, c.Map.GetField(coords[0], coords[1]))
		}
		return fs
	}
	return tc.th.GetFields()
}

func (tc *TownhallController) HandleClick(c *Controller, rf *renderer.RenderedField) bool {
	if tc.activeTrader != nil {
		return HandleClickForTrader(tc.activeTrader, c, rf)
	}
	for i := range tc.th.Household.Town.Roads {
		r := tc.th.Household.Town.Roads[i]
		if r.X == rf.F.X && r.Y == rf.F.Y {
			r.Allocated = false
			tc.th.Household.Town.Roads = append(tc.th.Household.Town.Roads[:i], tc.th.Household.Town.Roads[i+1:]...)
			return true
		}
	}
	if !rf.F.Allocated && rf.F.Road != nil && tc.th.FieldWithinDistance(rf.F) {
		tc.th.Household.Town.Roads = append(tc.th.Household.Town.Roads, rf.F)
		rf.F.Allocated = true
		return true
	}
	return false
}

func CreateOrderPanelForTownhall(x, y, sx, sy float64, th *social.Townhall, vc *economy.VehicleConstruction, m navigation.IMap) *gui.Panel {
	p := &gui.Panel{}
	l := p.AddTextLabel("", x, y+sy*2/3)
	var factories []*social.Factory
	for _, factory := range th.Household.Town.Factories {
		if economy.ConstructionCompatible(vc, factory.Household.Building.Plan.GetExtensions()) {
			factories = append(factories, factory)
		}
	}
	p.AddButton(OrderButton{
		b:         gui.ButtonGUI{Icon: "plus", X: x + sx, Y: y, SX: sy, SY: sy},
		factories: factories,
		vc:        vc,
		l:         l,
		m:         m,
	})
	p.AddTextLabel(fmt.Sprintf("$%v", factories[0].Price(vc)), x+sx+sy*2, y+sy*2/3)
	return p
}

func CreateTraderButtonForTownhall(x, y, sx, sy float64, th *social.Townhall, m navigation.IMap) gui.Button {
	return &gui.SimpleButton{
		ButtonGUI: gui.ButtonGUI{Icon: "plus", X: x + sx, Y: y, SX: sy, SY: sy},
		ClickImpl: func() {
			th.CreateTrader(m)
		},
	}
}

func CreateTraderButton(x, y float64, th *TownhallController, t *social.Trader) gui.Button {
	return &gui.SimpleButton{
		ButtonGUI: gui.ButtonGUI{Icon: "trader", X: x, Y: y, SX: IconS, SY: IconS},
		ClickImpl: func() {
			th.activeTrader = t
		},
		Highlight: func() bool {
			return t == th.activeTrader
		},
	}
}
