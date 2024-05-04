package controller

import (
	"github.com/tfriedel6/canvas"
	"medvil/model/economy"
	"medvil/model/navigation"
	"medvil/model/social"
	"medvil/renderer"
	"medvil/view/gui"
)

const FarmFieldUseTypeDisallocate uint8 = 255

type FarmController struct {
	householdPanel *gui.Panel
	farmPanel      *gui.Panel
	UseType        uint8
	farm           *social.Farm
	cp             *ControlPanel
}

func (fc *FarmController) GetUseType() uint8 {
	return fc.UseType
}

func (fc *FarmController) SetUseType(ut uint8) {
	fc.UseType = ut
}

func FarmToControlPanel(cp *ControlPanel, farm *social.Farm) {
	hp := &gui.Panel{X: 0, Y: ControlPanelDynamicPanelTop, SX: ControlPanelSX, SY: HouseholdControllerSY}
	fp := &gui.Panel{X: 0, Y: ControlPanelDynamicPanelTop + HouseholdControllerSY, SX: ControlPanelSX, SY: ControlPanelDynamicPanelSY - HouseholdControllerSY}
	HouseholdToControlPanel(cp, hp, farm.Household, "farm")
	fc := &FarmController{householdPanel: hp, farmPanel: fp, farm: farm, UseType: economy.FarmFieldUseTypeBarren, cp: cp}

	hcy := HouseholdControllerGUIBottomY * ControlPanelSY
	fp.AddTextLabel("Allocate farm land", 24, hcy-IconS/4.0)
	fp.AddButton(&LandUseButton{
		b:       gui.ButtonGUI{Icon: "artifacts/vegetable", X: float64(24 + IconW*0), Y: hcy, SX: IconS, SY: IconS},
		luc:     fc,
		useType: economy.FarmFieldUseTypeVegetables,
		cp:      cp,
		msg:     "Grow vegetables",
	})
	fp.AddButton(&LandUseButton{
		b:       gui.ButtonGUI{Icon: "artifacts/log", X: float64(24 + IconW*1), Y: hcy, SX: IconS, SY: IconS},
		luc:     fc,
		useType: economy.FarmFieldUseTypeForestry,
		cp:      cp,
		msg:     "Grow forests for log",
	})
	fp.AddButton(&LandUseButton{
		b:       gui.ButtonGUI{Icon: "artifacts/fruit", X: float64(24 + IconW*2), Y: hcy, SX: IconS, SY: IconS},
		luc:     fc,
		useType: economy.FarmFieldUseTypeOrchard,
		cp:      cp,
		msg:     "Grow orchards to produce fruit",
	})
	fp.AddButton(&LandUseButton{
		b:       gui.ButtonGUI{Icon: "artifacts/sheep", X: float64(24 + IconW*3), Y: hcy, SX: IconS, SY: IconS},
		luc:     fc,
		useType: economy.FarmFieldUseTypePasture,
		cp:      cp,
		msg:     "Raise sheep for meat and wool",
	})

	fp.AddButton(&LandUseButton{
		b:       gui.ButtonGUI{Icon: "artifacts/grain", X: float64(24 + IconW*0), Y: hcy + float64(IconH), SX: IconS, SY: IconS},
		luc:     fc,
		useType: economy.FarmFieldUseTypeWheat,
		cp:      cp,
		msg:     "Grow grain to make flour",
	})
	fp.AddButton(&LandUseButton{
		b:       gui.ButtonGUI{Icon: "artifacts/reed", X: float64(24 + IconW*1), Y: hcy + float64(IconH), SX: IconS, SY: IconS},
		luc:     fc,
		useType: economy.FarmFieldUseTypeReed,
		cp:      cp,
		msg:     "Grow reed for paper and thatching",
	})
	fp.AddButton(&LandUseButton{
		b:       gui.ButtonGUI{Icon: "artifacts/herb", X: float64(24 + IconW*2), Y: hcy + float64(IconH), SX: IconS, SY: IconS},
		luc:     fc,
		useType: economy.FarmFieldUseTypeHerb,
		cp:      cp,
		msg:     "Grow herbs to make medicine",
	})
	fp.AddButton(&LandUseButton{
		b:       gui.ButtonGUI{Icon: "clear_land", X: float64(24 + IconW*0), Y: hcy + float64(IconH*2), SX: IconS, SY: IconS},
		luc:     fc,
		useType: economy.FarmFieldUseTypeBarren,
		cp:      cp,
		msg:     "Clear land from trees to build buildings",
	})
	fp.AddButton(&LandUseButton{
		b:       gui.ButtonGUI{Icon: "cancel", X: float64(24 + IconW*1), Y: hcy + float64(IconH*2), SX: IconS, SY: IconS},
		luc:     fc,
		useType: FarmFieldUseTypeDisallocate,
		cp:      cp,
		msg:     "Stop farming",
	})
	fc.RefreshLandUseButtons()

	fp.AddButton(&gui.SimpleButton{
		ButtonGUI: gui.ButtonGUI{Icon: "tasks/calculate", X: 24 + IconS + gui.FontSize*10 + LargeIconD, Y: hcy - gui.FontSize/2.0, SX: LargeIconS, SY: LargeIconS, OnHoover: func() {
			cp.HelperMessage("Optimize land based on profit using paper.", true)
		}},
		Highlight: func() bool { return farm.AutoSwitch },
		ClickImpl: func() {
			farm.AutoSwitch = !farm.AutoSwitch
		}})

	cp.SetDynamicPanel(fc)
	cp.C.ClickHandler = fc
}

