package controller

import (
	"github.com/tfriedel6/canvas"
	"github.com/tfriedel6/canvas/backend/goglbackend"
	"image/color"
	"medvil/model/artifacts"
	"medvil/model/economy"
	"medvil/model/stats"
	"medvil/view/gui"
	"path/filepath"
	"strconv"
)

const DPoint = 2

type ChartsLabel struct {
	cp           *ControlPanel
	townSelector *gui.DropDown
	img          *canvas.Canvas
	state        uint8
	timeScale    uint8
	helperMsg    string
}

func icons(as []*artifacts.Artifact) []string {
	icons := make([]string, len(as))
	for i, a := range as {
		icons[i] = "icon/gui/artifacts/" + a.Name
	}
	return icons
}

func (l *ChartsLabel) Render(cv *canvas.Canvas) {
	l.Draw(l.img)
	cv.DrawImage(l.img, 24, ControlPanelSY*0.6, float64(l.img.Width()), float64(l.img.Height()))
}

type ElementLookup func(stats.HistoryElement) uint32

func (l *ChartsLabel) Draw(cv *canvas.Canvas) {
	ch := int(ControlPanelSY*0.3) / 4
	chf := float64(ControlPanelSY*0.3) / 4

	cv.ClearRect(0, 0, float64(l.img.Width()), float64(l.img.Height()))
	cv.SetFillStyle(filepath.FromSlash("texture/parchment.png"))
	cv.FillRect(0, 0, float64(l.img.Width()), float64(l.img.Height()))

	cv.SetLineWidth(2)
	switch l.state {
	case 1:
		l.drawChart(cv, "#22B", ch*1, []string{"icon/gui/person"}, stats.HistoryElement.GetPeople, false)
		l.drawChart(cv, "#22B", ch*2, []string{"icon/gui/death"}, stats.HistoryElement.GetDeaths, true)
		l.drawChart(cv, "#22B", ch*3, []string{"icon/gui/emigration"}, stats.HistoryElement.GetDepartures, true)
		l.drawChart(cv, "#22B", ch*4, []string{"icon/gui/coin"}, stats.HistoryElement.GetPoverty, false)
		l.helperMsg = "Population size, deaths, emigration and poverty"
	case 2:
		l.drawChart(cv, "#22B", ch*1, []string{"icon/gui/barrel"}, stats.HistoryElement.GetArtifacts, false)
		l.drawChart(cv, "#22B", ch*2, []string{"icon/gui/market", "icon/gui/barrel"}, stats.HistoryElement.GetExchangedQuantity, true)
		l.drawChart(cv, "#22B", ch*3, []string{"icon/gui/market", "icon/gui/coin"}, stats.HistoryElement.GetExchangedPrice, true)
		l.helperMsg = "Products and market transactions"
	case 3:
		l.drawChart(cv, "#22B", ch*1, icons(economy.FoodArtifacts), stats.HistoryElement.GetFoodPrice, false)
		l.drawChart(cv, "#22B", ch*2, icons(economy.HouseholdItems), stats.HistoryElement.GetHouseholdItemPrices, false)
		l.drawChart(cv, "#22B", ch*3, icons(economy.BuildingMaterials), stats.HistoryElement.GetBuildingMaterialsPrice, false)
		l.helperMsg = "Average price of food, building materials"
	case 4:
		l.drawChart(cv, "#22B", ch*1, []string{"icon/gui/tasks/transport"}, stats.HistoryElement.GetTransportTaskTime, true)
		l.drawChart(cv, "#22B", ch*2, []string{"icon/gui/tasks/exchange"}, stats.HistoryElement.GetExchangeTaskTime, true)
		l.drawChart(cv, "#22B", ch*3, []string{"icon/gui/tasks/ploughing"}, stats.HistoryElement.GetAgricultureTaskTime, true)
		l.drawChart(cv, "#22B", ch*4, []string{"icon/gui/tasks/milling"}, stats.HistoryElement.GetManufactureTaskTime, true)
		l.helperMsg = "Days spent on various tasks"
	case 5:
		l.drawCharts(cv, []string{"#872", "#96D", "#F11", "#D72", "#58F"}, ch*1, []string{"icon/gui/coin"},
			[]ElementLookup{
				stats.HistoryElement.GetFarmMoney,
				stats.HistoryElement.GetWorkshopMoney,
				stats.HistoryElement.GetMineMoney,
				stats.HistoryElement.GetTraderMoney,
				stats.HistoryElement.GetGovernmentMoney,
			}, false)
		l.drawCharts(cv, []string{"#872", "#96D", "#F11", "#D72", "#58F"}, ch*2, []string{"icon/gui/person"},
			[]ElementLookup{
				stats.HistoryElement.GetFarmPeople,
				stats.HistoryElement.GetWorkshopPeople,
				stats.HistoryElement.GetMinePeople,
				stats.HistoryElement.GetTraderPeople,
				stats.HistoryElement.GetGovernmentPeople,
			}, false)

		iconTop := chf*2 + 8
		w2 := ControlPanelSX / 2
		cv.DrawImage(filepath.FromSlash("icon/gui/farm.png"), 8, iconTop+IconS*0, IconS, IconS)
		l.drawLegend(cv, 8, iconTop+IconS*0, "#872")
		cv.DrawImage(filepath.FromSlash("icon/gui/workshop.png"), 8, iconTop+IconS*1, IconS, IconS)
		l.drawLegend(cv, 8, iconTop+IconS*1, "#96D")
		cv.DrawImage(filepath.FromSlash("icon/gui/mine.png"), 8, iconTop+IconS*2, IconS, IconS)
		l.drawLegend(cv, 8, iconTop+IconS*2, "#F11")
		cv.DrawImage(filepath.FromSlash("icon/gui/trader.png"), w2, iconTop+IconS*0, IconS, IconS)
		l.drawLegend(cv, w2, iconTop+IconS*0, "#D72")
		cv.DrawImage(filepath.FromSlash("icon/gui/town.png"), w2, iconTop+IconS*1, IconS, IconS)
		l.drawLegend(cv, w2, iconTop+IconS*1, "#58F")
		l.helperMsg = "Wealth and population of social classes"
	}
	l.CaptureClick(0, 0)
}

