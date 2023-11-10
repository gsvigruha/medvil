package buildings

import (
	"github.com/tfriedel6/canvas"
	"image/color"
	"math"
	"medvil/controller"
	"medvil/model/building"
	"medvil/model/materials"
	"medvil/renderer"
	"medvil/view/util"
	"path/filepath"
	"strconv"
)

func RoofMaterialName(r *building.RoofUnit) string {
	m := r.Roof.M
	shape := r.B.Shape
	broken := r.B.Broken
	if m == materials.GetMaterial("tile") {
		switch shape % 2 {
		case 0:
			return "tile_red"
		case 1:
			return "tile_darkred"
		}
	}
	if m == materials.GetMaterial("stone") && broken {
		return "stone_broken"
	}
	if m == materials.GetMaterial("brick") {
		return BrickMaterialName(shape)
	}
	if r.B.Plan.BuildingType == building.BuildingTypeMine && m == materials.GetMaterial("reed") {
		return "reed_old"
	}
	return m.Name
}

func RenderBuildingRoof(cv *canvas.Canvas, roof *building.RoofUnit, rf renderer.RenderedField, k int, c *controller.Controller) *renderer.RenderedBuildingRoof {
	if roof == nil {
		return nil
	}
	var roofPolygons []renderer.Polygon
	startL := 2 + c.Perspective
	if roof.Roof.RoofType == building.RoofTypeSplit {
		z := math.Min(math.Min(math.Min(rf.Z[0], rf.Z[1]), rf.Z[2]), rf.Z[3]) + float64(k*BuildingUnitHeight)*util.DZ
		midX := (rf.X[0] + rf.X[2]) / 2
		midY := (rf.Y[0] + rf.Y[2]) / 2

		for l := uint8(startL); l < 4+startL; l++ {
			rfIdx1 := (3 - (-c.Perspective + l)) % 4
			rfIdx2 := (2 - (-c.Perspective + l)) % 4
			if roof.Connected[l%4] {
				var suffix = ""
				if rfIdx1%2 == 0 {
					suffix = "_flipped"
				}

				sideMidX := (rf.X[rfIdx1] + rf.X[rfIdx2]) / 2
				sideMidY := (rf.Y[rfIdx1] + rf.Y[rfIdx2]) / 2
				rp1 := renderer.Polygon{Points: []renderer.Point{
					renderer.Point{X: rf.X[rfIdx1], Y: rf.Y[rfIdx1] - z},
					renderer.Point{X: sideMidX, Y: sideMidY - z - BuildingUnitHeight*util.DZ},
					renderer.Point{X: midX, Y: midY - z - BuildingUnitHeight*util.DZ},
				}}
				rp2 := renderer.Polygon{Points: []renderer.Point{
					renderer.Point{X: rf.X[rfIdx2], Y: rf.Y[rfIdx2] - z},
					renderer.Point{X: sideMidX, Y: sideMidY - z - BuildingUnitHeight*util.DZ},
					renderer.Point{X: midX, Y: midY - z - BuildingUnitHeight*util.DZ},
				}}
				roofPolygons = append(roofPolygons, rp1, rp2)

				if cv != nil {
					if !roof.Construction {
						cv.SetFillStyle(filepath.FromSlash("texture/building/" + RoofMaterialName(roof) + suffix + ".png"))
					} else {
						cv.SetFillStyle(filepath.FromSlash("texture/building/construction" + suffix + ".png"))
					}

					cv.SetStrokeStyle(color.RGBA{R: 192, G: 128, B: 64, A: 32})
					cv.SetLineWidth(3)

					util.RenderPolygon(cv, rp1, true)
					util.RenderPolygon(cv, rp2, true)
				}
			} else {
				var suffix = ""
				if rfIdx1%2 == 1 {
					suffix = "_flipped"
				}

				rp1 := renderer.Polygon{Points: []renderer.Point{
					renderer.Point{X: rf.X[rfIdx1], Y: rf.Y[rfIdx1] - z},
					renderer.Point{X: rf.X[rfIdx2], Y: rf.Y[rfIdx2] - z},
					renderer.Point{X: midX, Y: midY - z - BuildingUnitHeight*util.DZ},
				}}
				roofPolygons = append(roofPolygons, rp1)

				if cv != nil {
					if !roof.Construction {
						cv.SetFillStyle(filepath.FromSlash("texture/building/" + RoofMaterialName(roof) + suffix + ".png"))
					} else {
						cv.SetFillStyle(filepath.FromSlash("texture/building/construction" + suffix + ".png"))
					}
					cv.SetStrokeStyle(color.RGBA{R: 64, G: 32, B: 0, A: 32})
					cv.SetLineWidth(3)
					util.RenderPolygon(cv, rp1, true)
				}
			}
		}
	} else if roof.Roof.RoofType == building.RoofTypeFlat {
		z := float64(k*BuildingUnitHeight) * util.DZ
		if !roof.Construction {
			rp1 := renderer.Polygon{Points: []renderer.Point{
				renderer.Point{X: rf.X[0], Y: rf.Y[0] - rf.Z[0] - z},
				renderer.Point{X: rf.X[1], Y: rf.Y[1] - rf.Z[1] - z},
				renderer.Point{X: rf.X[2], Y: rf.Y[2] - rf.Z[2] - z},
				renderer.Point{X: rf.X[3], Y: rf.Y[3] - rf.Z[3] - z},
			}}
			roofPolygons = append(roofPolygons, rp1)
			if cv != nil {
				cv.SetFillStyle(filepath.FromSlash("texture/building/" + RoofMaterialName(roof) + "_flat.png"))
				util.RenderPolygon(cv, rp1, false)
				RenderRoofFence(cv, roof, rp1, c)

				if roof.B.Plan.BuildingType == building.BuildingTypeWorkshop || roof.B.Plan.BuildingType == building.BuildingTypeFactory {
					dirIdx0 := (c.Perspective + 0) % 4
					dirIdx1 := (c.Perspective + 1) % 4
					dirIdx2 := (c.Perspective + 2) % 4
					dirIdx3 := (c.Perspective + 3) % 4

					if roof.B.Shape%3 == 0 {
						px := (rf.X[dirIdx0]*3.0 + rf.X[dirIdx2]*1.0) / 4.0
						py := ((rf.Y[dirIdx0]-rf.Z[dirIdx0]-z)*3.0 + (rf.Y[dirIdx2]-rf.Z[dirIdx2]-z)*1.0) / 4.0
						cv.DrawImage(filepath.FromSlash("texture/building/plant_1.png"), px-16, py-32)
					} else if roof.B.Shape%3 == 1 {
						px := (rf.X[dirIdx3]*3.0 + rf.X[dirIdx1]*1.0) / 4.0
						py := ((rf.Y[dirIdx3]-rf.Z[dirIdx3]-z)*3.0 + (rf.Y[dirIdx1]-rf.Z[dirIdx1]-z)*1.0) / 4.0
						cv.DrawImage(filepath.FromSlash("texture/building/plant_2.png"), px-16, py-32)
					}
				}
			}
		}
	} else if roof.Roof.RoofType == building.RoofTypeRamp {
		z := math.Min(math.Min(math.Min(rf.Z[0], rf.Z[1]), rf.Z[2]), rf.Z[3]) + float64(k*BuildingUnitHeight)*util.DZ
		for l := uint8(startL); l < 4+startL; l++ {
			rfIdx1 := (3 - (-c.Perspective + l)) % 4
			rfIdx2 := (2 - (-c.Perspective + l)) % 4
			rfIdx3 := (1 - (-c.Perspective + l)) % 4
			rfIdx4 := (0 - (-c.Perspective + l)) % 4
			if roof.Connected[l%4] {
				var suffix = ""
				if rfIdx1%2 == 0 {
					suffix = "_flipped"
				}
				rp1 := renderer.Polygon{Points: []renderer.Point{
					renderer.Point{X: rf.X[rfIdx1], Y: rf.Y[rfIdx1] - z - BuildingUnitHeight*util.DZ},
					renderer.Point{X: rf.X[rfIdx1], Y: rf.Y[rfIdx1] - z},
					renderer.Point{X: rf.X[rfIdx4], Y: rf.Y[rfIdx4] - z},
				}}
				rp2 := renderer.Polygon{Points: []renderer.Point{
					renderer.Point{X: rf.X[rfIdx2], Y: rf.Y[rfIdx2] - z - BuildingUnitHeight*util.DZ},
					renderer.Point{X: rf.X[rfIdx2], Y: rf.Y[rfIdx2] - z},
					renderer.Point{X: rf.X[rfIdx3], Y: rf.Y[rfIdx3] - z},
				}}
				rp3 := renderer.Polygon{Points: []renderer.Point{
					renderer.Point{X: rf.X[rfIdx1], Y: rf.Y[rfIdx1] - z - BuildingUnitHeight*util.DZ},
					renderer.Point{X: rf.X[rfIdx2], Y: rf.Y[rfIdx2] - z - BuildingUnitHeight*util.DZ},
					renderer.Point{X: rf.X[rfIdx3], Y: rf.Y[rfIdx3] - z},
					renderer.Point{X: rf.X[rfIdx4], Y: rf.Y[rfIdx4] - z},
				}}
				roofPolygons = append(roofPolygons, rp1, rp2, rp3)

				if cv != nil {
					if !roof.Construction {
						cv.SetFillStyle(filepath.FromSlash("texture/building/" + RoofMaterialName(roof) + suffix + ".png"))
					} else {
						cv.SetFillStyle(filepath.FromSlash("texture/building/construction" + suffix + ".png"))
					}

					cv.SetStrokeStyle(color.RGBA{R: 192, G: 128, B: 64, A: 32})
					cv.SetLineWidth(3)

					util.RenderPolygon(cv, rp1, true)
					util.RenderPolygon(cv, rp2, true)

					if !roof.Construction {
						cv.SetFillStyle(filepath.FromSlash("texture/building/" + RoofMaterialName(roof) + "_flat.png"))
					} else {
						cv.SetFillStyle(filepath.FromSlash("texture/building/construction" + suffix + ".png"))
					}

					util.RenderPolygon(cv, rp3, true)
				}
			}
		}
	}
	if !roof.Construction && roof.B.HasExtension(building.Cooker) && cv != nil {
		RenderChimney(cv, rf, k, roof.Roof.RoofType == building.RoofTypeFlat, 255)
	}
	return &renderer.RenderedBuildingRoof{B: roof.Building(), Ps: roofPolygons}
}

