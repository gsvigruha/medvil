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
}

func FactoryToControlPanel(cp *ControlPanel, factory *social.Factory) {
	hp := &gui.Panel{X: 0, Y: ControlPanelDynamicPanelTop, SX: ControlPanelSX, SY: HouseholdControllerSY}
	fp := &gui.Panel{X: 0, Y: ControlPanelDynamicPanelTop + HouseholdControllerSY, SX: ControlPanelSX, SY: ControlPanelDynamicPanelSY - HouseholdControllerSY}
	HouseholdToControlPanel(hp, &factory.Household)
	fc := &FactoryController{factoryPanel: fp, householdPanel: hp, factory: factory}

	hcy := HouseholdControllerGUIBottomY * ControlPanelSY
	for i, vc := range economy.GetVehicleConstructions(factory.Household.Building.Plan.GetExtension()) {
		fp.AddPanel(CreateOrderPanelForFactory(10, float64(i*40)+hcy, 120, 20, factory, vc, cp.C.Map))
	}

	cp.SetDynamicPanel(fc)
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
	HouseholdToControlPanel(fc.householdPanel, &fc.factory.Household)
}

type OrderButton struct {
	b         gui.ButtonGUI
	factories []*social.Factory
	vc        *economy.VehicleConstruction
	l         *gui.TextLabel
	m         navigation.IMap
}

func (b OrderButton) Click() {
	factory := b.factories[rand.Intn(len(b.factories))]
	h := &factory.Household.Town.Townhall.Household
	order := factory.CreateOrder(b.vc, h)
	hx, hy, _ := social.GetRandomBuildingXY(h.Building, b.m, navigation.Field.BuildingNonExtension)
	fx, fy, _ := social.GetRandomBuildingXY(factory.Household.Building, b.m, navigation.Field.BuildingNonExtension)
	h.AddTask(&economy.FactoryPickupTask{
		PickupF:  b.m.GetField(fx, fy),
		DropoffF: b.m.GetField(hx, hy),
		Order:    order,
		TaskBase: economy.TaskBase{FieldCenter: true},
	})
}

func (b OrderButton) NumOrders() int {
	var o = 0
	for _, factory := range b.factories {
		o += factory.NumOrders(b.vc)
	}
	return o
}

func (b OrderButton) Render(cv *canvas.Canvas) {
	b.b.Render(cv)
	b.l.Text = fmt.Sprintf("%v %v", b.vc.Name, b.NumOrders())
}

func (b OrderButton) Contains(x float64, y float64) bool {
	return b.b.Contains(x, y)
}

func (b OrderButton) Enabled() bool {
	return b.b.Enabled()
}

func CreateOrderPanelForFactory(x, y, sx, sy float64, factory *social.Factory, vc *economy.VehicleConstruction, m navigation.IMap) *gui.Panel {
	p := &gui.Panel{}
	l := p.AddTextLabel("", x, y+sy*2/3)
	p.AddButton(OrderButton{
		b:         gui.ButtonGUI{Icon: "plus", X: x + sx, Y: y, SX: sy, SY: sy},
		factories: []*social.Factory{factory},
		vc:        vc,
		l:         l,
		m:         m,
	})
	p.AddTextLabel(fmt.Sprintf("$%v", factory.Price(vc)), x+sx+sy*2, y+sy*2/3)
	return p
}
