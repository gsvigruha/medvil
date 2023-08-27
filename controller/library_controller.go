package controller

import (
	"fmt"
	"io/ioutil"
	"medvil/maps"
	"medvil/view/gui"
	"path/filepath"
	"time"
)

func LibraryToControlPanel(cp *ControlPanel) {
	p := &gui.Panel{X: 0, Y: ControlPanelDynamicPanelTop, SX: ControlPanelSX, SY: HouseholdControllerSY}

	p.AddLargeTextLabel("New", 24, ControlPanelSY*0.15)
	nTop := ControlPanelSY * 0.15
	config := &maps.MapConfig{Size: 100, Hills: 5, Lakes: 5, Trees: 5, Resources: 5}

	p.AddPanel(gui.CreateNumberPaneFromVal(24, nTop+float64(IconS*1), 200, gui.FontSize, 100, 200, 50, "Map size %v", &config.Size).P)
	p.AddPanel(gui.CreateNumberPaneFromVal(24, nTop+float64(IconS*2), 200, gui.FontSize, 3, 10, 1, "Hills %v", &config.Hills).P)
	p.AddPanel(gui.CreateNumberPaneFromVal(24, nTop+float64(IconS*3), 200, gui.FontSize, 3, 10, 1, "Lakes %v", &config.Lakes).P)
	p.AddPanel(gui.CreateNumberPaneFromVal(24, nTop+float64(IconS*4), 200, gui.FontSize, 3, 10, 1, "Trees %v", &config.Trees).P)
	p.AddPanel(gui.CreateNumberPaneFromVal(24, nTop+float64(IconS*5), 200, gui.FontSize, 3, 10, 1, "Resources %v", &config.Resources).P)

	p.AddButton(gui.SimpleButton{
		ButtonGUI: gui.ButtonGUI{Icon: "plus", X: 24, Y: nTop + float64(IconS*7), SX: IconS, SY: IconS},
		ClickImpl: func() {
			go cp.C.NewMap(*config)
		}})

	p.AddLargeTextLabel("Load and save", 24, ControlPanelSY*0.45)

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

	lasTop := ControlPanelSY*0.45 + float64(IconH)
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
			go cp.C.Load(filesDropdown.GetSelectedValue())
			CPActionCancel(cp.C)
		}})

	savedButton := gui.SimpleButton{
		ButtonGUI: gui.ButtonGUI{Icon: "save", X: float64(24 + IconW*1), Y: lasTop + float64(IconH*2), SX: IconS, SY: IconS},
		ClickImpl: func() {
			go cp.C.Save(filesDropdown.GetSelectedValue())
		}}
	savedButton.Disabled = func() bool { return cp.C.Map == nil }
	p.AddButton(savedButton)

	plusButton := gui.SimpleButton{
		ButtonGUI: gui.ButtonGUI{Icon: "plus", X: float64(24 + IconW*2), Y: lasTop + float64(IconH*2), SX: IconS, SY: IconS},
		ClickImpl: func() {
			go cp.C.Save(time.Now().Format(time.RFC3339) + ".mdvl")
			CPActionCancel(cp.C)
		}}
	plusButton.Disabled = func() bool { return cp.C.Map == nil }
	p.AddButton(plusButton)

	p.AddLargeTextLabel("Settings", 24, ControlPanelSY*0.7)
	p.AddButton(gui.SimpleButton{
		ButtonGUI: gui.ButtonGUI{Icon: "warning", X: 24, Y: ControlPanelSY*0.7 + float64(IconH), SX: IconS, SY: IconS},
		Highlight: func() bool { return cp.C.ShowHouseIcons },
		ClickImpl: func() { cp.C.ShowHouseIcons = !cp.C.ShowHouseIcons }})

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
