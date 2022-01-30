package navigation

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
		if t != ot && absDistance(t.PX, ot.PX) < TravellerMinD && absDistance(t.PY, ot.PY) < TravellerMinD {
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
		if t.HasRoom(m, DirectionW) {
			t.PX--
		}
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
		if t.HasRoom(m, DirectionE) {
			t.PX++
		}
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
		if t.HasRoom(m, DirectionN) {
			t.PY--
		}
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
		if t.HasRoom(m, DirectionS) {
			t.PY++
		}
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
	if t.Path != nil {
		f := t.Path.F[0]
		if t.FY == f.Y {
			if t.PY < MaxPY/2 {
				t.MoveDown(m)
			} else if t.PY > MaxPY/2 {
				t.MoveUp(m)
			} else if t.FX > f.X {
				t.MoveLeft(m)
			} else if t.FX < f.X {
				t.MoveRight(m)
			}
		} else if t.FX == f.X {
			if t.PX < MaxPX/2 {
				t.MoveRight(m)
			} else if t.PX > MaxPX/2 {
				t.MoveLeft(m)
			} else if t.FY > f.Y {
				t.MoveUp(m)
			} else if t.FY < f.Y {
				t.MoveDown(m)
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
