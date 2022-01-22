package navigation

type Path struct {
	L []Location
}

func (p *Path) ConsumeElement() *Path {
	if len(p.L) > 1 {
		p.L = p.L[1:]
		return p
	}
	return nil
}

func (p *Path) LastLocation() Location {
	return p.L[len(p.L)-1]
}
