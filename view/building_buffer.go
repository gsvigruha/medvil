package view

import (
	"github.com/tfriedel6/canvas"
	"github.com/tfriedel6/canvas/backend/goglbackend"
	"medvil/controller"
	"medvil/model/building"
	"medvil/renderer"
	"time"
)

type BuildingImageCache struct {
	unitEntries map[*building.BuildingUnit]*CacheEntry
	roofEntries map[*building.RoofUnit]*CacheEntry
	ctx         *goglbackend.GLContext
}

func (ic *BuildingImageCache) RenderBuildingRoofOnBuffer(roof *building.RoofUnit, rf renderer.RenderedField, numUnits int, c *controller.Controller) *canvas.Canvas {
	t := time.Now().UnixNano()
	if ce, ok := ic.roofEntries[roof]; ok {
		return ce.cv
	} else {
		z := float64((numUnits+1)*BuildingUnitHeight)*DZ
		xMin, yMin, _, _ := rf.BoundingBox()
		bufferedRF := rf.Move(-xMin, -yMin+z)
		offscreen, _ := goglbackend.NewOffscreen(120, 125, true, ic.ctx)
		cv := canvas.New(offscreen)
		cv.ClearRect(0, 0, 120, 125)
		RenderBuildingRoof(cv, roof, bufferedRF, numUnits, c)
		ic.roofEntries[roof] = &CacheEntry{
			offscreen:   offscreen,
			cv:          cv,
			createdTime: t,
		}
		return cv
	}
}
