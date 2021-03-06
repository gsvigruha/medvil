package controller

import (
	"github.com/go-gl/glfw/v3.3/glfw"
	"medvil/model"
	"medvil/model/building"
	"medvil/model/navigation"
	"medvil/model/social"
	"medvil/model/time"
	"medvil/renderer"
)

const PerspectiveNE uint8 = 0
const PerspectiveSE uint8 = 1
const PerspectiveSW uint8 = 2
const PerspectiveNW uint8 = 3

type ClickHandler interface {
	HandleClick(c *Controller, rf *renderer.RenderedField) bool
	GetActiveFields(c *Controller, rf *renderer.RenderedField) []navigation.FieldWithContext
}

type Controller struct {
	X                         float64
	Y                         float64
	W                         int
	H                         int
	CenterX                   int
	CenterY                   int
	Perspective               uint8
	Map                       *model.Map
	Calendar                  *time.CalendarType
	RenderedFields            []*renderer.RenderedField
	RenderedBuildingParts     []renderer.RenderedBuildingPart
	TempRenderedFields        []*renderer.RenderedField
	TempRenderedBuildingParts []renderer.RenderedBuildingPart
	SelectedField             *navigation.Field
	ActiveTown                *social.Town
	SelectedFarm              *social.Farm
	SelectedWorkshop          *social.Workshop
	SelectedMine              *social.Mine
	SelectedFactory           *social.Factory
	SelectedTower             *social.Tower
	SelectedTownhall          *social.Townhall
	SelectedConstruction      *building.Construction
	SelectedMarketplace       *social.Marketplace
	ReverseReferences         *model.ReverseReferences
	ControlPanel              *ControlPanel
	Country                   *social.Country
	ClickHandler              ClickHandler
	TimeSpeed                 int
}

func (c *Controller) MoveCenter(dViewX, dViewY int) {
	var dCenterX, dCenterY = 0, 0
	switch c.Perspective {
	case PerspectiveNE:
		dCenterX, dCenterY = -dViewX+dViewY, -dViewX-dViewY
	case PerspectiveSE:
		dCenterX, dCenterY = -dViewX+dViewY, dViewX+dViewY
	case PerspectiveSW:
		dCenterX, dCenterY = dViewX-dViewY, dViewX+dViewY
	case PerspectiveNW:
		dCenterX, dCenterY = dViewX-dViewY, -dViewX-dViewY
	}
	c.CenterX += dCenterX
	c.CenterY += dCenterY
}

func (c *Controller) KeyboardCallback(wnd *glfw.Window, key glfw.Key, code int, action glfw.Action, mod glfw.ModifierKey) {
	if key == glfw.KeyEnter && action == glfw.Release {
		c.Perspective = (c.Perspective + 1) % 4
	}
	if action == glfw.Press {
		if key == glfw.KeyUp {
			c.MoveCenter(0, -2)
		}
		if key == glfw.KeyDown {
			c.MoveCenter(0, 2)
		}
		if key == glfw.KeyLeft {
			c.MoveCenter(-2, 0)
		}
		if key == glfw.KeyRight {
			c.MoveCenter(2, 0)
		}
	}
}

func (c *Controller) ShowBuildingController(bt building.BuildingType) {
	c.Reset()
	BuildingsToControlPanel(c.ControlPanel, bt)
}

func (c *Controller) ShowInfraController() {
	c.Reset()
	InfraToControlPanel(c.ControlPanel)
}

func (c *Controller) ShowNewTownController() {
	c.Reset()
	NewTownToControlPanel(c.ControlPanel, c.SelectedTownhall)
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
	c.ClickHandler = nil
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
	}
	return nil
}

