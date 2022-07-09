package navigation

import (
	"math/rand"
)

const MaxPX = 100
const MaxPY = 100

const TravellerMinD = 25

const MaxStuckCntr = 5

const TravellerTypePedestrian uint8 = 0
const TravellerTypeBoat uint8 = 1

const RoadBreakdownRate = 0.0001

type Traveller struct {
	FX        uint16
	FY        uint16
	FZ        uint8
	PX        uint8
	PY        uint8
	Direction uint8
	Motion    uint8
	Phase     uint8
	Visible   bool
	T         uint8
	path      *Path
	PE        PathElement
	lane      uint8
	stuckCntr uint8
	Vehicle   Vehicle
}

func (t *Traveller) consumePathElement() {
	if t.path != nil {
		path, removed := t.path.ConsumeElement()
		t.path = path
		t.PE = removed
	}
}

func absDistance(c1, c2 uint8) uint8 {
	if c1 > c2 {
		return c1 - c2
	}
	return c2 - c1
}

func (t *Traveller) BlockedBy(dir, opx, opy uint8) bool {
	if absDistance(t.PX, opx) < TravellerMinD && absDistance(t.PY, opy) < TravellerMinD {
		if (dir == DirectionW && t.PX > opx) ||
			(dir == DirectionE && t.PX < opx) ||
			(dir == DirectionN && t.PY > opy) ||
			(dir == DirectionS && t.PY < opy) {
			return true
		}
	}
	return false
}

func (t *Traveller) HasRoom(m IMap, dir uint8) bool {
	field := m.GetField(t.FX, t.FY)
	if (field.Plant != nil && field.Plant.T.TreeT != nil) || field.Animal != nil {
		// Trees are assumed to be in the middle of the field
		if t.BlockedBy(dir, MaxPX/2, MaxPY/2) {
			return false
		}
	}
	for _, ot := range field.Travellers {
		if t != ot && ot.Visible {
			if t.BlockedBy(dir, ot.PX, ot.PY) {
				return false
			}
		}
	}
	return true
}

func (t *Traveller) MoveLeft(m IMap) {
	if t.PX > 0 {
		t.PX--
	} else if t.FX > 0 {
		m.GetField(t.FX, t.FY).UnregisterTraveller(t)
		t.PX = MaxPX
		t.FX--
		m.GetField(t.FX, t.FY).RegisterTraveller(t)
		t.consumePathElement()
	}
	t.Direction = DirectionW
}

func (t *Traveller) MoveRight(m IMap) {
	if t.PX < MaxPX {
		t.PX++
	} else {
		m.GetField(t.FX, t.FY).UnregisterTraveller(t)
		t.PX = 0
		t.FX++
		m.GetField(t.FX, t.FY).RegisterTraveller(t)
		t.consumePathElement()
	}
	t.Direction = DirectionE
}

func (t *Traveller) MoveUp(m IMap) {
	if t.PY > 0 {
		t.PY--
	} else if t.FY > 0 {
		m.GetField(t.FX, t.FY).UnregisterTraveller(t)
		t.PY = MaxPY
		t.FY--
		m.GetField(t.FX, t.FY).RegisterTraveller(t)
		t.consumePathElement()
	}
	t.Direction = DirectionN
}

func (t *Traveller) MoveDown(m IMap) {
	if t.PY < MaxPY {
		t.PY++
	} else {
		m.GetField(t.FX, t.FY).UnregisterTraveller(t)
		t.PY = 0
		t.FY++
		m.GetField(t.FX, t.FY).RegisterTraveller(t)
		t.consumePathElement()
	}
	t.Direction = DirectionS
}

func (t *Traveller) MoveToDir(d uint8, m IMap) {
	switch d {
	case DirectionN:
		t.MoveUp(m)
	case DirectionS:
		t.MoveDown(m)
	case DirectionW:
		t.MoveLeft(m)
	case DirectionE:
		t.MoveRight(m)
	}
}

