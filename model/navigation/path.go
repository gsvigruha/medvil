package navigation

type PathElement interface {
	GetLocation() Location
	GetNeighbors(IMap) []PathElement
	GetSpeed() float64
	Walkable() bool
	Sailable() bool
}

type Path struct {
	P []PathElement
}

func (p *Path) ConsumeElement() (*Path, PathElement) {
	if len(p.P) > 1 {
		removed := p.P[0]
		p.P = p.P[1:]
		return p, removed
	}
	return nil, nil
}

func (p *Path) LastElement() PathElement {
	return p.P[len(p.P)-1]
}
