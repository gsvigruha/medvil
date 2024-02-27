package gui

type Popup struct {
	Message string
	ScreenX float64
	ScreenY float64
	yesFn   func()
	noFn    func()
}

func CreatePopup(message string, screenX float64, screenY float64, iconSize float64, yesFn func(), noFn func()) *Panel {
	w := screenX / 12
	h := screenY / 12
	p := &Panel{X: screenX/2 - w, Y: screenY/2 - h, SX: w * 2, SY: h * 2}
	p.AddTextLabel(message, screenX/2-EstimateWidth(message)*FontSize/2.0, screenY/2-iconSize)
	p.AddButton(&SimpleButton{
		ButtonGUI: ButtonGUI{Icon: "cancel", X: screenX/2 - w + iconSize*1, Y: screenY / 2, SX: iconSize, SY: iconSize},
		ClickImpl: noFn})
	p.AddButton(&SimpleButton{
		ButtonGUI: ButtonGUI{Icon: "click", X: screenX/2 + w - iconSize*2, Y: screenY / 2, SX: iconSize, SY: iconSize},
		ClickImpl: yesFn})
	return p
}
