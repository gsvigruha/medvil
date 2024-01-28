package controller

import (
	"github.com/tfriedel6/canvas"
	"github.com/tfriedel6/canvas/backend/goglbackend"
	"medvil/model/building"
	"medvil/view/gui"
	"path/filepath"
	"reflect"
	"strconv"
)

var ControlPanelSX = 300.0
var ControlPanelSY = 800.0
var ControlPanelDynamicPanelTop = 0.1
var ControlPanelDynamicPanelSY = 0.6

const CPButtonHighlightNone = 0
const CPButtonHighlightSmall = 1
const CPButtonHighlightLarge = 2

type ControlSubPanel interface {
	Panel
	GetHelperSuggestions() *gui.Suggestion
}

type ControlPanel struct {
	topPanel            *gui.Panel
	dynamicPanel        ControlSubPanel
	HelperPanel         *gui.Panel
	SelectedHelperPanel *gui.Panel
	dateLabel           *gui.TextLabel
	moneyLabel          *gui.TextLabel
	peopleLabel         *gui.TextLabel
	artifactsLabel      *gui.TextLabel
	buildingsLabel      *gui.TextLabel
	timeButton          *ControlPanelButton
	suggestion          *gui.Suggestion
	C                   *Controller
	buffer              *canvas.Canvas
	HelperBuffer        *canvas.Canvas
}

type ControlPanelButton struct {
	b         gui.ButtonGUI
	c         *Controller
	action    func(*Controller)
	highlight uint8
}

func (b ControlPanelButton) Click() {
	b.action(b.c)
}

func (b ControlPanelButton) Render(cv *canvas.Canvas) {
	if b.highlight == CPButtonHighlightSmall {
		cv.SetFillStyle("#48C")
		cv.FillRect(b.b.X, b.b.Y, b.b.SX, b.b.SY)
	} else if b.highlight == CPButtonHighlightLarge {
		cv.SetFillStyle("#8AD")
		cv.FillRect(b.b.X, b.b.Y, b.b.SX, b.b.SY)
	}
	b.b.Render(cv)
}

func (b *ControlPanelButton) SetHoover(h bool) {
	b.b.SetHoover(h)
}

func (b ControlPanelButton) Contains(x float64, y float64) bool {
	return b.b.Contains(x, y)
}

func (b ControlPanelButton) Enabled() bool {
	return b.b.Enabled()
}

func CPActionCancel(c *Controller) {
	c.Reset()
	c.ControlPanel.dynamicPanel = nil
}

func CPActionTimeScaleChange(c *Controller) {
	if c.TimeSpeed == 1 {
		c.TimeSpeed = 5
		c.ControlPanel.timeButton.highlight = CPButtonHighlightSmall
	} else if c.TimeSpeed == 5 {
		c.TimeSpeed = 20
		c.ControlPanel.timeButton.highlight = CPButtonHighlightLarge
	} else {
		c.TimeSpeed = 1
		c.ControlPanel.timeButton.highlight = CPButtonHighlightNone
	}
}

func (p *ControlPanel) Refresh() {
	p.HelperPanel.Clear()
	p.topPanel.CaptureMove(p.C.X, p.C.Y)
	p.dateLabel.Text = strconv.Itoa(
		int(p.C.Map.Calendar.Day)) + ", " +
		strconv.Itoa(int(p.C.Map.Calendar.Month)) + ", " +
		strconv.Itoa(int(p.C.Map.Calendar.Year))
	stats := p.C.Country.Stats()
	p.moneyLabel.Text = strconv.Itoa(int(stats.Global.Money))
	p.peopleLabel.Text = strconv.Itoa(int(stats.Global.People))
	p.artifactsLabel.Text = strconv.Itoa(int(stats.Global.Artifacts))
	p.buildingsLabel.Text = strconv.Itoa(int(stats.Global.Buildings))
	if p.dynamicPanel != nil {
		p.dynamicPanel.CaptureMove(p.C.X, p.C.Y)
		p.dynamicPanel.Refresh()
	}
}

func (p *ControlPanel) Clear() {
	if p.dynamicPanel != nil {
		p.dynamicPanel.Clear()
	}
}

