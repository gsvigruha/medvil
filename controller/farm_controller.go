package controller

import (
	"github.com/tfriedel6/canvas"
	"image/color"
	"medvil/model/economy"
	"medvil/model/social"
	"medvil/view/gui"
)

type FamController struct {
	UseType uint8
}

type LandUseButton struct {
	b       gui.ButtonGUI
	fc      *FamController
	useType uint8
}

func (b LandUseButton) Click() {
	b.fc.UseType = b.useType
}

func (b LandUseButton) Render(cv *canvas.Canvas) {
	b.b.Render(cv)
	if b.fc.UseType != b.useType {
		cv.SetFillStyle(color.RGBA{R: 64, G: 0, B: 0, A: 128})
		cv.FillRect(b.b.X, b.b.Y, 32, 32)
	}
}

func (b LandUseButton) Contains(x float64, y float64) bool {
	return b.b.Contains(x, y)
}

func FarmToControlPanel(p *gui.Panel, h *social.Farm) {
	HouseholdToControlPanel(p, &h.Household)
	fc := &FamController{}

	p.AddButton(LandUseButton{
		b:       gui.ButtonGUI{Texture: "terrain/grain", X: float64(10), Y: float64(HouseholdControllerGUIBottomY), SX: 32, SY: 32},
		fc:      fc,
		useType: economy.FarmFieldUseTypeWheat,
	})
}
