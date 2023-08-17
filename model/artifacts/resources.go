package artifacts

import (
	"encoding/json"
)

const InfiniteQuantity = 65535
const StorageFullThreshold = 0.9

type Artifacts struct {
	A        *Artifact
	Quantity uint16
}

type Order struct {
	A         *Artifact
	Quantity  uint16
	UnitPrice uint16
}

func (a Artifacts) Multiply(n uint16) Artifacts {
	return Artifacts{A: a.A, Quantity: a.Quantity * n}
}

func GetQuantity(as []Artifacts, oa *Artifact) uint16 {
	for _, a := range as {
		if a.A == oa && a.Quantity > 0 {
			return a.Quantity
		}
	}
	return false
}

func Multiply(as []Artifacts, n uint16) []Artifacts {
	r := make([]Artifacts, len(as))
	for i := range as {
		r[i] = as[i].Multiply(n)
	}
	return r
}

func Purchasable(as []Artifacts) []Artifacts {
	r := make([]Artifacts, 0, len(as))
	for i := range as {
		if as[i].A != Water {
			r = append(r, as[i])
		}
	}
	return r
}

func GetVolume(as []Artifacts) uint16 {
	var v uint16 = 0
	for _, a := range as {
		v += a.A.V * a.Quantity
	}
	return v
}

func (a Artifacts) Wrap() []Artifacts {
	return []Artifacts{a}
}

type Resources struct {
	Artifacts      map[*Artifact]uint16
	VolumeCapacity uint16
}

func (r *Resources) Init(volumeCapacity uint16) {
	r.Artifacts = make(map[*Artifact]uint16)
	r.VolumeCapacity = volumeCapacity
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

func (r *Resources) MarshalJSON() ([]byte, error) {
	var content map[string]interface{} = make(map[string]interface{})
	for a, q := range r.Artifacts {
		content[a.Name] = q
	}
	return json.Marshal(content)
}

func (r *Resources) AddAll(as []Artifacts) {
	for _, a := range as {
		r.Add(a.A, a.Quantity)
	}
}

func (r *Resources) AddResources(or Resources) {
	for a, q := range or.Artifacts {
		r.Add(a, q)
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

func (r *Resources) IsRealArtifact(a *Artifact) bool {
	if q, ok := r.Artifacts[a]; ok {
		return q < InfiniteQuantity
	}
	return false
}

func (r *Resources) HasRealArtifacts() bool {
	for _, q := range r.Artifacts {
		if q > 0 && q < InfiniteQuantity {
			return true
		}
	}
	return false
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
			if e < InfiniteQuantity {
				r.Artifacts[a] = e - q
			}
			return q
		} else {
			r.Artifacts[a] = 0
			return e
		}
	}
	return 0
}

func (r *Resources) Needs(as []Artifacts) []Artifacts {
	var remaining []Artifacts = nil
	for _, a := range as {
		if a.Quantity > 0 {
			if v, ok := r.Artifacts[a.A]; ok {
				if v < a.Quantity {
					remaining = append(remaining, Artifacts{A: a.A, Quantity: a.Quantity - v})
				}
			} else {
				remaining = append(remaining, Artifacts{A: a.A, Quantity: a.Quantity})
			}
		}
	}
	return remaining
}

func (r *Resources) GetAsManyAsPossible(as []Artifacts) []Artifacts {
	var result []Artifacts = nil
	for _, a := range as {
		if v, ok := r.Artifacts[a.A]; ok {
			if v <= a.Quantity {
				result = append(result, Artifacts{A: a.A, Quantity: v})
				r.Artifacts[a.A] = 0
			} else {
				result = append(result, Artifacts{A: a.A, Quantity: a.Quantity})
				r.Artifacts[a.A] = v - a.Quantity
			}
		}
	}
	return result
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
		if a.Quantity > 0 {
			if v, ok := r.Artifacts[a.A]; ok {
				if v < a.Quantity {
					return false
				}
			} else {
				return false
			}
		}
	}
	return true
}

func (r *Resources) HasAny(as []Artifacts) bool {
	if len(as) == 0 {
		return true
	}
	for _, a := range as {
		if v, ok := r.Artifacts[a.A]; ok {
			if v > 0 {
				return true
			}
		}
	}
	return false
}

func (r *Resources) IsEmpty() bool {
	if r.Artifacts == nil {
		r.Artifacts = make(map[*Artifact]uint16)
	}
	for _, q := range r.Artifacts {
		if q > 0 {
			return false
		}
	}
	return true
}

func (r *Resources) Volume() uint16 {
	if r.Artifacts == nil {
		r.Artifacts = make(map[*Artifact]uint16)
	}
	var v uint16 = 0
	for a, q := range r.Artifacts {
		v += a.V * q
	}
	return v
}

func (r *Resources) GetArtifacts() []*Artifact {
	var as []*Artifact
	for a, q := range r.Artifacts {
		if q > 0 {
			as = append(as, a)
		}
	}
	return as
}

func (r *Resources) UsedVolumeCapacity() float64 {
	return float64(r.Volume()) / float64(r.VolumeCapacity)
}

func (r *Resources) Full() bool {
	return r.UsedVolumeCapacity() >= StorageFullThreshold
}

func (r *Resources) NumArtifacts() uint32 {
	var n uint32
	for _, q := range r.Artifacts {
		n += uint32(q)
	}
	return n
}

func ArtifactsDiff(as1 []Artifacts, as2 []Artifacts) []Artifacts {
	r := Resources{}
	r.Init(0)
	r.AddAll(as2)
	return r.Needs(as1)
}
