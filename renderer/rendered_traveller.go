package renderer

import (
	"medvil/model/navigation"
)

type RenderedTraveller struct {
	X         float64
	Y         float64
	W         float64
	H         float64
	Traveller *navigation.Traveller
}

func (rt *RenderedTraveller) Contains(x float64, y float64) bool {
	return x >= rt.X-rt.W && x <= rt.X+rt.W && y >= rt.Y-rt.H && y <= rt.Y
}
