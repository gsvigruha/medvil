package renderer

type Point struct {
	X float64
	Y float64
}

type Polygon struct {
	Points []Point
}

func (p *Polygon) Contains(x float64, y float64) bool {
	var is = uint8(0)
	np := len(p.Points)
	for i := range p.Points {
		is += BtoI(RayIntersects(x, y, p.Points[i].X, p.Points[i].Y, p.Points[(i+1)%np].X, p.Points[(i+1)%np].Y))
	}
	return is%2 == 1
}
