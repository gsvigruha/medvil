package navigation

type Path struct {
	F []*Field
}

func (p *Path) ConsumeElement() *Path {
	if len(p.F) > 1 {
		p.F = p.F[1:]
		return p
	}
	return nil
}

func (p *Path) LastField() *Field {
	return p.F[len(p.F)-1]
}
