package controller

import (
	"fmt"
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

var TownhallControllerGUIBottomY = 0.75
var TownhallDocGUIY = 0.85

type TownhallControllerButton struct {
	tc        *TownhallController
	b         gui.ButtonGUI
	subPanel  *gui.Panel
	helperMsg string
}

func (b *TownhallControllerButton) Click() {
	b.tc.subPanel = b.subPanel
	if !b.tc.cp.C.ViewSettings.ShowSuggestions {
		b.tc.cp.SelectedHelperMessage(b.helperMsg)
	}
}

func (b *TownhallControllerButton) Render(cv *canvas.Canvas) {
	if b.tc.subPanel == b.subPanel {
		cv.SetFillStyle(gui.ButtonColorHighlight)
		cv.FillRect(b.b.X, b.b.Y, b.b.SX, b.b.SY)
	}
	b.b.Render(cv)
}

func (b *TownhallControllerButton) SetHoover(h bool) {
	b.b.SetHoover(h)
	if h {
		b.tc.cp.HelperMessage(b.helperMsg, true)
	}
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
		&TownhallControllerButton{
			tc: tc, subPanel: hp, b: gui.ButtonGUI{Icon: "town", X: float64(24 + LargeIconD*0), Y: top, SX: LargeIconS, SY: LargeIconS},
			helperMsg: "View status",
		},
		&TownhallControllerButton{
			tc: tc, subPanel: mp, b: gui.ButtonGUI{Icon: "taxes", X: float64(24 + LargeIconD*1), Y: top, SX: LargeIconS, SY: LargeIconS},
			helperMsg: "Adjust taxes and switch on/off activities.",
		},
		&TownhallControllerButton{
			tc: tc, subPanel: sp, b: gui.ButtonGUI{Icon: "barrel", X: float64(24 + LargeIconD*2), Y: top, SX: LargeIconS, SY: LargeIconS},
			helperMsg: "Store or offload goods.",
		},
		&TownhallControllerButton{
			tc: tc, subPanel: tp, b: gui.ButtonGUI{Icon: "trader", X: float64(24 + LargeIconD*3), Y: top, SX: LargeIconS, SY: LargeIconS,
				Disabled: func() bool { return !cp.C.ActiveSupplier.HasHousehold(building.BuildingTypeWorkshop) }},
			helperMsg: "Create traders to trade goods with other towns.",
		},
		&TownhallControllerButton{
			tc: tc, subPanel: ep, b: gui.ButtonGUI{Icon: "expedition", X: float64(24 + LargeIconD*4), Y: top, SX: LargeIconS, SY: LargeIconS,
				Disabled: func() bool { return !cp.C.ActiveSupplier.HasHousehold(building.BuildingTypeWorkshop) }},
			helperMsg: "Start expeditions to found new towns.",
		},
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

	HouseholdToControlPanel(tc.cp, tc.householdPanel, th.Household, "townhall")

	tpw := (ControlPanelSX - 30) / 2
	s := IconS / 2
	h := float64(LargeIconD / 3)
	tw := 24 + LargeIconD

	tp.AddImageLabel("farm", 24, top+h*2, LargeIconS, LargeIconS, gui.ImageLabelStyleRegular)
	tp.AddPanel(gui.CreateNumberPaneFromVal(tw, top+h*2, tpw-tw, s, 0, 100, 10, "tax rate %v", true, &th.Household.Town.Transfers.Farm.Rate).P)
	tp.AddPanel(gui.CreateNumberPaneFromVal(tw, top+h*3, tpw-tw, s, 0, 1000, 50, "subsidy %v", true, &th.Household.Town.Transfers.Farm.Threshold).P)

	tp.AddImageLabel("workshop", 24, top+h*5, LargeIconS, LargeIconS, gui.ImageLabelStyleRegular)
	tp.AddPanel(gui.CreateNumberPaneFromVal(tw, top+h*5, tpw-tw, s, 0, 100, 10, "tax rate %v", true, &th.Household.Town.Transfers.Workshop.Rate).P)
	tp.AddPanel(gui.CreateNumberPaneFromVal(tw, top+h*6, tpw-tw, s, 0, 1000, 50, "subsidy %v", true, &th.Household.Town.Transfers.Workshop.Threshold).P)

	tp.AddImageLabel("mine", 24+tpw, top+h*2, LargeIconS, LargeIconS, gui.ImageLabelStyleRegular)
	tp.AddPanel(gui.CreateNumberPaneFromVal(tw+tpw, top+h*2, tpw-tw, s, 0, 100, 10, "tax rate %v", true, &th.Household.Town.Transfers.Mine.Rate).P)
	tp.AddPanel(gui.CreateNumberPaneFromVal(tw+tpw, top+h*3, tpw-tw, s, 0, 1000, 50, "subsidy %v", true, &th.Household.Town.Transfers.Mine.Threshold).P)

	tp.AddImageLabel("factory", 24+tpw, top+h*5, LargeIconS, LargeIconS, gui.ImageLabelStyleRegular)
	tp.AddPanel(gui.CreateNumberPaneFromVal(tw+tpw, top+h*5, tpw-tw, s, 0, 100, 10, "tax rate %v", true, &th.Household.Town.Transfers.Factory.Rate).P)
	tp.AddPanel(gui.CreateNumberPaneFromVal(tw+tpw, top+h*6, tpw-tw, s, 0, 1000, 50, "subsidy %v", true, &th.Household.Town.Transfers.Factory.Threshold).P)

	tp.AddImageLabel("trader", 24, top+h*8, LargeIconS, LargeIconS, gui.ImageLabelStyleRegular)
	tp.AddPanel(gui.CreateNumberPaneFromVal(tw, top+h*8, tpw-tw, s, 0, 100, 10, "tax rate %v", true, &th.Household.Town.Transfers.Trader.Rate).P)
	tp.AddPanel(gui.CreateNumberPaneFromVal(tw, top+h*9, tpw-tw, s, 0, 1000, 50, "subsidy %v", true, &th.Household.Town.Transfers.Trader.Threshold).P)

	tp.AddPanel(gui.CreateNumberPaneFromVal(tw+tpw, top+h*8, tpw-tw, s, 0, 1000, 50, "military %v", true, &th.Household.Town.Transfers.Tower.Threshold).P)
	tp.AddPanel(gui.CreateNumberPaneFromVal(tw+tpw, top+h*9, tpw-tw, s, 0, 100, 10, "market %v", true, &th.Household.Town.Transfers.MarketFundingRate).P)

	tp.AddLargeTextLabel("Activities", 24, top+LargeIconD*4+s)
	tp.AddImageLabel("infra/cobble_road", 24, top+LargeIconD*5, LargeIconS, LargeIconS, gui.ImageLabelStyleRegular)
	tp.AddButton(&gui.SimpleButton{
		ButtonGUI: gui.ButtonGUI{Icon: "tasks/repair", X: 24 + LargeIconD, Y: top + LargeIconD*5, SX: LargeIconS, SY: LargeIconS, OnHoover: func() {
			tc.cp.HelperMessage("Start or stop repairing roads", true)
		}},
		Highlight: func() bool { return town.Settings.RoadRepairs },
		ClickImpl: func() {
			town.Settings.RoadRepairs = !town.Settings.RoadRepairs
		}})
	tp.AddTextLabel("Repair "+strconv.Itoa(len(town.Roads))+" roads", 24+LargeIconD*2, top+LargeIconD*5+LargeIconD/2)

	tp.AddImageLabel("infra/wall_small", 24, top+LargeIconD*6, LargeIconS, LargeIconS, gui.ImageLabelStyleRegular)
	tp.AddButton(&gui.SimpleButton{
		ButtonGUI: gui.ButtonGUI{Icon: "tasks/repair", X: 24 + LargeIconD, Y: top + LargeIconD*6, SX: LargeIconS, SY: LargeIconS, OnHoover: func() {
			tc.cp.HelperMessage("Start or stop repairing walls", true)
		}},
		Highlight: func() bool { return town.Settings.WallRepairs },
		ClickImpl: func() {
			town.Settings.WallRepairs = !town.Settings.WallRepairs
		}})
	tp.AddTextLabel("Repair "+strconv.Itoa(len(town.Walls))+" walls", 24+LargeIconD*2, top+LargeIconD*6+LargeIconD/2)

	tp.AddImageLabel("market", 24, top+LargeIconD*7, LargeIconS, LargeIconS, gui.ImageLabelStyleRegular)
	tp.AddButton(&gui.SimpleButton{
		ButtonGUI: gui.ButtonGUI{Icon: "trader", X: 24 + LargeIconD, Y: top + LargeIconD*7, SX: LargeIconS, SY: LargeIconS, OnHoover: func() {
			tc.cp.HelperMessage("Enable or disable trading with this city", true)
		}},
		Highlight: func() bool { return town.Settings.Trading },
		ClickImpl: func() {
			town.Settings.Trading = !town.Settings.Trading
		}})
	tp.AddTextLabel("This city has "+strconv.Itoa(len(town.Townhall.Traders))+" traders", 24+LargeIconD*2, top+LargeIconD*7+LargeIconD/2)

	tp.AddButton(&gui.SimpleButton{
		ButtonGUI: gui.ButtonGUI{Icon: "barrel", X: 24 + LargeIconD, Y: top + LargeIconD*8, SX: LargeIconS, SY: LargeIconS, OnHoover: func() {
			tc.cp.HelperMessage("Start or stop collecting nearby abandoned items", true)
		}},
		Highlight: func() bool { return town.Settings.ArtifactCollection },
		ClickImpl: func() {
			town.Settings.ArtifactCollection = !town.Settings.ArtifactCollection
		}})
	tp.AddTextLabel("Storage is "+strconv.Itoa(int(th.Household.Resources.UsedVolumeCapacity()*100.0))+"% full", 24+LargeIconD*2, top+LargeIconD*8+LargeIconD/2)

	tp.AddButton(&gui.SimpleButton{
		ButtonGUI: gui.ButtonGUI{Icon: "coin", X: 24 + LargeIconD, Y: top + LargeIconD*9, SX: LargeIconS, SY: LargeIconS, OnHoover: func() {
			tc.cp.HelperMessage("Start or stop minting gold coins", true)
		}},
		Highlight: func() bool { return town.Settings.Coinage },
		ClickImpl: func() {
			town.Settings.Coinage = !town.Settings.Coinage
		}})
	tp.AddTextLabel(""+strconv.Itoa(int(town.Stats.Global.Money))+" coins in circulation", 24+LargeIconD*2, top+LargeIconD*9+LargeIconD/2)

	tp.AddButton(&gui.SimpleButton{
		ButtonGUI: gui.ButtonGUI{Icon: "town", X: 24 + LargeIconD, Y: top + LargeIconD*10, SX: LargeIconS, SY: LargeIconS, OnHoover: func() {
			tc.cp.HelperMessage("Start or stop relying on supplies", true)
		}},
		Highlight: func() bool { return town.Settings.UseSupplier },
		ClickImpl: func() {
			town.Settings.UseSupplier = !town.Settings.UseSupplier
		}})
	if town.Supplier != nil {
		tp.AddTextLabel(town.Supplier.GetName()+" supplies "+town.GetName(), 24+LargeIconD*2, top+LargeIconD*10+LargeIconD/2)
	} else {
		tp.AddTextLabel("No supplier", 24+LargeIconD*2, top+LargeIconD*10+LargeIconD/2)
	}

	if tc.cp.C.ViewSettings.ShowSuggestions {
		AddLines(tp, []string{
			"Set taxes and subsidies for different social",
			"classes. Taxes are used to fund the townhall",
			"market, as well as paying for subsidies.",
			"",
			"If taxes are high, citizens become less happy.",
		}, 24, TownhallDocGUIY*ControlPanelSY)
	}

	var aI = 0
	for _, a := range artifacts.All {
		var q uint16 = 0
		if storageQ, ok := th.Household.Resources.Artifacts[a]; ok {
			q = storageQ
		}
		ArtifactStorageToControlPanel(sp, tc.cp, th.StorageTarget, aI, a, q, ControlPanelSY*0.175, false)
		aI++
	}

	traderTop := top + ControlPanelSY*0.20
	{
		for i, vc := range social.GetVehicleConstructions(th.Household.Town.Factories, func(vc *economy.VehicleConstruction) bool { return vc.Output.Trader }) {
			tc.traderPanel.AddPanel(CreateTraderButtonForTownhall(24, float64(i)*LargeIconD+top, tc, vc, tc.cp.C.Map))
		}
		for i := 0; i < th.Household.NumTasks("create_trader", economy.EmptyTag); i++ {
			tc.traderPanel.AddImageLabel("tasks/factory_pickup", float64(24+i*IconW), top+ControlPanelSY*0.15, IconS, IconS, gui.ImageLabelStyleRegular)
		}
		for i, t := range th.Traders {
			tc.traderPanel.AddButton(SelectTraderButton(float64(24+i*IconW), traderTop, tc, t))
		}
		if tc.activeTrader != nil {
			MoneyToControlPanel(tc.cp, tc.traderPanel, th.Household, tc.activeTrader, 24, 10, traderTop+float64(IconH)+IconS)
			for i, task := range tc.activeTrader.Tasks {
				TaskToControlPanel(tc.cp, tc.traderPanel, i, traderTop+float64(IconH*3)+IconS, task, IconW)
			}
			if tc.activeTrader.SourceExchange != nil {
				tc.traderPanel.AddImageLabel("market", 24, traderTop+float64(IconH*5), IconS, IconS, gui.ImageLabelStyleRegular)
				tc.traderPanel.AddTextLabel(tc.activeTrader.SourceExchange.Town.Name, 24+float64(IconW), traderTop+float64(IconH)*5.5)
			}
			if tc.activeTrader.TargetExchange != nil {
				tc.traderPanel.AddImageLabel("market", 24, traderTop+float64(IconH*6), IconS, IconS, gui.ImageLabelStyleRegular)
				tc.traderPanel.AddTextLabel(tc.activeTrader.TargetExchange.Town.Name, 24+float64(IconW), traderTop+float64(IconH)*6.5)
			}
		}
	}
	if tc.cp.C.ViewSettings.ShowSuggestions {
		AddLines(tc.traderPanel, []string{
			"Create traders to move goods between cities.",
			"Each trader has to be assigned a destination",
			"market to trade with. Traders select the most",
			"profitable goods and can reduce shortages.",
		}, 24, TownhallDocGUIY*ControlPanelSY)
	}

	{
		for i, vc := range social.GetVehicleConstructions(th.Household.Town.Factories, func(vc *economy.VehicleConstruction) bool { return vc.Output.Expedition }) {
			tc.expeditionPanel.AddPanel(CreateExpeditionButtonForTownhall(24, float64(i)*LargeIconD+top, tc, vc, tc.cp.C.Map))
		}
		for i := 0; i < th.Household.NumTasks("create_expedition", economy.EmptyTag); i++ {
			tc.expeditionPanel.AddImageLabel("tasks/factory_pickup", float64(24+i*IconW), top+ControlPanelSY*0.15, IconS, IconS, gui.ImageLabelStyleRegular)
		}
		for i, e := range th.Expeditions {
			tc.expeditionPanel.AddButton(SelectExpeditionButton(float64(24+i*IconW), traderTop, tc, e))
		}
		if tc.activeExpedition != nil {
			MoneyToControlPanel(tc.cp, tc.expeditionPanel, th.Household, tc.activeExpedition, 24, 10, traderTop+float64(IconH)+IconS)
			for i, task := range tc.activeExpedition.Tasks {
				TaskToControlPanel(tc.cp, tc.expeditionPanel, i, traderTop+float64(IconH*3)+IconS, task, IconW)
			}
		}
	}
	if tc.cp.C.ViewSettings.ShowSuggestions {
		AddLines(tc.expeditionPanel, []string{
			"Create expedition to fund faraway cities.",
			"Expeditions can store large amounts of good",
			"and travel to far distances that way.",
			"",
			"It is recommended to load them with food and",
			"building materials that are needed to establish",
			"brand new cities.",
		}, 24, TownhallDocGUIY*ControlPanelSY)
	}
}

func abs(val int) int {
	if val < 0 {
		return -val
	}
	return val
}

func ArtifactStorageToControlPanel(p *gui.Panel, cp *ControlPanel, st map[*artifacts.Artifact]int, i int, a *artifacts.Artifact, q uint16, top float64, offload bool) {
	rowH := int(IconS * 2)
	xI := i % IconRowMaxButtons
	yI := i / IconRowMaxButtons
	w := int(float64(IconW) * float64(IconRowMax) / float64(IconRowMaxButtons))
	p.AddImageLabel("artifacts/"+a.Name, float64(24+xI*w), top+float64(yI*rowH), IconS, IconS, gui.ImageLabelStyleRegular)
	p.AddTextLabel(strconv.Itoa(int(q)), float64(24+xI*w), top+float64(yI*rowH+IconH+4))
	var offloadButton *gui.SimpleButton
	if offload {
		var icon = "arrow_small_up"
		if st[a] < 0 {
			icon = "arrow_small_down"
		}
		offloadButton = &gui.SimpleButton{
			ButtonGUI: gui.ButtonGUI{Icon: icon, X: float64(24+xI*w) + IconS, Y: top + float64(yI*rowH), SX: IconS / 2.0, SY: IconS,
				Disabled: func() bool { return st[a] == 0 },
			},
			ClickImpl: func() {
				if st[a] < 0 {
					st[a] = -st[a]
					offloadButton.Icon = "arrow_small_up"
					cp.HelperMessage("Transport goods to the destination town", true)
				} else if st[a] > 0 {
					st[a] = -st[a]
					offloadButton.Icon = "arrow_small_down"
					cp.HelperMessage("Transport goods back", true)
				}
			},
		}
		p.AddButton(offloadButton)
	}
	p.AddPanel(gui.CreateNumberPanel(float64(24+xI*w), top+float64(yI*rowH+IconH+4), float64(IconW+8), gui.FontSize*1.5, 0, 250, 6, "%v", false,
		func() int { return abs(st[a]) },
		func(v int) {
			if st[a] < 0 {
				st[a] = -v
			} else {
				st[a] = v
			}
		}).P)
}

func (tc *TownhallController) CaptureMove(x, y float64) {
	tc.topPanel.CaptureMove(x, y)
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
	tc.CaptureMove(tc.cp.C.X, tc.cp.C.Y)
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
	if tc.activeExpedition != nil {
		return HandleClickForExpedition(tc.activeExpedition, c, rf)
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

func (tc *TownhallController) GetHelperSuggestions() *gui.Suggestion {
	suggestion := GetHouseholdHelperSuggestions(tc.th.Household)
	if suggestion != nil {
		return suggestion
	}
	top := 15 + IconS + LargeIconD
	if int(tc.th.Household.Money) < int(tc.th.Household.Town.Stats.Global.Money)/10 {
		return &gui.Suggestion{
			Message: "Your townhall needs more money. You can either\nincrease tax rates or reduce the subsidies your town\ngives out for poor households.",
			Icon:    "coin", X: float64(24 + LargeIconD*2), Y: top + LargeIconD/2.0,
		}
	} else if len(tc.th.Household.Town.Country.Towns) > 1 && len(tc.th.Traders) == 0 && len(tc.th.Household.Town.Factories) > 0 {
		return &gui.Suggestion{
			Message: "Create traders to trade with other towns.\nYou can direct them by clicking on other town's marketplaces.",
			Icon:    "trader", X: float64(24 + LargeIconD*4), Y: top + LargeIconD/2.0,
		}
	}
	return nil
}

func CreateTraderButtonForTownhall(x, y float64, tc *TownhallController, vc *economy.VehicleConstruction, m navigation.IMap) *gui.Panel {
	p := &gui.Panel{}
	p.AddImageLabel("vehicles/"+vc.Name, 24, y, LargeIconS, LargeIconS, gui.ImageLabelStyleRegular)
	p.AddButton(&gui.SimpleButton{
		ButtonGUI: gui.ButtonGUI{Icon: "plus", X: x + 240, Y: y + float64(IconH/4), SX: IconS, SY: IconS, OnHoover: func() {
			tc.cp.HelperMessage("Create a new trader", true)
		}},
		ClickImpl: func() {
			h := tc.th.Household
			factory := social.PickFactory(h.Town.Factories, vc.BuildingExtensionType, tc.th.Household.Town.Marketplace.Building, m)
			if factory != nil {
				order := factory.CreateOrder(vc, h)
				if order != nil {
					h.AddTask(&economy.CreateTraderTask{
						Townhall: tc.th,
						PickupD:  factory.Household.Destination(building.NonExtension),
						Order:    order,
					})
				}
			}
		},
	})
	if tc.th.Household.Town.Marketplace != nil {
		p.AddTextLabel(fmt.Sprintf("$%v", social.VehiclePrice(tc.th.Household.Town.Marketplace, vc)), 24+x+float64(IconW)*2, y+float64(LargeIconD)/2)
	}
	return p
}

func SelectTraderButton(x, y float64, th *TownhallController, t *social.Trader) gui.Button {
	return &gui.SimpleButton{
		ButtonGUI: gui.ButtonGUI{Icon: "trader", X: x, Y: y, SX: IconS, SY: IconS, OnHoover: func() {
			th.cp.HelperMessage("Select trader", true)
		}},
		ClickImpl: func() {
			th.activeTrader = t
		},
		Highlight: func() bool {
			return t == th.activeTrader
		},
	}
}

func CreateExpeditionButtonForTownhall(x, y float64, tc *TownhallController, vc *economy.VehicleConstruction, m navigation.IMap) *gui.Panel {
	p := &gui.Panel{}
	p.AddImageLabel("vehicles/"+vc.Name, 24, y, LargeIconS, LargeIconS, gui.ImageLabelStyleRegular)
	p.AddButton(&gui.SimpleButton{
		ButtonGUI: gui.ButtonGUI{Icon: "plus", X: x + 240, Y: y + float64(IconH/4), SX: IconS, SY: IconS, OnHoover: func() {
			tc.cp.HelperMessage("Create a new expedition", true)
		}},
		ClickImpl: func() {
			h := tc.th.Household
			factory := social.PickFactory(h.Town.Factories, vc.BuildingExtensionType, nil, m)
			if factory != nil {
				order := factory.CreateOrder(vc, h)
				if order != nil {
					h.AddTask(&economy.CreateExpeditionTask{
						PickupD:  factory.Household.Destination(building.NonExtension),
						Order:    order,
						Townhall: tc.th,
					})
				}
			}
		},
	})
	if tc.th.Household.Town.Marketplace != nil {
		p.AddTextLabel(fmt.Sprintf("$%v", social.VehiclePrice(tc.th.Household.Town.Marketplace, vc)), 24+x+float64(IconW)*2, y+float64(LargeIconD)/2)
	}
	return p
}

func SelectExpeditionButton(x, y float64, tc *TownhallController, e *social.Expedition) gui.Button {
	return &gui.SimpleButton{
		ButtonGUI: gui.ButtonGUI{Icon: "vehicles/" + e.Vehicle.T.Name, X: x, Y: y, SX: IconS, SY: IconS, OnHoover: func() {
			tc.cp.HelperMessage("Select expedition", true)
		}},
		ClickImpl: func() {
			tc.activeExpedition = e
		},
		Highlight: func() bool {
			return e == tc.activeExpedition
		},
	}
}
