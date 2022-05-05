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

func (p Polygon) Move(dx float64, dy float64) Polygon {
	np := Polygon{Points: make([]Point, len(p.Points))}
	for i := range p.Points {
		np.Points[i] = Point{X: p.Points[i].X + dx, Y: p.Points[i].Y + dy}
	}
	return np
}
