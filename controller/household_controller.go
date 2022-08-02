package controller

import (
	"github.com/tfriedel6/canvas"
	"medvil/model/artifacts"
	"medvil/model/economy"
	"medvil/model/social"
	"medvil/model/vehicles"
	"medvil/view/gui"
	"strconv"
)

const IconH = 40
const IconW = 40
const IconRowMax = 7
const PersonGUIY = 110
const ArtifactsGUIY = 350
const TaskGUIY = 440
const MaxNumTasks = 24
const VehicleGUIY = 600
const HouseholdControllerSY = 550
const HouseholdControllerGUIBottomY = ControlPanelDynamicPanelTop + HouseholdControllerSY

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

func personIconW(h *social.Household) int {
	var numPeople = int(h.TargetNumPeople)
	if len(h.People) > numPeople {
		numPeople = len(h.People)
	}
	var w = IconW
	if numPeople > 6 {
		w = 6 * IconW / numPeople
	}
	return w
}

func HouseholdToControlPanel(p *gui.Panel, h *social.Household) {
	MoneyToControlPanel(p, h.Town, &h.Money, 100, 10, 80)
	piw := personIconW(h)
	for i, person := range h.People {
		PersonToControlPanel(p, i, person, piw)
	}
	for i := len(h.People); i < int(h.TargetNumPeople); i++ {
		p.AddImageLabel("person", float64(10+i*piw), PersonGUIY, 32, 32, gui.ImageLabelStyleDisabled)
	}
	p.AddButton(HouseholdControllerButton{
		b: gui.ButtonGUI{Icon: "plus", X: ControlPanelSX - 40, Y: PersonGUIY, SX: 16, SY: 16},
		h: h, action: IncreaseHouseholdTargetNumPeople})
	p.AddButton(HouseholdControllerButton{
		b: gui.ButtonGUI{Icon: "minus", X: ControlPanelSX - 40, Y: PersonGUIY + 16, SX: 16, SY: 16},
		h: h, action: DecreaseHouseholdTargetNumPeople})
	p.AddScaleLabel("heating", 10, ArtifactsGUIY, 32, 32, 4, h.Heating, false)
	p.AddScaleLabel("barrel", 50, ArtifactsGUIY, 32, 32, 4, h.Resources.UsedVolumeCapacity(), false)
	var aI = 2
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
		TaskToControlPanel(p, i%IconRowMax, float64(TaskGUIY+i/IconRowMax*IconH), task, IconW)
	}
	for i, vehicle := range h.Vehicles {
		VehicleToControlPanel(p, i, VehicleGUIY, vehicle, IconW)
	}
}

func PersonToControlPanel(p *gui.Panel, i int, person *social.Person, w int) {
	p.AddImageLabel("person", float64(10+i*w), PersonGUIY, 32, 32, gui.ImageLabelStyleRegular)
	if person.Equipment.Weapon() {
		p.AddImageLabel("tasks/swordsmith", float64(10+i*w)+16, PersonGUIY+16, 24, 24, gui.ImageLabelStyleRegular)
	} else if person.Equipment.Tool() {
		p.AddImageLabel("tasks/toolsmith", float64(10+i*w)+16, PersonGUIY+16, 24, 24, gui.ImageLabelStyleRegular)
	}
	p.AddScaleLabel("food", float64(10+i*w), PersonGUIY+IconH, 32, 32, 4, float64(person.Food)/float64(social.MaxPersonState), false)
	p.AddScaleLabel("drink", float64(10+i*w), PersonGUIY+IconH*2, 32, 32, 4, float64(person.Water)/float64(social.MaxPersonState), false)
	p.AddScaleLabel("health", float64(10+i*w), PersonGUIY+IconH*3, 32, 32, 4, float64(person.Health)/float64(social.MaxPersonState), false)
	p.AddScaleLabel("happiness", float64(10+i*w), PersonGUIY+IconH*4, 32, 32, 4, float64(person.Happiness)/float64(social.MaxPersonState), false)
	if person.Task != nil {
		TaskToControlPanel(p, i, PersonGUIY+IconH*5, person.Task, w)
	}
}

func ArtifactsToControlPanel(p *gui.Panel, i int, a *artifacts.Artifact, q uint16, top float64) {
	xI := i % IconRowMax
	yI := i / IconRowMax
	p.AddImageLabel("artifacts/"+a.Name, float64(10+xI*IconW), top+float64(yI)*IconH, 32, 32, gui.ImageLabelStyleRegular)
	p.AddTextLabel(strconv.Itoa(int(q)), float64(10+xI*IconW), top+float64(yI)*IconH+IconH+4)
}

func TaskToControlPanel(p *gui.Panel, i int, y float64, task economy.Task, w int) {
	var style uint8 = gui.ImageLabelStyleHighlight
	if task.Blocked() {
		style = gui.ImageLabelStyleDisabled
	}
	p.AddImageLabel("tasks/"+task.Name(), float64(10+i*w), y, 32, 32, style)
}

func VehicleToControlPanel(p *gui.Panel, i int, y float64, vehicle *vehicles.Vehicle, w int) {
	var style uint8 = gui.ImageLabelStyleHighlight
	if !vehicle.InUse {
		style = gui.ImageLabelStyleDisabled
	}
	p.AddImageLabel("vehicles/"+vehicle.T.Name, float64(10+i*w), y, 32, 32, style)
}