func (fc *FarmController) RefreshLandUseButtons() {
	landDist := fc.farm.GetLandDistribution()
	for _, b := range fc.farmPanel.Buttons {
		if lub, ok := b.(*LandUseButton); ok {
			lub.cnt = landDist[lub.useType]
		}
	}
}

func (fc *FarmController) CaptureMove(x, y float64) {
	fc.householdPanel.CaptureMove(x, y)
	fc.farmPanel.CaptureMove(x, y)
}

func (fc *FarmController) CaptureClick(x, y float64) {
	fc.householdPanel.CaptureClick(x, y)
	fc.farmPanel.CaptureClick(x, y)
}

func (fc *FarmController) Render(cv *canvas.Canvas) {
	fc.householdPanel.Render(cv)
	fc.farmPanel.Render(cv)
}

func (fc *FarmController) Clear() {}

func (fc *FarmController) Refresh() {
	fc.householdPanel.Clear()
	HouseholdToControlPanel(fc.cp, fc.householdPanel, fc.farm.Household, "farm")
	fc.CaptureMove(fc.cp.C.X, fc.cp.C.Y)
}

func (fc *FarmController) GetActiveFields(c *Controller, rf *renderer.RenderedField) []navigation.FieldWithContext {
	fields := fc.farm.GetFields()
	if fc.farm.FieldUsableFor(fc.cp.C.Map, rf.F, fc.UseType) && !rf.F.Allocated && fc.UseType != FarmFieldUseTypeDisallocate &&
		(fc.farm.FieldWithinDistance(rf.F) || fc.UseType == economy.FarmFieldUseTypeBarren && fc.farm.FieldWithinDistanceClearing(rf.F)) {
		fields = append(fields, social.FarmLand{
			X:       rf.F.X,
			Y:       rf.F.Y,
			UseType: fc.UseType,
			F:       rf.F,
		})
	} else if fc.UseType != FarmFieldUseTypeDisallocate {
		fields = append(fields, &navigation.BlockedField{F: rf.F})
	}
	return fields
}

func (fc *FarmController) HandleClick(c *Controller, rf *renderer.RenderedField) bool {
	for i := range fc.farm.Land {
		l := &fc.farm.Land[i]
		if l.F.X == rf.F.X && l.F.Y == rf.F.Y {
			if fc.UseType == FarmFieldUseTypeDisallocate {
				// Disallocate land
				fc.farm.Land = append(fc.farm.Land[:i], fc.farm.Land[i+1:]...)
				rf.F.Allocated = false
			} else if fc.farm.FieldUsableFor(c.Map, l.F, fc.UseType) {
				l.UseType = fc.UseType
			}
			fc.RefreshLandUseButtons()
			return true
		}
	}
	if fc.UseType != FarmFieldUseTypeDisallocate && !rf.F.Allocated && fc.farm.FieldUsableFor(c.Map, rf.F, fc.UseType) &&
		(fc.farm.FieldWithinDistance(rf.F) || fc.UseType == economy.FarmFieldUseTypeBarren && fc.farm.FieldWithinDistanceClearing(rf.F)) {
		fc.farm.Land = append(fc.farm.Land, social.FarmLand{
			X:       rf.F.X,
			Y:       rf.F.Y,
			UseType: fc.UseType,
			F:       rf.F,
		})
		rf.F.Allocated = true
		fc.RefreshLandUseButtons()
		return true
	}
	return false
}

