package artifacts

import (
	"encoding/json"
)

type Artifacts struct {
	A        *Artifact
	Quantity uint16
}

type Resources struct {
	Artifacts map[*Artifact]uint16
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
