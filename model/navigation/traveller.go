package navigation

const MaxPX = 100
const MaxPY = 100

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
}

func (t *Traveller) consumePathElement() {
	if t.Path != nil {
		t.Path = t.Path.ConsumeElement()
	}
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
}

func (t *Traveller) Move(m IMap) {
	if t.Path != nil {
		l := t.Path.L[0]
		if t.FX > l.X {
			t.MoveLeft(m)
			t.Direction = DirectionW
		} else if t.FX < l.X {
			t.MoveRight(m)
			t.Direction = DirectionE
		} else if t.FY > l.Y {
			t.MoveUp(m)
			t.Direction = DirectionN
		} else if t.FY < l.Y {
			t.MoveDown(m)
			t.Direction = DirectionS
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
