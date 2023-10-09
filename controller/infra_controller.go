package controller

import (
	"github.com/tfriedel6/canvas"
	"image/color"
	"medvil/model/building"
	"medvil/model/economy"
	"medvil/model/navigation"
	"medvil/model/social"
	"medvil/renderer"
	"medvil/view/gui"
)

type InfraType uint8

const InfraTypeNone = 0
const InfraTypeDirtRoad = 1
const InfraTypeCobbleRoad = 2
const InfraTypeCanal = 3
const InfraTypeAqueduct = 4
const InfraTypeBridge = 5
const InfraTypeStoneWall1 = 11
const InfraTypeStoneWall2 = 12
const InfraTypeStoneWallRamp = 14
const InfraTypeStoneTower1 = 15
const InfraTypeStoneTower2 = 16
const InfraTypeGateNS = 21
const InfraTypeGateEW = 22
const InfraTypeLevelForBuilding = 31
const InfraTypeLevelForRoad = 32
const InfraTypeFountain = 41
const InfraTypeObelisk = 42
const InfraTypeOakTree = 51
const InfraTypeAppleTree = 52

const InfraPanelTop = 0.1

type InfraController struct {
	it InfraType
	cp *ControlPanel
}

type InfraBuildButton struct {
	b   *gui.ButtonGUI
	it  InfraType
	msg string
	ic  *InfraController
}

func (b *InfraBuildButton) SetHoover(h bool) {
	b.b.SetHoover(h)
}

func (b InfraBuildButton) Click() {
	b.ic.it = b.it
	b.ic.cp.HelperMessage(b.msg)
}

func (b InfraBuildButton) Render(cv *canvas.Canvas) {
	b.b.Render(cv)
	if b.ic.it != b.it {
		cv.SetFillStyle(color.RGBA{R: 64, G: 0, B: 0, A: 128})
		cv.FillRect(b.b.X, b.b.Y, LargeIconS, LargeIconS)
	}
}

func (b InfraBuildButton) Contains(x float64, y float64) bool {
	return b.b.Contains(x, y)
}

func (b InfraBuildButton) Enabled() bool {
	return b.b.Enabled()
}

func (ic *InfraController) CheckField(c *Controller, rf *renderer.RenderedField) bool {
	if !c.ActiveSupplier.FieldWithinDistance(rf.F) {
		return false
	}
	if ic.it == InfraTypeDirtRoad {
		return rf.F.RoadCompatible()
	} else if ic.it == InfraTypeCobbleRoad {
		return rf.F.RoadCompatible() || (rf.F.Road != nil && rf.F.Road.T == building.DirtRoadType)
	} else if ic.it == InfraTypeCanal {
		return rf.F.Buildable()
	} else if ic.it == InfraTypeBridge {
		return c.Map.Shore(rf.F.X, rf.F.Y)
	} else if ic.it == InfraTypeStoneWall1 || ic.it == InfraTypeStoneWall2 {
		return rf.F.WallCompatible()
	} else if ic.it == InfraTypeStoneTower1 || ic.it == InfraTypeStoneTower2 {
		return rf.F.Buildable()
	} else if ic.it == InfraTypeStoneWallRamp {
		return navigation.IsRampPossible(c.Map, rf.F.X, rf.F.Y)
	} else if ic.it == InfraTypeGateNS {
		return c.Map.IsBuildingPossible(rf.F.X, rf.F.Y, building.SmallGate, building.DirectionN)
	} else if ic.it == InfraTypeGateEW {
		return c.Map.IsBuildingPossible(rf.F.X, rf.F.Y, building.SmallGate, building.DirectionE)
	} else if ic.it == InfraTypeLevelForBuilding {
		return navigation.FieldCanBeLeveledForBuilding(*rf.F, c.Map)
	} else if ic.it == InfraTypeLevelForRoad {
		return navigation.FieldCanBeLeveledForRoad(*rf.F, c.Map)
	} else if ic.it == InfraTypeFountain || ic.it == InfraTypeObelisk {
		return rf.F.StatueCompatible()
	} else if ic.it == InfraTypeOakTree || ic.it == InfraTypeAppleTree {
		return rf.F.Plantable()
	}
	return false
}

func (ic *InfraController) GetActiveFields(c *Controller, rf *renderer.RenderedField) []navigation.FieldWithContext {
	if ic.CheckField(c, rf) {
		return []navigation.FieldWithContext{rf.F}
	}
	return nil
}

