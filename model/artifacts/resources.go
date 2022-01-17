package artifacts

type Artifacts struct {
	A        *Artifact
	Quantity uint16
}

type Resources struct {
	Artifacts map[*Artifact]uint16
}

func (r *Resources) Add(a *Artifact, q uint16) {
	if r.Artifacts == nil {
		r.Artifacts = make(map[*Artifact]uint16)
	}
	if e, ok := r.Artifacts[a]; ok {
		r.Artifacts[a] = e + q
	} else {
		r.Artifacts[a] = q
	}
}

func (r *Resources) IsEmpty() bool {
	return len(r.Artifacts) == 0
}
