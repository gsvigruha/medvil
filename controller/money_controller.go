package controller

import (
	"github.com/tfriedel6/canvas"
	"image/color"
	"medvil/model/social"
	"medvil/view/gui"
	"strconv"
)

type MoneyControllerButton struct {
	b            gui.ButtonGUI
	sourceWallet *uint32
	targetWallet *uint32
	amount       uint32
}

func (b MoneyControllerButton) Click() {
	if *b.sourceWallet >= b.amount {
		*b.sourceWallet -= b.amount
		*b.targetWallet += b.amount
	}
}

func (b MoneyControllerButton) Render(cv *canvas.Canvas) {
	b.b.Render(cv)
	if *b.sourceWallet == 0 {
		cv.SetFillStyle(color.RGBA{R: 0, G: 0, B: 0, A: 64})
		cv.FillRect(b.b.X, b.b.Y, b.b.SX, b.b.SY)
	}
}

func (b MoneyControllerButton) Contains(x float64, y float64) bool {
	return b.b.Contains(x, y)
}

func MoneyToControlPanel(p *gui.Panel, town *social.Town, targetWallet *uint32, amount uint32, x, y float64) {
	p.AddTextLabel("$ "+strconv.Itoa(int(*targetWallet)), x, y)
	if town != nil {
		p.AddButton(MoneyControllerButton{
			b:            gui.ButtonGUI{Icon: "plus", X: x + 60, Y: y - 10, SX: 12, SY: 12},
			sourceWallet: &town.Townhall.Household.Money,
			targetWallet: targetWallet,
			amount:       amount,
		})
	}
}
