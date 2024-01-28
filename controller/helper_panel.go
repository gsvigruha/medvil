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
	y := float64(IconH) / 4.0
	p.AddTextLabel(task.Description(), 24, y+IconS+gui.FontSize*1.2)
	p.AddImageLabel("tasks/"+economy.IconName(task), 24, y, IconS, IconS, style)
	switch v := task.(type) {
	case *economy.BuyTask:
		ArtifactsToHelperPanel(p, v.Goods, 0, y)
		x := float64(24 + IconW + len(v.Goods)*IconW)
		p.AddImageLabel("coin", x, y, IconS, IconS, gui.ImageLabelStyleRegular)
		p.AddTextLabel(strconv.Itoa(int(v.MaxPrice)), x+IconS*0.75, y+IconS)
	case *economy.SellTask:
		ArtifactsToHelperPanel(p, v.Goods, 0, y)
	case *economy.ExchangeTask:
		p.AddImageLabel("tasks/buy", 24, y, IconS, IconS, style)
		ArtifactsToHelperPanel(p, v.GoodsToBuy, 0, y)
		sellSX := float64((len(v.GoodsToBuy)+1)*IconW + 24)
		p.AddImageLabel("tasks/sell", sellSX, y, IconS, IconS, style)
		ArtifactsToHelperPanel(p, v.GoodsToSell, sellSX, y)
		sellSX = sellSX + float64((len(v.GoodsToSell)+1)*IconW+24)
		if v.Vehicle != nil {
			p.AddImageLabel("vehicles/"+v.Vehicle.T.Name, sellSX, y, IconS, IconS, style)
		}
	case *economy.TransportTask:
		x := float64(24 + IconW)
		p.AddImageLabel("artifacts/"+v.A.Name, x, y, IconS, IconS, gui.ImageLabelStyleRegular)
		p.AddTextLabel(strconv.Itoa(int(v.ActualQuantity))+"/"+strconv.Itoa(int(v.ActualQuantity+v.TargetQuantity)), x+IconS*0.75, y+IconS)
	case *economy.RepairTask:
		ArtifactsToHelperPanel(p, v.Repairable.RepairCost(), 0, y)
	case *economy.FactoryPickupTask:
		p.AddImageLabel("vehicles/"+v.Order.Name(), 24+float64(IconW), y, IconS, IconS, style)
	case *economy.TradeTask:
		ArtifactsToHelperPanel(p, v.Goods, 0, y)
	case *economy.MiningTask:
		p.AddImageLabel(MineUseTypeIcon(v.UseType), 24+float64(IconW), y, IconS, IconS, style)
	}
}

func ArtifactsToHelperPanel(p *gui.Panel, as []artifacts.Artifacts, sx, y float64) {
	for i, as := range as {
		x := float64(24+IconW+i*IconW) + sx
		p.AddImageLabel("artifacts/"+as.A.Name, x, y, IconS, IconS, gui.ImageLabelStyleRegular)
		p.AddTextLabel(strconv.Itoa(int(as.Quantity)), x+IconS*0.75, y+IconS)
	}
}

func ArtifactToHelperPanel(p *gui.Panel, a *artifacts.Artifact) {
	y := float64(IconH) / 2.0
	p.AddImageLabel("artifacts/"+a.Name, 24, y, IconS, IconS, gui.ImageLabelStyleRegular)
	p.AddTextLabel(a.Description, 24+float64(IconW), y+IconS-gui.FontSize/2.0)
}
