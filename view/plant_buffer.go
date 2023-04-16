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

const DefaultPlantBufferW = 200
const DefaultPlantBufferH = 300

func getPlantBufferSize(p *terrain.Plant) (float64, float64) {
	if p.T.TreeT == &terrain.Oak {
		return 150, 200
	} else if p.T.TreeT == &terrain.Apple {
		return 120, 150
	} else {
		return DefaultPlantBufferW, DefaultPlantBufferH
	}
}

type PlantImageCache struct {
	entries map[*terrain.Plant]*CacheEntry
	ctx     *goglbackend.GLContext
}

func (ic *PlantImageCache) RenderPlantOnBuffer(p *terrain.Plant, rf renderer.RenderedField, c *controller.Controller) *canvas.Canvas {
	t := time.Now().UnixNano()
	plantBufferW, plantBufferH := getPlantBufferSize(p)
	if ce, ok := ic.entries[p]; ok {
		if t-ce.createdTime > int64(PlantRenderBufferTimeMs)*1000*1000 {
			ce.cv.ClearRect(0, 0, plantBufferW, plantBufferH)
			RenderPlant(ce.cv, p, rf, c)
			ce.createdTime = t - int64(rand.Intn(PlantRenderBufferTimeMs/2)*1000*1000) + int64(PlantRenderBufferTimeMs/4*1000*1000)
		}
		return ce.cv
	} else {
		offscreen, _ := goglbackend.NewOffscreen(int(plantBufferW), int(plantBufferH), true, ic.ctx)
		cv := canvas.New(offscreen)
		cv.ClearRect(0, 0, plantBufferW, plantBufferH)
		RenderPlant(cv, p, rf, c)
		ic.entries[p] = &CacheEntry{
			offscreen:   offscreen,
			cv:          cv,
			createdTime: t - int64(rand.Intn(PlantRenderBufferTimeMs/2)*1000*1000) + int64(PlantRenderBufferTimeMs/4*1000*1000),
		}
		return cv
	}
}
