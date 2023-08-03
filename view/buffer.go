package view

import (
	"github.com/tfriedel6/canvas"
	"github.com/tfriedel6/canvas/backend/goglbackend"
	"log"
	"os"
	"time"
)

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
			entries: make(map[string]*CacheEntry),
			ctx:     ctx,
		},
		Bic: &BuildingImageCache{
			unitEntries:      make(map[string]*CacheEntry),
			roofEntries:      make(map[string]*CacheEntry),
			extensionEntries: make(map[string]*CacheEntry),
			ctx:              ctx,
		},
	}
}

func (ic *ImageCache) Clean() {
	t := time.Now().UnixNano()
	for k, v := range ic.Pic.entries {
		if t-v.createdTime > int64(PlantRenderBufferTimeMs)*1000*1000 {
			v.offscreen.Delete()
			delete(ic.Pic.entries, k)
		}
	}
	for k, v := range ic.Fic.entries {
		if t-v.createdTime > int64(RenderBufferTimeMs)*10000*1000 {
			v.offscreen.Delete()
			delete(ic.Fic.entries, k)
		}
	}
	for k, v := range ic.Bic.roofEntries {
		if t-v.createdTime > int64(RenderBufferTimeMs)*10000*1000 {
			v.offscreen.Delete()
			delete(ic.Bic.roofEntries, k)
		}
	}
	for k, v := range ic.Bic.unitEntries {
		if t-v.createdTime > int64(RenderBufferTimeMs)*10000*1000 {
			v.offscreen.Delete()
			delete(ic.Bic.unitEntries, k)
		}
	}
	for k, v := range ic.Bic.extensionEntries {
		if t-v.createdTime > int64(RenderBufferTimeMs)*10000*1000 {
			v.offscreen.Delete()
			delete(ic.Bic.extensionEntries, k)
		}
	}
	if os.Getenv("MEDVIL_VERBOSE") == "1" {
		log.Printf("Buffer sizes: Plant %s, Fields %s, Building units %s, Building roofs %s, Building extensions %s",
			len(ic.Pic.entries), len(ic.Fic.entries), len(ic.Bic.roofEntries), len(ic.Bic.unitEntries), len(ic.Bic.extensionEntries))
	}
}
