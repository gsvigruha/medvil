package controller

import (
	"github.com/tfriedel6/canvas"
	"math"
	"medvil/model/artifacts"
	"medvil/model/building"
	"medvil/model/economy"
	"medvil/view/gui"
)

const SplashPageGameplay = 0
const SplashPageControls = 1
const SplashPageEconomics = 2

type Splash struct {
	p            *gui.Panel
	gp           *gui.Panel
	cp           *gui.Panel
	ep           *gui.Panel
	selectedA    *artifacts.Artifact
	controlPanel *ControlPanel
	page         int
}

func (s *Splash) CaptureClick(x float64, y float64) {
	s.p.CaptureClick(x, y)
	if s.ep != nil {
		s.ep.CaptureClick(x, y)
	}
}

func (s *Splash) CaptureMove(x float64, y float64) {
	s.p.CaptureMove(x, y)
	if s.ep != nil {
		s.ep.CaptureMove(x, y)
	}
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
		ButtonGUI: gui.ButtonGUI{Icon: "barrel", X: px + LargeIconS*2 + 24, Y: py + 24, SX: IconS, SY: IconS},
		ClickImpl: func() {
			s.page = SplashPageEconomics
			s.addEconomicsLabels(s.controlPanel, s.ep)
		},
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
		s.selectedA = artifacts.GetArtifact("bread")
		s.addEconomicsLabels(cp, ep)
		s.ep = ep
	}

	s.p = p
	s.controlPanel = cp
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

type EconomicChainNode struct {
	A               *artifacts.Artifact
	Parents         []*EconomicChainNode
	ProducerIcon    string
	ProducerSubIcon string
	Description     string
	SubDescription  string
}

func setupWorkshopType(et *building.BuildingExtensionType, nodes map[string]*EconomicChainNode) {
	manufactures := economy.GetManufactureNamesForET(et)
	for _, mName := range manufactures {
		m := economy.GetManufacture(mName)
		for _, output := range m.Outputs {
			if _, ok := nodes[output.A.Name]; !ok {
				nodes[output.A.Name] = &EconomicChainNode{
					A: output.A, ProducerIcon: "workshop", ProducerSubIcon: "building/" + et.Name,
					Description:    HelperMessageForBuildingType(building.BuildingTypeWorkshop),
					SubDescription: et.Description,
				}
				for _, input := range m.Inputs {
					nodes[output.A.Name].Parents = append(nodes[output.A.Name].Parents, nodes[input.A.Name])
				}
			}
		}
	}
}

func (s *Splash) renderNodes(a *artifacts.Artifact, nodes map[string]*EconomicChainNode, i int, px, py float64) {
	var x = s.ep.X + s.ep.SX - 24 - LargeIconD*2 - LargeIconD*2*float64(i)
	var y = py
	dy := s.ep.SY / math.Pow(2, float64(i))
	if len(nodes[a.Name].Parents) > 0 {
		y = y - float64(len(nodes[a.Name].Parents)-1)/2*dy - LargeIconD
	} else {
		s.ep.AddLabel(&gui.ArrowLabel{SX: x + LargeIconD + 4, SY: y + float64(IconH)/2, EX: px - 4, EY: py + float64(IconH)/2})
	}
	if nodes[a.Name].ProducerSubIcon != "" {
		s.ep.AddLabel(&gui.ImageLabel{Icon: nodes[a.Name].ProducerIcon, X: x - LargeIconD/2, Y: y, SX: LargeIconS, SY: LargeIconS,
			Style: gui.ImageLabelStyleRegular, OnHoover: func() {
				s.controlPanel.HelperMessage(nodes[a.Name].Description, false)
			}})
		s.ep.AddLabel(&gui.ImageLabel{Icon: nodes[a.Name].ProducerSubIcon, X: x + LargeIconD/2, Y: y, SX: LargeIconS, SY: LargeIconS,
			Style: gui.ImageLabelStyleRegular, OnHoover: func() {
				s.controlPanel.HelperMessage(nodes[a.Name].SubDescription, false)
			}})
	} else {
		s.ep.AddLabel(&gui.ImageLabel{Icon: nodes[a.Name].ProducerIcon, X: x, Y: y, SX: LargeIconS, SY: LargeIconS,
			Style: gui.ImageLabelStyleRegular, OnHoover: func() {
				s.controlPanel.HelperMessage(nodes[a.Name].Description, false)
			}})
	}
	y = y + LargeIconD
	for j, node := range nodes[a.Name].Parents {
		aJ := node.A
		s.ep.AddLabel(&gui.ImageLabel{Icon: "artifacts/" + node.A.Name, X: x, Y: y, SX: IconS, SY: IconS,
			Style: gui.ImageLabelStyleRegular, OnHoover: func() {
				ArtifactToHelperPanel(s.controlPanel.GetHelperPanel(true), aJ)
			}})
		s.ep.AddLabel(&gui.ArrowLabel{SX: x + float64(IconW) + 4, SY: y + float64(IconH)/2, EX: px - 4, EY: py + float64(IconH)/2 + float64(j-len(nodes[a.Name].Parents)/2)*10})
		s.renderNodes(aJ, nodes, i+1, x, y)
		y = y + dy
	}
}