func (l *ChartsLabel) drawLegend(cv *canvas.Canvas, x, y float64, c string) {
	cv.SetLineWidth(3.0)
	cv.SetStrokeStyle(c)
	cv.BeginPath()
	cv.MoveTo(x+float64(IconW), y+float64(IconH)/2)
	cv.LineTo(x+float64(IconW)+IconS, y+float64(IconH)/2)
	cv.ClosePath()
	cv.Stroke()
}

func (l *ChartsLabel) CaptureClick(x float64, y float64) {
	l.cp.HelperMessage(l.helperMsg)
}

func (l *ChartsLabel) drawChart(cv *canvas.Canvas, c string, y int, icons []string, fn ElementLookup, sum bool) {
	l.drawCharts(cv, []string{c}, y, icons, []ElementLookup{fn}, sum)
}

func (l *ChartsLabel) drawCharts(cv *canvas.Canvas, cs []string, y int, icons []string, fns []ElementLookup, sum bool) {
	var s *stats.History
	if l.townSelector.Selected == 0 {
		s = l.cp.C.Map.Countries[0].History
	} else {
		s = l.cp.C.Map.Countries[0].Towns[l.townSelector.Selected-1].History
	}

	maxPoints := (int(ControlPanelSX) - 48) / DPoint
	numPoints := len(s.Elements) / int(l.timeScale)
	var dPoint = float64(DPoint)
	var startIdx = 0
	if numPoints > maxPoints {
		startIdx = (numPoints - maxPoints) * int(l.timeScale)
	} else if numPoints < maxPoints {
		dPoint = DPoint * float64(maxPoints) / float64(numPoints)
	}

	var max uint32 = 0
	for _, fn := range fns {
		var scaleCntr uint32 = 0
		var scaleAggr uint32 = 0
		for i := startIdx; i < len(s.Elements); i++ {
			he := s.Elements[i]
			scaleAggr += fn(he)
			scaleCntr++
			if scaleCntr == uint32(l.timeScale) {
				var val uint32
				if sum {
					val = scaleAggr
				} else {
					val = scaleAggr / scaleCntr
				}
				if val > max {
					max = val
				}
				scaleCntr = 0
				scaleAggr = 0
			}
		}
		if max == 0 {
			max = 1
		}
	}

	uh := float64(ControlPanelSY*0.3) / 4 * 0.25
	lh := float64(ControlPanelSY*0.3) / 4 * 0.75

	cv.SetStrokeStyle(color.RGBA{R: 128, G: 64, B: 255, A: 64})
	for i := 0; i < int(float64(maxPoints)*dPoint/20)+1; i++ {
		cv.BeginPath()
		cv.MoveTo(float64(i*20), float64(y))
		cv.LineTo(float64(i*20), float64(y)-lh)
		cv.ClosePath()
		cv.Stroke()
	}

	for i := 0; i <= 5; i++ {
		dh := lh / 5
		cv.BeginPath()
		cv.MoveTo(float64(0), float64(y)-float64(i)*dh)
		cv.LineTo(float64(int(ControlPanelSX)), float64(y)-float64(i)*dh)
		cv.ClosePath()
		cv.Stroke()
	}

	for i, fn := range fns {
		c := cs[i]
		cv.SetStrokeStyle(c)
		cv.BeginPath()
		var scaleCntr uint32 = 0
		var scaleAggr uint32 = 0
		for i := startIdx; i < len(s.Elements); i++ {
			he := s.Elements[i]
			scaleAggr += fn(he)
			scaleCntr++
			if scaleCntr == uint32(l.timeScale) {
				var val uint32
				if sum {
					val = scaleAggr
				} else {
					val = scaleAggr / scaleCntr
				}
				scaleCntr = 0
				scaleAggr = 0
				cv.LineTo(float64(i-startIdx)*dPoint/float64(l.timeScale), float64(y)-float64(val)*lh/float64(max))
			}
		}
		cv.MoveTo(0, 0)
		cv.ClosePath()
		cv.Stroke()
	}

	for i, icon := range icons {
		cv.DrawImage(filepath.FromSlash(icon+".png"), float64(i)*IconS+8, float64(y)-lh-uh, IconS, IconS)
	}

	cv.SetFillStyle("#22B")
	cv.SetFont(gui.Font, gui.FontSize)
	text := strconv.Itoa(int(max))
	cv.FillText(text, ControlPanelSX-60-float64(len(text))*gui.FontSize*0.5, float64(y)-lh-gui.FontSize*0.25)
}

