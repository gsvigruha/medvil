package view

import (
	"github.com/tfriedel6/canvas"
	"image/color"
	"medvil/controller"
)

func RenderActiveBuildingPlanBase(cv *canvas.Canvas, c *controller.Controller) {
	fields := c.GetActiveFields()
	if fields != nil {
		for _, f := range fields {
			for _, rf := range c.RenderedFields {
				if rf.F.X == f.X && rf.F.Y == f.Y {
					cv.SetStrokeStyle(color.RGBA{R: 0, G: 192, B: 0, A: 255})
					cv.SetLineWidth(2)
					rf.Draw(cv)
					cv.Stroke()
					break
				}
			}
		}
	}
}