func (s *Splash) addEconomicsLabels(cp *ControlPanel, p *gui.Panel) {
	p.Clear()
	p.AddLargeTextLabel("Economic production chain", p.X+24, p.Y+24)
	var nodes map[string]*EconomicChainNode = make(map[string]*EconomicChainNode)

	var y = p.Y + 24 + LargeIconD
	p.AddImageLabel("person", p.X+p.SX-24-LargeIconD*2, p.Y+24+LargeIconD, LargeIconS, LargeIconS, gui.ImageLabelStyleRegular)
	y = y + LargeIconD
	for _, aName := range []string{"vegetable", "fruit", "bread", "meat", "log", "textile", "leather", "beer", "medicine", "tools"} {
		a := artifacts.GetArtifact(aName)
		p.AddButton(&gui.SimpleButton{
			ButtonGUI: gui.ButtonGUI{Icon: "artifacts/" + a.Name, X: p.X + p.SX - 24 - LargeIconD*2, Y: y, SX: IconS, SY: IconS},
			ClickImpl: func() {
				s.selectedA = a
				s.addEconomicsLabels(s.controlPanel, s.ep)
			},
			Highlight: func() bool {
				return s.selectedA == a
			},
		})
		y = y + float64(IconH)
	}
	y = p.Y + 24 + LargeIconD
	p.AddImageLabel("house", p.X+p.SX-24-LargeIconD, y, LargeIconS, LargeIconS, gui.ImageLabelStyleRegular)
	y = y + LargeIconD
	for _, aName := range []string{"cube", "brick", "board", "tile", "thatch", "paper", "sword", "shield"} {
		a := artifacts.GetArtifact(aName)
		p.AddButton(&gui.SimpleButton{
			ButtonGUI: gui.ButtonGUI{Icon: "artifacts/" + a.Name, X: p.X + p.SX - 24 - LargeIconD, Y: y, SX: IconS, SY: IconS},
			ClickImpl: func() {
				s.selectedA = a
				s.addEconomicsLabels(s.controlPanel, s.ep)
			},
			Highlight: func() bool {
				return s.selectedA == a
			},
		})
		y = y + float64(IconH)
	}

	for _, aName := range []string{"log", "vegetable", "fruit", "herb", "sheep", "wool", "grain", "reed"} {
		a := artifacts.GetArtifact(aName)
		nodes[a.Name] = &EconomicChainNode{A: a, ProducerIcon: "farm", Description: HelperMessageForBuildingType(building.BuildingTypeFarm)}
	}
	for _, aName := range []string{"clay", "stone", "iron_ore", "gold_ore"} {
		a := artifacts.GetArtifact(aName)
		nodes[a.Name] = &EconomicChainNode{A: a, ProducerIcon: "mine", Description: HelperMessageForBuildingType(building.BuildingTypeMine)}
	}
	water := artifacts.GetArtifact("water")
	nodes[water.Name] = &EconomicChainNode{A: water, ProducerIcon: "house", Description: "Any building can collect water"}

	setupWorkshopType(building.Workshop, nodes)
	setupWorkshopType(building.WaterMillWheel, nodes)
	setupWorkshopType(building.Kiln, nodes)
	setupWorkshopType(building.Forge, nodes)
	setupWorkshopType(building.Cooker, nodes)

	if s.page == SplashPageEconomics && s.selectedA != nil {
		x := s.ep.X + s.ep.SX - 24 - LargeIconD*3
		y := s.ep.Y + s.ep.SY/2
		s.ep.AddLabel(&gui.ImageLabel{Icon: "artifacts/" + s.selectedA.Name, X: x, Y: y, SX: IconS, SY: IconS,
			Style: gui.ImageLabelStyleRegular, OnHoover: func() {
				ArtifactToHelperPanel(s.controlPanel.GetHelperPanel(true), s.selectedA)
			}})
		s.renderNodes(s.selectedA, nodes, 2, x, y)
	}
}
