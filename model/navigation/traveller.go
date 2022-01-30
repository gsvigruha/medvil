package navigation

import (
	"math/rand"
)

const MaxPX = 100
const MaxPY = 100

const TravellerMinD = 25

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
	Path      *Path
	Visible   bool
	Lane      uint8
}

func (t *Traveller) consumePathElement() {
	if t.Path != nil {
		t.Path = t.Path.ConsumeElement()
	}
}

func absDistance(c1, c2 uint8) uint8 {
	if c1 > c2 {
		return c1 - c2
	}
	return c2 - c1
}

func (t *Traveller) HasRoom(m IMap, dir uint8) bool {
	for _, ot := range m.GetField(t.FX, t.FY).Travellers {
		if t != ot && ot.Visible && absDistance(t.PX, ot.PX) < TravellerMinD && absDistance(t.PY, ot.PY) < TravellerMinD {
			if (dir == DirectionW && t.PX > ot.PX) ||
				(dir == DirectionE && t.PX < ot.PX) ||
				(dir == DirectionN && t.PY > ot.PY) ||
				(dir == DirectionS && t.PY < ot.PY) {
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

func (t *Traveller) Move(m IMap) {
	if t.Lane == 0 {
		t.Lane = uint8(rand.Intn(3) + 1)
	}
	if t.Path != nil {
		f := t.Path.F[0]
		var d1 uint8 = DirectionNone
		var d2 uint8 = DirectionNone
		if t.FY == f.Y {
			if t.PY < MaxPY/4*t.Lane {
				d1 = DirectionS
			} else if t.PY > MaxPY/4*t.Lane {
				d1 = DirectionN
			} 
			if t.FX > f.X {
				d2 = DirectionW
			} else if t.FX < f.X {
				d2 = DirectionE
			}
		} else if t.FX == f.X {
			if t.PX < MaxPX/4*t.Lane {
				d1 = DirectionE
			} else if t.PX > MaxPX/4*t.Lane {
				d1 = DirectionW
			}
			if t.FY > f.Y {
				d2 = DirectionN
			} else if t.FY < f.Y {
				d2 = DirectionS
			}
		}
		if d1 != DirectionNone && t.HasRoom(m, d1) {
			switch d1 {
			case DirectionN:
				t.MoveUp(m)
			case DirectionS:
				t.MoveDown(m)
			case DirectionW:
				t.MoveLeft(m)
			case DirectionE:
				t.MoveRight(m)
			}
		} else if d2 != DirectionNone && t.HasRoom(m, d2) {
			switch d2 {
			case DirectionN:
				t.MoveUp(m)
			case DirectionS:
				t.MoveDown(m)
			case DirectionW:
				t.MoveLeft(m)
			case DirectionE:
				t.MoveRight(m)
			}
		} else {
			t.Lane = uint8(rand.Intn(3) + 1)
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
