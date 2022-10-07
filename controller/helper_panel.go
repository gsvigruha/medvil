package controller

import (
	"medvil/model/artifacts"
	"medvil/model/economy"
	"medvil/view/gui"
	"strconv"
)

func TaskToHelperPanel(p *gui.Panel, task economy.Task) {
	var style uint8 = gui.ImageLabelStyleHighlight
	if task.Blocked() {
		style = gui.ImageLabelStyleDisabled
	}
	y := ControlPanelSY * 0.95
	p.AddImageLabel("tasks/"+task.Name(), 10, y, IconS, IconS, style)
	switch v := task.(type) {
	case *economy.BuyTask:
		ArtifactsToHelperPanel(p, v.Goods, 0)
		x := float64(10 + IconW + len(v.Goods)*IconW)
		p.AddImageLabel("coin", x, y, IconS, IconS, gui.ImageLabelStyleRegular)
		p.AddTextLabel(strconv.Itoa(int(v.MaxPrice)), x+IconS*0.75, y+IconS)
	case *economy.SellTask:
		ArtifactsToHelperPanel(p, v.Goods, 0)
	case *economy.ExchangeTask:
		ArtifactsToHelperPanel(p, v.GoodsToBuy, 0)
		ArtifactsToHelperPanel(p, v.GoodsToSell, float64(len(v.GoodsToBuy)*IconW+10))
	case *economy.TransportTask:
		x := float64(10 + IconW)
		p.AddImageLabel("artifacts/"+v.A.Name, x, y, IconS, IconS, gui.ImageLabelStyleRegular)
		p.AddTextLabel(strconv.Itoa(int(v.Quantity)), x+IconS*0.75, y+IconS)
	}
}

func ArtifactsToHelperPanel(p *gui.Panel, as []artifacts.Artifacts, sx float64) {
	y := ControlPanelSY * 0.95
	for i, as := range as {
		x := float64(10+IconW+i*IconW) + sx
		p.AddImageLabel("artifacts/"+as.A.Name, x, y, IconS, IconS, gui.ImageLabelStyleRegular)
		p.AddTextLabel(strconv.Itoa(int(as.Quantity)), x+IconS*0.75, y+IconS)
	}
}
