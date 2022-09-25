package buildings

import (
	"github.com/tfriedel6/canvas"
	"image/color"
	"math"
	"medvil/controller"
	"medvil/model/building"
	"medvil/model/materials"
	"medvil/renderer"
	"strconv"
)

func RenderBuildingRoof(cv *canvas.Canvas, roof *building.RoofUnit, rf renderer.RenderedField, k int, c *controller.Controller) *renderer.RenderedBuildingRoof {
	if roof == nil {
		return nil
	}
	var roofPolygons []renderer.Polygon
	startL := 2 + c.Perspective
	if roof.Roof.RoofType == building.RoofTypeSplit {
		z := math.Min(math.Min(math.Min(rf.Z[0], rf.Z[1]), rf.Z[2]), rf.Z[3]) + float64(k*BuildingUnitHeight)*DZ
		midX := (rf.X[0] + rf.X[2]) / 2
		midY := (rf.Y[0] + rf.Y[2]) / 2

		for l := uint8(startL); l < 4+startL; l++ {
			rfIdx1 := (3 - (-c.Perspective + l)) % 4
			rfIdx2 := (2 - (-c.Perspective + l)) % 4
			if roof.Elevated[l%4] {
				var suffix = ""
				if rfIdx1%2 == 0 {
					suffix = "_flipped"
				}

				sideMidX := (rf.X[rfIdx1] + rf.X[rfIdx2]) / 2
				sideMidY := (rf.Y[rfIdx1] + rf.Y[rfIdx2]) / 2
				rp1 := renderer.Polygon{Points: []renderer.Point{
					renderer.Point{X: rf.X[rfIdx1], Y: rf.Y[rfIdx1] - z},
					renderer.Point{X: sideMidX, Y: sideMidY - z - BuildingUnitHeight*DZ},
					renderer.Point{X: midX, Y: midY - z - BuildingUnitHeight*DZ},
				}}
				rp2 := renderer.Polygon{Points: []renderer.Point{
					renderer.Point{X: rf.X[rfIdx2], Y: rf.Y[rfIdx2] - z},
					renderer.Point{X: sideMidX, Y: sideMidY - z - BuildingUnitHeight*DZ},
					renderer.Point{X: midX, Y: midY - z - BuildingUnitHeight*DZ},
				}}
				roofPolygons = append(roofPolygons, rp1, rp2)

				if cv != nil {
					if !roof.Construction {
						cv.SetFillStyle("texture/building/" + RoofMaterialName(roof.Roof.M, roof.B.Shape) + suffix + ".png")
					} else {
						cv.SetFillStyle("texture/building/construction" + suffix + ".png")
					}

					cv.SetStrokeStyle(color.RGBA{R: 192, G: 128, B: 64, A: 32})
					cv.SetLineWidth(3)

					RenderPolygon(cv, rp1, true)
					RenderPolygon(cv, rp2, true)
				}
			} else {
				var suffix = ""
				if rfIdx1%2 == 1 {
					suffix = "_flipped"
				}

				rp1 := renderer.Polygon{Points: []renderer.Point{
					renderer.Point{X: rf.X[rfIdx1], Y: rf.Y[rfIdx1] - z},
					renderer.Point{X: rf.X[rfIdx2], Y: rf.Y[rfIdx2] - z},
					renderer.Point{X: midX, Y: midY - z - BuildingUnitHeight*DZ},
				}}
				roofPolygons = append(roofPolygons, rp1)

				if cv != nil {
					if !roof.Construction {
						cv.SetFillStyle("texture/building/" + RoofMaterialName(roof.Roof.M, roof.B.Shape) + suffix + ".png")
					} else {
						cv.SetFillStyle("texture/building/construction" + suffix + ".png")
					}
					cv.SetStrokeStyle(color.RGBA{R: 64, G: 32, B: 0, A: 32})
					cv.SetLineWidth(3)
					RenderPolygon(cv, rp1, true)
				}
			}
		}
	} else if roof.Roof.RoofType == building.RoofTypeFlat {
		z := float64(k*BuildingUnitHeight) * DZ
		if !roof.Construction {
			rp1 := renderer.Polygon{Points: []renderer.Point{
				renderer.Point{X: rf.X[0], Y: rf.Y[0] - rf.Z[0] - z},
				renderer.Point{X: rf.X[1], Y: rf.Y[1] - rf.Z[1] - z},
				renderer.Point{X: rf.X[2], Y: rf.Y[2] - rf.Z[2] - z},
				renderer.Point{X: rf.X[3], Y: rf.Y[3] - rf.Z[3] - z},
			}}
			roofPolygons = append(roofPolygons, rp1)
			if cv != nil {
				cv.SetFillStyle("texture/building/" + RoofMaterialName(roof.Roof.M, roof.B.Shape) + "_flat.png")
				RenderPolygon(cv, rp1, false)
			}
		}
	} else if roof.Roof.RoofType == building.RoofTypeRamp {
		z := math.Min(math.Min(math.Min(rf.Z[0], rf.Z[1]), rf.Z[2]), rf.Z[3]) + float64(k*BuildingUnitHeight)*DZ
		for l := uint8(startL); l < 4+startL; l++ {
			rfIdx1 := (3 - (-c.Perspective + l)) % 4
			rfIdx2 := (2 - (-c.Perspective + l)) % 4
			rfIdx3 := (1 - (-c.Perspective + l)) % 4
			rfIdx4 := (0 - (-c.Perspective + l)) % 4
			if roof.Elevated[l%4] {
				var suffix = ""
				if rfIdx1%2 == 0 {
					suffix = "_flipped"
				}
				rp1 := renderer.Polygon{Points: []renderer.Point{
					renderer.Point{X: rf.X[rfIdx1], Y: rf.Y[rfIdx1] - z - BuildingUnitHeight*DZ},
					renderer.Point{X: rf.X[rfIdx1], Y: rf.Y[rfIdx1] - z},
					renderer.Point{X: rf.X[rfIdx4], Y: rf.Y[rfIdx4] - z},
				}}
				rp2 := renderer.Polygon{Points: []renderer.Point{
					renderer.Point{X: rf.X[rfIdx2], Y: rf.Y[rfIdx2] - z - BuildingUnitHeight*DZ},
					renderer.Point{X: rf.X[rfIdx2], Y: rf.Y[rfIdx2] - z},
					renderer.Point{X: rf.X[rfIdx3], Y: rf.Y[rfIdx3] - z},
				}}
				rp3 := renderer.Polygon{Points: []renderer.Point{
					renderer.Point{X: rf.X[rfIdx1], Y: rf.Y[rfIdx1] - z - BuildingUnitHeight*DZ},
					renderer.Point{X: rf.X[rfIdx2], Y: rf.Y[rfIdx2] - z - BuildingUnitHeight*DZ},
					renderer.Point{X: rf.X[rfIdx3], Y: rf.Y[rfIdx3] - z},
					renderer.Point{X: rf.X[rfIdx4], Y: rf.Y[rfIdx4] - z},
				}}
				roofPolygons = append(roofPolygons, rp1, rp2, rp3)

				if cv != nil {
					if !roof.Construction {
						cv.SetFillStyle("texture/building/" + RoofMaterialName(roof.Roof.M, roof.B.Shape) + suffix + ".png")
					} else {
						cv.SetFillStyle("texture/building/construction" + suffix + ".png")
					}

					cv.SetStrokeStyle(color.RGBA{R: 192, G: 128, B: 64, A: 32})
					cv.SetLineWidth(3)

					RenderPolygon(cv, rp1, true)
					RenderPolygon(cv, rp2, true)

					if !roof.Construction {
						cv.SetFillStyle("texture/building/" + RoofMaterialName(roof.Roof.M, roof.B.Shape) + "_flat.png")
					} else {
						cv.SetFillStyle("texture/building/construction" + suffix + ".png")
					}

					RenderPolygon(cv, rp3, true)
				}
			}
		}
	}
	return &renderer.RenderedBuildingRoof{B: roof.Building(), Ps: roofPolygons}
}
