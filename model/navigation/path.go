package navigation

type PathElement struct {
	F *Field
	Z uint8
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
