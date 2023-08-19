package controller

import (
	"github.com/tfriedel6/canvas"
	"github.com/tfriedel6/canvas/backend/goglbackend"
	"medvil/model/stats"
	"medvil/view/gui"
)

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
	cv.SetLineWidth(2)
	drawChart(cv, "#B00", 100, l.s, stats.HistoryElement.GetDeaths)
	drawChart(cv, "#808", 200, l.s, stats.HistoryElement.GetDepartures)
	drawChart(cv, "#DDD", 300, l.s, stats.HistoryElement.GetPeople)
	drawChart(cv, "#DDD", 400, l.s, stats.HistoryElement.GetArtifacts)
	drawChart(cv, "#FF0", 500, l.s, stats.HistoryElement.GetExchangedNum)
	drawChart(cv, "#FF0", 600, l.s, stats.HistoryElement.GetExchangedPrice)
}

func (l *ChartsLabel) CaptureClick(x float64, y float64) {

}

func drawChart(cv *canvas.Canvas, color string, y int, s *stats.History, fn func(stats.HistoryElement) uint32) {
	var max uint32 = 0
	for _, he := range s.Elements {
		if max < fn(he) {
			max = fn(he)
		}
	}
	if max == 0 {
		max = 1
	}
	cv.SetStrokeStyle(color)
	cv.BeginPath()
	for i, he := range s.Elements {
		cv.LineTo(float64(i*2), float64(y-int(fn(he)*95/max)))
	}
	cv.MoveTo(0, 0)
	cv.ClosePath()
	cv.Stroke()
}

func DrawStats(cp *ControlPanel, p *gui.Panel) {
	if cp.C.Map != nil && cp.C.Map.Countries[0] != nil && cp.C.Map.Countries[0].History != nil {
		offscreen, _ := goglbackend.NewOffscreen(int(ControlPanelSX)-48, 600, true, cp.C.ctx)
		cv := canvas.New(offscreen)
		cl := &ChartsLabel{cp: cp, s: cp.C.Map.Countries[0].History, img: cv}
		cl.Draw(cv)
		p.AddLabel(cl)
	}
}