func (t *Traveller) Move(m IMap) {
	if t.path != nil {
		t.Motion = MotionWalk
		pe := t.path.P[0]
		l := pe.GetLocation()
		f := m.GetField(t.FX, t.FY)
		var steps = 1
		if rand.Float64() < f.GetSpeed()-1.0 {
			steps = 2
			if rand.Float64() < RoadBreakdownRate {
				f.Road.Broken = true
			}
		}
		for i := 0; i < steps; i++ {
			var dirToLane uint8 = DirectionNone
			var dirToNextField uint8 = DirectionNone
			if t.FY == l.Y {
				if t.PY < MaxPY/4*t.lane {
					dirToLane = DirectionS
				} else if t.PY > MaxPY/4*t.lane {
					dirToLane = DirectionN
				}
				if t.FX > l.X {
					dirToNextField = DirectionW
				} else if t.FX < l.X {
					dirToNextField = DirectionE
				}
			} else if t.FX == l.X {
				if t.PX < MaxPX/4*t.lane {
					dirToLane = DirectionE
				} else if t.PX > MaxPX/4*t.lane {
					dirToLane = DirectionW
				}
				if t.FY > l.Y {
					dirToNextField = DirectionN
				} else if t.FY < l.Y {
					dirToNextField = DirectionS
				}
			}
			if dirToLane != DirectionNone && t.HasRoom(m, dirToLane) {
				// Move towards the lane if possible
				t.MoveToDir(dirToLane, m)
				t.stuckCntr = 0
			} else if dirToNextField != DirectionNone && t.HasRoom(m, dirToNextField) {
				// Move towards the next field if in correct lane and possible
				t.MoveToDir(dirToNextField, m)
				t.stuckCntr = 0
			} else if t.stuckCntr < MaxStuckCntr {
				// Try to pick a different lane a few times
				t.lane = uint8(rand.Intn(3) + 1)
				t.stuckCntr++
			} else {
				// Move towards any available space
				for i := uint8(0); i < 4; i++ {
					d := (dirToNextField + i) % 4
					if t.HasRoom(m, d) {
						t.MoveToDir(d, m)
						break
					}
				}
			}
		}
		t.IncPhase()
		if t.Vehicle != nil {
			t.Vehicle.GetTraveller().MoveWith(m, t)
		}
	}
}

func (t *Traveller) ResetPhase() {
	t.Phase = 0
}

func (t *Traveller) IncPhase() {
	t.Phase++
	if t.Phase >= 128 {
		t.Phase = 0
	}
}

func (t *Traveller) EnsurePath(f *Field, m IMap) bool {
	if t.path == nil || t.path.LastElement().GetLocation() != f.GetLocation() {
		t.path = m.ShortPath(Location{X: t.FX, Y: t.FY, Z: t.FZ}, Location{X: f.X, Y: f.Y, Z: 0}, t.TravellerType())
		t.lane = uint8(rand.Intn(3) + 1)
		t.stuckCntr = 0
	}
	return t.path != nil
}

func (t *Traveller) TravellerType() uint8 {
	if t.Vehicle != nil {
		return t.Vehicle.TravellerType()
	}
	return TravellerTypePedestrian
}

func (t *Traveller) UseVehicle(v Vehicle) {
	t.Vehicle = v
	t.SyncTo(v.GetTraveller())
	v.SetInUse(true)
	t.Visible = false
}

func (t *Traveller) ExitVehicle() {
	if t.Vehicle != nil {
		t.Vehicle.GetTraveller().PX = 50
		t.Vehicle.GetTraveller().PY = 50
		t.Vehicle.SetInUse(false)
		t.Vehicle = nil
		t.Visible = true
	}
}

func (t *Traveller) MoveWith(m IMap, ot *Traveller) {
	oFX, oFY := t.FX, t.FY
	t.SyncTo(ot)
	t.FX = ot.FX
	t.FY = ot.FY
	if oFX != t.FX || oFY != t.FY {
		m.GetField(oFX, oFY).UnregisterTraveller(t)
		m.GetField(t.FX, t.FY).RegisterTraveller(t)
	}
}

func (t *Traveller) SyncTo(ot *Traveller) {
	t.PX = ot.PX
	t.PY = ot.PY
	t.Phase = ot.Phase
	t.Direction = ot.Direction
}

func (t *Traveller) IsOnFieldCenter() bool {
	return t.PX > MaxPX/4 && t.PY > MaxPY/4 && t.PX < MaxPX*3/4 && t.PY < MaxPY*3/4
}
