package controller

import (
	"fmt"
	"io/ioutil"
	"medvil/view/gui"
	"path/filepath"
	"time"
)

func LibraryToControlPanel(cp *ControlPanel) {
	p := &gui.Panel{X: 0, Y: ControlPanelDynamicPanelTop, SX: ControlPanelSX, SY: HouseholdControllerSY}

	files, err := ioutil.ReadDir("saved/")
	if err != nil {
		fmt.Println(err)
	}

	var savedGames []string
	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".mdvl" {
			savedGames = append(savedGames, file.Name())
		}
	}

	var icons []string
	for _, _ = range savedGames {
		icons = append(icons, "library")
	}

	filesDropdown := &gui.DropDown{
		X:        float64(10),
		Y:        ControlPanelSY * 0.15,
		SX:       IconS + gui.FontSize*16,
		SY:       IconS,
		Options:  savedGames,
		Icons:    icons,
		Selected: -1,
	}
	p.AddDropDown(filesDropdown)

	p.AddButton(gui.SimpleButton{
		ButtonGUI: gui.ButtonGUI{Icon: "load", X: float64(10 + IconW*0), Y: ControlPanelSY*0.15 + float64(IconH*2), SX: IconS, SY: IconS},
		ClickImpl: func() {
			cp.C.Load(filesDropdown.GetSelectedValue())
			CPActionCancel(cp.C)
		}})
	p.AddButton(gui.SimpleButton{
		ButtonGUI: gui.ButtonGUI{Icon: "save", X: float64(10 + IconW*1), Y: ControlPanelSY*0.15 + float64(IconH*2), SX: IconS, SY: IconS},
		ClickImpl: func() { cp.C.Save(filesDropdown.GetSelectedValue()) }})
	p.AddButton(gui.SimpleButton{
		ButtonGUI: gui.ButtonGUI{Icon: "plus", X: float64(10 + IconW*2), Y: ControlPanelSY*0.15 + float64(IconH*2), SX: IconS, SY: IconS},
		ClickImpl: func() {
			cp.C.Save(time.Now().Format(time.RFC3339) + ".mdvl")
			CPActionCancel(cp.C)
		}})

	cp.SetDynamicPanel(p)
}