func RenderChimney(cv *canvas.Canvas, rf renderer.RenderedField, k int, flatRoof bool, phase uint8) {
	z := math.Min(math.Min(math.Min(rf.Z[0], rf.Z[1]), rf.Z[2]), rf.Z[3]) + float64(k*BuildingUnitHeight)*DZ
	if flatRoof {
		z -= float64(BuildingUnitHeight) * DZ
	}
	midX := (rf.X[0] + rf.X[2]) / 2
	midY := (rf.Y[0] + rf.Y[2]) / 2
	h := 8.0
	rp1 := renderer.Polygon{Points: []renderer.Point{
		renderer.Point{X: midX, Y: midY - z - BuildingUnitHeight*DZ + 12},
		renderer.Point{X: midX, Y: midY - z - BuildingUnitHeight*DZ - h},
		renderer.Point{X: midX - 9, Y: midY - z - BuildingUnitHeight*DZ - h - 6},
		renderer.Point{X: midX - 9, Y: midY - z - BuildingUnitHeight*DZ + 6},
	}}
	cv.SetFillStyle(filepath.FromSlash("texture/building/stone_terra.png"))
	util.RenderPolygon(cv, rp1, true)

	rp2 := renderer.Polygon{Points: []renderer.Point{
		renderer.Point{X: midX, Y: midY - z - BuildingUnitHeight*DZ + 12},
		renderer.Point{X: midX, Y: midY - z - BuildingUnitHeight*DZ - h},
		renderer.Point{X: midX + 9, Y: midY - z - BuildingUnitHeight*DZ - h - 6},
		renderer.Point{X: midX + 9, Y: midY - z - BuildingUnitHeight*DZ + 6},
	}}
	cv.SetFillStyle(filepath.FromSlash("texture/building/stone_terra_flipped.png"))
	util.RenderPolygon(cv, rp2, true)

	rp3 := renderer.Polygon{Points: []renderer.Point{
		renderer.Point{X: midX + 9, Y: midY - z - BuildingUnitHeight*DZ - h - 6},
		renderer.Point{X: midX, Y: midY - z - BuildingUnitHeight*DZ - h},
		renderer.Point{X: midX - 9, Y: midY - z - BuildingUnitHeight*DZ - h - 6},
		renderer.Point{X: midX, Y: midY - z - BuildingUnitHeight*DZ - h - 12},
	}}
	cv.SetFillStyle(color.RGBA{R: 0, G: 0, B: 0, A: 224})
	util.RenderPolygon(cv, rp3, true)

	if phase != 255 {
		cv.DrawImage(filepath.FromSlash("texture/building/smoke_"+strconv.Itoa(int(phase/3))+".png"), midX-16, midY-z-BuildingUnitHeight*DZ-h-52, 32, 48)
	}
}
