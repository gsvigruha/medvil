package controller

import (
	"github.com/tfriedel6/canvas"
	"medvil/model/artifacts"
	"medvil/model/social"
	"medvil/view/gui"
	"strconv"
)

type MarketplaceController struct {
	mp          *gui.Panel
	marketplace *social.Marketplace
}

func MarketplaceToControlPanel(cp *ControlPanel, m *social.Marketplace) {
	mp := &gui.Panel{X: 0, Y: ControlPanelDynamicPanelTop, SX: ControlPanelSX, SY: HouseholdControllerSY}
	mc := &MarketplaceController{mp: mp, marketplace: m}
	MarketplaceToMarketPanel(mp, m)
	cp.SetDynamicPanel(mc)
}

func MarketplaceToMarketPanel(mp *gui.Panel, m *social.Marketplace) {
	MoneyToControlPanel(mp, m.Town, &m.Money, 100, 10, float64(IconH+50))
	var aI = 0
	for _, a := range artifacts.All {
		if q, ok := m.Storage.Artifacts[a]; ok {
			ArtifactsToMarketPanel(mp, aI, a, q, m.Prices[a])
			aI++
		}
	}
}

func ArtifactsToMarketPanel(mp *gui.Panel, i int, a *artifacts.Artifact, q uint16, p uint32) {
	rowH := IconH + int(IconS)
	xI := i % IconRowMax
	yI := i / IconRowMax
	mp.AddImageLabel("artifacts/"+a.Name, float64(10+xI*IconW), ArtifactsGUIY*ControlPanelSY+float64(yI*rowH), IconS, IconS, gui.ImageLabelStyleRegular)
	mp.AddTextLabel(strconv.Itoa(int(q)), float64(10+xI*IconW), ArtifactsGUIY*ControlPanelSY+float64(yI*rowH+IconH+4))
	mp.AddTextLabel("$"+strconv.Itoa(int(p)), float64(10+xI*IconW), ArtifactsGUIY*ControlPanelSY+float64(yI*rowH+IconH+4)+gui.FontSize)
}

func (mc *MarketplaceController) CaptureClick(x, y float64) {
	mc.mp.CaptureClick(x, y)
}

func (mc *MarketplaceController) Render(cv *canvas.Canvas) {
	mc.mp.Render(cv)
}

func (mc *MarketplaceController) Clear() {}

func (mc *MarketplaceController) Refresh() {
	mc.mp.Clear()
	MarketplaceToMarketPanel(mc.mp, mc.marketplace)
}
