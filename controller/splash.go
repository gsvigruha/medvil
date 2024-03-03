package controller

import (
	"github.com/tfriedel6/canvas"
	"medvil/view/gui"
)

const SplashPageGameplay = 0
const SplashPageControls = 1
const SplashPageEconomics = 2

type Splash struct {
	p    *gui.Panel
	gp   *gui.Panel
	cp   *gui.Panel
	ep   *gui.Panel
	page int
}

func (s *Splash) CaptureClick(x float64, y float64) {
	s.p.CaptureClick(x, y)
}

func (s *Splash) CaptureMove(x float64, y float64) {
	s.p.CaptureMove(x, y)
}

func (s *Splash) Setup(w, h int) {
	px := float64(w) * 0.3
	pw := float64(w) * 0.4
	py := float64(h) * 0.2
	ph := float64(h) * 0.6
	p := &gui.Panel{X: px, Y: py, SX: pw, SY: ph}
	p.AddButton(&gui.SimpleButton{
		ButtonGUI: gui.ButtonGUI{Icon: "house", X: px + LargeIconS*0 + 24, Y: py + 24, SX: IconS, SY: IconS},
		ClickImpl: func() { s.page = SplashPageGameplay },
	})
	p.AddButton(&gui.SimpleButton{
		ButtonGUI: gui.ButtonGUI{Icon: "settings", X: px + LargeIconS*1 + 24, Y: py + 24, SX: IconS, SY: IconS},
		ClickImpl: func() { s.page = SplashPageControls },
	})
	p.AddButton(&gui.SimpleButton{
		ButtonGUI: gui.ButtonGUI{Icon: "coin", X: px + LargeIconS*2 + 24, Y: py + 24, SX: IconS, SY: IconS},
		ClickImpl: func() { s.page = SplashPageEconomics },
	})

	spy := py + LargeIconD + 24
	sph := ph - LargeIconD - 24

	{
		cp := &gui.Panel{X: px, Y: spy, SX: pw, SY: sph}
		cp.AddLargeTextLabel("Controls", px+24, spy+24)
		lines := []string{
			"Keyboard",
			"",
			"- Move the camera with arrows or the A/S/D/W keys",
			"- Rotate with Enter or Q/E",
			"- Escape quits the game after saving to latest_autosave",
			"- Tab changes the game speed",
			"- Pause the game with P",
			"",
			"Mouse",
			"- You can click on any building to configure them.",
			"- Any tile or person to view their status - but you can not directly control the citizens.",
		}
		AddLines(cp, lines, px+24, spy+24+LargeIconD)
		s.cp = cp
	}

	{
		gp := &gui.Panel{X: px, Y: spy, SX: pw, SY: sph}
		gp.AddLargeTextLabel("Gameplay", px+24, spy+24)
		lines := []string{
			"The purpose of the game is to expand your city while developing an increasingly complex economy.",
			"You can create buildings, add people to them and configure their function.",
			"You cannot control people directly. Your population automatically expands until your buildings are full.",
			"",
			"First, build farms and configure them to be self sustaining by growing trees and vegetables.",
			"",
			"Second, build workshops and configure them to be butchers, mills, bakers or sew clothes.",
			"At this stage your town can sustain a small population.",
			"",
			"Third, mine rocks and clay so you can produce building materials with more complex workshops.",
			"This will enable you expanding your town.",
			"",
			"In the later stages of development, you can start mining silver to make vehicles,",
			"or gold to increase your money supply. You can also establish new towns and establish",
			"trade routes between them.",
		}
		AddLines(gp, lines, px+24, spy+24+LargeIconD)
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
	} else if s.page == SplashPageControls {
		s.cp.Render(cv)
	}
}