func (p *ControlPanel) GetHelperPanel(clear bool) *gui.Panel {
	if clear {
		p.HelperPanel.Clear()
	}
	return p.HelperPanel
}

func (p *ControlPanel) SetDims(size uint8) {
	c := p.C
	if size == SizeSmall {
		ControlPanelSX = 300.0
		ControlPanelSY = float64(c.H)
		IconS = 24.0
		IconW = 30
		IconH = 30
		LargeIconS = 36.0
		LargeIconD = 40.0
		gui.FontSize = 12.0
		ScaleBuildingControllerElements(0.75)
	} else if size == SizeMedium {
		ControlPanelSX = 400.0
		ControlPanelSY = float64(c.H)
		IconS = 32.0
		IconW = 40
		IconH = 40
		LargeIconS = 48.0
		LargeIconD = 52.0
		gui.FontSize = 16.0
		ScaleBuildingControllerElements(1.0)
	} else if size == SizeLarge {
		ControlPanelSX = 600.0
		ControlPanelSY = float64(c.H)
		IconS = 48.0
		IconW = 60
		IconH = 60
		LargeIconS = 72.0
		LargeIconD = 80.0
		gui.FontSize = 24.0
		ScaleBuildingControllerElements(1.5)
	}
}

func (p *ControlPanel) SetupDims(width, height int) {
	c := p.C
	c.H = height
	c.W = width
	if c.ViewSettings.Size == SizeAuto {
		if c.H < 1000 {
			p.SetDims(SizeSmall)
		} else if c.H < 1500 {
			p.SetDims(SizeMedium)
		} else {
			p.SetDims(SizeLarge)
		}
	} else {
		p.SetDims(c.ViewSettings.Size)
	}

	{
		offscreen, _ := goglbackend.NewOffscreen(int(ControlPanelSX), int(ControlPanelSY), false, c.ctx)
		p.buffer = canvas.New(offscreen)
	}

	p.topPanel.SX = ControlPanelSX
	p.topPanel.SY = ControlPanelSY
	p.HelperPanel.Y = ControlPanelSY * 0.95
	p.HelperPanel.SX = ControlPanelSX
	p.HelperPanel.SY = float64(IconH) * 2.0
	{
		offscreen, _ := goglbackend.NewOffscreen(int(p.HelperPanel.SX), int(p.HelperPanel.SY), false, c.ctx)
		p.HelperBuffer = canvas.New(offscreen)
	}

	if c.Map != nil {
		p.GenerateButtons()
	}
}

func (p *ControlPanel) Setup(c *Controller, ctx *goglbackend.GLContext) {
	p.C = c

	p.topPanel = &gui.Panel{X: 0, Y: 0, SX: ControlPanelSX, SY: ControlPanelSY}
	p.HelperPanel = &gui.Panel{X: 0, Y: ControlPanelSY * 0.95, SX: ControlPanelSX, SY: ControlPanelSY * 0.05}
	p.SelectedHelperPanel = &gui.Panel{X: 0, Y: ControlPanelSY * 0.95, SX: ControlPanelSX, SY: ControlPanelSY * 0.05}

	p.SetupDims(p.C.W, p.C.H)
}

