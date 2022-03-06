package view

import (
	"github.com/tfriedel6/canvas"
	"github.com/tfriedel6/canvas/backend/goglbackend"
	"math/rand"
	"medvil/controller"
	"medvil/model/terrain"
	"medvil/renderer"
	"time"
)

const PlantBufferW = 200
const PlantBufferH = 300

type PlantImageCache struct {
	entries map[*terrain.Plant]*CacheEntry
	ctx     *goglbackend.GLContext
}

func (ic *PlantImageCache) RenderPlantOnBuffer(p *terrain.Plant, rf renderer.RenderedField, c *controller.Controller) *canvas.Canvas {
	t := time.Now().UnixNano()
	if ce, ok := ic.entries[p]; ok {
		if t-ce.createdTime > 1000*1000*1000 {
			ce.cv.ClearRect(0, 0, PlantBufferW, PlantBufferH)
			RenderPlant(ce.cv, p, rf, c)
			ce.createdTime = t - int64(rand.Intn(500)*1000*1000) + int64(250*1000*1000)
		}
		return ce.cv
	} else {
		offscreen, _ := goglbackend.NewOffscreen(PlantBufferW, PlantBufferH, true, ic.ctx)
		cv := canvas.New(offscreen)
		cv.ClearRect(0, 0, PlantBufferW, PlantBufferH)
		RenderPlant(cv, p, rf, c)
		ic.entries[p] = &CacheEntry{
			offscreen:   offscreen,
			cv:          cv,
			createdTime: t - int64(rand.Intn(500)*1000*1000) + int64(250*1000*1000),
		}
		return cv
	}
}
