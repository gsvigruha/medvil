package controller

import (
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/tfriedel6/canvas/backend/goglbackend"
	"medvil/maps"
	"medvil/model"
	"medvil/model/building"
	"medvil/model/navigation"
	"medvil/model/social"
	"medvil/renderer"
	"sync"
)

const PerspectiveNE uint8 = 0
const PerspectiveSE uint8 = 1
const PerspectiveSW uint8 = 2
const PerspectiveNW uint8 = 3

const MaxRenderCnt = 10

type Window interface {
	GetGLFWWindow() *glfw.Window
	Close()
	SetController(*Controller)
}

type ClickHandler interface {
	HandleClick(c *Controller, rf *renderer.RenderedField) bool
	GetActiveFields(c *Controller, rf *renderer.RenderedField) []navigation.FieldWithContext
}

type ViewSettings struct {
	ShowHouseIcons       bool
	ShowAllocatedFields  bool
	ShowLabels           bool
	ShowTutorialMessages bool
}

var DefaultViewSettings = ViewSettings{
	ShowHouseIcons:       true,
	ShowAllocatedFields:  false,
	ShowLabels:           true,
	ShowTutorialMessages: true,
}

type Controller struct {
	X                         float64
	Y                         float64
	Window                    Window
	DX                        int
	DY                        int
	W                         int
	H                         int
	CenterX                   int
	CenterY                   int
	Perspective               uint8
	ViewSettings              ViewSettings
	Map                       *model.Map
	MapLock                   sync.Mutex
	RenderedFields            []*renderer.RenderedField
	RenderedBuildingParts     []renderer.RenderedBuildingPart
	RenderedTravellers        []*renderer.RenderedTraveller
	TempRenderedFields        []*renderer.RenderedField
	TempRenderedBuildingParts []renderer.RenderedBuildingPart
	TempRenderedTravellers    []*renderer.RenderedTraveller
	SelectedField             *navigation.Field
	ActiveSupplier            social.Supplier
	SelectedFarm              *social.Farm
	SelectedWorkshop          *social.Workshop
	SelectedMine              *social.Mine
	SelectedFactory           *social.Factory
	SelectedTower             *social.Tower
	SelectedTownhall          *social.Townhall
	SelectedConstruction      *building.Construction
	SelectedMarketplace       *social.Marketplace
	SelectedTraveller         *navigation.Traveller
	SelectedTrader            *social.Trader
	SelectedExpedition        *social.Expedition
	ReverseReferences         *model.ReverseReferences
	ControlPanel              *ControlPanel
	Country                   *social.Country
	ClickHandler              ClickHandler
	TimeSpeed                 int
	RenderCnt                 int
	ctx                       *goglbackend.GLContext
}

func (c *Controller) MoveCenter(dViewX, dViewY int) {
	var dCenterX, dCenterY = 0, 0
	switch c.Perspective {
	case PerspectiveNE:
		dCenterX, dCenterY = -dViewX+dViewY, -dViewX-dViewY
	case PerspectiveSE:
		dCenterX, dCenterY = dViewX+dViewY, -dViewX+dViewY
	case PerspectiveSW:
		dCenterX, dCenterY = dViewX-dViewY, dViewX+dViewY
	case PerspectiveNW:
		dCenterX, dCenterY = -dViewX-dViewY, dViewX-dViewY
	}
	c.CenterX += dCenterX
	c.CenterY += dCenterY
}

func (c *Controller) KeyboardCallback(wnd *glfw.Window, key glfw.Key, code int, action glfw.Action, mod glfw.ModifierKey) {
	if key == glfw.KeyEnter && action == glfw.Release {
		c.Perspective = (c.Perspective + 1) % 4
	}
	if action == glfw.Press {
		if key == glfw.KeyQ {
			c.Perspective = (c.Perspective + 1) % 4
		}
		if key == glfw.KeyE {
			c.Perspective = (c.Perspective - 1) % 4
		}
		if key == glfw.KeyUp || key == glfw.KeyW {
			c.DY = -1
		}
		if key == glfw.KeyDown || key == glfw.KeyS {
			c.DY = 1
		}
		if key == glfw.KeyLeft || key == glfw.KeyA {
			c.DX = -1
		}
		if key == glfw.KeyRight || key == glfw.KeyD {
			c.DX = 1
		}
		if key == glfw.KeyTab {
			CPActionTimeScaleChange(c)
		}
		if key == glfw.KeyL {
			c.Load(GetLatestFile())
		}
		if key == glfw.KeyEscape {
			c.Save("latest_autosave.mdvl")
			c.Window.Close()
		}
	}
	if action == glfw.Release {
		if key == glfw.KeyUp || key == glfw.KeyW {
			c.DY = 0
		}
		if key == glfw.KeyDown || key == glfw.KeyS {
			c.DY = 0
		}
		if key == glfw.KeyLeft || key == glfw.KeyA {
			c.DX = 0
		}
		if key == glfw.KeyRight || key == glfw.KeyD {
			c.DX = 0
		}
	}
}

