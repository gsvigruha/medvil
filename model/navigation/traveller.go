package navigation

import (
	"math/rand"
)

const MaxPX = 100
const MaxPY = 100

const TravellerMinD = 25

const MaxStuckCntr = 5

const TravellerTypePedestrian uint8 = 0

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
	path      *Path
	lane      uint8
	stuckCntr uint8
}

func (t *Traveller) consumePathElement() {
	if t.path != nil {
		t.path = t.path.ConsumeElement()
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
	if field.Plant != nil && field.Plant.T.TreeT != nil {
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
	} else {
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
	} else {
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
		f := t.path.F[0]
		var dirToLane uint8 = DirectionNone
		var dirToNextField uint8 = DirectionNone
		if t.FY == f.Y {
			if t.PY < MaxPY/4*t.lane {
				dirToLane = DirectionS
			} else if t.PY > MaxPY/4*t.lane {
				dirToLane = DirectionN
			}
			if t.FX > f.X {
				dirToNextField = DirectionW
			} else if t.FX < f.X {
				dirToNextField = DirectionE
			}
		} else if t.FX == f.X {
			if t.PX < MaxPX/4*t.lane {
				dirToLane = DirectionE
			} else if t.PX > MaxPX/4*t.lane {
				dirToLane = DirectionW
			}
			if t.FY > f.Y {
				dirToNextField = DirectionN
			} else if t.FY < f.Y {
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
		t.IncPhase()
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

func (t *Traveller) EnsurePath(f *Field, travellerType uint8, m IMap) {
	if t.path == nil || t.path.LastField() != f {
		t.path = m.ShortPath(t.FX, t.FY, f.X, f.Y, travellerType)
		t.lane = uint8(rand.Intn(3) + 1)
		t.stuckCntr = 0
	}
}