func (ic *InfraController) HandleClick(c *Controller, rf *renderer.RenderedField) bool {
	if c.ActiveSupplier == nil {
		return false
	}
	if ic.CheckField(c, rf) {
		if ic.it == InfraTypeLevelForBuilding {
			c.Map.AddLevelingTask(c.ActiveSupplier, rf.F.X, rf.F.Y, economy.TerraformTaskTypeLevelForBuilding)
			return true
		} else if ic.it == InfraTypeLevelForRoad {
			c.Map.AddLevelingTask(c.ActiveSupplier, rf.F.X, rf.F.Y, economy.TerraformTaskTypeLevelForRoad)
			return true
		}
		if activeTown, ok := c.ActiveSupplier.(*social.Town); ok {
			if ic.it == InfraTypeDirtRoad {
				c.Map.AddRoadConstruction(activeTown, rf.F.X, rf.F.Y, building.DirtRoadType)
			} else if ic.it == InfraTypeCobbleRoad {
				c.Map.AddRoadConstruction(activeTown, rf.F.X, rf.F.Y, building.CobbleRoadType)
			} else if ic.it == InfraTypeCanal {
				c.Map.AddInfraConstruction(activeTown, rf.F.X, rf.F.Y, building.CanalType)
			} else if ic.it == InfraTypeBridge {
				c.Map.AddRoadConstruction(activeTown, rf.F.X, rf.F.Y, building.BridgeRoadType)
			} else if ic.it == InfraTypeStoneWall1 {
				c.Map.AddBuildingConstruction(activeTown, rf.F.X, rf.F.Y, building.StoneWall1Type, building.DirectionNone)
			} else if ic.it == InfraTypeStoneWall2 {
				c.Map.AddBuildingConstruction(activeTown, rf.F.X, rf.F.Y, building.StoneWall2Type, building.DirectionNone)
			} else if ic.it == InfraTypeStoneTower1 {
				c.Map.AddBuildingConstruction(activeTown, rf.F.X, rf.F.Y, building.Tower1Type, building.DirectionNone)
			} else if ic.it == InfraTypeStoneTower2 {
				c.Map.AddBuildingConstruction(activeTown, rf.F.X, rf.F.Y, building.Tower2Type, building.DirectionNone)
			} else if ic.it == InfraTypeStoneWallRamp {
				c.Map.AddWallRampConstruction(activeTown, rf.F.X, rf.F.Y)
			} else if ic.it == InfraTypeGateNS {
				c.Map.AddBuildingConstruction(activeTown, rf.F.X, rf.F.Y, building.SmallGate, building.DirectionN)
			} else if ic.it == InfraTypeGateEW {
				c.Map.AddBuildingConstruction(activeTown, rf.F.X, rf.F.Y, building.SmallGate, building.DirectionE)
			} else if ic.it == InfraTypeFountain {
				c.Map.AddStatueConstruction(activeTown, rf.F.X, rf.F.Y, building.FountainType)
			} else if ic.it == InfraTypeObelisk {
				c.Map.AddStatueConstruction(activeTown, rf.F.X, rf.F.Y, building.ObeliskType)
			} else if ic.it == InfraTypeOakTree {
				activeTown.Townhall.Household.AddTask(&economy.AgriculturalTask{T: economy.AgriculturalTaskPlantingOakTree, F: rf.F, Start: *c.Map.Calendar})
			} else if ic.it == InfraTypeAppleTree {
				activeTown.Townhall.Household.AddTask(&economy.AgriculturalTask{T: economy.AgriculturalTaskPlantingAppleTree, F: rf.F, Start: *c.Map.Calendar})
			}
			return true
		}
	}
	return false
}