func (c *Controller) MouseButtonCallback(wnd *glfw.Window, button glfw.MouseButton, action glfw.Action, mod glfw.ModifierKey) {
	if action == glfw.Press && button == glfw.MouseButton1 {
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
		for i := range c.RenderedBuildingParts {
			rbp := c.RenderedBuildingParts[i]
			if rbp.Contains(c.X, c.Y) {
				c.Reset()
				c.SelectedTownhall = c.ReverseReferences.BuildingToTownhall[rbp.GetBuilding()]
				if c.SelectedTownhall != nil {
					c.ActiveTown = c.SelectedTownhall.Household.Town
					TownhallToControlPanel(c.ControlPanel, c.SelectedTownhall)
				}
				c.SelectedMarketplace = c.ReverseReferences.BuildingToMarketplace[rbp.GetBuilding()]
				if c.SelectedMarketplace != nil {
					c.ActiveTown = c.SelectedMarketplace.Town
					MarketplaceToControlPanel(c.ControlPanel, c.SelectedMarketplace)
				}
				c.SelectedFarm = c.ReverseReferences.BuildingToFarm[rbp.GetBuilding()]
				if c.SelectedFarm != nil {
					c.ActiveTown = c.SelectedFarm.Household.Town
					FarmToControlPanel(c.ControlPanel, c.SelectedFarm)
				}
				c.SelectedWorkshop = c.ReverseReferences.BuildingToWorkshop[rbp.GetBuilding()]
				if c.SelectedWorkshop != nil {
					c.ActiveTown = c.SelectedWorkshop.Household.Town
					WorkshopToControlPanel(c.ControlPanel, c.SelectedWorkshop)
				}
				c.SelectedMine = c.ReverseReferences.BuildingToMine[rbp.GetBuilding()]
				if c.SelectedMine != nil {
					c.ActiveTown = c.SelectedMine.Household.Town
					MineToControlPanel(c.ControlPanel, c.SelectedMine)
				}
				c.SelectedFactory = c.ReverseReferences.BuildingToFactory[rbp.GetBuilding()]
				if c.SelectedFactory != nil {
					c.ActiveTown = c.SelectedFactory.Household.Town
					FactoryToControlPanel(c.ControlPanel, c.SelectedFactory)
				}
				c.SelectedTower = c.ReverseReferences.BuildingToTower[rbp.GetBuilding()]
				if c.SelectedTower != nil {
					c.ActiveTown = c.SelectedTower.Household.Town
					TowerToControlPanel(c.ControlPanel, c.SelectedTower)
				}
				c.SelectedConstruction = c.ReverseReferences.BuildingToConstruction[rbp.GetBuilding()]
				if c.SelectedConstruction != nil {
					ConstructionToControlPanel(c.ControlPanel, c.SelectedConstruction)
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

func (c *Controller) MouseMoveCallback(wnd *glfw.Window, x float64, y float64) {
	c.X = x
	c.Y = y
}

func (c *Controller) MouseScrollCallback(wnd *glfw.Window, x float64, y float64) {
}

func Link(wnd *glfw.Window, Map *model.Map) *Controller {
	W, H := wnd.GetSize()
	Calendar := &time.CalendarType{
		Year:  1000,
		Month: 1,
		Day:   1,
		Hour:  0,
	}
	controlPanel := &ControlPanel{}
	C := &Controller{H: H, W: W, Calendar: Calendar, ControlPanel: controlPanel, Map: Map, Country: Map.Countries[0], TimeSpeed: 1, ActiveTown: Map.Countries[0].Towns[0]}
	controlPanel.Setup(C)
	wnd.SetKeyCallback(C.KeyboardCallback)
	wnd.SetMouseButtonCallback(C.MouseButtonCallback)
	wnd.SetCursorPosCallback(C.MouseMoveCallback)
	wnd.SetScrollCallback(C.MouseScrollCallback)
	return C
}

func (c *Controller) SwapRenderedObjects() {
	c.RenderedFields = c.TempRenderedFields
	c.RenderedBuildingParts = c.TempRenderedBuildingParts
	c.TempRenderedFields = []*renderer.RenderedField{}
	c.TempRenderedBuildingParts = []renderer.RenderedBuildingPart{}
}

func (c *Controller) AddRenderedField(rf *renderer.RenderedField) {
	c.TempRenderedFields = append(c.TempRenderedFields, rf)
}

func (c *Controller) AddRenderedBuildingPart(rbp renderer.RenderedBuildingPart) {
	c.TempRenderedBuildingParts = append(c.TempRenderedBuildingParts, rbp)
}
