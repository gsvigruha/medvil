package gui

import (
	"github.com/tfriedel6/canvas"
)

type Panel struct {
	X       float64
	Y       float64
	SX      float64
	SY      float64
	Buttons []*Button
	Labels  []Label
}

func (p *Panel) Render(cv *canvas.Canvas) {
	cv.SetFillStyle("#210")
	cv.FillRect(p.X, p.Y, p.SX, p.SY)
	for i := range p.Buttons {
		p.Buttons[i].Render(cv)
	}
	for i := range p.Labels {
		p.Labels[i].Render(cv)
	}
}

func (p *Panel) CaptureButton(x float64, y float64) *Button {
	for i := range p.Buttons {
		if p.Buttons[i].Contains(x, y) {
			return p.Buttons[i]
		}
	}
	return nil
}

func (p *Panel) Clear() {
	p.Buttons = []*Button{}
	p.Labels = []Label{}
}

func (p *Panel) AddTextLabel(text string, x float64, y float64) {
	p.Labels = append(p.Labels, &TextLabel{Text: text, X: x, Y: y})
}

func (p *Panel) AddImageLabel(icon string, x, y, sx, sy float64, style uint8) {
	p.Labels = append(p.Labels, &ImageLabel{Icon: icon, X: x, Y: y, SX: sx, SY: sy, Style: style})
}

func (p *Panel) AddScaleLabel(icon string, x, y, sx, sy, scaleW, scale float64) {
	p.Labels = append(p.Labels, &ScaleLabel{Icon: icon, X: x, Y: y, SX: sx, SY: sy, ScaleW: scaleW, Scale: scale})
}

func (p *Panel) AddTextureLabel(texture string, x, y, sx, sy float64) {
	p.Labels = append(p.Labels, &TextureLabel{Texture: texture, X: x, Y: y, SX: sx, SY: sy})
}
