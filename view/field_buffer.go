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
	if f.Terrain.T == terrain.Grass {
		if c.Calendar.Season() == 3 {
			cv.SetFillStyle("texture/terrain/" + f.Terrain.T.Name + "_winter_" + strconv.Itoa(int(f.Terrain.Shape)) + ".png")
		} else if c.Calendar.Season() == 2 {
			cv.SetFillStyle("texture/terrain/" + f.Terrain.T.Name + "_fall_" + strconv.Itoa(int(f.Terrain.Shape)) + ".png")
		} else {
			cv.SetFillStyle("texture/terrain/" + f.Terrain.T.Name + "_" + strconv.Itoa(int(f.Terrain.Shape)) + ".png")
		}
	} else {
		cv.SetFillStyle("texture/terrain/" + f.Terrain.T.Name + ".png")
	}

	rf.Draw(cv)
	cv.Fill()

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
}

func (ic *FieldImageCache) RenderFieldOnBuffer(f *navigation.Field, rf renderer.RenderedField, c *controller.Controller) *canvas.Canvas {
	key := f.CacheKey() + "#" + strconv.Itoa(int(c.Perspective)) + "#" + strconv.Itoa(int(c.Calendar.Season()))
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