func (p *ControlPanel) GenerateButtons() {
	p.topPanel.Clear()
	c := p.C
	ih := 4.0
	th := IconS/2 - gui.FontSize/2
	p.dateLabel = p.topPanel.AddTextLabel("", ControlPanelSX*0.03, th+gui.FontSize)
	p.topPanel.AddImageLabel("coin", ControlPanelSX*0.25, ih, IconS, IconS, gui.ImageLabelStyleRegular)
	p.moneyLabel = p.topPanel.AddTextLabel("", ControlPanelSX*0.25+IconS, th+gui.FontSize)
	p.topPanel.AddImageLabel("person", ControlPanelSX*0.5, ih, IconS, IconS, gui.ImageLabelStyleRegular)
	p.peopleLabel = p.topPanel.AddTextLabel("", ControlPanelSX*0.5+IconS, th+gui.FontSize)
	p.topPanel.AddImageLabel("barrel", ControlPanelSX*0.65, ih, IconS, IconS, gui.ImageLabelStyleRegular)
	p.artifactsLabel = p.topPanel.AddTextLabel("", ControlPanelSX*0.65+IconS, th+gui.FontSize)
	p.topPanel.AddImageLabel("house", ControlPanelSX*0.85, ih, IconS, IconS, gui.ImageLabelStyleRegular)
	p.buildingsLabel = p.topPanel.AddTextLabel("", ControlPanelSX*0.85+IconS, th+gui.FontSize)

	iconTop := 15 + IconS
	p.topPanel.AddButton(&gui.SimpleButton{
		ButtonGUI: gui.ButtonGUI{Icon: "house", X: float64(24 + LargeIconD*0), Y: iconTop, SX: LargeIconS, SY: LargeIconS, OnHoover: func() {
			p.HelperMessage("Create buildings", true)
		}},
		Highlight: func() bool { return p.IsBuildingType() },
		ClickImpl: func() {
			suggestion := p.GetHelperSuggestions()
			if suggestion == nil {
				c.ShowBuildingController()
			} else {
				if suggestion.Icon == "farm" {
					c.ShowBuildingControllerForType(building.BuildingTypeFarm)
				} else if suggestion.Icon == "workshop" {
					c.ShowBuildingControllerForType(building.BuildingTypeWorkshop)
				} else if suggestion.Icon == "factory" {
					c.ShowBuildingControllerForType(building.BuildingTypeFactory)
				} else if suggestion.Icon == "mine" {
					c.ShowBuildingControllerForType(building.BuildingTypeMine)
				} else if suggestion.Icon == "tower" {
					c.ShowBuildingControllerForType(building.BuildingTypeTower)
				} else if suggestion.Icon == "townhall" {
					c.ShowBuildingControllerForType(building.BuildingTypeTownhall)
				} else if suggestion.Icon == "market" {
					c.ShowBuildingControllerForType(building.BuildingTypeMarket)
				} else {
					c.ShowBuildingController()
				}
			}
		}})
	p.topPanel.AddButton(&gui.SimpleButton{
		ButtonGUI: gui.ButtonGUI{Icon: "infra", X: float64(24 + LargeIconD*1), Y: iconTop, SX: LargeIconS, SY: LargeIconS, OnHoover: func() {
			p.HelperMessage("Create infrastructure like roads or walls", true)
		}},
		Highlight: func() bool { return p.IsInfraType() },
		ClickImpl: func() { c.ShowInfraController() }})
	p.topPanel.AddButton(&gui.SimpleButton{
		ButtonGUI: gui.ButtonGUI{Icon: "demolish", X: float64(24 + LargeIconD*2), Y: iconTop, SX: LargeIconS, SY: LargeIconS, OnHoover: func() {
			p.HelperMessage("Demolish buildings", true)
		},
			Disabled: func() bool { return c.GetActiveTownhall() == nil }},
		Highlight: func() bool { return p.IsDynamicPanelType("DemolishController") },
		ClickImpl: func() { c.ShowDemolishController() }})
	p.topPanel.AddButton(&gui.SimpleButton{
		ButtonGUI: gui.ButtonGUI{Icon: "map", X: float64(24 + LargeIconD*3), Y: iconTop, SX: LargeIconS, SY: LargeIconS, OnHoover: func() {
			p.HelperMessage("View map and charts", true)
		}},
		Highlight: func() bool { return p.IsDynamicPanelType("MapController") },
		ClickImpl: func() { c.ShowMapController() }})
	p.topPanel.AddButton(&gui.SimpleButton{
		ButtonGUI: gui.ButtonGUI{Icon: "library", X: float64(24 + LargeIconD*4), Y: iconTop, SX: LargeIconS, SY: LargeIconS, OnHoover: func() {
			p.HelperMessage("Load, save and settings", true)
		}},
		Highlight: func() bool { return p.IsDynamicPanelType("LibraryController") },
		ClickImpl: func() { c.ShowLibraryController() }})
	p.topPanel.AddButton(&ControlPanelButton{
		b: gui.ButtonGUI{Icon: "cancel", X: float64(24 + LargeIconD*5), Y: iconTop, SX: LargeIconS, SY: LargeIconS, OnHoover: func() {
			p.HelperMessage("Cancel", true)
		}},
		c: c, action: CPActionCancel})
	p.timeButton = &ControlPanelButton{
		b: gui.ButtonGUI{Icon: "time", X: float64(24 + LargeIconD*6), Y: iconTop, SX: LargeIconS, SY: LargeIconS, OnHoover: func() {
			p.HelperMessage("Speed up game", true)
		}},
		c: c, action: CPActionTimeScaleChange}

	p.topPanel.AddButton(p.timeButton)
}

