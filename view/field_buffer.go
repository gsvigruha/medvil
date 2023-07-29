package view

import (
	"github.com/tfriedel6/canvas"
	"github.com/tfriedel6/canvas/backend/goglbackend"
	"image/color"
	"medvil/controller"
	"medvil/model/navigation"
	"medvil/model/terrain"
	"medvil/renderer"
	"strconv"
	"time"
)

type FieldImageCache struct {
	entries map[string]*CacheEntry
	ctx     *goglbackend.GLContext
}

func renderField(cv *canvas.Canvas, c *controller.Controller, f *navigation.Field, rf renderer.RenderedField) {
	if f.Terrain.T == terrain.Grass || f.Terrain.T.Object {
		if c.Map.Calendar.Season() == 3 {
			cv.SetFillStyle("texture/terrain/grass_winter_" + strconv.Itoa(int(f.Terrain.Shape)) + ".png")
		} else {
			cv.SetFillStyle("texture/terrain/grass_" + strconv.Itoa(int(f.Terrain.Shape)) + ".png")
		}
	} else {
		cv.SetFillStyle("texture/terrain/" + f.Terrain.T.Name + ".png")
	}

	rf.Draw(cv)
	cv.Fill()

	if f.Terrain.T.Object {
		cv.SetFillStyle("texture/terrain/" + f.Terrain.T.Name + ".png")
		rf.Draw(cv)
		cv.Fill()
	}

	if (f.SE + f.SW) > (f.NE + f.NW) {
		slope := (f.SE + f.SW) - (f.NE + f.NW)
		cv.SetFillStyle(color.RGBA{R: 255, G: 255, B: 255, A: slope * 4})
		rf.Draw(cv)
		cv.Fill()
	} else if (f.SE + f.SW) < (f.NE + f.NW) {
		slope := (f.NE + f.NW) - (f.SE + f.SW)
		cv.SetFillStyle(color.RGBA{R: 0, G: 0, B: 0, A: slope * 16})
		rf.Draw(cv)
		cv.Fill()
	}

	for i := uint8(0); i < 4; i++ {
		idx1 := (3 - (-c.Perspective + i)) % 4
		idx2 := (2 - (-c.Perspective + i)) % 4
		idx4 := (0 - (-c.Perspective + i)) % 4
		if f.Surroundings[(i-1)%4] == navigation.SurroundingGrass {
			if c.Map.Calendar.Season() == 3 {
				cv.SetFillStyle("texture/terrain/grass_winter_" + strconv.Itoa(int(f.Terrain.Shape)) + ".png")
			} else {
				cv.SetFillStyle("texture/terrain/grass_" + strconv.Itoa(int(f.Terrain.Shape)) + ".png")
			}
		} else if f.Surroundings[(i-1)%4] == navigation.SurroundingWater {
			cv.SetFillStyle("texture/terrain/water.png")
		} else if f.Surroundings[(i-1)%4] == navigation.SurroundingDarkSlope {
			cv.SetFillStyle(color.RGBA{R: 0, G: 0, B: 0, A: 8})
		} else {
			continue
		}
		cv.BeginPath()
		cv.LineTo(rf.X[idx1], rf.Y[idx1]-rf.Z[idx1])
		cv.LineTo((rf.X[idx1]+rf.X[idx2])/2, (rf.Y[idx1]-rf.Z[idx1]+rf.Y[idx2]-rf.Z[idx2])/2)
		cv.QuadraticCurveTo(
			(4*rf.X[idx1]+rf.X[idx2]+rf.X[idx4])/6,
			(4*rf.Y[idx1]-4*rf.Z[idx1]+rf.Y[idx2]-rf.Z[idx2]+rf.Y[idx4]-rf.Z[idx4])/6,
			(rf.X[idx1]+rf.X[idx4])/2,
			(rf.Y[idx1]-rf.Z[idx1]+rf.Y[idx4]-rf.Z[idx4])/2)
		cv.ClosePath()
		cv.Fill()
	}
}

func (ic *FieldImageCache) RenderFieldOnBuffer(f *navigation.Field, rf renderer.RenderedField, c *controller.Controller) *canvas.Canvas {
	key := f.CacheKey() + "#" + strconv.Itoa(int(c.Perspective)) + "#" + strconv.Itoa(int(c.Map.Calendar.Season()))
	t := time.Now().UnixNano()
	if ce, ok := ic.entries[key]; ok {
		return ce.cv
	} else {
		xMin, yMin, xMax, yMax := rf.BoundingBox()
		bufferedRF := rf.Move(-xMin, -yMin)
		w := xMax - xMin
		h := yMax - yMin

		offscreen, _ := goglbackend.NewOffscreen(int(w), int(h), true, ic.ctx)
		cv := canvas.New(offscreen)
		cv.ClearRect(0, 0, w, h)
		renderField(cv, c, f, bufferedRF)
		ic.entries[key] = &CacheEntry{
			offscreen:   offscreen,
			cv:          cv,
			createdTime: t,
		}
		return cv
	}
}
