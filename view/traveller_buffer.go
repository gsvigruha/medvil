package view

import (
	"github.com/tfriedel6/canvas"
	"github.com/tfriedel6/canvas/backend/goglbackend"
	"medvil/controller"
	"medvil/model/building"
	"medvil/model/navigation"
	"medvil/renderer"
	"medvil/view/buildings"
	"strconv"
	"time"
)

type TravellerImageCache struct {
	entries map[string]*CacheEntry
	ctx     *goglbackend.GLContext
}

func (tc *TravellerImageCache) RenderTravellerOnBuffer(
	cv *canvas.Canvas, t *navigation.Traveller, f *navigation.Field, w, h int, c *controller.Controller) *canvas.Canvas {
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
		return 48, 48
	} else if t == navigation.TravellerTypeTradingCart {
		return 48, 48
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
		travellerImg := ic.Tic.RenderTravellerOnBuffer(cv, t, rf.F, w, h, c)
		cv.DrawImage(travellerImg, x-float64(w/2), y-z-float64(h*3/4), float64(w), float64(h))
		c.AddRenderedTraveller(&renderer.RenderedTraveller{X: x, Y: y - z + float64(h/4), H: float64(h), W: float64(w), Traveller: t})
	}
}
