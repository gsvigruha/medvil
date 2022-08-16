package controller

import (
	"fmt"
	"github.com/tfriedel6/canvas"
	"math/rand"
	"medvil/model/economy"
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

	for i, vc := range economy.GetVehicleConstructions(factory.Household.Building.Plan.GetExtension()) {
		fp.AddPanel(CreateOrderPanelForFactory(10, float64(i*40+600), 60, 20, factory, vc))
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
	factory   *social.Factory
	factories []*social.Factory
	vc        *economy.VehicleConstruction
	l         *gui.TextLabel
}

func (b OrderButton) Click() {
	if b.factory != nil {
		b.factory.CreateOrder(b.vc, &b.factory.Household.Town.Townhall.Household)
	} else {
		factory := b.factories[rand.Intn(len(b.factories))]
		factory.CreateOrder(b.vc, &factory.Household.Town.Townhall.Household)
	}
}

func (b OrderButton) NumOrders() int {
	if b.factory != nil {
		return b.factory.NumOrders(b.vc)
	} else {
		var o = 0
		for _, factory := range b.factories {
			o += factory.NumOrders(b.vc)
		}
		return o
	}
}

func (b OrderButton) Render(cv *canvas.Canvas) {
	b.b.Render(cv)
	b.l.Text = fmt.Sprintf("%v %v", b.vc.Name, b.NumOrders())
}

func (b OrderButton) Contains(x float64, y float64) bool {
	return b.b.Contains(x, y)
}

func CreateOrderPanelForFactory(x, y, sx, sy float64, factory *social.Factory, vc *economy.VehicleConstruction) *gui.Panel {
	p := &gui.Panel{}
	l := p.AddTextLabel("", x, y+sy*2/3)
	p.AddButton(OrderButton{
		b:       gui.ButtonGUI{Icon: "plus", X: x + sx, Y: y, SX: sy, SY: sy},
		factory: factory,
		vc:      vc,
		l:       l,
	})
	p.AddTextLabel(fmt.Sprintf("$%v", factory.Price(vc)), x+sx+sy*2, y+sy*2/3)
	return p
}
