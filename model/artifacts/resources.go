package artifacts

import (
	"encoding/json"
)

type Artifacts struct {
	A        *Artifact
	Quantity uint16
}

type Resources struct {
	Artifacts      map[*Artifact]uint16
	VolumeCapacity uint16
}

func (r *Resources) UnmarshalJSON(data []byte) error {
	r.Artifacts = make(map[*Artifact]uint16)
	var j map[string]json.RawMessage
	if err := json.Unmarshal(data, &j); err != nil {
		return err
	}
	for k, v := range j {
		var quantity uint16
		e := json.Unmarshal(v, &quantity)
		if e != nil {
			panic("Error unmarshalling json")
		}
		r.Artifacts[GetArtifact(k)] = quantity
	}
	return nil
}

func (r *Resources) AddAll(as []Artifacts) {
	for _, a := range as {
		r.Add(a.A, a.Quantity)
	}
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
	for _, q := range r.Artifacts {
		if q > 0 {
			return false
		}
	}
	return true
}

func (r *Resources) Get(a *Artifact) uint16 {
	if r.Artifacts == nil {
		r.Artifacts = make(map[*Artifact]uint16)
	}
	if q, ok := r.Artifacts[a]; ok {
		return q
	}
	return 0
}

func (r *Resources) RemoveAll(as []Artifacts) bool {
	if !r.Has(as) {
		return false
	}
	for _, a := range as {
		r.Remove(a.A, a.Quantity)
	}
	return true
}

func (r *Resources) Remove(a *Artifact, q uint16) uint16 {
	if r.Artifacts == nil {
		r.Artifacts = make(map[*Artifact]uint16)
	}
	if e, ok := r.Artifacts[a]; ok {
		if e >= q {
			r.Artifacts[a] = e - q
			return q
		} else {
			r.Artifacts[a] = 0
			return q - e
		}
	}
	return 0
}

func (r *Resources) Needs(as []Artifacts) []Artifacts {
	var remaining []Artifacts = nil
	for _, a := range as {
		if v, ok := r.Artifacts[a.A]; ok {
			if v < a.Quantity {
				remaining = append(remaining, Artifacts{A: a.A, Quantity: a.Quantity - v})
			}
		} else {
			remaining = append(remaining, Artifacts{A: a.A, Quantity: a.Quantity})
		}
	}
	return remaining
}

func (r *Resources) HasArtifact(a *Artifact) bool {
	if q, ok := r.Artifacts[a]; ok {
		if q > 0 {
			return true
		}
	}
	return false
}

func (r *Resources) Has(as []Artifacts) bool {
	for _, a := range as {
		if v, ok := r.Artifacts[a.A]; ok {
			if v < a.Quantity {
				return false
			}
		} else {
			return false
		}
	}
	return true
}

func (r *Resources) Volume() uint16 {
	var v uint16 = 0
	for a, q := range r.Artifacts {
		v += a.V * q
	}
	return v
}

func (r *Resources) UsedVolumeCapacity() float64 {
	return float64(r.Volume()) / float64(r.VolumeCapacity)
}
