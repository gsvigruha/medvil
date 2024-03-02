package controller

import (
	"github.com/tfriedel6/canvas"
	"medvil/view/gui"
)

const SplashPageGameplay = 0
const SplashPageEconomics = 1

type Splash struct {
	p    *gui.Panel
	gp   *gui.Panel
	ep   *gui.Panel
	page int
}

func (s *Splash) Setup(w, h int) {
	px := float64(w) * 0.3
	pw := float64(w) * 0.4
	py := float64(h) * 0.3
	ph := float64(h) * 0.4
	p := &gui.Panel{X: px, Y: py, SX: pw, SY: ph}
	p.AddButton(&gui.SimpleButton{
		ButtonGUI: gui.ButtonGUI{Icon: "cancel", X: px, Y: py, SX: IconS, SY: IconS},
		ClickImpl: func() { s.page = SplashPageGameplay },
	})
	p.AddButton(&gui.SimpleButton{
		ButtonGUI: gui.ButtonGUI{Icon: "cancel", X: px + LargeIconS, Y: py, SX: IconS, SY: IconS},
		ClickImpl: func() { s.page = SplashPageEconomics },
	})

	spy := py + LargeIconD
	sph := ph - LargeIconD

	{
		gp := &gui.Panel{X: px, Y: spy, SX: pw, SY: sph}
		gp.AddLargeTextLabel("Gameplay", px+24, spy+24)
		s.gp = gp
	}

	{
		ep := &gui.Panel{X: px, Y: spy, SX: pw, SY: sph}
		ep.AddLargeTextLabel("Economics", px+24, spy+24)
		s.ep = ep
	}

	s.p = p
}

func (s *Splash) Render(cv *canvas.Canvas) {
	s.p.Render(cv)
	if s.page == SplashPageGameplay {
		s.gp.Render(cv)
	} else if s.page == SplashPageEconomics {
		s.ep.Render(cv)
	}
}