func (fc *FarmController) GetHelperSuggestions() *gui.Suggestion {
	suggestion := GetHouseholdHelperSuggestions(fc.farm.Household)
	if suggestion != nil {
		return suggestion
	}
	hcy := HouseholdControllerGUIBottomY * ControlPanelSY
	if fc.UseType == 0 {
		return &gui.Suggestion{
			Message: "Select land cultivation method, then allocate land\nfor various purposes like growing vegetables, grain,\ntrees and sheep by clicking on the land.",
			Icon:    "farm_mixed", X: float64(24 + IconW*4), Y: hcy + float64(IconH),
		}
	}

	landDist := fc.farm.GetLandDistribution()
	if landDist[economy.FarmFieldUseTypeVegetables] < 2 {
		return &gui.Suggestion{
			Message: ("It's recommended to allocate some land to grow vegetables\nin order to make the farm self sustaining.\n" +
				"The farmers will sell excess vegetables on the market,\nso it can be used to feed other villagers."),
			Icon: "artifacts/vegetable", X: float64(24 + IconW*4), Y: hcy + float64(IconH),
		}
	}
	if landDist[economy.FarmFieldUseTypePasture] < 2 {
		return &gui.Suggestion{
			Message: ("Sheeps are useful, they produce meat and materials for clothes.\nIt takes 3 years to raise sheep. The villagers will\n" +
				"sell the sheep at the marketplace. you will need a butchershop,\na certain type of workshop, to produce meat and leather."),
			Icon: "artifacts/sheep", X: float64(24 + IconW*4), Y: hcy + float64(IconH),
		}
	}
	if landDist[economy.FarmFieldUseTypeForestry] < 4 {
		return &gui.Suggestion{
			Message: ("Make sure to grow some trees for firewood and building materials.\nTrees grow slowly and don't need much work, so it's best\n" +
				"to allocate a bit more land for them."),
			Icon: "artifacts/log", X: float64(24 + IconW*4), Y: hcy + float64(IconH),
		}
	}
	if len(fc.farm.Land) > int(fc.farm.Household.TargetNumPeople)*3 {
		return &gui.Suggestion{
			Message: ("Be careful allocating too much land for one farm.\nThe villagers might not be able to cultivate all the\n" +
				"land before winter. You can either release\nthe land or add more villagers to this farm."),
			Icon: "warning", X: float64(24 + IconW*4), Y: hcy + float64(IconH),
		}
	}
	/*
		if fc.UseType == economy.FarmFieldUseTypeWheat {
			return &gui.Suggestion{
				Message: ("Grow wheat to produce grain. It will need to be\nturned into flour using waterwheel mills,\nthen baked as bread."),
				Icon:    "artifacts/grain", X: float64(24 + IconW*4), Y: hcy + float64(IconH),
			}
		}
		if fc.UseType == economy.FarmFieldUseTypeHerb {
			return &gui.Suggestion{
				Message: ("Grow herbs to produce medicine."),
				Icon:    "artifacts/herb", X: float64(24 + IconW*4), Y: hcy + float64(IconH),
			}
		}
		if fc.UseType == economy.FarmFieldUseTypeReed {
			return &gui.Suggestion{
				Message: ("Grow reed to produce thatch for roofs or paper."),
				Icon:    "artifacts/reed", X: float64(24 + IconW*4), Y: hcy + float64(IconH),
			}
		}
		if fc.UseType == economy.FarmFieldUseTypeOrchard {
			return &gui.Suggestion{
				Message: ("Grow an orchard to produce fruits."),
				Icon:    "artifacts/fruit", X: float64(24 + IconW*4), Y: hcy + float64(IconH),
			}
		}
		if fc.UseType == economy.FarmFieldUseTypeForestry {
			return &gui.Suggestion{
				Message: "Grow trees for firewood and building materials.",
				Icon:    "artifacts/log", X: float64(24 + IconW*4), Y: hcy + float64(IconH),
			}
		}
		if fc.UseType == economy.FarmFieldUseTypeVegetables {
			return &gui.Suggestion{
				Message: "Grow vegetables for food.",
				Icon:    "artifacts/vegetable", X: float64(24 + IconW*4), Y: hcy + float64(IconH),
			}
		}
		if fc.UseType == economy.FarmFieldUseTypePasture {
			return &gui.Suggestion{
				Message: "Raise sheep for food and textile.",
				Icon:    "artifacts/sheep", X: float64(24 + IconW*4), Y: hcy + float64(IconH),
			}
		}
		if fc.UseType == economy.FarmFieldUseTypeBarren {
			return &gui.Suggestion{
				Message: ("Clear land in order to build houses on them."),
				Icon:    "clear_land", X: float64(24 + IconW*4), Y: hcy + float64(IconH),
			}
		}
	*/
	return nil
}
