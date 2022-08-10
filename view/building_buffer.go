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
var BuildingBufferH = DY*3 + DZ*BuildingUnitHeight
var BuildingExtensionBufferH = 200.0

type BuildingImageCache struct {
	unitEntries      map[string]*CacheEntry
	roofEntries      map[string]*CacheEntry
	extensionEntries map[string]*CacheEntry
	ctx              *goglbackend.GLContext
}

func (ic *BuildingImageCache) RenderBuildingRoofOnBuffer(
	roof *building.RoofUnit,
	rf renderer.RenderedField,
	numUnits int,
	c *controller.Controller) (*canvas.Canvas, renderer.RenderedBuildingRoof, float64, float64) {

	t := time.Now().UnixNano()
	key := roof.CacheKey() + "#" + strconv.Itoa(int(c.Perspective)) + "#" + rf.F.CacheKey()
	z := float64((numUnits+1)*BuildingUnitHeight) * DZ
	xMin, yMin, _, _ := rf.BoundingBox()
	bufferedRF := rf.Move(-xMin, -yMin+z)

	if ce, ok := ic.roofEntries[key]; ok {
		return ce.cv, RenderBuildingRoof(nil, roof, bufferedRF, numUnits, c).Move(xMin, yMin-z), xMin, yMin - z
	} else {
		offscreen, _ := goglbackend.NewOffscreen(int(BuildingBufferW), int(BuildingBufferH), true, ic.ctx)
		cv := canvas.New(offscreen)
		cv.ClearRect(0, 0, BuildingBufferW, BuildingBufferH)
		rbr := RenderBuildingRoof(cv, roof, bufferedRF, numUnits, c)
		ic.roofEntries[key] = &CacheEntry{
			offscreen:   offscreen,
			cv:          cv,
			createdTime: t,
		}
		return cv, rbr.Move(xMin, yMin-z), xMin, yMin - z
	}
}

func (ic *BuildingImageCache) RenderBuildingUnitOnBuffer(
	unit *building.BuildingUnit,
	rf renderer.RenderedField,
	numUnits int,
	c *controller.Controller) (*canvas.Canvas, renderer.RenderedBuildingUnit, float64, float64) {

	t := time.Now().UnixNano()
	key := unit.CacheKey() + "#" + strconv.Itoa(int(c.Perspective)) + "#" + rf.F.CacheKey()
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

func (ic *BuildingImageCache) RenderBuildingExtensionOnBuffer(
	extension *building.ExtensionUnit,
	rf renderer.RenderedField,
	numUnits int,
	c *controller.Controller) (*canvas.Canvas, float64, float64) {

	t := time.Now().UnixNano()
	phase := c.Calendar.Hour % BuildingAnimationMaxPhase
	key := extension.CacheKey() + "#" + strconv.Itoa(int(c.Perspective)) + "#" + strconv.Itoa(int(phase))
	z := BuildingExtensionBufferH / 2
	xMin, yMin, _, _ := rf.BoundingBox()
	bufferedRF := rf.Move(-xMin, -yMin+z)

	if ce, ok := ic.extensionEntries[key]; ok {
		return ce.cv, xMin, yMin - z
	} else {
		offscreen, _ := goglbackend.NewOffscreen(int(BuildingBufferW), int(BuildingExtensionBufferH), true, ic.ctx)
		cv := canvas.New(offscreen)
		cv.ClearRect(0, 0, BuildingBufferW, BuildingExtensionBufferH)
		RenderBuildingExtension(cv, extension, bufferedRF, numUnits, phase, c)
		ic.extensionEntries[key] = &CacheEntry{
			offscreen:   offscreen,
			cv:          cv,
			createdTime: t,
		}
		return cv, xMin, yMin - z
	}
}
