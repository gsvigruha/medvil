package controller

import (
	"fmt"
	"github.com/tfriedel6/canvas"
	"math/rand"
	"medvil/model/economy"
	"medvil/model/navigation"
	"medvil/model/social"
	"medvil/view/gui"
)

type FactoryController struct {
	householdPanel *gui.Panel
	factoryPanel   *gui.Panel
	factory        *social.Factory
	cp             *ControlPanel
}

func FactoryToControlPanel(cp *ControlPanel, factory *social.Factory) {
	hp := &gui.Panel{X: 0, Y: ControlPanelDynamicPanelTop, SX: ControlPanelSX, SY: HouseholdControllerSY}
	fp := &gui.Panel{X: 0, Y: ControlPanelDynamicPanelTop + HouseholdControllerSY, SX: ControlPanelSX, SY: ControlPanelDynamicPanelSY - HouseholdControllerSY}
	HouseholdToControlPanel(cp, hp, factory.Household, "factory")
	fc := &FactoryController{factoryPanel: fp, householdPanel: hp, factory: factory, cp: cp}

	hcy := HouseholdControllerGUIBottomY * ControlPanelSY
	for i, vc := range economy.GetVehicleConstructions(factory.Household.Building.Plan.GetExtensions()) {
		fp.AddPanel(CreateOrderPanelForFactory(cp, 10, float64(i*IconH)+hcy, factory, vc, cp.C.Map))
	}

	cp.SetDynamicPanel(fc)
}

func (fc *FactoryController) CaptureMove(x, y float64) {
	fc.householdPanel.CaptureMove(x, y)
	fc.factoryPanel.CaptureMove(x, y)
}

func (fc *FactoryController) CaptureClick(x, y float64) {
	fc.householdPanel.CaptureClick(x, y)
	fc.factoryPanel.CaptureClick(x, y)
}

func (fc *FactoryController) Render(cv *canvas.Canvas) {
	fc.householdPanel.Render(cv)
	fc.factoryPanel.Render(cv)
}

func (fc *FactoryController) Clear() {}

func (fc *FactoryController) Refresh() {
	fc.householdPanel.Clear()
	HouseholdToControlPanel(fc.cp, fc.householdPanel, fc.factory.Household, "factory")
	fc.CaptureMove(fc.cp.C.X, fc.cp.C.Y)
}

func (fc *FactoryController) GetHelperSuggestions() *gui.Suggestion {
	return nil
}

type OrderButton struct {
	factories []*social.Factory
	vc        *economy.VehicleConstruction
	icon      *gui.ImageLabel
	price     *gui.TextLabel
	orders    *gui.TextLabel
	m         navigation.IMap
}

func (b OrderButton) Click() {
	factory := b.factories[rand.Intn(len(b.factories))]
	h := factory.Household.Town.Townhall.Household
	order := factory.CreateOrder(b.vc, h)
	if order != nil {
		hx, hy, _ := social.GetRandomBuildingXY(h.Building, b.m, navigation.Field.BuildingNonExtension)
		fx, fy, _ := social.GetRandomBuildingXY(factory.Household.Building, b.m, navigation.Field.BuildingNonExtension)
		h.AddTask(&economy.FactoryPickupTask{
			PickupD:  b.m.GetField(fx, fy),
			DropoffD: b.m.GetField(hx, hy),
			Order:    order,
			TaskBase: economy.TaskBase{FieldCenter: true},
		})
	}
}

func (b OrderButton) NumOrders() int {
	var o = 0
	for _, factory := range b.factories {
		o += factory.NumOrders(b.vc)
	}
	return o
}

func (b OrderButton) Render(cv *canvas.Canvas) {
	b.price.Text = fmt.Sprintf("$%v", b.factories[0].Price(b.vc))
	b.orders.Text = fmt.Sprintf("%v", b.NumOrders())
}

func (b OrderButton) SetHoover(h bool) {}

func (b OrderButton) Contains(x float64, y float64) bool {
	return false
}

func (b OrderButton) Enabled() bool {
	return true
}

func CreateOrderPanelForFactory(cp *ControlPanel, x, y float64, factory *social.Factory, vc *economy.VehicleConstruction, m navigation.IMap) *gui.Panel {
	p := &gui.Panel{}
	orders := p.AddTextLabel("", 24+x+float64(IconW)*2, y)
	price := p.AddTextLabel("", 24+x+float64(IconW)*3, y)
	icon := p.AddImageLabel("vehicles/"+vc.Name, 24, y-float64(IconH)/2, IconS, IconS, gui.ImageLabelStyleRegular)
	p.AddButton(OrderButton{
		factories: []*social.Factory{factory},
		vc:        vc,
		icon:      icon,
		price:     price,
		orders:    orders,
		m:         m,
	})
	for i, as := range vc.Inputs {
		ArtifactsToControlPanel(cp, p, i+5, as.A, as.Quantity, y-float64(IconH)*2/3)
	}
	return p
}
