package view

import (
	"os"
	"strconv"
)

var PlantRenderBufferTimeMs = 10000
var RenderBufferTimeMs = 10000

func init() {
	if val, exists := os.LookupEnv("MEDVIL_PLANT_RENDER_BUFFER_TIME_MS"); exists {
		if time, err := strconv.Atoi(val); err == nil {
			PlantRenderBufferTimeMs = time
		}
	}
	if val, exists := os.LookupEnv("MEDVIL_RENDER_BUFFER_TIME_MS"); exists {
		if time, err := strconv.Atoi(val); err == nil {
			RenderBufferTimeMs = time
		}
	}
}
