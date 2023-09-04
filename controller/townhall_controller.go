package controller

import (
	"fmt"
	"github.com/tfriedel6/canvas"
	"medvil/model/artifacts"
	"medvil/model/building"
	"medvil/model/economy"
	"medvil/model/navigation"
	"medvil/model/social"
	"medvil/model/vehicles"
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
	if b.tc.subPanel == b.subPanel {
		cv.SetFillStyle(gui.ButtonColorHighlight)
		cv.FillRect(b.b.X, b.b.Y, b.b.SX, b.b.SY)
	}
	b.b.Render(cv)
}

func (b *TownhallControllerButton) Contains(x float64, y float64) bool {
	return b.b.Contains(x, y)
}

func (b *TownhallControllerButton) Enabled() bool {
	return b.b.Enabled()
}

type TownhallController struct {
	cp               *ControlPanel
	topPanel         *gui.Panel
	householdPanel   *gui.Panel
	taxPanel         *gui.Panel
	storagePanel     *gui.Panel
	traderPanel      *gui.Panel
	expeditionPanel  *gui.Panel
	buttons          []*TownhallControllerButton
	subPanel         *gui.Panel
	th               *social.Townhall
	activeTrader     *social.Trader
	activeExpedition *social.Expedition
}

func TownhallToControlPanel(cp *ControlPanel, th *social.Townhall) {
	top := 15 + IconS + LargeIconD
	topPanel := &gui.Panel{X: 0, Y: ControlPanelDynamicPanelTop, SX: ControlPanelSX, SY: HouseholdControllerSY}
	hp := &gui.Panel{X: 0, Y: ControlPanelDynamicPanelTop, SX: ControlPanelSX, SY: HouseholdControllerSY}
	mp := &gui.Panel{X: 0, Y: ControlPanelDynamicPanelTop, SX: ControlPanelSX, SY: HouseholdControllerSY}
	sp := &gui.Panel{X: 0, Y: ControlPanelDynamicPanelTop, SX: ControlPanelSX, SY: HouseholdControllerSY}
	tp := &gui.Panel{X: 0, Y: ControlPanelDynamicPanelTop, SX: ControlPanelSX, SY: HouseholdControllerSY}
	ep := &gui.Panel{X: 0, Y: ControlPanelDynamicPanelTop, SX: ControlPanelSX, SY: HouseholdControllerSY}

	tc := &TownhallController{cp: cp, th: th, topPanel: topPanel, householdPanel: hp, taxPanel: mp, storagePanel: sp, traderPanel: tp, expeditionPanel: ep}
	tc.buttons = []*TownhallControllerButton{
		&TownhallControllerButton{tc: tc, subPanel: hp, b: gui.ButtonGUI{Icon: "town", X: float64(24 + LargeIconD*0), Y: top, SX: LargeIconS, SY: LargeIconS}},
		&TownhallControllerButton{tc: tc, subPanel: mp, b: gui.ButtonGUI{Icon: "taxes", X: float64(24 + LargeIconD*1), Y: top, SX: LargeIconS, SY: LargeIconS}},
		&TownhallControllerButton{tc: tc, subPanel: sp, b: gui.ButtonGUI{Icon: "barrel", X: float64(24 + LargeIconD*2), Y: top, SX: LargeIconS, SY: LargeIconS}},
		&TownhallControllerButton{tc: tc, subPanel: tp, b: gui.ButtonGUI{Icon: "trader", X: float64(24 + LargeIconD*3), Y: top, SX: LargeIconS, SY: LargeIconS}},
		&TownhallControllerButton{tc: tc, subPanel: ep, b: gui.ButtonGUI{Icon: "expedition", X: float64(24 + LargeIconD*4), Y: top, SX: LargeIconS, SY: LargeIconS}},
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
	tp.AddLargeTextLabel("Finances", 24, 15+IconS+LargeIconD*3)
	top := 15 + IconS + LargeIconD*3

	HouseholdToControlPanel(tc.cp, tc.householdPanel, th.Household)

	tpw := (ControlPanelSX - 30) / 2
	s := IconS / 2
	h := float64(LargeIconD / 3)
	tw := 24 + LargeIconD

	tp.AddImageLabel("farm", 24, top+h*2, LargeIconS, LargeIconS, gui.ImageLabelStyleRegular)
	tp.AddPanel(gui.CreateNumberPaneFromVal(tw, top+h*2, tpw-tw, s, 0, 100, 10, "tax rate %v", &th.Household.Town.Transfers.Farm.Rate).P)
	tp.AddPanel(gui.CreateNumberPaneFromVal(tw, top+h*3, tpw-tw, s, 0, 1000, 50, "subsidy %v", &th.Household.Town.Transfers.Farm.Threshold).P)

	tp.AddImageLabel("workshop", 24, top+h*5, LargeIconS, LargeIconS, gui.ImageLabelStyleRegular)
	tp.AddPanel(gui.CreateNumberPaneFromVal(tw, top+h*5, tpw-tw, s, 0, 100, 10, "tax rate %v", &th.Household.Town.Transfers.Workshop.Rate).P)
	tp.AddPanel(gui.CreateNumberPaneFromVal(tw, top+h*6, tpw-tw, s, 0, 1000, 50, "subsidy %v", &th.Household.Town.Transfers.Workshop.Threshold).P)

	tp.AddImageLabel("mine", 24+tpw, top+h*2, LargeIconS, LargeIconS, gui.ImageLabelStyleRegular)
	tp.AddPanel(gui.CreateNumberPaneFromVal(tw+tpw, top+h*2, tpw-tw, s, 0, 100, 10, "tax rate %v", &th.Household.Town.Transfers.Mine.Rate).P)
	tp.AddPanel(gui.CreateNumberPaneFromVal(tw+tpw, top+h*3, tpw-tw, s, 0, 1000, 50, "subsidy %v", &th.Household.Town.Transfers.Mine.Threshold).P)

	tp.AddImageLabel("factory", 24+tpw, top+h*5, LargeIconS, LargeIconS, gui.ImageLabelStyleRegular)
	tp.AddPanel(gui.CreateNumberPaneFromVal(tw+tpw, top+h*5, tpw-tw, s, 0, 100, 10, "tax rate %v", &th.Household.Town.Transfers.Factory.Rate).P)
	tp.AddPanel(gui.CreateNumberPaneFromVal(tw+tpw, top+h*6, tpw-tw, s, 0, 1000, 50, "subsidy %v", &th.Household.Town.Transfers.Factory.Threshold).P)

	tp.AddImageLabel("trader", 24, top+h*8, LargeIconS, LargeIconS, gui.ImageLabelStyleRegular)
	tp.AddPanel(gui.CreateNumberPaneFromVal(tw, top+h*8, tpw-tw, s, 0, 100, 10, "tax rate %v", &th.Household.Town.Transfers.Trader.Rate).P)
	tp.AddPanel(gui.CreateNumberPaneFromVal(tw, top+h*9, tpw-tw, s, 0, 1000, 50, "subsidy %v", &th.Household.Town.Transfers.Trader.Threshold).P)

	tp.AddPanel(gui.CreateNumberPaneFromVal(tw+tpw, top+h*8, tpw-s, s, 0, 100, 50, "military %v", &th.Household.Town.Transfers.Tower.Threshold).P)
	tp.AddPanel(gui.CreateNumberPaneFromVal(tw+tpw, top+h*9, tpw-s, s, 0, 100, 10, "market %v", &th.Household.Town.Transfers.MarketFundingRate).P)

	tp.AddLargeTextLabel("Activities", 24, top+LargeIconD*4+s)
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

	tp.AddButton(gui.SimpleButton{
		ButtonGUI: gui.ButtonGUI{Icon: "coin", X: 24 + LargeIconD, Y: top + LargeIconD*9, SX: LargeIconS, SY: LargeIconS},
		Highlight: func() bool { return town.Settings.Coinage },
		ClickImpl: func() {
			town.Settings.Coinage = !town.Settings.Coinage
			tc.cp.HelperMessage("Start or stop minting gold coins")
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

	for i, vc := range social.GetVehicleConstructions(th.Household.Town.Factories, func(vc *economy.VehicleConstruction) bool { return vc.Output.Trader }) {
		tc.traderPanel.AddPanel(CreateTraderButtonForTownhall(24, float64(i)*LargeIconD+top, th, vc, tc.cp.C.Map))
	}
	for i := 0; i < th.Household.NumTasks("create_trader", ""); i++ {
		tc.traderPanel.AddImageLabel("tasks/factory_pickup", float64(24+i*IconW), top+ControlPanelSY*0.15, IconS, IconS, gui.ImageLabelStyleRegular)
	}

	traderTop := top + ControlPanelSY*0.20
	for i, t := range th.Traders {
		tc.traderPanel.AddButton(CreateTraderButton(float64(24+i*IconW), traderTop, tc, t))
	}
	if tc.activeTrader != nil {
		MoneyToControlPanel(tc.traderPanel, th.Household.Town, &tc.activeTrader.Money, 24, 10, traderTop+float64(IconH)+IconS)
		for i, task := range tc.activeTrader.Tasks {
			TaskToControlPanel(tc.cp, tc.traderPanel, i, traderTop+float64(IconH*3)+IconS, task, IconW)
		}
	}

	for i, vc := range social.GetVehicleConstructions(th.Household.Town.Factories, func(vc *economy.VehicleConstruction) bool { return vc.Output.Expedition }) {
		tc.expeditionPanel.AddPanel(CreateExpeditionButtonForTownhall(24, float64(i)*LargeIconD+top, th, vc, tc.cp.C.Map))
	}
	for i, v := range th.Household.GetVehicles(func(v *vehicles.Vehicle) bool { return v.T.Expedition }) {
		tc.expeditionPanel.AddButton(CreateExpeditionButton(float64(24+i*IconW), top+ControlPanelSY*0.15, tc, v))
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
	tc.expeditionPanel.Clear()
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

func CreateTraderButtonForTownhall(x, y float64, th *social.Townhall, vc *economy.VehicleConstruction, m navigation.IMap) *gui.Panel {
	p := &gui.Panel{}
	p.AddImageLabel("vehicles/"+vc.Name, 24, y, LargeIconS, LargeIconS, gui.ImageLabelStyleRegular)
	p.AddButton(&gui.SimpleButton{
		ButtonGUI: gui.ButtonGUI{Icon: "plus", X: x + 240, Y: y + float64(IconH/4), SX: IconS, SY: IconS},
		ClickImpl: func() {
			h := th.Household
			factory := social.PickFactory(h.Town.Factories, vc.BuildingExtensionType, th.Household, m)
			order := factory.CreateOrder(vc, h)
			if order != nil {
				h.AddTask(&economy.CreateTraderTask{
					Townhall: th,
					PickupD:  factory.Household.Destination(building.NonExtension),
					Order:    order,
				})
			}
		},
	})
	if th.Household.Town.Marketplace != nil {
		p.AddTextLabel(fmt.Sprintf("$%v", social.VehiclePrice(th.Household.Town.Marketplace, vc)), 24+x+float64(IconW)*2, y+float64(LargeIconD)/2)
	}
	return p
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

func CreateExpeditionButtonForTownhall(x, y float64, th *social.Townhall, vc *economy.VehicleConstruction, m navigation.IMap) *gui.Panel {
	p := &gui.Panel{}
	p.AddImageLabel("vehicles/"+vc.Name, 24, y, LargeIconS, LargeIconS, gui.ImageLabelStyleRegular)
	p.AddButton(&gui.SimpleButton{
		ButtonGUI: gui.ButtonGUI{Icon: "plus", X: x + 240, Y: y + float64(IconH/4), SX: IconS, SY: IconS},
		ClickImpl: func() {
			h := th.Household
			factory := social.PickFactory(h.Town.Factories, vc.BuildingExtensionType, th.Household, m)
			order := factory.CreateOrder(vc, h)
			if order != nil {
				h.AddTask(&economy.FactoryPickupTask{
					PickupD:  factory.Household.Destination(building.NonExtension),
					DropoffD: h.Destination(vc.BuildingExtensionType),
					Order:    order,
					TaskBase: economy.TaskBase{FieldCenter: true},
				})
			}
		},
	})
	if th.Household.Town.Marketplace != nil {
		p.AddTextLabel(fmt.Sprintf("$%v", social.VehiclePrice(th.Household.Town.Marketplace, vc)), 24+x+float64(IconW)*2, y+float64(LargeIconD)/2)
	}
	return p
}

func CreateExpeditionButton(x, y float64, th *TownhallController, v *vehicles.Vehicle) gui.Button {
	expedition := &social.Expedition{Vehicle: v}
	return &gui.SimpleButton{
		ButtonGUI: gui.ButtonGUI{Icon: "vehicles/" + v.T.Name, X: x, Y: y, SX: IconS, SY: IconS},
		ClickImpl: func() {
			th.activeExpedition = expedition
		},
		Highlight: func() bool {
			return expedition == th.activeExpedition
		},
	}
}
