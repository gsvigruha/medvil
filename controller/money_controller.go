package controller

import (
	"github.com/tfriedel6/canvas"
	"image/color"
	"medvil/model/economy"
	"medvil/view/gui"
	"strconv"
)

type MoneyControllerButton struct {
	b            gui.ButtonGUI
	sourceWallet economy.Wallet
	targetWallet economy.Wallet
	amount       uint32
}

func (b *MoneyControllerButton) SetHoover(h bool) {
	b.b.SetHoover(h)
}

func (b MoneyControllerButton) Click() {
	if b.sourceWallet.GetMoney() >= b.amount {
		b.sourceWallet.Spend(b.amount)
		b.targetWallet.Earn(b.amount)
	}
}

func (b MoneyControllerButton) Render(cv *canvas.Canvas) {
	b.b.Render(cv)
	if b.sourceWallet.GetMoney() == 0 {
		cv.SetFillStyle(color.RGBA{R: 0, G: 0, B: 0, A: 64})
		cv.FillRect(b.b.X, b.b.Y, b.b.SX, b.b.SY)
	}
}

func (b MoneyControllerButton) Contains(x float64, y float64) bool {
	return b.b.Contains(x, y)
}

func (b MoneyControllerButton) Enabled() bool {
	return b.b.Enabled()
}

func MoneyToControlPanel(cp *ControlPanel, p *gui.Panel, srcWallet economy.Wallet, targetWallet economy.Wallet, amount uint32, x, y float64) {
	p.AddImageLabel("coin", x, y-gui.FontSize*0.8, gui.FontSize, gui.FontSize, gui.ImageLabelStyleRegular)
	p.AddTextLabel(strconv.Itoa(int(targetWallet.GetMoney())), x+gui.FontSize+4, y)
	if srcWallet != nil {
		p.AddButton(&MoneyControllerButton{
			b: gui.ButtonGUI{Icon: "plus", X: x + gui.FontSize*4, Y: y - gui.FontSize*0.8, SX: gui.FontSize, SY: gui.FontSize, OnHoover: func() {
				cp.HelperMessage("Send money from the townhall")
			}},
			sourceWallet: srcWallet,
			targetWallet: targetWallet,
			amount:       amount,
		})
	}
}
