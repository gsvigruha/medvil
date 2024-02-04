package view

import (
	"github.com/tfriedel6/canvas"
	"image/color"
	"medvil/controller"
	"medvil/model/navigation"
	"medvil/model/social"
	"medvil/renderer"
	"medvil/view/buildings"
	"medvil/view/gui"
	"path/filepath"
)

var coinI = filepath.FromSlash("icon/gui/coin.png")
var foodI = filepath.FromSlash("icon/gui/food.png")
var drinkI = filepath.FromSlash("icon/gui/drink.png")
var happinessI = filepath.FromSlash("icon/gui/happiness.png")
var healthI = filepath.FromSlash("icon/gui/health.png")
var clothesI = filepath.FromSlash("icon/gui/artifacts/clothes.png")
var heatingI = filepath.FromSlash("icon/gui/heating.png")
var repairI = filepath.FromSlash("icon/gui/tasks/repair.png")
var personI = filepath.FromSlash("icon/gui/person.png")
var woodI = filepath.FromSlash("texture/wood.png")
var warnI = filepath.FromSlash("icon/gui/warning_slim.png")

func iconsFromHousehold(h *social.Household, moneyThreshold int, icons *[]string) {
	if int(h.Money) < moneyThreshold {
		*icons = append(*icons, coinI)
	}

	if len(h.People) > 0 {
		var food = 0
		var water = 0
		var happiness = 0
		var health = 0
		for _, person := range h.People {
			food += int(person.Food)
			water += int(person.Water)
			happiness += int(person.Happiness)
			health += int(person.Health)
		}
		if food/len(h.People) < 25 {
			*icons = append(*icons, foodI)
		}
		if water/len(h.People) < 25 {
			*icons = append(*icons, drinkI)
		}
		if happiness/len(h.People) < 25 {
			*icons = append(*icons, happinessI)
		}
		if health/len(h.People) < 25 {
			*icons = append(*icons, healthI)
		}
		if !h.HasEnoughClothes() {
			*icons = append(*icons, clothesI)
		}
		if h.GetHeating() < 100 {
			*icons = append(*icons, heatingI)
		}
		if h.Building.Broken {
			*icons = append(*icons, repairI)
		}
	} else {
		*icons = append(*icons, personI)
	}
}

func DrawHouseholdIcons(cv *canvas.Canvas, rf renderer.RenderedField, f *navigation.Field, k int, c *controller.Controller) {
	b := f.Building.GetBuilding()
	if b.X != f.X || b.Y != f.Y {
		return
	}

	z := (float64(k) + 0.5) * buildings.BuildingUnitHeight * buildings.DZ
	midX, midY := rf.MidScreenPoint()

	var icons []string

	farm := c.ReverseReferences.BuildingToFarm[b]
	if farm != nil {
		iconsFromHousehold(farm.Household, farm.Household.Town.Transfers.Farm.Threshold, &icons)
	}
	workshop := c.ReverseReferences.BuildingToWorkshop[b]
	if workshop != nil {
		iconsFromHousehold(workshop.Household, workshop.Household.Town.Transfers.Workshop.Threshold, &icons)
	}
	mine := c.ReverseReferences.BuildingToMine[b]
	if mine != nil {
		iconsFromHousehold(mine.Household, mine.Household.Town.Transfers.Mine.Threshold, &icons)
	}
	factory := c.ReverseReferences.BuildingToFactory[b]
	if factory != nil {
		iconsFromHousehold(factory.Household, factory.Household.Town.Transfers.Factory.Threshold, &icons)
	}
	townhall := c.ReverseReferences.BuildingToTownhall[b]
	if townhall != nil {
		iconsFromHousehold(townhall.Household, int(townhall.Household.Town.Stats.Global.Money)/10, &icons)
	}
	market := c.ReverseReferences.BuildingToMarketplace[b]
	if market != nil {
		if int(market.Money) < int(market.Town.Stats.Global.Money)/10 {
			icons = append(icons, coinI)
		}
	}

	if len(icons) > 0 {
		var s float64
		if b == c.HooveredBuilding {
			s = controller.IconS
		} else {
			s = controller.IconS * 0.6
		}

		cv.SetFillStyle(color.RGBA{R: 0, G: 0, B: 0, A: 192})
		left := midX - float64(len(icons))*s/2
		cv.FillRect(left-s*0.15, midY-z-s, float64(len(icons)+1)*s, s)
		cv.DrawImage(warnI, left-s*0.15, midY-z-s, s, s)
		for i, icon := range icons {
			cv.DrawImage(icon, left+float64(i)*s+s/2, midY-z-s, s, s)
		}
	}
}

func DrawLabels(cv *canvas.Canvas, rf renderer.RenderedField, f *navigation.Field, k int, c *controller.Controller) {
	if f.Building.GetBuilding().X != f.X || f.Building.GetBuilding().Y != f.Y {
		return
	}

	z := float64(k+1) * buildings.BuildingUnitHeight * buildings.DZ
	midX, midY := rf.MidScreenPoint()

	townhall := c.ReverseReferences.BuildingToTownhall[f.Building.GetBuilding()]
	if townhall != nil {
		name := townhall.Household.Town.Name
		if name != "" {
			dx := gui.EstimateWidth(name) * gui.FontSize / 2.0
			y := midY - z - 10
			dy := gui.FontSize
			if c.ActiveSupplier == townhall.Household.Town {
				cv.SetStrokeStyle(color.RGBA{R: 0, G: 192, B: 0, A: 255})
				cv.SetLineWidth(4.0)
				cv.StrokeRect(midX-dx-8, y-dy-2, dx*2+16, dy+10)
			}
			cv.SetFillStyle(woodI)
			cv.FillRect(midX-dx-8, y-dy-2, dx*2+16, dy+10)
			cv.SetFillStyle("#FED")
			cv.SetFont(gui.Font, gui.FontSize)
			cv.FillText(name, midX-dx, y)
		}
	}
}
