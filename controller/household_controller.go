package controller

import (
	"github.com/tfriedel6/canvas"
	"medvil/model/artifacts"
	"medvil/model/economy"
	"medvil/model/social"
	"medvil/view/gui"
	"strconv"
)

const IconH = 40
const IconW = 40
const IconRowMax = 6
const PersonGUIY = 110
const ArtifactsGUIY = 350
const TaskGUIY = 400
const MaxNumTasks = 24
const HouseholdControllerGUIBottomY = 580

type HouseholdControllerButton struct {
	b      gui.ButtonGUI
	h      *social.Household
	action func(*social.Household)
}

func (b HouseholdControllerButton) Click() {
	b.action(b.h)
}

func (b HouseholdControllerButton) Render(cv *canvas.Canvas) {
	b.b.Render(cv)
}

func (b HouseholdControllerButton) Contains(x float64, y float64) bool {
	return b.b.Contains(x, y)
}

func IncreaseHouseholdTargetNumPeople(h *social.Household) {
	h.IncTargetNumPeople()
}

func DecreaseHouseholdTargetNumPeople(h *social.Household) {
	h.DecTargetNumPeople()
}

func HouseholdToControlPanel(p *gui.Panel, h *social.Household) {
	MoneyToControlPanel(p, h.Town, &h.Money, 100, 10, 80)
	for i, person := range h.People {
		PersonToControlPanel(p, i, person)
	}
	for i := uint16(len(h.People)); i < h.TargetNumPeople; i++ {
		p.AddImageLabel("person", float64(10+i*IconW), PersonGUIY, 32, 32, gui.ImageLabelStyleDisabled)
	}
	p.AddButton(HouseholdControllerButton{
		b: gui.ButtonGUI{Icon: "plus", X: ControlPanelSX - 40, Y: PersonGUIY, SX: 16, SY: 16},
		h: h, action: IncreaseHouseholdTargetNumPeople})
	p.AddButton(HouseholdControllerButton{
		b: gui.ButtonGUI{Icon: "minus", X: ControlPanelSX - 40, Y: PersonGUIY + 16, SX: 16, SY: 16},
		h: h, action: DecreaseHouseholdTargetNumPeople})
	p.AddScaleLabel("barrel", 10, ArtifactsGUIY, 32, 32, 4, h.Resources.UsedVolumeCapacity(), false)
	var aI = 1
	for _, a := range artifacts.All {
		if q, ok := h.Resources.Artifacts[a]; ok {
			ArtifactsToControlPanel(p, aI, a, q, ArtifactsGUIY)
			aI++
		}
	}
	for i, task := range h.Tasks {
		if i >= MaxNumTasks {
			break
		}
		TaskToControlPanel(p, i%IconRowMax, float64(TaskGUIY+i/IconRowMax*IconH), task)
	}
}

func PersonToControlPanel(p *gui.Panel, i int, person *social.Person) {
	p.AddImageLabel("person", float64(10+i*IconW), PersonGUIY, 32, 32, gui.ImageLabelStyleRegular)
	p.AddScaleLabel("food", float64(10+i*IconW), PersonGUIY+IconH, 32, 32, 4, float64(person.Food)/float64(social.MaxPersonState), false)
	p.AddScaleLabel("drink", float64(10+i*IconW), PersonGUIY+IconH*2, 32, 32, 4, float64(person.Water)/float64(social.MaxPersonState), false)
	p.AddScaleLabel("health", float64(10+i*IconW), PersonGUIY+IconH*3, 32, 32, 4, float64(person.Health)/float64(social.MaxPersonState), false)
	p.AddScaleLabel("happiness", float64(10+i*IconW), PersonGUIY+IconH*4, 32, 32, 4, float64(person.Happiness)/float64(social.MaxPersonState), false)
	if person.Task != nil {
		TaskToControlPanel(p, i, PersonGUIY+IconH*5, person.Task)
	}
}

func ArtifactsToControlPanel(p *gui.Panel, i int, a *artifacts.Artifact, q uint16, top float64) {
	p.AddImageLabel("artifacts/"+a.Name, float64(10+i*IconW), top, 32, 32, gui.ImageLabelStyleRegular)
	p.AddTextLabel(strconv.Itoa(int(q)), float64(10+i*IconW), top+IconH+4)
}

func TaskToControlPanel(p *gui.Panel, i int, y float64, task economy.Task) {
	var style uint8 = gui.ImageLabelStyleHighlight
	if task.Blocked() {
		style = gui.ImageLabelStyleDisabled
	}
	p.AddImageLabel("tasks/"+task.Name(), float64(10+i*IconW), y, 32, 32, style)
}