func (p *ControlPanel) SetDynamicPanel(dp ControlSubPanel) {
	p.Clear()
	p.dynamicPanel = dp
}

func (p *ControlPanel) CaptureClick(x, y float64) {
	p.topPanel.CaptureClick(x, y)
	if p.dynamicPanel != nil {
		p.dynamicPanel.CaptureClick(x, y)
	}
}

func (p *ControlPanel) CaptureMove(x, y float64) {
	p.topPanel.CaptureMove(x, y)
	if p.dynamicPanel != nil {
		p.dynamicPanel.CaptureMove(x, y)
	}
}

func (p *ControlPanel) Render(cv *canvas.Canvas, c *Controller) {
	if c.RenderCnt == 0 {
		p.topPanel.Render(p.buffer)
		if p.dynamicPanel != nil {
			p.dynamicPanel.Render(p.buffer)
		}
		if !p.HelperPanel.IsEmpty() {
			p.HelperBuffer.SetFillStyle(filepath.FromSlash("texture/wood.png"))
			p.HelperBuffer.FillRect(0, 0, p.HelperPanel.SX, p.HelperPanel.SY)
			p.HelperBuffer.SetStrokeStyle("#DDD")
			p.HelperBuffer.SetLineWidth(2)
			p.HelperBuffer.StrokeRect(1, 1, p.HelperPanel.SX-2, p.HelperPanel.SY-2)
			p.HelperPanel.Render(p.HelperBuffer)
		}
		if !p.SelectedHelperPanel.IsEmpty() {
			p.SelectedHelperPanel.Render(p.buffer)
		}
		p.GetSuggestion()
	}
	cv.DrawImage(p.buffer, 0, 0, ControlPanelSX, ControlPanelSY)
	if p.suggestion != nil && p.C.ViewSettings.ShowSuggestions {
		p.suggestion.Render(cv, LargeIconS, LargeIconD)
	}
	if !p.HelperPanel.IsEmpty() {
		cv.DrawImage(p.HelperBuffer, c.X, c.Y, p.HelperPanel.SX, p.HelperPanel.SY)
	}
}

func (p *ControlPanel) IsDynamicPanelType(typeName string) bool {
	if p.dynamicPanel == nil {
		return false
	}
	return reflect.TypeOf(p.dynamicPanel).String() == ("*controller." + typeName)
}

func (p *ControlPanel) IsBuildingTypeOf(bt building.BuildingType) bool {
	if p.C.ClickHandler == nil {
		return false
	}
	if bc, ok := p.C.ClickHandler.(*BuildingsController); ok {
		return bc.Plan.BuildingType == bt
	}
	return false
}

func (p *ControlPanel) IsBuildingType() bool {
	if p.C.ClickHandler == nil {
		return false
	}
	_, ok := p.C.ClickHandler.(*BuildingsController)
	return ok
}

func (p *ControlPanel) IsInfraType() bool {
	if p.C.ClickHandler == nil {
		return false
	}
	_, ok := p.C.ClickHandler.(*InfraController)
	return ok
}

func (p *ControlPanel) SelectedHelperMessage(msg string) {
	p.SelectedHelperPanel.Clear()
	p.SelectedHelperPanel.AddTextLabel(msg, 24, ControlPanelSY*0.95+IconS-gui.FontSize/2.0)
}

func (p *ControlPanel) HelperMessage(msg string, actionable bool) {
	hp := p.GetHelperPanel(true)
	if actionable {
		hp.AddImageLabel("arrow_small_right", 24, float64(IconH)/2.0, IconS, IconS, gui.ImageLabelStyleRegular)
	} else {
		hp.AddImageLabel("help", 24, float64(IconH)/2.0, IconS, IconS, gui.ImageLabelStyleRegular)
	}
	hp.AddTextLabel(msg, 24+float64(IconW), float64(IconH)/2.0+IconS-gui.FontSize/2.0)
}

