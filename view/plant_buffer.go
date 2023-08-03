package view

import (
	"github.com/tfriedel6/canvas"
	"github.com/tfriedel6/canvas/backend/goglbackend"
	"math/rand"
	"medvil/controller"
	"medvil/model/terrain"
	"medvil/renderer"
	"strconv"
	"time"
)

const DefaultPlantBufferW = 200
const DefaultPlantBufferH = 300

func getPlantBufferSize(p *terrain.Plant) (float64, float64) {
	if !p.IsTree() {
		return DX * 2, DY*2 + 50
	} else if p.T.TreeT == &terrain.Oak {
		return 150, 200
	} else if p.T.TreeT == &terrain.Apple {
		return 120, 150
	} else {
		return DefaultPlantBufferW, DefaultPlantBufferH
	}
}

type PlantImageCache struct {
	entries map[string]*CacheEntry
	ctx     *goglbackend.GLContext
}

func (ic *PlantImageCache) RenderPlantOnBuffer(p *terrain.Plant, rf renderer.RenderedField, c *controller.Controller) *canvas.Canvas {
	keyByShape := p.CacheKey(c.Map.Calendar) + "#" + strconv.Itoa(int(c.Map.Calendar.Month)) + "#" + strconv.Itoa(int(c.Map.Calendar.Day))
	keyByLoc := strconv.Itoa(int(rf.F.X)) + "#" + strconv.Itoa(int(rf.F.Y))
	t := time.Now().UnixNano()
	nt := t - int64(rand.Intn(PlantRenderBufferTimeMs/2)*1000*1000) + int64(PlantRenderBufferTimeMs/4*1000*1000)
	plantBufferW, plantBufferH := getPlantBufferSize(p)
	xMin, yMin, _, _ := rf.BoundingBox()
	bufferedRF := rf.Move(-xMin, -yMin+plantBufferH)
	if ce, ok := ic.entries[keyByShape]; ok {
		ic.entries[keyByLoc] = ce
		return ce.cv
	}
	if ce, ok := ic.entries[keyByLoc]; ok {
		if t-ce.createdTime > int64(PlantRenderBufferTimeMs)*1000*1000 {
			ce.cv.ClearRect(0, 0, plantBufferW, plantBufferH)
			RenderPlant(ce.cv, p, bufferedRF, c)
			ce.createdTime = nt
		}
		return ce.cv
	}
	offscreen, _ := goglbackend.NewOffscreen(int(plantBufferW), int(plantBufferH), true, ic.ctx)
	cv := canvas.New(offscreen)
	cv.ClearRect(0, 0, plantBufferW, plantBufferH)
	RenderPlant(cv, p, bufferedRF, c)
	e := &CacheEntry{
		offscreen:   offscreen,
		cv:          cv,
		createdTime: nt,
	}
	ic.entries[keyByShape] = e
	ic.entries[keyByLoc] = e
	return cv
}
