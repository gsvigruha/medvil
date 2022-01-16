package renderer

import (
	"math"
)

func RayIntersects(x float64, y float64, lx1 float64, ly1 float64, lx2 float64, ly2 float64) bool {
	if lx1 != lx2 {
		a := (ly2 - ly1) / (lx2 - lx1)
		b := ly1 - a*lx1
		xi := (y - b) / a
		return xi >= math.Min(lx1, lx2) && xi <= math.Max(lx1, lx2) && xi >= x
	} else {
		xi := lx1
		yi := y
		return yi >= math.Min(ly1, ly2) && yi <= math.Max(ly1, ly2) && xi >= x
	}
}

func BtoI(b bool) uint8 {
	if b {
		return 1
	}
	return 0
}