func (c *Controller) ShowBuildingController() {
	c.Reset()
	BuildingsToControlPanel(c.ControlPanel)
}

func (c *Controller) ShowLibraryController() {
	c.Reset()
	LibraryToControlPanel(c.ControlPanel)
}

func (c *Controller) ShowMapController() {
	c.Reset()
	MapToControlPanel(c.ControlPanel)
}

func (c *Controller) ShowInfraController() {
	c.Reset()
	InfraToControlPanel(c.ControlPanel)
}

func (c *Controller) ShowDemolishController() {
	c.Reset()
	DemolishToControlPanel(c.ControlPanel, c.GetActiveTownhall())
}

func (c *Controller) GetActiveTownhall() *social.Townhall {
	if c.SelectedTownhall != nil {
		return c.SelectedTownhall
	}
	if c.ActiveSupplier != nil {
		if town, ok := c.ActiveSupplier.(*social.Town); ok {
			return town.Townhall
		}
	}
	return nil
}

func (c *Controller) Refresh() {
	c.ReverseReferences = c.Map.ReverseReferences()
	c.ControlPanel.Refresh()
}

func (c *Controller) Reset() {
	c.SelectedField = nil
	c.SelectedFarm = nil
	c.SelectedMine = nil
	c.SelectedFactory = nil
	c.SelectedTower = nil
	c.SelectedWorkshop = nil
	c.SelectedTraveller = nil
	c.SelectedTrader = nil
	c.SelectedExpedition = nil
	c.ClickHandler = nil
	c.ControlPanel.GetHelperPanel()
}

func (c *Controller) CaptureRenderedField(x, y float64) *renderer.RenderedField {
	for i := range c.RenderedFields {
		rf := c.RenderedFields[i]
		if rf.Contains(c.X, c.Y) {
			return rf
		}
	}
	return nil
}

func (c *Controller) GetActiveFields() []navigation.FieldWithContext {
	if c.ClickHandler != nil {
		rf := c.CaptureRenderedField(c.X, c.Y)
		if rf != nil {
			return c.ClickHandler.GetActiveFields(c, rf)
		}
	} else if c.SelectedTraveller != nil {
		return c.SelectedTraveller.GetPathFields(c.Map)
	}
	return nil
}

func (c *Controller) MouseButtonCallback(wnd *glfw.Window, button glfw.MouseButton, action glfw.Action, mod glfw.ModifierKey) {
	if action == glfw.Press && button == glfw.MouseButton1 {
		c.RenderCnt = MaxRenderCnt
		c.ControlPanel.CaptureClick(c.X, c.Y)
		if c.X < ControlPanelSX {
			return
		}
		rf := c.CaptureRenderedField(c.X, c.Y)
		if c.ClickHandler != nil && rf != nil {
			if c.ClickHandler.HandleClick(c, rf) {
				return
			}
		}
		var rbp renderer.RenderedBuildingPart
		for _, rbpI := range c.RenderedBuildingParts {
			if rbpI.Contains(c.X, c.Y) {
				rbp = rbpI
			}
		}
		if rbp != nil {
			c.Reset()
			c.SelectedTownhall = c.ReverseReferences.BuildingToTownhall[rbp.GetBuilding()]
			if c.SelectedTownhall != nil {
				c.ActiveSupplier = c.SelectedTownhall.Household.Town
				TownhallToControlPanel(c.ControlPanel, c.SelectedTownhall)
			}
			c.SelectedMarketplace = c.ReverseReferences.BuildingToMarketplace[rbp.GetBuilding()]
			if c.SelectedMarketplace != nil {
				c.ActiveSupplier = c.SelectedMarketplace.Town
				MarketplaceToControlPanel(c.ControlPanel, c.SelectedMarketplace)
			}
			c.SelectedFarm = c.ReverseReferences.BuildingToFarm[rbp.GetBuilding()]
			if c.SelectedFarm != nil {
				c.ActiveSupplier = c.SelectedFarm.Household.Town
				FarmToControlPanel(c.ControlPanel, c.SelectedFarm)
			}
			c.SelectedWorkshop = c.ReverseReferences.BuildingToWorkshop[rbp.GetBuilding()]
			if c.SelectedWorkshop != nil {
				c.ActiveSupplier = c.SelectedWorkshop.Household.Town
				WorkshopToControlPanel(c.ControlPanel, c.SelectedWorkshop)
			}
			c.SelectedMine = c.ReverseReferences.BuildingToMine[rbp.GetBuilding()]
			if c.SelectedMine != nil {
				c.ActiveSupplier = c.SelectedMine.Household.Town
				MineToControlPanel(c.ControlPanel, c.SelectedMine)
			}
			c.SelectedFactory = c.ReverseReferences.BuildingToFactory[rbp.GetBuilding()]
			if c.SelectedFactory != nil {
				c.ActiveSupplier = c.SelectedFactory.Household.Town
				FactoryToControlPanel(c.ControlPanel, c.SelectedFactory)
			}
			c.SelectedTower = c.ReverseReferences.BuildingToTower[rbp.GetBuilding()]
			if c.SelectedTower != nil {
				c.ActiveSupplier = c.SelectedTower.Household.Town
				TowerToControlPanel(c.ControlPanel, c.SelectedTower)
			}
			c.SelectedConstruction = c.ReverseReferences.BuildingToConstruction[rbp.GetBuilding()]
			if c.SelectedConstruction != nil {
				ConstructionToControlPanel(c.ControlPanel, c.SelectedConstruction)
			}
			return
		}
		for i := range c.RenderedTravellers {
			rt := c.RenderedTravellers[i]
			if rt.Contains(c.X, c.Y) {
				c.Reset()
				c.SelectedTraveller = rt.Traveller
				person := c.ReverseReferences.TravellerToPerson[c.SelectedTraveller]
				if person != nil {
					PersonToControlPanel(c.ControlPanel, person)
				}
				trader := c.ReverseReferences.TravellerToTrader[c.SelectedTraveller]
				if trader != nil {
					c.SelectedTrader = trader
					TraderToControlPanel(c.ControlPanel, trader)
				}
				expedition := c.ReverseReferences.TravellerToExpedition[c.SelectedTraveller]
				if expedition != nil {
					c.SelectedExpedition = expedition
					c.ActiveSupplier = expedition
					ExpeditionToControlPanel(c.ControlPanel, expedition)
				}
				return
			}
		}
		if c.ClickHandler == nil && rf != nil {
			c.Reset()
			c.SelectedField = rf.F
			FieldToControlPanel(c.ControlPanel, c.SelectedField)
			return
		}
	}
}

