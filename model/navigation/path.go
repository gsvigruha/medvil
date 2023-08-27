package navigation

type PathElement interface {
	GetLocation() Location
	GetNeighbors(IMap) []PathElement
	GetSpeed() float64
	Walkable() bool
	Sailable() bool
	TravellerVisible() bool
	Crowded() bool
	LocationXY() (uint16, uint16)
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
	if len(p.P) == 1 {
		return nil, p.P[0]
	}
	return nil, nil
}

func (p *Path) LastElement() PathElement {
	return p.P[len(p.P)-1]
}
