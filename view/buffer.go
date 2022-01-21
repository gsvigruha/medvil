package view

import (
	"github.com/tfriedel6/canvas"
	"github.com/tfriedel6/canvas/backend/goglbackend"
	"medvil/controller"
	"medvil/model/terrain"
	"medvil/renderer"
	"time"
)

type CacheEntry struct {
	offscreen *goglbackend.GoGLBackendOffscreen
	cv        *canvas.Canvas
	createdTime int64
}

type ImageCache struct {
	entries map[*terrain.Plant]*CacheEntry
	ctx     *goglbackend.GLContext
}

func NewImageCache(ctx *goglbackend.GLContext) *ImageCache {
	return &ImageCache{
		entries: make(map[*terrain.Plant]*CacheEntry),
		ctx:     ctx,
	}
}

func (ic *ImageCache) Clean() {
	t := time.Now().UnixNano()
	for k, v := range ic.entries {
		if t - v.createdTime > 1000*1000*1000 {
			v.offscreen.Delete()
			delete(ic.entries, k)
		}
	}
}

func (ic *ImageCache) RenderPlantOnBuffer(p *terrain.Plant, rf renderer.RenderedField, c *controller.Controller) *canvas.Canvas {
	t := time.Now().UnixNano()
	if ce, ok := ic.entries[p]; ok {
		if t - ce.createdTime > 300*1000*1000 {
			ce.cv.ClearRect(0, 0, 120, 300)
			RenderPlant(ce.cv, p, rf, c)
			ce.createdTime = t
		}
		return ce.cv
	} else {
		offscreen, _ := goglbackend.NewOffscreen(120, 300, true, ic.ctx)
		cv := canvas.New(offscreen)
		cv.ClearRect(0, 0, 120, 300)
		RenderPlant(cv, p, rf, c)
		ic.entries[p] = &CacheEntry{
			offscreen: offscreen,
			cv: cv,
			createdTime: t,
		}
		return cv
	}

}