func DrawStats(cp *ControlPanel, p *gui.Panel) {
	if cp.C.Map != nil && cp.C.Map.Countries[0] != nil && cp.C.Map.Countries[0].History != nil {
		ch := int(ControlPanelSY * 0.3)
		offscreen, _ := goglbackend.NewOffscreen(int(ControlPanelSX)-48, ch, true, cp.C.ctx)
		cv := canvas.New(offscreen)
		cl := &ChartsLabel{cp: cp, img: cv, timeScale: 1}
		cl.Draw(cv)
		p.AddLabel(cl)

		p.AddButton(&gui.SimpleButton{
			ButtonGUI: gui.ButtonGUI{Icon: "person", X: float64(24 + LargeIconD*0), Y: ControlPanelSY * 0.5, SX: LargeIconS, SY: LargeIconS},
			ClickImpl: func() { cl.state = 1 }})
		p.AddButton(&gui.SimpleButton{
			ButtonGUI: gui.ButtonGUI{Icon: "tasks/exchange", X: float64(24 + LargeIconD*1), Y: ControlPanelSY * 0.5, SX: LargeIconS, SY: LargeIconS},
			ClickImpl: func() { cl.state = 2 }})
		p.AddButton(&gui.SimpleButton{
			ButtonGUI: gui.ButtonGUI{Icon: "coin", X: float64(24 + LargeIconD*2), Y: ControlPanelSY * 0.5, SX: LargeIconS, SY: LargeIconS},
			ClickImpl: func() { cl.state = 3 }})
		p.AddButton(&gui.SimpleButton{
			ButtonGUI: gui.ButtonGUI{Icon: "tasks/transport", X: float64(24 + LargeIconD*3), Y: ControlPanelSY * 0.5, SX: LargeIconS, SY: LargeIconS},
			ClickImpl: func() { cl.state = 4 }})
		p.AddButton(&gui.SimpleButton{
			ButtonGUI: gui.ButtonGUI{Icon: "classes", X: float64(24 + LargeIconD*4), Y: ControlPanelSY * 0.5, SX: LargeIconS, SY: LargeIconS},
			ClickImpl: func() { cl.state = 5 }})

		p.AddButton(&gui.SimpleButton{
			ButtonGUI: gui.ButtonGUI{Icon: "time", X: float64(24 + LargeIconD*6), Y: ControlPanelSY * 0.5, SX: LargeIconS, SY: LargeIconS},
			ClickImpl: func() {
				switch cl.timeScale {
				case 1:
					cl.timeScale = 3
				case 3:
					cl.timeScale = 12
				case 12:
					cl.timeScale = 1
				}
			}})
		p.AddDynamicTextLabel(func() string {
			switch cl.timeScale {
			case 1:
				return "M"
			case 3:
				return "Q"
			case 12:
				return "Y"
			}
			return ""
		}, float64(24+LargeIconD*6.5), ControlPanelSY*0.5+LargeIconD*0.7)

		var names []string = []string{"Country"}
		var icons []string = []string{"town"}
		for _, town := range cp.C.Map.Countries[0].Towns {
			names = append(names, town.Name)
			icons = append(icons, "town")
		}
		cl.townSelector = &gui.DropDown{
			X:        float64(24),
			Y:        ControlPanelSY*0.5 + LargeIconD + float64(IconH/8),
			SX:       IconS + gui.FontSize*12,
			SY:       IconS,
			Options:  names,
			Icons:    icons,
			Selected: 0,
		}
		p.AddDropDown(cl.townSelector)
	}
}
