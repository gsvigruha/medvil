package view

import (
	"github.com/tfriedel6/canvas"
	"github.com/tfriedel6/canvas/backend/goglbackend"
	"medvil/model/terrain"
	"time"
)

const MaxCacheDeletePerIteration = 10

type CacheEntry struct {
	offscreen   *goglbackend.GoGLBackendOffscreen
	cv          *canvas.Canvas
	createdTime int64
}

type ImageCache struct {
	Pic *PlantImageCache
	Fic *FieldImageCache
	Bic *BuildingImageCache
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
		Bic: &BuildingImageCache{
			unitEntries: make(map[string]*CacheEntry),
			roofEntries: make(map[string]*CacheEntry),
			ctx:         ctx,
		},
	}
}

func (ic *ImageCache) Clean() {
	t := time.Now().UnixNano()
	var i = 0
	for k, v := range ic.Pic.entries {
		if i < MaxCacheDeletePerIteration && t-v.createdTime > 1000*1000*1000 {
			v.offscreen.Delete()
			delete(ic.Pic.entries, k)
			i++
		}
	}
	for k, v := range ic.Fic.entries {
		if i < MaxCacheDeletePerIteration && t-v.createdTime > 10000*1000*1000 {
			v.offscreen.Delete()
			delete(ic.Fic.entries, k)
			i++
		}
	}
	for k, v := range ic.Bic.roofEntries {
		if i < MaxCacheDeletePerIteration && t-v.createdTime > 10000*1000*1000 {
			v.offscreen.Delete()
			delete(ic.Bic.roofEntries, k)
			i++
		}
	}
	for k, v := range ic.Bic.unitEntries {
		if i < MaxCacheDeletePerIteration && t-v.createdTime > 10000*1000*1000 {
			v.offscreen.Delete()
			delete(ic.Bic.unitEntries, k)
			i++
		}
	}
}
