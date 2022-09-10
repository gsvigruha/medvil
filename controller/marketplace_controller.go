package controller

import (
	"github.com/tfriedel6/canvas"
	"medvil/model/artifacts"
	"medvil/model/social"
	"medvil/view/gui"
	"strconv"
)

const MPIconH = 60

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
	MoneyToControlPanel(mp, m.Town, &m.Money, 100, 10, 80)
	var aI = 0
	for _, a := range artifacts.All {
		if q, ok := m.Storage.Artifacts[a]; ok {
			ArtifactsToMarketPanel(mp, aI, a, q, m.Prices[a])
			aI++
		}
	}
}

func ArtifactsToMarketPanel(mp *gui.Panel, i int, a *artifacts.Artifact, q uint16, p uint32) {
	xI := i % IconRowMax
	yI := i / IconRowMax
	mp.AddImageLabel("artifacts/"+a.Name, float64(10+xI*IconW), ArtifactsGUIY+float64(yI*MPIconH), 32, 32, gui.ImageLabelStyleRegular)
	mp.AddTextLabel(strconv.Itoa(int(q)), float64(10+xI*IconW), ArtifactsGUIY+float64(yI*MPIconH+IconH+4))
	mp.AddTextLabel("$"+strconv.Itoa(int(p)), float64(10+xI*IconW), ArtifactsGUIY+float64(yI*MPIconH+IconH+16))
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
