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

var BuildingBufferW = DX * 2
var BuildingBufferH = DY*2 + DZ*BuildingUnitHeight

type BuildingImageCache struct {
	unitEntries map[string]*CacheEntry
	roofEntries map[string]*CacheEntry
	ctx         *goglbackend.GLContext
}

func (ic *BuildingImageCache) RenderBuildingRoofOnBuffer(
	roof *building.RoofUnit,
	rf renderer.RenderedField,
	numUnits int,
	c *controller.Controller) (*canvas.Canvas, float64, float64) {

	t := time.Now().UnixNano()
	key := roof.CacheKey() + "#" + strconv.Itoa(int(c.Perspective))
	z := float64((numUnits+1)*BuildingUnitHeight) * DZ
	xMin, yMin, _, _ := rf.BoundingBox()
	bufferedRF := rf.Move(-xMin, -yMin+z)

	if ce, ok := ic.roofEntries[key]; ok {
		return ce.cv, xMin, yMin - z
	} else {
		offscreen, _ := goglbackend.NewOffscreen(int(BuildingBufferW), int(BuildingBufferH), true, ic.ctx)
		cv := canvas.New(offscreen)
		cv.ClearRect(0, 0, BuildingBufferW, BuildingBufferH)
		RenderBuildingRoof(cv, roof, bufferedRF, numUnits, c)
		ic.roofEntries[key] = &CacheEntry{
			offscreen:   offscreen,
			cv:          cv,
			createdTime: t,
		}
		return cv, xMin, yMin - z
	}
}

func (ic *BuildingImageCache) RenderBuildingUnitOnBuffer(
	unit *building.BuildingUnit,
	rf renderer.RenderedField,
	numUnits int,
	c *controller.Controller) (*canvas.Canvas, renderer.RenderedBuildingUnit, float64, float64) {

	t := time.Now().UnixNano()
	key := unit.CacheKey() + "#" + strconv.Itoa(int(c.Perspective))
	z := float64((numUnits+1)*BuildingUnitHeight) * DZ
	xMin, yMin, _, _ := rf.BoundingBox()
	bufferedRF := rf.Move(-xMin, -yMin+z)

	if ce, ok := ic.unitEntries[key]; ok {
		return ce.cv, RenderBuildingUnit(nil, unit, bufferedRF, numUnits, c).Move(xMin, yMin-z), xMin, yMin - z
	} else {
		offscreen, _ := goglbackend.NewOffscreen(int(BuildingBufferW), int(BuildingBufferH), true, ic.ctx)
		cv := canvas.New(offscreen)
		cv.ClearRect(0, 0, BuildingBufferW, BuildingBufferH)
		rbu := RenderBuildingUnit(cv, unit, bufferedRF, numUnits, c)
		ic.unitEntries[key] = &CacheEntry{
			offscreen:   offscreen,
			cv:          cv,
			createdTime: t,
		}
		return cv, rbu.Move(xMin, yMin-z), xMin, yMin - z
	}
}
