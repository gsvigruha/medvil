package navigation

type PathElement interface {
	GetLocation() Location
	GetNeighbors(IMap) []PathElement
	GetSpeed() float64
	Walkable() bool
}

type Path struct {
	P []PathElement
}

func (p *Path) ConsumeElement() *Path {
	if len(p.P) > 1 {
		p.P = p.P[1:]
		return p
	}
	return nil
}

func (p *Path) LastElement() PathElement {
	return p.P[len(p.P)-1]
}
