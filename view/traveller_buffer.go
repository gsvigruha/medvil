package view

import (
	"github.com/tfriedel6/canvas"
	"github.com/tfriedel6/canvas/backend/goglbackend"
	"image/color"
	"medvil/controller"
	"medvil/model/building"
	"medvil/model/navigation"
	"medvil/renderer"
	"medvil/view/buildings"
	"medvil/view/gui"
	"strconv"
	"time"
)

type TravellerImageCache struct {
	entries map[string]*CacheEntry
	ctx     *goglbackend.GLContext
}

func (tc *TravellerImageCache) RenderTravellerOnBuffer(t *navigation.Traveller, f *navigation.Field, w, h int, c *controller.Controller) *canvas.Canvas {
	key := t.CacheKey(c.Perspective) + "#" + strconv.FormatBool(tallPlant(f))
	person := t.Person
	if person != nil {
		key = key + "#" + person.CacheKey()
	}
	if ce, ok := tc.entries[key]; ok {
		return ce.cv
	} else {
		offscreen, _ := goglbackend.NewOffscreen(w, h, true, tc.ctx)
		cv := canvas.New(offscreen)
		cv.ClearRect(0, 0, float64(w), float64(h))
		DrawTraveller(cv, t, float64(w/2), float64(h*3/4), f, c)
		tc.entries[key] = &CacheEntry{
			offscreen:   offscreen,
			cv:          cv,
			createdTime: time.Now().UnixNano(),
		}
		return cv
	}
}

func getZByDir(bpe *navigation.BuildingPathElement, dir uint8) float64 {
	if bpe.BC.Connection(dir) == building.ConnectionTypeUpperLevel {
		return float64(bpe.GetLocation().Z) * buildings.DZ * buildings.BuildingUnitHeight
	} else if bpe.BC.Connection(dir) == building.ConnectionTypeLowerLevel {
		return float64(bpe.GetLocation().Z-1) * buildings.DZ * buildings.BuildingUnitHeight
	}
	return 0
}

func travellerWH(t uint8) (int, int) {
	if t == navigation.TravellerTypePedestrian {
		return 16, 48
	} else if t == navigation.TravellerTypeBoat {
		return 48, 48
	} else if t == navigation.TravellerTypeTradingBoat {
		return 48, 48
	} else if t == navigation.TravellerTypeExpeditionBoat {
		return 96, 96
	} else if t == navigation.TravellerTypeCart {
		return 36, 48
	} else if t == navigation.TravellerTypeTradingCart {
		return 36, 48
	} else if t == navigation.TravellerTypeExpeditionCart {
		return 96, 96
	}
	return 0, 0
}

func RenderTravellers(ic *ImageCache, cv *canvas.Canvas, travellers []*navigation.Traveller, show func(*navigation.Traveller) bool, rf renderer.RenderedField, c *controller.Controller) {
	for i := range travellers {
		t := travellers[i]
		px := float64(t.PX)
		py := float64(t.PY)
		x, y := GetScreenXY(t, rf, c)
		if !show(t) {
			continue
		}
		if !t.Visible {
			continue
		}
		var z = 0.0
		if t.GetPathElement() != nil && t.GetPathElement().GetLocation().Z > 0 {
			if bpe, ok := t.GetPathElement().(*navigation.BuildingPathElement); ok {
				z1 := getZByDir(bpe, t.Direction)
				z2 := getZByDir(bpe, building.OppDir(t.Direction))
				switch t.Direction {
				case navigation.DirectionN:
					z = (z1*(MaxPY-py) + z2*py) / MaxPY
				case navigation.DirectionS:
					z = (z1*py + z2*(MaxPY-py)) / MaxPY
				case navigation.DirectionW:
					z = (z1*(MaxPX-px) + z2*px) / MaxPX
				case navigation.DirectionE:
					z = (z1*px + z2*(MaxPX-px)) / MaxPX
				}
			}
		}
		w, h := travellerWH(t.T)
		travellerImg := ic.Tic.RenderTravellerOnBuffer(t, rf.F, w, h, c)
		rx := x - float64(w/2)
		ry := y - z - float64(h*3/4)
		cv.DrawImage(travellerImg, rx, ry, float64(w), float64(h))
		rt := &renderer.RenderedTraveller{X: rx, Y: ry, W: float64(w), H: float64(h), Traveller: t}
		c.AddRenderedTraveller(rt)
		if c.ViewSettings.ShowLabels && c.ReverseReferences.TravellerToExpedition[t] != nil {
			expedition := c.ReverseReferences.TravellerToExpedition[t]
			if expedition != nil {
				name := expedition.Name
				if name != "" {
					dx := gui.EstimateWidth(name) * gui.FontSize / 2.0
					dy := gui.FontSize
					if c.ActiveSupplier == expedition {
						cv.SetStrokeStyle(color.RGBA{R: 0, G: 192, B: 0, A: 255})
						cv.SetLineWidth(4.0)
						cv.StrokeRect(x-dx-8, y-z-float64(h)-dy-2, dx*2+16, dy+10)
					}
					cv.SetFillStyle("texture/wood.png")
					cv.FillRect(x-dx-8, y-z-float64(h)-dy-2, dx*2+16, dy+10)
					cv.SetFillStyle("#FED")
					cv.SetFont("texture/font/Go-Regular.ttf", gui.FontSize)
					cv.FillText(name, x-dx, y-z-float64(h))
				}
			}
		}
		if t == c.SelectedTraveller {
			rt.Draw(cv)
		}
	}
}
