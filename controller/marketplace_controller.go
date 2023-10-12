package controller

import (
	"github.com/tfriedel6/canvas"
	"medvil/model/artifacts"
	"medvil/model/social"
	"medvil/view/gui"
	"strconv"
)

const IconRowMaxButtons = 7

var MarketplaceGUIY = 0.175

type MarketplaceController struct {
	mp          *gui.Panel
	marketplace *social.Marketplace
	cp          *ControlPanel
}

func MarketplaceToControlPanel(cp *ControlPanel, m *social.Marketplace) {
	mp := &gui.Panel{X: 0, Y: ControlPanelDynamicPanelTop, SX: ControlPanelSX, SY: HouseholdControllerSY}
	mc := &MarketplaceController{mp: mp, marketplace: m, cp: cp}
	MarketplaceToMarketPanel(mp, m)
	cp.SetDynamicPanel(mc)
}

func MarketplaceToMarketPanel(mp *gui.Panel, m *social.Marketplace) {
	MoneyToControlPanel(mp, m.Town.Townhall.Household, m, 100, 24, LargeIconD*2+float64(IconH)+24)
	mp.AddTextLabel("marketplace / "+m.Town.Name, 200, LargeIconD*2+float64(IconH)+24)
	var aI = 0
	for _, a := range artifacts.All {
		if q, ok := m.Storage.Artifacts[a]; ok {
			ArtifactsToMarketPanel(mp, aI, a, q, m)
			aI++
		}
	}
}

func ArtifactsToMarketPanel(mp *gui.Panel, i int, a *artifacts.Artifact, q uint16, m *social.Marketplace) {
	rowH := int(IconS * 2)
	xI := i % IconRowMaxButtons
	yI := i / IconRowMaxButtons
	w := int(float64(IconW) * float64(IconRowMax) / float64(IconRowMaxButtons))
	mp.AddImageLabel("artifacts/"+a.Name, float64(24+xI*w), MarketplaceGUIY*ControlPanelSY+float64(yI*rowH), IconS, IconS, gui.ImageLabelStyleRegular)
	mp.AddTextLabel(strconv.Itoa(int(q)), float64(24+xI*w), MarketplaceGUIY*ControlPanelSY+float64(yI*rowH+IconH+4))
	mp.AddPanel(gui.CreateNumberPanel(float64(24+xI*w), MarketplaceGUIY*ControlPanelSY+float64(yI*rowH+IconH+4), float64(IconW+8), gui.FontSize*1.5, 0, 1000, 10, "$%v", false,
		func() int { return int(m.Prices[a]) },
		func(v int) {
			if uint32(v) > m.Prices[a] {
				m.IncPrice(a)
			} else if uint32(v) < m.Prices[a] {
				m.DecPrice(a)
			}
		}).P)
}

func (mc *MarketplaceController) CaptureMove(x, y float64) {
	mc.mp.CaptureMove(x, y)
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
	mc.CaptureMove(mc.cp.C.X, mc.cp.C.Y)
}

func (mc *MarketplaceController) GetHelperSuggestions() *gui.Suggestion {
	return nil
}
