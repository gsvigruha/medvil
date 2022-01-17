package controller

import (
	"medvil/model/artifacts"
	"medvil/model/economy"
	"medvil/model/social"
	"medvil/view/gui"
	"strconv"
)

const IconH = 40
const IconW = 40
const IconRowMax = 6
const PersonGUIY = 50
const ArtifactsGUIY = 200
const TaskGUIY = 300

func HouseholdToControlPanel(p *gui.Panel, h *social.Household) {
	p.AddTextLabel("money "+strconv.Itoa(int(h.Money)), 10, 50)
	for i, person := range h.People {
		PersonToControlPanel(p, i, person)
	}
	var aI = 0
	for a, q := range h.Resources.Artifacts {
		ArtifactsToControlPanel(p, aI, a, q)
		aI++
	}
	for i, task := range h.Tasks {
		TaskToControlPanel(p, i%IconRowMax, float64(TaskGUIY+i/IconRowMax*IconH), task)
	}
}

func PersonToControlPanel(p *gui.Panel, i int, person *social.Person) {
	p.AddScaleLabel("food", float64(10+i*IconW), PersonGUIY, 32, 32, 4, float64(person.Food)/float64(social.MaxPersonState))
	p.AddScaleLabel("drink", float64(10+i*IconW), PersonGUIY+IconH, 32, 32, 4, float64(person.Water)/float64(social.MaxPersonState))
	if person.Task != nil {
		TaskToControlPanel(p, i, PersonGUIY+IconH*2, person.Task)
	}
}

func ArtifactsToControlPanel(p *gui.Panel, i int, a *artifacts.Artifact, q uint16) {
	p.AddImageLabel("artifacts/"+a.Name, float64(10+i*IconW), ArtifactsGUIY, 32, 32, gui.ImageLabelStyleRegular)
	p.AddTextLabel(strconv.Itoa(int(q)), float64(10+i*IconW), ArtifactsGUIY+IconH)
}

func TaskToControlPanel(p *gui.Panel, i int, y float64, task economy.Task) {
	var style uint8 = gui.ImageLabelStyleHighlight
	if task.Blocked() {
		style = gui.ImageLabelStyleDisabled
	}
	p.AddImageLabel("tasks/"+task.Name(), float64(10+i*IconW), y, 32, 32, style)
}
