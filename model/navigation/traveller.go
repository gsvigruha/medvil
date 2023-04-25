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
const TravellerTypeCart uint8 = 2
const TravellerTypeTradingBoat uint8 = 3
const TravellerTypeTradingCart uint8 = 4

const RoadBreakdownRate = 0.0001

type PathType uint8

const PathTypePedestrian PathType = 0
const PathTypeBoat PathType = 1

type PathComp struct {
	path      *Path
	pe        PathElement
	computing bool
	pc        chan *Path
}

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
	Vehicle   Vehicle
	pc        PathComp `ser:"false"`
	Lane      uint8
	StuckCntr uint8
}

func (t *Traveller) consumePathElement() {
	if t.pc.path != nil {
		path, removed := t.pc.path.ConsumeElement()
		t.pc.path = path
		t.pc.pe = removed
		if t.pc.pe != nil {
			t.FZ = t.pc.pe.GetLocation().Z
			t.Visible = t.pc.pe.TravellerVisible()
		}
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

func (t *Traveller) MoveLeft(stayInField bool, m IMap) {
	if t.PX > 0 {
		t.PX--
	} else if t.FX > 0 && !stayInField {
		m.GetField(t.FX, t.FY).UnregisterTraveller(t)
		t.PX = MaxPX
		t.FX--
		m.GetField(t.FX, t.FY).RegisterTraveller(t)
		t.consumePathElement()
	}
	t.Direction = DirectionW
}

func (t *Traveller) MoveRight(stayInField bool, m IMap) {
	if t.PX < MaxPX {
		t.PX++
	} else if !stayInField {
		m.GetField(t.FX, t.FY).UnregisterTraveller(t)
		t.PX = 0
		t.FX++
		m.GetField(t.FX, t.FY).RegisterTraveller(t)
		t.consumePathElement()
	}
	t.Direction = DirectionE
}

func (t *Traveller) MoveUp(stayInField bool, m IMap) {
	if t.PY > 0 {
		t.PY--
	} else if t.FY > 0 && !stayInField {
		m.GetField(t.FX, t.FY).UnregisterTraveller(t)
		t.PY = MaxPY
		t.FY--
		m.GetField(t.FX, t.FY).RegisterTraveller(t)
		t.consumePathElement()
	}
	t.Direction = DirectionN
}

func (t *Traveller) MoveDown(stayInField bool, m IMap) {
	if t.PY < MaxPY {
		t.PY++
	} else if !stayInField {
		m.GetField(t.FX, t.FY).UnregisterTraveller(t)
		t.PY = 0
		t.FY++
		m.GetField(t.FX, t.FY).RegisterTraveller(t)
		t.consumePathElement()
	}
	t.Direction = DirectionS
}

func (t *Traveller) Jump(m IMap) {
	t.consumePathElement()
}

func (t *Traveller) MoveToDir(d uint8, stayInField bool, m IMap) {
	switch d {
	case DirectionN:
		t.MoveUp(stayInField, m)
	case DirectionS:
		t.MoveDown(stayInField, m)
	case DirectionW:
		t.MoveLeft(stayInField, m)
	case DirectionE:
		t.MoveRight(stayInField, m)
	}
}

func (t *Traveller) Move(m IMap) {
	t.Motion = MotionWalk
	var steps = 1
	{
		f := m.GetField(t.FX, t.FY)
		if rand.Float64() < f.GetSpeed()-1.0 {
			steps = 2
			if rand.Float64() < RoadBreakdownRate {
				f.Road.Broken = true
			}
		}
	}
	for i := 0; i < steps; i++ {
		if t.pc.path != nil {
			pe := t.pc.path.P[0]
			l := pe.GetLocation()
			var dirToLane uint8 = DirectionNone
			var dirToNextField uint8 = DirectionNone
			if t.FX == l.X && t.FY == l.Y {
				t.Jump(m)
			} else if t.FY == l.Y {
				if t.PY < MaxPY/4*t.Lane {
					dirToLane = DirectionS
				} else if t.PY > MaxPY/4*t.Lane {
					dirToLane = DirectionN
				}
				if t.FX > l.X {
					dirToNextField = DirectionW
				} else if t.FX < l.X {
					dirToNextField = DirectionE
				}
			} else if t.FX == l.X {
				if t.PX < MaxPX/4*t.Lane {
					dirToLane = DirectionE
				} else if t.PX > MaxPX/4*t.Lane {
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
				t.MoveToDir(dirToLane, true, m)
				t.StuckCntr = 0
			} else if dirToNextField != DirectionNone && t.HasRoom(m, dirToNextField) {
				// Move towards the next field if in correct lane and possible
				t.MoveToDir(dirToNextField, false, m)
				t.StuckCntr = 0
			} else if t.StuckCntr < MaxStuckCntr {
				// Try to pick a different lane a few times
				t.Lane = uint8(rand.Intn(3) + 1)
				t.StuckCntr++
			} else {
				// Move towards any available space
				for i := uint8(0); i < 4; i++ {
					d := (dirToNextField + i) % 4
					if t.HasRoom(m, d) {
						t.MoveToDir(d, true, m)
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

func (t *Traveller) IsAtDestination(dest Destination) bool {
	return dest.Check(t.pc.pe)
}

func (t *Traveller) EnsurePath(dest Destination, m IMap) (bool, bool) {
	if t.pc.path == nil || !dest.Check(t.pc.path.LastElement()) {
		if t.pc.pc == nil {
			t.pc.pc = make(chan *Path)
		}
		if !t.pc.computing {
			t.pc.computing = true
			go func(c chan *Path) {
				c <- m.ShortPath(Location{X: t.FX, Y: t.FY, Z: t.FZ}, dest, t.PathType())
			}(t.pc.pc)
		} else {
			select {
			case t.pc.path = <-t.pc.pc:
				t.pc.computing = false
				t.Lane = uint8(rand.Intn(3) + 1)
				t.StuckCntr = 0
			}
		}
	}
	return t.pc.path != nil, t.pc.computing
}

func (t *Traveller) PathType() PathType {
	if t.Vehicle != nil {
		return t.Vehicle.PathType()
	}
	return PathTypePedestrian
}

func (t *Traveller) UseVehicle(v Vehicle) {
	t.Vehicle = v
	t.SyncTo(v.GetTraveller())
	v.SetInUse(true)
}

func (t *Traveller) ExitVehicle() {
	if t.Vehicle != nil {
		t.Vehicle.GetTraveller().PX = 50
		t.Vehicle.GetTraveller().PY = 50
		t.Vehicle.SetInUse(false)
		t.Vehicle = nil
	}
}

func (t *Traveller) MoveWith(m IMap, ot *Traveller) {
	oFX, oFY := t.FX, t.FY
	t.SyncTo(ot)
	t.FX = ot.FX
	t.FY = ot.FY
	t.FZ = ot.FZ
	t.pc.pe = ot.pc.pe
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

func (t *Traveller) SetHome(home bool) {
	if home {
		t.Visible = false
	}
	if t.Vehicle != nil {
		t.Vehicle.SetHome(home)
	}
}

func (t *Traveller) GetPathFields(m IMap) []FieldWithContext {
	var fs []FieldWithContext
	if t.pc.pe != nil {
		l := t.pc.pe.GetLocation()
		fs = append(fs, m.GetField(l.X, l.Y))
	}
	if t.pc.path != nil {
		for _, pe := range t.pc.path.P {
			l := pe.GetLocation()
			fs = append(fs, m.GetField(l.X, l.Y))
		}
	}
	return fs
}

func (t *Traveller) GetPathElement() PathElement {
	return t.pc.pe
}

func (t *Traveller) InitPathElement(pe PathElement) {
	t.pc.pe = pe
}
