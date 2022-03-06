package view

import (
	"github.com/tfriedel6/canvas"
	"github.com/tfriedel6/canvas/backend/goglbackend"
	"medvil/controller"
	"medvil/model/building"
	"medvil/renderer"
	"strconv"
	"time"
)

type BuildingImageCache struct {
	unitEntries map[string]*CacheEntry
	roofEntries map[string]*CacheEntry
	ctx         *goglbackend.GLContext
}

func (ic *BuildingImageCache) RenderBuildingRoofOnBuffer(roof *building.RoofUnit, rf renderer.RenderedField, numUnits int, c *controller.Controller) *canvas.Canvas {
	t := time.Now().UnixNano()
	key := roof.CacheKey() + "#" + strconv.Itoa(int(c.Perspective))
	if ce, ok := ic.roofEntries[key]; ok {
		return ce.cv
	} else {
		z := float64((numUnits+1)*BuildingUnitHeight) * DZ
		xMin, yMin, _, _ := rf.BoundingBox()
		bufferedRF := rf.Move(-xMin, -yMin+z)
		offscreen, _ := goglbackend.NewOffscreen(120, 125, true, ic.ctx)
		cv := canvas.New(offscreen)
		cv.ClearRect(0, 0, 120, 125)
		RenderBuildingRoof(cv, roof, bufferedRF, numUnits, c)
		ic.roofEntries[key] = &CacheEntry{
			offscreen:   offscreen,
			cv:          cv,
			createdTime: t,
		}
		return cv
	}
}

func (ic *BuildingImageCache) RenderBuildingUnitOnBuffer(
	unit *building.BuildingUnit,
	rf renderer.RenderedField,
	numUnits int,
	c *controller.Controller) (*canvas.Canvas, renderer.RenderedBuildingUnit) {

	t := time.Now().UnixNano()
	key := unit.CacheKey() + "#" + strconv.Itoa(int(c.Perspective))
	z := float64((numUnits+1)*BuildingUnitHeight) * DZ
	xMin, yMin, _, _ := rf.BoundingBox()
	bufferedRF := rf.Move(-xMin, -yMin+z)
	if ce, ok := ic.unitEntries[key]; ok {
		return ce.cv, RenderBuildingUnit(nil, unit, bufferedRF, numUnits, c)
	} else {
		offscreen, _ := goglbackend.NewOffscreen(120, 180, true, ic.ctx)
		cv := canvas.New(offscreen)
		cv.ClearRect(0, 0, 120, 180)
		rbu := RenderBuildingUnit(cv, unit, bufferedRF, numUnits, c)
		ic.unitEntries[key] = &CacheEntry{
			offscreen:   offscreen,
			cv:          cv,
			createdTime: t,
		}
		return cv, rbu
	}
}
