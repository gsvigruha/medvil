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

	p.AddTextLabel("Load and save", 24, ControlPanelSY*0.15)

	lasTop := ControlPanelSY*0.15 + LargeIconD
	filesDropdown := &gui.DropDown{
		X:        24,
		Y:        lasTop,
		SX:       IconS + gui.FontSize*16,
		SY:       IconS,
		Options:  savedGames,
		Icons:    icons,
		Selected: -1,
	}
	p.AddDropDown(filesDropdown)

	p.AddButton(gui.SimpleButton{
		ButtonGUI: gui.ButtonGUI{Icon: "load", X: float64(24 + IconW*0), Y: lasTop + float64(IconH*2), SX: IconS, SY: IconS},
		ClickImpl: func() {
			cp.C.Load(filesDropdown.GetSelectedValue())
			CPActionCancel(cp.C)
		}})

	savedButton := gui.SimpleButton{
		ButtonGUI: gui.ButtonGUI{Icon: "save", X: float64(24 + IconW*1), Y: lasTop + float64(IconH*2), SX: IconS, SY: IconS},
		ClickImpl: func() { cp.C.Save(filesDropdown.GetSelectedValue()) }}
	savedButton.Disabled = func() bool { return cp.C.Map == nil }
	p.AddButton(savedButton)

	plusButton := gui.SimpleButton{
		ButtonGUI: gui.ButtonGUI{Icon: "plus", X: float64(24 + IconW*2), Y: lasTop + float64(IconH*2), SX: IconS, SY: IconS},
		ClickImpl: func() {
			cp.C.Save(time.Now().Format(time.RFC3339) + ".mdvl")
			CPActionCancel(cp.C)
		}}
	plusButton.Disabled = func() bool { return cp.C.Map == nil }
	p.AddButton(plusButton)

	cp.SetDynamicPanel(p)
}

func GetLatestFile() string {
	files, err := ioutil.ReadDir("saved/")
	if err != nil {
		fmt.Println(err)
	}
	var latestModTime time.Time
	var latestFile string
	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".mdvl" && file.ModTime().After(latestModTime) {
			latestModTime = file.ModTime()
			latestFile = file.Name()
		}
	}
	return latestFile
}
