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
)

func iconsFromHousehold(h *social.Household, moneyThreshold int, icons *[]string) {
	if int(h.Money) < moneyThreshold {
		*icons = append(*icons, "icon/gui/profitable.png")
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
			*icons = append(*icons, "icon/gui/food.png")
		}
		if water/len(h.People) < 25 {
			*icons = append(*icons, "icon/gui/drink.png")
		}
		if happiness/len(h.People) < 25 {
			*icons = append(*icons, "icon/gui/happiness.png")
		}
		if health/len(h.People) < 25 {
			*icons = append(*icons, "icon/gui/health.png")
		}
		if !h.HasEnoughClothes() {
			*icons = append(*icons, "icon/gui/artifacts/clothes.png")
		}
		if h.GetHeating() < 100 {
			*icons = append(*icons, "icon/gui/heating.png")
		}
		if h.Building.Broken {
			*icons = append(*icons, "icon/gui/tasks/repair.png")
		}
	}
}

func DrawHouseholdIcons(cv *canvas.Canvas, rf renderer.RenderedField, f *navigation.Field, k int, c *controller.Controller) {
	if f.Building.GetBuilding().X != f.X || f.Building.GetBuilding().Y != f.Y {
		return
	}

	z := float64(k+1) * buildings.BuildingUnitHeight * buildings.DZ
	midX, midY := rf.MidScreenPoint()

	var icons []string

	farm := c.ReverseReferences.BuildingToFarm[f.Building.GetBuilding()]
	if farm != nil {
		iconsFromHousehold(farm.Household, farm.Household.Town.Transfers.Farm.Threshold, &icons)
	}
	workshop := c.ReverseReferences.BuildingToWorkshop[f.Building.GetBuilding()]
	if workshop != nil {
		iconsFromHousehold(workshop.Household, workshop.Household.Town.Transfers.Workshop.Threshold, &icons)
	}
	mine := c.ReverseReferences.BuildingToMine[f.Building.GetBuilding()]
	if mine != nil {
		iconsFromHousehold(mine.Household, mine.Household.Town.Transfers.Mine.Threshold, &icons)
	}
	factory := c.ReverseReferences.BuildingToFactory[f.Building.GetBuilding()]
	if factory != nil {
		iconsFromHousehold(factory.Household, factory.Household.Town.Transfers.Factory.Threshold, &icons)
	}
	townhall := c.ReverseReferences.BuildingToTownhall[f.Building.GetBuilding()]
	if townhall != nil {
		iconsFromHousehold(townhall.Household, int(townhall.Household.Town.Stats.Global.Money)/10, &icons)
	}
	market := c.ReverseReferences.BuildingToMarketplace[f.Building.GetBuilding()]
	if market != nil {
		if int(market.Money) < int(market.Town.Stats.Global.Money)/10 {
			icons = append(icons, "icon/gui/profitable.png")
		}
	}

	for i, icon := range icons {
		cv.DrawImage(icon, midX-float64(len(icons))*24/2+float64(i)*24, midY-z, 24, 24)
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
			cv.SetFillStyle("texture/wood.png")
			cv.FillRect(midX-dx-8, y-dy-2, dx*2+16, dy+10)
			cv.SetFillStyle("#FED")
			cv.SetFont("texture/font/Go-Regular.ttf", gui.FontSize)
			cv.FillText(name, midX-dx, y)
		}
	}
}