func InfraToControlPanel(cp *ControlPanel) {
	p := &gui.Panel{X: 0, Y: ControlPanelDynamicPanelTop, SX: ControlPanelSX, SY: ControlPanelDynamicPanelSY}
	ic := &InfraController{it: InfraTypeNone, cp: cp}

	top := InfraPanelTop * ControlPanelSY

	if cp.C.ActiveSupplier != nil && cp.C.ActiveSupplier.BuildHousesEnabled() {

		p.AddButton(&InfraBuildButton{
			b:  &gui.ButtonGUI{Texture: "terrain/grass", X: float64(24 + LargeIconD*0), Y: top, SX: LargeIconS, SY: LargeIconS},
			it: InfraTypeNone,
			ic: ic,
		})

		p.AddButton(&InfraBuildButton{
			b:   &gui.ButtonGUI{Icon: "infra/dirt_road", X: float64(24 + LargeIconD*1), Y: top, SX: LargeIconS, SY: LargeIconS},
			it:  InfraTypeDirtRoad,
			msg: "Build dirt road. Speeds up commute with 50%.",
			ic:  ic,
		})

		p.AddButton(&InfraBuildButton{
			b:   &gui.ButtonGUI{Icon: "infra/cobble_road", X: float64(24 + LargeIconD*2), Y: top, SX: LargeIconS, SY: LargeIconS},
			it:  InfraTypeCobbleRoad,
			msg: "Build cobble road. Speeds up commute with 100%.",
			ic:  ic,
		})

		p.AddButton(&InfraBuildButton{
			b:   &gui.ButtonGUI{Icon: "infra/canal", X: float64(24 + LargeIconD*3), Y: top, SX: LargeIconS, SY: LargeIconS},
			it:  InfraTypeCanal,
			msg: "Extend water with canals for drinking and transport.",
			ic:  ic,
		})

		p.AddButton(&InfraBuildButton{
			b:   &gui.ButtonGUI{Icon: "infra/bridge", X: float64(24 + LargeIconD*4), Y: top, SX: LargeIconS, SY: LargeIconS},
			it:  InfraTypeBridge,
			msg: "Build bridges. People can cross small rivers.",
			ic:  ic,
		})

		p.AddButton(&InfraBuildButton{
			b:   &gui.ButtonGUI{Icon: "infra/wall_small", X: float64(24 + LargeIconD*0), Y: top + float64(LargeIconD*1), SX: LargeIconS, SY: LargeIconS},
			it:  InfraTypeStoneWall1,
			msg: "Build a small wall.",
			ic:  ic,
		})

		p.AddButton(&InfraBuildButton{
			b:   &gui.ButtonGUI{Icon: "infra/wall_large", X: float64(24 + LargeIconD*1), Y: top + float64(LargeIconD*1), SX: LargeIconS, SY: LargeIconS},
			it:  InfraTypeStoneWall2,
			msg: "Build a large wall.",
			ic:  ic,
		})

		p.AddButton(&InfraBuildButton{
			b:   &gui.ButtonGUI{Icon: "infra/tower_1", X: float64(24 + LargeIconD*2), Y: top + float64(LargeIconD*1), SX: LargeIconS, SY: LargeIconS},
			it:  InfraTypeStoneTower1,
			msg: "Build a small tower.",
			ic:  ic,
		})

		p.AddButton(&InfraBuildButton{
			b:   &gui.ButtonGUI{Icon: "infra/tower_2", X: float64(24 + LargeIconD*3), Y: top + float64(LargeIconD*1), SX: LargeIconS, SY: LargeIconS},
			it:  InfraTypeStoneTower2,
			msg: "Build a large tower.",
			ic:  ic,
		})

		p.AddButton(&InfraBuildButton{
			b:   &gui.ButtonGUI{Icon: "infra/ramp", X: float64(24 + LargeIconD*4), Y: top + float64(LargeIconD*1), SX: LargeIconS, SY: LargeIconS},
			it:  InfraTypeStoneWallRamp,
			msg: "Build a ramp to make walls accessible.",
			ic:  ic,
		})

		p.AddButton(&InfraBuildButton{
			b:   &gui.ButtonGUI{Icon: "infra/gate_ns", X: float64(24 + LargeIconD*0), Y: top + float64(LargeIconD*2), SX: LargeIconS, SY: LargeIconS},
			it:  InfraTypeGateNS,
			msg: "Build a north-south gate over water or land.",
			ic:  ic,
		})

		p.AddButton(&InfraBuildButton{
			b:   &gui.ButtonGUI{Icon: "infra/gate_ew", X: float64(24 + LargeIconD*1), Y: top + float64(LargeIconD*2), SX: LargeIconS, SY: LargeIconS},
			it:  InfraTypeGateEW,
			msg: "Build an east-west gate over water or land.",
			ic:  ic,
		})
	}

	p.AddButton(&InfraBuildButton{
		b:   &gui.ButtonGUI{Icon: "infra/terraform_building", X: float64(24 + LargeIconD*0), Y: top + float64(LargeIconD*3), SX: LargeIconS, SY: LargeIconS},
		it:  InfraTypeLevelForBuilding,
		msg: "Terraform hills in order to build houses on it.",
		ic:  ic,
	})

	p.AddButton(&InfraBuildButton{
		b:   &gui.ButtonGUI{Icon: "infra/terraform_road", X: float64(24 + LargeIconD*1), Y: top + float64(LargeIconD*3), SX: LargeIconS, SY: LargeIconS},
		it:  InfraTypeLevelForRoad,
		msg: "Terraform hills in order to build roads on it.",
		ic:  ic,
	})

	if cp.C.ActiveSupplier != nil && cp.C.ActiveSupplier.BuildHousesEnabled() {
		p.AddButton(&InfraBuildButton{
			b:   &gui.ButtonGUI{Icon: "infra/fountain", X: float64(24 + LargeIconD*0), Y: top + float64(LargeIconD*4), SX: LargeIconS, SY: LargeIconS},
			it:  InfraTypeFountain,
			msg: "Statues make your population happy.",
			ic:  ic,
		})

		p.AddButton(&InfraBuildButton{
			b:   &gui.ButtonGUI{Icon: "infra/obelisk", X: float64(24 + LargeIconD*1), Y: top + float64(LargeIconD*4), SX: LargeIconS, SY: LargeIconS},
			it:  InfraTypeObelisk,
			msg: "Statues make your population happy.",
			ic:  ic,
		})

		p.AddButton(&InfraBuildButton{
			b:   &gui.ButtonGUI{Icon: "infra/oak_tree", X: float64(24 + LargeIconD*0), Y: top + float64(LargeIconD*5), SX: LargeIconS, SY: LargeIconS},
			it:  InfraTypeOakTree,
			msg: "Trees make your population happy and healthy.",
			ic:  ic,
		})

		p.AddButton(&InfraBuildButton{
			b:   &gui.ButtonGUI{Icon: "infra/apple_tree", X: float64(24 + LargeIconD*1), Y: top + float64(LargeIconD*5), SX: LargeIconS, SY: LargeIconS},
			it:  InfraTypeAppleTree,
			msg: "Trees make your population happy and healthy.",
			ic:  ic,
		})
	}

	cp.SetDynamicPanel(p)
	cp.C.ClickHandler = ic
}
