package controller

import (
	"fmt"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/tfriedel6/canvas"
	"io/ioutil"
	"medvil/maps"
	"medvil/model/stats"
	"medvil/view/gui"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type LibraryController struct {
	p                      *gui.Panel
	fileTextField          *gui.TextLabel
	historyLengthTextField *gui.TextLabel
}

func (lc *LibraryController) CaptureKey(key glfw.Key) {
	if (key >= glfw.KeyA && key <= glfw.KeyZ) || key == glfw.KeySpace || (key >= glfw.Key0 && key <= glfw.Key9) {
		if len(lc.fileTextField.Text) < 20 {
			lc.fileTextField.Text = lc.fileTextField.Text + strings.ToLower(string(key))
		}
	} else if key == glfw.KeyBackspace {
		if len(lc.fileTextField.Text) > 0 {
			lc.fileTextField.Text = lc.fileTextField.Text[:len(lc.fileTextField.Text)-1]
		}
	}
}

func (lc *LibraryController) CaptureMove(x, y float64) {
	lc.p.CaptureMove(x, y)
}

func (lc *LibraryController) CaptureClick(x, y float64) {
	lc.p.CaptureClick(x, y)
}

func (lc *LibraryController) Render(cv *canvas.Canvas) {
	lc.p.Render(cv)
}

func (lc *LibraryController) Clear() {}

func (lc *LibraryController) Refresh() {
}

func (lc *LibraryController) GetHelperSuggestions() *gui.Suggestion {
	return nil
}

func LibraryToControlPanel(cp *ControlPanel) {
	p := &gui.Panel{X: 0, Y: ControlPanelDynamicPanelTop, SX: ControlPanelSX, SY: HouseholdControllerSY}
	lc := &LibraryController{p: p}

	p.AddLargeTextLabel("New", 24, ControlPanelSY*0.15)
	nTop := ControlPanelSY * 0.15
	config := &maps.MapConfig{Size: 100, Hills: 5, Lakes: 5, Trees: 5, Resources: 5}

	w := ControlPanelSX / 2
	p.AddPanel(gui.CreateNumberPaneFromVal(24, nTop+float64(IconS*1), w, gui.FontSize, 100, 200, 50, "Map size %v", true, &config.Size).P)
	p.AddPanel(gui.CreateNumberPaneFromVal(24, nTop+float64(IconS*2), w, gui.FontSize, 3, 10, 1, "Hills %v", true, &config.Hills).P)
	p.AddPanel(gui.CreateNumberPaneFromVal(24, nTop+float64(IconS*3), w, gui.FontSize, 3, 10, 1, "Lakes %v", true, &config.Lakes).P)
	p.AddPanel(gui.CreateNumberPaneFromVal(24, nTop+float64(IconS*4), w, gui.FontSize, 3, 10, 1, "Trees %v", true, &config.Trees).P)
	p.AddPanel(gui.CreateNumberPaneFromVal(24, nTop+float64(IconS*5), w, gui.FontSize, 3, 10, 1, "Resources %v", true, &config.Resources).P)

	p.AddButton(&gui.SimpleButton{
		ButtonGUI: gui.ButtonGUI{Icon: "plus", X: 24, Y: nTop + float64(IconS*6.5), SX: IconS, SY: IconS, OnHoover: func() {
			cp.HelperMessage("Start a new game")
		}},
		ClickImpl: func() {
			go cp.C.NewMap(*config)
		}})
	p.AddTextLabel("Start a new game", 24+float64(LargeIconD), nTop+float64(IconS*7.16))

	p.AddLargeTextLabel("Load and save", 24, ControlPanelSY*0.45)

	files, err := ioutil.ReadDir(filepath.FromSlash("saved/"))
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
		X:        float64(24 + IconW*2),
		Y:        lasTop,
		SX:       IconS + gui.FontSize*12,
		SY:       IconS,
		Options:  savedGames,
		Icons:    icons,
		Selected: -1,
	}
	p.AddDropDown(filesDropdown)

	p.AddButton(&gui.SimpleButton{
		ButtonGUI: gui.ButtonGUI{Icon: "load", X: float64(24 + IconW*0), Y: lasTop, SX: IconS, SY: IconS, OnHoover: func() {
			cp.HelperMessage("Load game")
		}},
		ClickImpl: func() {
			go cp.C.Load(filesDropdown.GetSelectedValue())
			CPActionCancel(cp.C)
		}})

	saveButton := &gui.SimpleButton{
		ButtonGUI: gui.ButtonGUI{Icon: "save", X: float64(24 + IconW*1), Y: lasTop, SX: IconS, SY: IconS, OnHoover: func() {
			cp.HelperMessage("Save game to an existing file")
		}},
		ClickImpl: func() {
			go cp.C.Save(filesDropdown.GetSelectedValue())
		}}
	saveButton.Disabled = func() bool {
		return cp.C.Map == nil || strings.HasPrefix(filesDropdown.GetSelectedValue(), "example")
	}
	p.AddButton(saveButton)

	plusButton := &gui.SimpleButton{
		ButtonGUI: gui.ButtonGUI{Icon: "plus_save", X: float64(24 + IconW*1), Y: lasTop + float64(IconH*1), SX: IconS, SY: IconS, OnHoover: func() {
			cp.HelperMessage("Save game to a new file")
		}},
		ClickImpl: func() {
			go cp.C.Save(lc.fileTextField.Text + ".mdvl")
			CPActionCancel(cp.C)
		}}
	plusButton.Disabled = func() bool { return cp.C.Map == nil }
	p.AddButton(plusButton)

	lc.fileTextField = p.AddEditableTextLabel(float64(24+IconW*2), lasTop+float64(IconH*1), IconS+gui.FontSize*12, IconS)

	settingsTop := ControlPanelSY * 0.65
	p.AddLargeTextLabel("Settings", 24, settingsTop)
	p.AddButton(&gui.SimpleButton{
		ButtonGUI: gui.ButtonGUI{Icon: "warning_slim", X: 24 + float64(IconW)*0, Y: settingsTop + float64(IconH), SX: IconS, SY: IconS, OnHoover: func() {
			cp.HelperMessage("Show warning icons for households")
		}},
		Highlight: func() bool { return cp.C.ViewSettings.ShowHouseIcons },
		ClickImpl: func() {
			cp.C.ViewSettings.ShowHouseIcons = !cp.C.ViewSettings.ShowHouseIcons
			cp.C.SaveSettings()
		}})
	p.AddButton(&gui.SimpleButton{
		ButtonGUI: gui.ButtonGUI{Icon: "flag", X: 24 + float64(IconW)*1, Y: settingsTop + float64(IconH), SX: IconS, SY: IconS, OnHoover: func() {
			cp.HelperMessage("Show flags for allocated land")
		}},
		Highlight: func() bool { return cp.C.ViewSettings.ShowAllocatedFields },
		ClickImpl: func() {
			cp.C.ViewSettings.ShowAllocatedFields = !cp.C.ViewSettings.ShowAllocatedFields
			cp.C.SaveSettings()
		}})
	p.AddButton(&gui.SimpleButton{
		ButtonGUI: gui.ButtonGUI{Icon: "label", X: 24 + float64(IconW)*2, Y: settingsTop + float64(IconH), SX: IconS, SY: IconS, OnHoover: func() {
			cp.HelperMessage("Show town and expedition labels")
		}},
		Highlight: func() bool { return cp.C.ViewSettings.ShowLabels },
		ClickImpl: func() {
			cp.C.ViewSettings.ShowLabels = !cp.C.ViewSettings.ShowLabels
			cp.C.SaveSettings()
		}})
	p.AddButton(&gui.SimpleButton{
		ButtonGUI: gui.ButtonGUI{Icon: "help", X: 24 + float64(IconW)*3, Y: settingsTop + float64(IconH), SX: IconS, SY: IconS, OnHoover: func() {
			cp.HelperMessage("Show suggestions and descriptions")
		}},
		Highlight: func() bool { return cp.C.ViewSettings.ShowSuggestions },
		ClickImpl: func() {
			cp.C.ViewSettings.ShowSuggestions = !cp.C.ViewSettings.ShowSuggestions
			cp.C.SaveSettings()
		}})

	if cp.C.Map == nil {
		p.AddButton(&gui.SimpleButton{
			ButtonGUI: gui.ButtonGUI{Icon: "icon_size", X: 24, Y: settingsTop + float64(IconH)*2, SX: IconS, SY: IconS, OnHoover: func() {
				cp.HelperMessage("Adjust icon sizes")
			}},
			Highlight: func() bool { return cp.C.ViewSettings.Size == SizeAuto },
			ClickImpl: func() {
				cp.C.ViewSettings.Size = (cp.C.ViewSettings.Size + 1) % 4
				cp.C.SaveSettings()
				cp.SetupDims(cp.C.W, cp.C.H)
				cp.C.ShowLibraryController()
			}})
		p.AddDynamicTextLabel(func() string { return IconSizeStr(cp.C.ViewSettings.Size) }, 24, settingsTop+float64(IconH)*3)
		p.AddButton(&gui.SimpleButton{
			ButtonGUI: gui.ButtonGUI{Icon: "size", X: 24 + float64(IconW), Y: settingsTop + float64(IconH)*2, SX: IconS, SY: IconS, OnHoover: func() {
				cp.HelperMessage("Resolution: " + ResolutionStr(cp.C.ViewSettings.Resolution) + " (applied after restart)")
			}},
			Highlight: func() bool { return false },
			ClickImpl: func() {
				cp.C.ViewSettings.Resolution = (cp.C.ViewSettings.Resolution + 1) % 3
				cp.C.SaveSettings()
			}})
		p.AddDynamicTextLabel(func() string { return ResolutionStr(cp.C.ViewSettings.Resolution) }, 24+float64(IconW)-gui.FontSize*0.4, settingsTop+float64(IconH)*3)
		p.AddButton(&gui.SimpleButton{
			ButtonGUI: gui.ButtonGUI{Icon: "screen", X: 24 + float64(IconW)*2, Y: settingsTop + float64(IconH)*2, SX: IconS, SY: IconS, OnHoover: func() {
				cp.HelperMessage("Full screen (applied after restart)")
			}},
			Highlight: func() bool { return cp.C.ViewSettings.FullScreen },
			ClickImpl: func() {
				cp.C.ViewSettings.FullScreen = !cp.C.ViewSettings.FullScreen
				cp.C.SaveSettings()
			}})
		p.AddButton(&gui.SimpleButton{
			ButtonGUI: gui.ButtonGUI{Icon: "chart", X: 24 + float64(IconW)*3, Y: settingsTop + float64(IconH)*2, SX: IconS, SY: IconS, OnHoover: func() {
				cp.HelperMessage("Chart history length")
			}},
			Highlight: func() bool { return false },
			ClickImpl: func() {
				if stats.MaxHistory == 120 {
					stats.MaxHistory = 1200
				} else if stats.MaxHistory == 1200 {
					stats.MaxHistory = 2400
				} else {
					stats.MaxHistory = 120
				}
				cp.C.SaveSettings()
			}})
		p.AddDynamicTextLabel(func() string { return strconv.Itoa(stats.MaxHistory/12) + "Y" }, 24+float64(IconW)*3-gui.FontSize/2, settingsTop+float64(IconH)*3)
	}

	p.AddLargeTextLabel("Quit", 24, ControlPanelSY*0.85)
	p.AddButton(&gui.SimpleButton{
		ButtonGUI: gui.ButtonGUI{Icon: "cancel_gold", X: 24, Y: ControlPanelSY*0.85 + float64(IconH)/2, SX: IconS, SY: IconS,
			OnHoover: func() {
				cp.HelperMessage("Stop game in progress")
			},
			Disabled: func() bool { return cp.C.Map == nil }},
		ClickImpl: func() {
			cp.C.Save("latest_autosave.mdvl")
			cp.C.Map = nil
			cp.C.ShowLibraryController()
		}})
	p.AddButton(&gui.SimpleButton{
		ButtonGUI: gui.ButtonGUI{Icon: "quit", X: 24 + float64(IconW), Y: ControlPanelSY*0.85 + float64(IconH)/2, SX: IconS, SY: IconS, OnHoover: func() {
			cp.HelperMessage("Quit Medville")
		}},
		ClickImpl: func() {
			cp.C.Window.Close()
		}})

	cp.SetDynamicPanel(lc)
	cp.C.KeyHandler = lc
}

func GetLatestFile() string {
	files, err := ioutil.ReadDir(filepath.FromSlash("saved/"))
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
