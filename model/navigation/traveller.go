package navigation

import (
	"math/rand"
	"medvil/model/building"
	"strconv"
)

const MaxPX = 100
const MaxPY = 100

const TravellerMinD = 25

const MaxStuckCntr = 5

const TravellerTypePedestrianM uint8 = 0
const TravellerTypePedestrianF uint8 = 7
const TravellerTypeBoat uint8 = 1
const TravellerTypeCart uint8 = 2
const TravellerTypeTradingBoat uint8 = 3
const TravellerTypeTradingCart uint8 = 4
const TravellerTypeExpeditionBoat uint8 = 5
const TravellerTypeExpeditionCart uint8 = 6

const RoadBreakdownRate = 0.0001

type PathType uint8

const PathTypePedestrian PathType = 0
const PathTypeBoat PathType = 1
const PathTypeCart PathType = 2

type Person interface {
}

type PathComp struct {
	Path      *Path
	PE        PathElement
	computing bool       `ser:"false"`
	pc        chan *Path `ser:"false"`
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
	Person    Person
	PathComp  PathComp
	Lane      uint8
	StuckCntr uint8
}

func (t *Traveller) consumePathElement() {
	if t.PathComp.Path != nil {
		path, removed := t.PathComp.Path.ConsumeElement()
		t.PathComp.Path = path
		t.PathComp.PE = removed
		if t.PathComp.PE != nil {
			t.FZ = t.PathComp.PE.GetLocation().Z
			t.Visible = t.PathComp.PE.TravellerVisible()
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
	if (field.Plant != nil && field.Plant.T.TreeT != nil) || field.Animal != nil || field.Statue != nil {
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
	if !field.Building.Empty() {
		if ext, ok := field.Building.BuildingComponents[0].(*building.ExtensionUnit); ok {
			if ext.Direction == DirectionN && t.BlockedBy(dir, MaxPX/2, MaxPY/4) {
				return false
			}
			if ext.Direction == DirectionS && t.BlockedBy(dir, MaxPX/2, MaxPY*3/4) {
				return false
			}
			if ext.Direction == DirectionW && t.BlockedBy(dir, MaxPX/4, MaxPY/2) {
				return false
			}
			if ext.Direction == DirectionE && t.BlockedBy(dir, MaxPX*3/4, MaxPY/2) {
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

func (t *Traveller) MoveToCenter(m IMap) {
	if t.PX < MaxPX/2 {
		t.MoveToDir(DirectionE, true, m)
	} else if t.PX > MaxPX/2 {
		t.MoveToDir(DirectionW, true, m)
	} else if t.PY < MaxPY/2 {
		t.MoveToDir(DirectionS, true, m)
	} else if t.PY > MaxPY/2 {
		t.MoveToDir(DirectionN, true, m)
	}
}

func (t *Traveller) Move(m IMap) {
	t.Motion = MotionWalk
	if t.Vehicle != nil && t.Vehicle.Water() {
		t.Motion = MotionPaddle
	}
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
		if t.PathComp.Path != nil {
			pe := t.PathComp.Path.P[0]
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
	return dest.Check(t.PathComp.PE)
}

func (t *Traveller) EnsurePath(dest Destination, m IMap) (bool, bool) {
	if t.PathComp.Path == nil || !dest.Check(t.PathComp.Path.LastElement()) {
		if t.PathComp.pc == nil {
			t.PathComp.pc = make(chan *Path)
		}
		if !t.PathComp.computing {
			t.PathComp.computing = true
			go func(c chan *Path) {
				c <- m.ShortPath(Location{X: t.FX, Y: t.FY, Z: t.FZ}, dest, t.PathType())
			}(t.PathComp.pc)
		} else {
			select {
			case t.PathComp.Path = <-t.PathComp.pc:
				t.PathComp.computing = false
				if t.Lane == 0 {
					if t.T == TravellerTypeExpeditionCart || t.T == TravellerTypeExpeditionBoat {
						t.Lane = 2
					} else {
						t.Lane = uint8(rand.Intn(3) + 1)
					}
				}
				t.StuckCntr = 0
			}
		}
	}
	return t.PathComp.Path != nil, t.PathComp.computing
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
		t.Vehicle.GetTraveller().PX = MaxPX / 2
		t.Vehicle.GetTraveller().PY = MaxPY / 2
		t.PX = MaxPX / 2
		t.PY = MaxPY / 2
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
	t.PathComp.PE = ot.PathComp.PE
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
	if t.PathComp.PE != nil {
		l := t.PathComp.PE.GetLocation()
		fs = append(fs, m.GetField(l.X, l.Y))
	}
	if t.PathComp.Path != nil {
		for _, pe := range t.PathComp.Path.P {
			l := pe.GetLocation()
			fs = append(fs, m.GetField(l.X, l.Y))
		}
	}
	return fs
}

func (t *Traveller) GetPathElement() PathElement {
	return t.PathComp.PE
}

func (t *Traveller) InitPathElement(pe PathElement) {
	t.PathComp.PE = pe
}

func (t *Traveller) DrawingPhase() uint8 {
	return (t.Phase / 2) % 8
}

func (t *Traveller) CacheKey(perspective uint8) string {
	dirIdx := (perspective - t.Direction) % 4
	key := strconv.Itoa(int(dirIdx)) + "#" +
		strconv.Itoa(int(t.Motion)) + "#" +
		strconv.Itoa(int(t.DrawingPhase())) + "#" +
		strconv.Itoa(int(t.T))
	if t.Vehicle != nil {
		key = key + "#v" + strconv.FormatBool(t.Vehicle.Water())
	}
	return key
}
