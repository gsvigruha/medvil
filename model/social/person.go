package social

import (
	"medvil/model/economy"
	"medvil/model/navigation"
	"medvil/model/time"
)

const MaxPX = 10
const MaxPY = 10

type Person struct {
	FX        uint16
	FY        uint16
	FZ        uint8
	PX        uint8
	PY        uint8
	Hunger    uint8
	Thirst    uint8
	Happiness uint8
	Health    uint8
	Household *Household
	Task      economy.Task
	IsHome    bool
}

func (p *Person) MoveLeft() {
	if p.PX > 0 {
		p.PX--
	} else {
		p.PX = MaxPX
		p.FX--
	}
}

func (p *Person) MoveRight() {
	if p.PX < MaxPX {
		p.PX++
	} else {
		p.PX = 0
		p.FX++
	}
}

func (p *Person) MoveUp() {
	if p.PY > 0 {
		p.PY--
	} else {
		p.PY = MaxPY
		p.FY--
	}
}

func (p *Person) MoveDown() {
	if p.PY < MaxPY {
		p.PY++
	} else {
		p.PY = 0
		p.FY++
	}
}

func (p *Person) ElapseTime(Calendar *time.CalendarType, Map navigation.IMap) {
	if p.Task != nil {
		if p.FX == p.Task.Location().X && p.FY == p.Task.Location().Y {
			p.Task.Tick()
		} else {
			if p.FX > p.Task.Location().X {
				p.MoveLeft()
			} else if p.FX < p.Task.Location().X {
				p.MoveRight()
			} else if p.FY > p.Task.Location().Y {
				p.MoveDown()
			} else if p.FY < p.Task.Location().Y {
				p.MoveUp()
			}
		}
	}
}
