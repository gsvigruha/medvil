package view

import (
	"github.com/tfriedel6/canvas"
	"medvil/controller"
	"medvil/model/navigation"
	"medvil/renderer"
	"path/filepath"
	"strconv"
)

func RenderAnimal(cv *canvas.Canvas, rf renderer.RenderedField, f *navigation.Field, c *controller.Controller) {
	var phase = 0
	if c.TimeSpeed == 1 {
		phase = (int((c.Map.Calendar.Hour+(c.Map.Calendar.Day%2*24))/12) + int(f.Terrain.Shape)) % 4
	}
	maturity := float64(f.Animal.AgeYears(c.Map.Calendar)) / float64(f.Animal.T.MaturityAgeYears)
	s := 32 + maturity*32
	cv.DrawImage(filepath.FromSlash("texture/terrain/"+f.Animal.T.Name+"_"+strconv.Itoa(phase)+".png"), rf.X[0]-s/2, rf.Y[2]-64, s, s)
}