func (c *Controller) RenderTick() {
	c.RenderCnt++
	if c.RenderCnt >= MaxRenderCnt {
		c.RenderCnt = 0
	}
	c.MoveCenter(c.DX, c.DY)
}

func (c *Controller) MouseMoveCallback(wnd *glfw.Window, x float64, y float64) {
	w, h := wnd.GetSize()
	fbw, fbh := wnd.GetFramebufferSize()
	c.X = x * float64(fbw) / float64(w)
	c.Y = y * float64(fbh) / float64(h)
}

func (c *Controller) MouseScrollCallback(wnd *glfw.Window, x float64, y float64) {
}

func Link(window Window, ctx *goglbackend.GLContext) *Controller {
	wnd := window.GetGLFWWindow()
	W, H := wnd.GetFramebufferSize()
	controlPanel := &ControlPanel{}
	c := &Controller{H: H, W: W, Window: window, ControlPanel: controlPanel, TimeSpeed: 1, ViewSettings: DefaultViewSettings}
	controlPanel.Setup(c, ctx)
	wnd.SetKeyCallback(c.KeyboardCallback)
	wnd.SetMouseButtonCallback(c.MouseButtonCallback)
	wnd.SetCursorPosCallback(c.MouseMoveCallback)
	wnd.SetScrollCallback(c.MouseScrollCallback)
	c.ctx = ctx
	window.SetController(c)
	return c
}

func (c *Controller) SwapRenderedObjects() {
	c.RenderedFields = c.TempRenderedFields
	c.RenderedBuildingParts = c.TempRenderedBuildingParts
	c.RenderedTravellers = c.TempRenderedTravellers
	c.TempRenderedFields = []*renderer.RenderedField{}
	c.TempRenderedBuildingParts = []renderer.RenderedBuildingPart{}
	c.TempRenderedTravellers = []*renderer.RenderedTraveller{}
}

func (c *Controller) AddRenderedField(rf *renderer.RenderedField) {
	c.TempRenderedFields = append(c.TempRenderedFields, rf)
}

func (c *Controller) AddRenderedBuildingPart(rbp renderer.RenderedBuildingPart) {
	c.TempRenderedBuildingParts = append(c.TempRenderedBuildingParts, rbp)
}

func (c *Controller) AddRenderedTraveller(rt *renderer.RenderedTraveller) {
	c.TempRenderedTravellers = append(c.TempRenderedTravellers, rt)
}

func (c *Controller) Save(fileName string) {
	if c.Map == nil {
		return
	}
	c.MapLock.Lock()
	maps.Serialize(c.Map, "saved/"+fileName)
	c.MapLock.Unlock()
}

func (c *Controller) Load(fileName string) {
	m := maps.Deserialize("saved/" + fileName).(*model.Map)
	c.MapLock.Lock()
	c.Map = m
	c.LinkMap()
	c.MapLock.Unlock()
}

func (c *Controller) NewMap(config maps.MapConfig) {
	m := maps.NewMap(config)
	c.MapLock.Lock()
	c.Map = m
	c.LinkMap()
	c.MapLock.Unlock()
}

func (c *Controller) LinkMap() {
	c.Country = c.Map.Countries[0]
	town := c.Map.Countries[0].Towns[0]
	c.CenterX = int(town.Townhall.Household.Building.X)
	c.CenterY = int(town.Townhall.Household.Building.Y)
	c.ActiveSupplier = town
	c.ControlPanel.GenerateButtons()
	c.Refresh()
}
