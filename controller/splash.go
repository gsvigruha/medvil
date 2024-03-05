package controller

import (
	"github.com/tfriedel6/canvas"
	"medvil/model/artifacts"
	"medvil/model/building"
	"medvil/model/economy"
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

func (s *Splash) Setup(cp *ControlPanel, w, h int) {
	px := float64(w) * 0.3
	pw := float64(w) * 0.4
	py := float64(h) * 0.15
	ph := float64(h) * 0.7
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
			"You cannot control people directly, but you can observe how they feel and what are they doing.",
			"Your population expands automatically until your buildings are full.",
			"",
			"At first, build farms and configure them to be self sustaining by growing trees and vegetables.",
			"",
			"Second, build workshops and configure them to be butchers, mills, bakers or sew clothes.",
			"At this stage your town can sustain a small population.",
			"",
			"Third, mine rocks and clay so you can produce building materials with more complex workshops.",
			"This will enable you expanding your town and your villagers can repair their buildings.",
			"",
			"In the later stages of development, you can start mining silver to make vehicles.",
			"Mining gold and minting coins increases your money supply, which is great when you are expanding.",
			"You can also establish new towns and establish trade routes between them.",
		}
		AddLines(gp, lines, px+24, spy+24+LargeIconD)
		s.gp = gp
	}

	{
		ep := &gui.Panel{X: px, Y: spy, SX: pw, SY: sph}
		ep.AddLargeTextLabel("Economics", px+24, spy+24)
		s.addEconomicsLabels(cp, ep)
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

func (s *Splash) addEconomicsLabels(cp *ControlPanel, p *gui.Panel) {
	var y = p.Y + 24 + LargeIconD
	p.AddImageLabel("farm", p.X+24, y, LargeIconS, LargeIconS, gui.ImageLabelStyleRegular)
	y = y + LargeIconD
	for _, aName := range []string{"log", "vegetable", "fruit", "herb", "sheep", "grain", "reed"} {
		a := artifacts.GetArtifact(aName)
		p.AddLabel(&gui.ImageLabel{Icon: "artifacts/" + a.Name, X: p.X + 24, Y: y, SX: IconS, SY: IconS,
			Style: gui.ImageLabelStyleRegular, OnHoover: func() {
				ArtifactToHelperPanel(cp.GetHelperPanel(true), a)
			}})
		y = y + float64(IconH)
	}
	p.AddImageLabel("mine", p.X+24, y, LargeIconS, LargeIconS, gui.ImageLabelStyleRegular)
	y = y + LargeIconD
	for _, aName := range []string{"clay", "stone", "iron_ore", "gold_ore"} {
		a := artifacts.GetArtifact(aName)
		p.AddLabel(&gui.ImageLabel{Icon: "artifacts/" + a.Name, X: p.X + 24, Y: y, SX: IconS, SY: IconS,
			Style: gui.ImageLabelStyleRegular, OnHoover: func() {
				ArtifactToHelperPanel(cp.GetHelperPanel(true), a)
			}})
		y = y + float64(IconH)
	}

	y = p.Y + 24 + LargeIconD
	p.AddImageLabel("workshop", p.X+p.SX/2+24, y, LargeIconS, LargeIconS, gui.ImageLabelStyleRegular)
	y = y + LargeIconD
	for _, mName := range economy.GetManufactureNamesForET(building.WaterMillWheel) {
		m := economy.GetManufacture(mName)
		p.AddLabel(&gui.ImageLabel{Icon: "tasks/" + m.Name, X: p.X + p.SX/2 + 24, Y: y, SX: IconS, SY: IconS, Style: gui.ImageLabelStyleRegular})
		y = y + float64(IconH)
	}
	p.AddImageLabel("person", p.X+p.SX-24-LargeIconD, p.Y+24+LargeIconD, LargeIconS, LargeIconS, gui.ImageLabelStyleRegular)
}
