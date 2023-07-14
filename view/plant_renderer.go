package view

import (
	//"fmt"
	"github.com/tfriedel6/canvas"
	//"image/color"
	"math"
	"math/rand"
	"medvil/controller"
	"medvil/model/terrain"
	"medvil/model/time"
	"medvil/renderer"
)

func DrawBranch(cv *canvas.Canvas, plant *terrain.Plant, r *rand.Rand, sx float64, sy float64, width float64, length float64, angle float64, i uint8, prevSeasonPhase uint8, c *controller.Controller) {
	maturity := plant.Maturity(c.Map.Calendar)
	ex := sx + math.Cos(angle)*length
	ey := sy + math.Sin(angle)*length

	dx := math.Cos(angle+math.Pi/2) * width * math.Max(maturity, 0.2)
	dy := math.Sin(angle+math.Pi/2) * width * math.Max(maturity, 0.2)

	cv.SetFillStyle("texture/terrain/tree_bark.png")
	cv.BeginPath()
	cv.LineTo(sx-dx, sy-dy)
	cv.LineTo(ex-dx, ey-dy)
	cv.LineTo(ex+dx, ey+dy)
	cv.LineTo(sx+dx, sy+dy)
	cv.ClosePath()
	cv.Fill()

	var seasonPhase = uint8(r.Intn(int(30 / plant.T.TreeT.BranchingIterations)))
	if prevSeasonPhase >= seasonPhase {
		seasonPhase = prevSeasonPhase - seasonPhase
	}

	if c.Map.Calendar.Season() != time.Winter {
		if i > plant.T.TreeT.LeavesMinIterarion {
			dxL := math.Cos(angle+math.Pi/2) * plant.T.TreeT.LeavesSize
			dyL := math.Sin(angle+math.Pi/2) * plant.T.TreeT.LeavesSize
			var draw = false
			if (c.Map.Calendar.Month == 3 && seasonPhase <= c.Map.Calendar.Day) ||
				(c.Map.Calendar.Month == 4 && seasonPhase > c.Map.Calendar.Day) {
				if plant.T.TreeT.Blooms {
					cv.SetFillStyle("texture/terrain/leaves_blooming.png")
				} else {
					cv.SetFillStyle("texture/terrain/leaves_v2.png")
				}
				draw = true
			}
			if (c.Map.Calendar.Month == 4 && seasonPhase <= c.Map.Calendar.Day) ||
				(c.Map.Calendar.Month > 4 && c.Map.Calendar.Month < 9) ||
				(c.Map.Calendar.Month == 9 && seasonPhase > c.Map.Calendar.Day) {
				cv.SetFillStyle("texture/terrain/leaves_v2.png")
				draw = true
			}
			if (c.Map.Calendar.Month == 9 && seasonPhase <= c.Map.Calendar.Day) ||
				(c.Map.Calendar.Month == 10) ||
				(c.Map.Calendar.Month == 11 && seasonPhase > c.Map.Calendar.Day) {
				cv.SetFillStyle("texture/terrain/leaves_colored.png")
				draw = true
			}
			if draw {
				cv.BeginPath()
				cv.LineTo(sx-dxL, sy-dyL)
				cv.LineTo(ex-dxL+dyL, ey-dyL-dxL)
				cv.LineTo(ex+dxL+dyL, ey+dyL-dxL)
				cv.LineTo(sx+dxL, sy+dyL)
				cv.ClosePath()
				cv.Fill()
			}
			if plant.T.TreeT.Blooms {
				if (c.Map.Calendar.Month == 6 && seasonPhase <= c.Map.Calendar.Day) ||
					(c.Map.Calendar.Month == 7) ||
					(c.Map.Calendar.Month == 8 && seasonPhase > c.Map.Calendar.Day) {
					cv.SetFillStyle("texture/terrain/fruit.png")
					cv.BeginPath()
					cv.LineTo(sx-dxL, sy-dyL)
					cv.LineTo(ex-dxL+dyL, ey-dyL-dxL)
					cv.LineTo(ex+dxL+dyL, ey+dyL-dxL)
					cv.LineTo(sx+dxL, sy+dyL)
					cv.ClosePath()
					cv.Fill()
				}
			}

		}
	}

	if (c.Map.Calendar.Month == 12 && seasonPhase < c.Map.Calendar.Day) ||
		(c.Map.Calendar.Month == 1) ||
		(c.Map.Calendar.Month == 2 && seasonPhase > c.Map.Calendar.Day) {
		if angle < -math.Pi/2-math.Pi/4 || angle > -math.Pi/2+math.Pi/4 {
			cv.SetFillStyle("texture/terrain/snow_patches.png")
			cv.BeginPath()
			cv.LineTo(sx-dx, sy-dy)
			cv.LineTo(ex-dx, ey-dy)
			cv.LineTo(ex+dx, ey+dy)
			cv.LineTo(sx+dx, sy+dy)
			cv.ClosePath()
			cv.Fill()
		}
	}

	maxI := uint8(math.Max(2.0, math.Ceil(float64(plant.T.TreeT.BranchingIterations)*maturity)))
	if i < maxI {
		var prevAngleD = 0.0
		nextAngles := make([]float64, len(plant.T.TreeT.BranchAngles))
		for branchI, nextAngle := range plant.T.TreeT.BranchAngles {
			var angleD = 2.0*r.Float64()*nextAngle - nextAngle
			if (prevAngleD < 0 && angleD < 0) || (prevAngleD > 0 && angleD > 0) {
				angleD = -angleD
			}
			nextAngles[branchI] = angle + angleD
			prevAngleD = angleD
		}

		for branchI, _ := range plant.T.TreeT.BranchAngles {
			nextWidth := width * plant.T.TreeT.BranchWidthD[branchI]
			nextLength := length * plant.T.TreeT.BranchLengthD[branchI]
			DrawBranch(cv, plant, r, ex, ey, nextWidth, nextLength, nextAngles[branchI], i+1, seasonPhase, c)
		}
	}
}

func RenderTree(cv *canvas.Canvas, plant *terrain.Plant, rf renderer.RenderedField, c *controller.Controller) {
	r := rand.New(rand.NewSource(int64(plant.Shape)))
	midX, midY := rf.MidPoint()
	DrawBranch(cv, plant, r, midX, midY, plant.T.TreeT.BranchWidth0, plant.T.TreeT.BranchLength0, -math.Pi/2, 0, 30, c)
}

func RenderRegularPlant(cv *canvas.Canvas, plant *terrain.Plant, rf renderer.RenderedField, c *controller.Controller) {
	cv.DrawImage("texture/terrain/"+plant.T.Name+".png", rf.X[1], rf.Y[2]-108, 120, 108)
}

func RenderPlant(cv *canvas.Canvas, plant *terrain.Plant, rf renderer.RenderedField, c *controller.Controller) {
	if plant.T.TreeT != nil {
		RenderTree(cv, plant, rf, c)
	} else {
		RenderRegularPlant(cv, plant, rf, c)
	}
}
