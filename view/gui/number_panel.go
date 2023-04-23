package gui

import (
	"fmt"
	"github.com/tfriedel6/canvas"
)

type NumberPanel struct {
	P          *Panel
	TextLabel  *TextLabel
	UpButton   NumberPanelButton
	DownButton NumberPanelButton
	format     string
	min        int
	max        int
	get        func() int
	set        func(int)
}

type NumberPanelButton struct {
	b   ButtonGUI
	np  *NumberPanel
	inc int
}

func (b NumberPanelButton) Click() {
	if b.np.get()+b.inc >= b.np.min && b.np.get()+b.inc <= b.np.max {
		b.np.set(b.np.get() + b.inc)
	}
}

func (b NumberPanelButton) Render(cv *canvas.Canvas) {
	b.np.TextLabel.Text = fmt.Sprintf(b.np.format, b.np.get())
	b.b.Render(cv)
}

func (b NumberPanelButton) Contains(x float64, y float64) bool {
	return b.b.Contains(x, y)
}

func (b NumberPanelButton) Enabled() bool {
	return b.b.Enabled()
}

func CreateNumberPaneFromVal(x, y, sx, sy float64, min, max, inc int, format string, val *int) *NumberPanel {
	return CreateNumberPanel(x, y, sx, sy, min, max, inc, format, func() int { return *val }, func(v int) { *val = v })
}

func CreateNumberPanel(x, y, sx, sy float64, min, max, inc int, format string, get func() int, set func(int)) *NumberPanel {
	p := &Panel{}
	np := &NumberPanel{P: p, format: format, min: min, max: max, get: get, set: set}
	np.TextLabel = p.AddTextLabel("", x, y+sy*2/3)
	si := sy / 2
	np.UpButton = NumberPanelButton{
		b:   ButtonGUI{Icon: "plus", X: x + sx - si, Y: y, SX: si, SY: si},
		np:  np,
		inc: inc,
	}
	p.AddButton(np.UpButton)
	np.DownButton = NumberPanelButton{
		b:   ButtonGUI{Icon: "minus", X: x + sx - si, Y: y + si, SX: si, SY: si},
		np:  np,
		inc: -inc,
	}
	p.AddButton(np.DownButton)
	return np
}
