package controller

import (
	"github.com/tfriedel6/canvas"
	"github.com/tfriedel6/canvas/backend/goglbackend"
	"image/color"
	"medvil/model/artifacts"
	"medvil/model/social"
	"medvil/model/stats"
	"medvil/view/gui"
	"strconv"
)

const DPoint = 2

type ChartsLabel struct {
	cp        *ControlPanel
	s         *stats.History
	img       *canvas.Canvas
	state     uint8
	timeScale uint8
	helperMsg string
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
	cv.DrawImage(l.img, 24, ControlPanelSY*0.55, float64(l.img.Width()), float64(l.img.Height()))
}

func (l *ChartsLabel) Draw(cv *canvas.Canvas) {
	cv.ClearRect(0, 0, float64(l.img.Width()), float64(l.img.Height()))
	cv.SetLineWidth(2)
	switch l.state {
	case 1:
		l.drawChart(cv, "#DDD", 130, []string{"icon/gui/person"}, stats.HistoryElement.GetPeople, false)
		l.drawChart(cv, "#B00", 260, []string{"icon/gui/person"}, stats.HistoryElement.GetDeaths, true)
		l.drawChart(cv, "#808", 390, []string{"icon/gui/person"}, stats.HistoryElement.GetDepartures, true)
		l.helperMsg = "Population, deaths and departures"
	case 2:
		l.drawChart(cv, "#DDD", 130, []string{"icon/gui/barrel"}, stats.HistoryElement.GetArtifacts, false)
		l.drawChart(cv, "#FF0", 260, []string{"icon/gui/market", "icon/gui/barrel"}, stats.HistoryElement.GetExchangedNum, true)
		l.drawChart(cv, "#FF0", 390, []string{"icon/gui/market", "icon/gui/coin"}, stats.HistoryElement.GetExchangedPrice, true)
		l.helperMsg = "Products and market transactions"
	case 3:
		l.drawChart(cv, "#D82", 130, icons(social.FoodArtifacts), stats.HistoryElement.GetFoodPrice, false)
		l.drawChart(cv, "#660", 260, icons(social.BuildingMaterials), stats.HistoryElement.GetHouseholdItemPrices, false)
		l.drawChart(cv, "#D42", 390, icons(social.HouseholdItems), stats.HistoryElement.GetBuildingMaterialsPrice, false)
		l.helperMsg = "Average price of food, building materials"
	}
	l.CaptureClick(0, 0)
}

func (l *ChartsLabel) CaptureClick(x float64, y float64) {
	l.cp.HelperMessage(l.helperMsg)
}

func (l *ChartsLabel) drawChart(cv *canvas.Canvas, c string, y int, icons []string, fn func(stats.HistoryElement) uint32, sum bool) {
	maxPoints := (int(ControlPanelSX) - 48) / DPoint
	var startIdx = 0
	if len(l.s.Elements)/int(l.timeScale) > maxPoints {
		startIdx = len(l.s.Elements)/int(l.timeScale) - maxPoints
	}

	var max uint32 = 0
	var scaleCntr uint32 = 0
	var scaleAggr uint32 = 0
	for i := startIdx; i < len(l.s.Elements); i++ {
		he := l.s.Elements[i]
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

	cv.SetStrokeStyle(color.RGBA{R: 192, G: 192, B: 192, A: 128})
	for i := 0; i < maxPoints*DPoint/20+1; i++ {
		cv.BeginPath()
		cv.MoveTo(float64(i*20), float64(y))
		cv.LineTo(float64(i*20), float64(y-100))
		cv.ClosePath()
		cv.Stroke()
	}

	for i := 0; i <= 5; i++ {
		cv.BeginPath()
		cv.MoveTo(float64(0), float64(y-i*20))
		cv.LineTo(float64(int(ControlPanelSX)), float64(y-i*20))
		cv.ClosePath()
		cv.Stroke()
	}

	cv.SetStrokeStyle(c)
	cv.BeginPath()
	scaleCntr = 0
	scaleAggr = 0
	for i := startIdx; i < len(l.s.Elements); i++ {
		he := l.s.Elements[i]
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
			cv.LineTo(float64((i-startIdx)/int(l.timeScale)*DPoint), float64(y-int(val*100/max)))
		}
	}
	cv.MoveTo(0, 0)
	cv.ClosePath()
	cv.Stroke()

	for i, icon := range icons {
		cv.DrawImage(icon+".png", float64(i*32), float64(y-128), 32, 32)
	}

	cv.SetFillStyle(c)
	cv.SetFont("texture/font/Go-Regular.ttf", gui.FontSize)
	text := strconv.Itoa(int(max))
	cv.FillText(text, ControlPanelSX-60-float64(len(text))*gui.FontSize*0.5, float64(y-104))
}

func DrawStats(cp *ControlPanel, p *gui.Panel) {
	if cp.C.Map != nil && cp.C.Map.Countries[0] != nil && cp.C.Map.Countries[0].History != nil {
		offscreen, _ := goglbackend.NewOffscreen(int(ControlPanelSX)-48, 520, true, cp.C.ctx)
		cv := canvas.New(offscreen)
		cl := &ChartsLabel{cp: cp, s: cp.C.Map.Countries[0].History, img: cv, timeScale: 1}
		cl.Draw(cv)
		p.AddLabel(cl)

		p.AddButton(gui.SimpleButton{
			ButtonGUI: gui.ButtonGUI{Icon: "person", X: float64(24 + LargeIconD*0), Y: ControlPanelSY * 0.5, SX: LargeIconS, SY: LargeIconS},
			ClickImpl: func() { cl.state = 1 }})
		p.AddButton(gui.SimpleButton{
			ButtonGUI: gui.ButtonGUI{Icon: "tasks/exchange", X: float64(24 + LargeIconD*1), Y: ControlPanelSY * 0.5, SX: LargeIconS, SY: LargeIconS},
			ClickImpl: func() { cl.state = 2 }})
		p.AddButton(gui.SimpleButton{
			ButtonGUI: gui.ButtonGUI{Icon: "coin", X: float64(24 + LargeIconD*2), Y: ControlPanelSY * 0.5, SX: LargeIconS, SY: LargeIconS},
			ClickImpl: func() { cl.state = 3 }})

		p.AddButton(gui.SimpleButton{
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
	}
}
