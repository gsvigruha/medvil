package view

import (
	"github.com/tfriedel6/canvas"
	"github.com/tfriedel6/canvas/backend/goglbackend"
	"medvil/controller"
	"medvil/model/terrain"
	"medvil/renderer"
	"time"
)

const BufferW = 200
const BufferH = 300

type CacheEntry struct {
	offscreen   *goglbackend.GoGLBackendOffscreen
	cv          *canvas.Canvas
	createdTime int64
}

type ImageCache struct {
	Pic *PlantImageCache
	Fic *FieldImageCache
}

type PlantImageCache struct {
	entries map[*terrain.Plant]*CacheEntry
	ctx     *goglbackend.GLContext
}

func NewImageCache(ctx *goglbackend.GLContext) *ImageCache {
	return &ImageCache{
		Fic: &FieldImageCache{
			entries: make(map[string]*CacheEntry),
			ctx:     ctx,
		},
		Pic: &PlantImageCache{
			entries: make(map[*terrain.Plant]*CacheEntry),
			ctx:     ctx,
		},
	}
}

func (ic *ImageCache) Clean() {
	t := time.Now().UnixNano()
	for k, v := range ic.Pic.entries {
		if t-v.createdTime > 1000*1000*1000 {
			v.offscreen.Delete()
			delete(ic.Pic.entries, k)
		}
	}
}

func (ic *PlantImageCache) RenderPlantOnBuffer(p *terrain.Plant, rf renderer.RenderedField, c *controller.Controller) *canvas.Canvas {
	t := time.Now().UnixNano()
	if ce, ok := ic.entries[p]; ok {
		if t-ce.createdTime > 300*1000*1000 {
			ce.cv.ClearRect(0, 0, BufferW, BufferH)
			RenderPlant(ce.cv, p, rf, c)
			ce.createdTime = t
		}
		return ce.cv
	} else {
		offscreen, _ := goglbackend.NewOffscreen(BufferW, BufferH, true, ic.ctx)
		cv := canvas.New(offscreen)
		cv.ClearRect(0, 0, BufferW, BufferH)
		RenderPlant(cv, p, rf, c)
		ic.entries[p] = &CacheEntry{
			offscreen:   offscreen,
			cv:          cv,
			createdTime: t,
		}
		return cv
	}

}