func (p *ControlPanel) GetSuggestion() {
	if p.dynamicPanel != nil {
		p.suggestion = p.dynamicPanel.GetHelperSuggestions()
	} else {
		p.suggestion = p.GetHelperSuggestions()
	}
}

func (p *ControlPanel) GetHelperSuggestions() *gui.Suggestion {
	if p.C.Map != nil {
		if len(p.C.Map.Countries[0].Towns[0].Farms) == 0 && len(p.C.Map.Countries[0].Towns[0].Constructions) == 0 {
			return &gui.Suggestion{
				Message: "Build farms. Your town needs farms to produce\ngrain, sheep, textile and logs.",
				Icon:    "farm", X: float64(24 + LargeIconD*1), Y: IconS + 15 + LargeIconD/2.0,
			}
		} else if len(p.C.Map.Countries[0].Towns[0].Workshops) == 0 && len(p.C.Map.Countries[0].Towns[0].Constructions) == 0 {
			return &gui.Suggestion{
				Message: "Build workshops. Workshops can turn raw materials\ninto food or other goods like building materials.",
				Icon:    "workshop", X: float64(24 + LargeIconD*1), Y: IconS + 15 + LargeIconD/2.0,
			}
		} else if len(p.C.Map.Countries[0].Towns[0].Mines) == 0 && len(p.C.Map.Countries[0].Towns[0].Constructions) == 0 {
			return &gui.Suggestion{
				Message: "Build mines to get metals and raw building materials.\nClay and stone are used to produce bricks for roads and houses.\nGold is used as a currency, iron is needed for weapons and vehicles.",
				Icon:    "mine", X: float64(24 + LargeIconD*1), Y: IconS + 15 + LargeIconD/2.0,
			}
		} else if len(p.C.Map.Countries[0].Towns[0].Roads) == 0 && len(p.C.Map.Countries[0].Towns[0].Constructions) == 0 {
			return &gui.Suggestion{
				Message: "Build roads and bridges to make commuting faster for your\n villagers. This will make your economy more efficient.\nHowever, roads need to be maintained by the townhall.",
				Icon:    "infra/cobble_road", X: float64(24 + LargeIconD*2), Y: IconS + 15 + LargeIconD/2.0,
			}
		} else if len(p.C.Map.Countries[0].Towns[0].Factories) == 0 && len(p.C.Map.Countries[0].Towns[0].Constructions) == 0 {
			return &gui.Suggestion{
				Message: "Build factories to create vehicles. Vehicles make it more\nefficient for your villagers to transport goods to and\nfrom the market. They can also be used to create traders\n and to launch expeditions.",
				Icon:    "factory", X: float64(24 + LargeIconD*1), Y: IconS + 15 + LargeIconD/2.0,
			}
		} else if len(p.C.Map.Countries[0].Towns) == 1 && p.C.Map.Countries[0].Towns[0].Stats.Global.People > 80 {
			return &gui.Suggestion{
				Message: "Establish a new town by building a new townhall.\nYou can extract materials from distant lands and trade.",
				Icon:    "townhall", X: float64(24 + LargeIconD*1), Y: IconS + 15 + LargeIconD/2.0,
			}
		} else if p.C.GetActiveTownhall() != nil && p.C.GetActiveTownhall().Household.Town.Marketplace == nil && len(p.C.GetActiveTownhall().Household.Town.Constructions) == 0 {
			return &gui.Suggestion{
				Message: "Build a marketplace for your new town. The villagers need a place\nto trade goods with each other.",
				Icon:    "market", X: float64(24 + LargeIconD*1), Y: IconS + 15 + LargeIconD/2.0,
			}
		} else if len(p.C.Map.Countries[0].Towns[0].Towers) == 0 && len(p.C.Map.Countries[0].Towns[0].Constructions) == 0 {
			return &gui.Suggestion{
				Message: "Build towers and walls to protect your town from the outlaws.\nThey live in small villages with wooden buildings and\nsteal your crops if you build too close to them.",
				Icon:    "tower", X: float64(24 + LargeIconD*1), Y: IconS + 15 + LargeIconD/2.0,
			}
		}
	}
	return nil
}
