package view

import (
	"github.com/tfriedel6/canvas"
	"image/color"
	"medvil/controller"
	"medvil/model/artifacts"
)

func RenderActiveFields(cv *canvas.Canvas, c *controller.Controller) {
	fields := c.GetActiveFields()
	if fields != nil {
		for _, f := range fields {
			for _, rf := range c.RenderedFields {
				if rf.F.X == f.Field().X && rf.F.Y == f.Field().Y {
					cv.SetStrokeStyle(color.RGBA{R: 0, G: 192, B: 0, A: 255})
					cv.SetLineWidth(2)
					rf.Draw(cv)
					cv.Stroke()
					if f.Context() != "" {
						midX := (rf.X[0] + rf.X[2]) / 2
						midY := (rf.Y[0] + rf.Y[2]) / 2
						a := artifacts.GetArtifact(f.Context())
						cv.DrawImage("icon/gui/artifacts/"+a.Name+".png", midX-16, midY-32, 32, 32)
					}
					break
				}
			}
		}
	}
}
