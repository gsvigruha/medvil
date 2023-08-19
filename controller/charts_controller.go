package controller

import (
	"github.com/tfriedel6/canvas"
	"github.com/tfriedel6/canvas/backend/goglbackend"
	"image/color"
	"medvil/model/stats"
	"medvil/view/gui"
	"strconv"
)

const DPoint = 2

type ChartsLabel struct {
	cp  *ControlPanel
	s   *stats.History
	img *canvas.Canvas
	cnt int
}

func (l *ChartsLabel) Render(cv *canvas.Canvas) {
	if l.cnt >= 10 {
		l.Draw(l.img)
		l.cnt = 0
	}
	cv.DrawImage(l.img, 24, ControlPanelSY*0.5, float64(l.img.Width()), float64(l.img.Height()))
	l.cnt++
}

func (l *ChartsLabel) Draw(cv *canvas.Canvas) {
	cv.ClearRect(0, 0, float64(l.img.Width()), float64(l.img.Height()))
	cv.SetLineWidth(2)
	drawChart(cv, "#B00", 120, l.s, stats.HistoryElement.GetDeaths)
	drawChart(cv, "#808", 240, l.s, stats.HistoryElement.GetDepartures)
	drawChart(cv, "#DDD", 360, l.s, stats.HistoryElement.GetPeople)
	drawChart(cv, "#DDD", 480, l.s, stats.HistoryElement.GetArtifacts)
	drawChart(cv, "#FF0", 600, l.s, stats.HistoryElement.GetExchangedNum)
	drawChart(cv, "#FF0", 720, l.s, stats.HistoryElement.GetExchangedPrice)
}

func (l *ChartsLabel) CaptureClick(x float64, y float64) {

}

func drawChart(cv *canvas.Canvas, c string, y int, s *stats.History, fn func(stats.HistoryElement) uint32) {
	maxPoints := (int(ControlPanelSX) - 48) / DPoint
	var startIdx = 0
	if len(s.Elements) > maxPoints {
		startIdx = len(s.Elements) - maxPoints
	}

	var max uint32 = 0
	for i := startIdx; i < len(s.Elements); i++ {
		he := s.Elements[i]
		if max < fn(he) {
			max = fn(he)
		}
	}
	if max == 0 {
		max = 1
	}

	cv.SetStrokeStyle(color.RGBA{R: 192, G: 192, B: 192, A: 128})
	for i := 0; i < maxPoints*DPoint/20; i++ {
		cv.BeginPath()
		cv.MoveTo(float64(i*20), float64(y))
		cv.LineTo(float64(i*20), float64(y-100))
		cv.ClosePath()
		cv.Stroke()
	}

	for i := 0; i <= 5; i++ {
		cv.BeginPath()
		cv.MoveTo(float64(24), float64(y-i*20))
		cv.LineTo(float64(int(ControlPanelSX)-24), float64(y-i*20))
		cv.ClosePath()
		cv.Stroke()
	}

	cv.SetStrokeStyle(c)
	cv.BeginPath()
	for i := startIdx; i < len(s.Elements); i++ {
		he := s.Elements[i]
		cv.LineTo(float64((i-startIdx)*DPoint), float64(y-int(fn(he)*100/max)))
	}
	cv.MoveTo(0, 0)
	cv.ClosePath()
	cv.Stroke()

	cv.SetFillStyle(c)
	cv.SetFont("texture/font/Go-Regular.ttf", gui.FontSize)
	cv.FillText(strconv.Itoa(int(max)), ControlPanelSX-120, float64(y-80))
}

func DrawStats(cp *ControlPanel, p *gui.Panel) {
	if cp.C.Map != nil && cp.C.Map.Countries[0] != nil && cp.C.Map.Countries[0].History != nil {
		offscreen, _ := goglbackend.NewOffscreen(int(ControlPanelSX)-48, 720, true, cp.C.ctx)
		cv := canvas.New(offscreen)
		cl := &ChartsLabel{cp: cp, s: cp.C.Map.Countries[0].History, img: cv}
		cl.Draw(cv)
		p.AddLabel(cl)
	}
}
