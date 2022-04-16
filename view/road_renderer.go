package view

import (
	"github.com/tfriedel6/canvas"
	"medvil/controller"
	"medvil/model/navigation"
	"medvil/renderer"
)

func RenderRoad(cv *canvas.Canvas, rf renderer.RenderedField, f *navigation.Field, c *controller.Controller) {
	if f.Construction || f.Road.Construction {
		cv.SetFillStyle("texture/infra/construction.png")
	} else {
		cv.SetFillStyle("texture/infra/" + f.Road.T.Name + ".png")
	}
	cv.BeginPath()
	for i := uint8(0); i < 4; i++ {
		idx1 := (3 - (-c.Perspective + i)) % 4
		idx2 := (2 - (-c.Perspective + i)) % 4
		idx3 := (1 - (-c.Perspective + i)) % 4
		idx4 := (0 - (-c.Perspective + i)) % 4
		if !f.Construction {
			leftEdge := f.Road.EdgeConnections[(i-1)%4]
			corner := f.Road.CornerConnections[(i-1)%4]
			rightEdge := f.Road.EdgeConnections[i]
			if leftEdge {
				cv.LineTo((rf.X[idx1]*7+rf.X[idx4])/8, (rf.Y[idx1]*7+rf.Y[idx4])/8-(rf.Z[idx1]*7+rf.Z[idx4])/8)
			}
			if leftEdge && corner && rightEdge {
				cv.LineTo(rf.X[idx1], rf.Y[idx1]-rf.Z[idx1])
			} else {
				cv.LineTo((rf.X[idx1]*7+rf.X[idx3])/8, (rf.Y[idx1]*7+rf.Y[idx3])/8-(rf.Z[idx1]*7+rf.Z[idx3])/8)
			}
			if rightEdge {
				cv.LineTo((rf.X[idx1]*7+rf.X[idx2])/8, (rf.Y[idx1]*7+rf.Y[idx2])/8-(rf.Z[idx1]*7+rf.Z[idx2])/8)
			}
		} else {
			cv.LineTo(rf.X[idx1], rf.Y[idx1]-rf.Z[idx1])
		}
	}
	cv.ClosePath()
	cv.Fill()
}
