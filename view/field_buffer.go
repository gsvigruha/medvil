package view

import (
	"github.com/tfriedel6/canvas"
	"github.com/tfriedel6/canvas/backend/goglbackend"
	"image/color"
	"medvil/model/navigation"
	"medvil/renderer"
	"time"
)

type FieldImageCache struct {
	entries map[string]*CacheEntry
	ctx     *goglbackend.GLContext
}

func NewFieldImageCache(ctx *goglbackend.GLContext) *FieldImageCache {
	return &FieldImageCache{
		entries: make(map[string]*CacheEntry),
		ctx:     ctx,
	}
}

func (ic *FieldImageCache) Clean() {
	t := time.Now().UnixNano()
	for k, v := range ic.entries {
		if t-v.createdTime > 1000*1000*1000 {
			v.offscreen.Delete()
			delete(ic.entries, k)
		}
	}
}

func renderField(cv *canvas.Canvas, f *navigation.Field, rf renderer.RenderedField) {
	offsetX, offsetY := rf.Offset()
	rf2 := rf.Move(-offsetX, -offsetY)
	cv.SetFillStyle("texture/terrain/" + f.Terrain.T.Name + ".png")

	rf2.Draw(cv)
	cv.Fill()

	if (f.SE + f.SW) > (f.NE + f.NW) {
		slope := (f.SE + f.SW) - (f.NE + f.NW)
		cv.SetFillStyle(color.RGBA{R: 255, G: 255, B: 255, A: slope * 4})
		rf2.Draw(cv)
		cv.Fill()
	} else if (f.SE + f.SW) < (f.NE + f.NW) {
		slope := (f.NE + f.NW) - (f.SE + f.SW)
		cv.SetFillStyle(color.RGBA{R: 0, G: 0, B: 0, A: slope * 16})
		rf2.Draw(cv)
		cv.Fill()
	}
}

func (ic *FieldImageCache) RenderFieldOnBuffer(f *navigation.Field, rf renderer.RenderedField) *canvas.Canvas {
	t := time.Now().UnixNano()
	if ce, ok := ic.entries[f.CacheKey()]; ok {
		if t-ce.createdTime > 300*1000*1000 {
			ce.cv.ClearRect(0, 0, BufferW, BufferH)
			renderField(ce.cv, f, rf)
			ce.createdTime = t
		}
		return ce.cv
	} else {
		offscreen, _ := goglbackend.NewOffscreen(BufferW, BufferH, true, ic.ctx)
		cv := canvas.New(offscreen)
		cv.ClearRect(0, 0, BufferW, BufferH)
		renderField(cv, f, rf)
		ic.entries[f.CacheKey()] = &CacheEntry{
			offscreen:   offscreen,
			cv:          cv,
			createdTime: t,
		}
		return cv
	}

}
