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
	val        *int
}

type NumberPanelButton struct {
	b   ButtonGUI
	np  *NumberPanel
	inc int
}

func (b NumberPanelButton) Click() {
	if *b.np.val > b.np.min && *b.np.val < b.np.max {
		*b.np.val += b.inc
	}
}

func (b NumberPanelButton) Render(cv *canvas.Canvas) {
	b.np.TextLabel.Text = fmt.Sprintf(b.np.format, *b.np.val)
	b.b.Render(cv)
}

func (b NumberPanelButton) Contains(x float64, y float64) bool {
	return b.b.Contains(x, y)
}

func CreateNumberPanel(x, y, sx, sy float64, min, max, inc int, format string, val *int) *NumberPanel {
	p := &Panel{}
	np := &NumberPanel{P: p, val: val, format: format, min: min, max: max}
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
