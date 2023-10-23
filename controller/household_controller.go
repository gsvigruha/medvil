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

const IconRowMax = 9

var PersonGUIY = 0.175
var ArtifactsGUIY = 0.45
var TaskGUIY = 0.6

const MaxNumTasks = 24

var VehicleGUIY = 0.7
var HouseholdControllerSY = 0.7
var HouseholdControllerGUIBottomY = 0.75

type HouseholdControllerButton struct {
	b      *gui.ButtonGUI
	h      *social.Household
	action func(*social.Household)
}

func (b *HouseholdControllerButton) SetHoover(h bool) {
	b.b.SetHoover(h)
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

func (b HouseholdControllerButton) Enabled() bool {
	return b.b.Enabled()
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

func taskIconW(h *social.Household) (int, int) {
	var numTasks = len(h.Tasks)
	var w = IconW
	var n = IconRowMax
	if numTasks > IconRowMax*2 {
		if numTasks >= MaxNumTasks {
			w = IconRowMax * 2 * IconW / MaxNumTasks
			n = MaxNumTasks / 2
		} else {
			w = IconRowMax * 2 * IconW / numTasks
			n = (numTasks + 1) / 2
		}
	}
	return w, n
}

func HouseholdToControlPanel(cp *ControlPanel, p *gui.Panel, h *social.Household, name string) {
	if h.Town.Townhall.Household == h && h.Town.Supplier != nil {
		MoneyToControlPanel(p, h.Town.Supplier.GetHome(), h, 100, 24, LargeIconD*2+float64(IconH)+24)
	} else {
		MoneyToControlPanel(p, h.Town.Townhall.Household, h, 100, 24, LargeIconD*2+float64(IconH)+24)
	}
	p.AddTextLabel(name+" / "+h.Town.Name, 200, LargeIconD*2+float64(IconH)+24)

	piw := personIconW(h)
	for i, person := range h.People {
		PersonToPanel(cp, p, i, person, piw, PersonGUIY*ControlPanelSY)
	}
	for i := len(h.People); i < int(h.TargetNumPeople); i++ {
		p.AddImageLabel("person", float64(24+i*piw), PersonGUIY*ControlPanelSY, IconS, IconS, gui.ImageLabelStyleDisabled)
	}
	s := IconS / 2
	p.AddButton(&HouseholdControllerButton{
		b: &gui.ButtonGUI{Icon: "plus", X: ControlPanelSX - 24 - s, Y: PersonGUIY * ControlPanelSY, SX: s, SY: s, OnHoover: func() {
			cp.HelperMessage("Add people to this household")
		}},
		h: h, action: IncreaseHouseholdTargetNumPeople})
	p.AddButton(&HouseholdControllerButton{
		b: &gui.ButtonGUI{Icon: "minus", X: ControlPanelSX - 24 - s, Y: PersonGUIY*ControlPanelSY + s, SX: s, SY: s, OnHoover: func() {
			cp.HelperMessage("Remove people from this household")
		}},
		h: h, action: DecreaseHouseholdTargetNumPeople})
	p.AddScaleLabel("heating", 24, ArtifactsGUIY*ControlPanelSY, IconS, IconS, 4, float64(h.GetHeating())/100, false)
	p.AddScaleLabel("barrel", 24+float64(IconW), ArtifactsGUIY*ControlPanelSY, IconS, IconS, 4, h.Resources.UsedVolumeCapacity(), false)
	var aI = 2
	for _, a := range artifacts.All {
		if q, ok := h.Resources.Artifacts[a]; ok {
			ArtifactsToControlPanel(cp, p, aI, a, q, ArtifactsGUIY*ControlPanelSY)
			aI++
		}
	}
	tiw, tirm := taskIconW(h)
	for i, task := range h.Tasks {
		if i >= MaxNumTasks {
			tasksStr := strconv.Itoa(len(h.Tasks))
			p.AddTextLabel(tasksStr, ControlPanelSX-24-float64(len(tasksStr))*gui.FontSize*0.5, TaskGUIY*ControlPanelSY+float64(IconH*2))
			break
		}
		TaskToControlPanel(cp, p, i%tirm, TaskGUIY*ControlPanelSY+float64(i/tirm*IconH), task, tiw)
	}
	p.AddButton(&gui.SimpleButton{
		ButtonGUI: gui.ButtonGUI{Icon: "vehicles/boat", X: 24, Y: VehicleGUIY * ControlPanelSY, SX: IconS, SY: IconS},
		Highlight: func() bool { return h.IsBoatEnabled() },
		ClickImpl: func() {
			h.BoatEnabled = !h.BoatEnabled
			cp.HelperMessage("Start or stop using boats and waterways")
		}})
	for i, vehicle := range h.Vehicles {
		VehicleToControlPanel(p, i, VehicleGUIY*ControlPanelSY, vehicle, IconW)
	}
}

func PersonToPanel(cp *ControlPanel, p *gui.Panel, i int, person *social.Person, w int, top float64) {
	p.AddImageLabel("person", float64(24+i*w), top, IconS, IconS, gui.ImageLabelStyleRegular)
	if person.Equipment.Weapon {
		p.AddImageLabel("tasks/swordsmith", float64(24+i*w)+16, top+16, 24, 24, gui.ImageLabelStyleRegular)
	} else if person.Equipment.Tool {
		p.AddImageLabel("tasks/toolsmith", float64(24+i*w)+16, top+16, 24, 24, gui.ImageLabelStyleRegular)
	}
	p.AddScaleLabel("food", float64(24+i*w), top+float64(IconH), IconS, IconS, 4, float64(person.Food)/float64(social.MaxPersonState), false)
	p.AddScaleLabel("drink", float64(24+i*w), top+float64(IconH*2), IconS, IconS, 4, float64(person.Water)/float64(social.MaxPersonState), false)
	p.AddScaleLabel("health", float64(24+i*w), top+float64(IconH*3), IconS, IconS, 4, float64(person.Health)/float64(social.MaxPersonState), false)
	p.AddScaleLabel("happiness", float64(24+i*w), top+float64(IconH*4), IconS, IconS, 4, float64(person.Happiness)/float64(social.MaxPersonState), false)
	if person.Task != nil {
		TaskToControlPanel(cp, p, i, top+float64(IconH*5), person.Task, w)
	}
}

func ArtifactQStr(q uint16) string {
	var qStr = strconv.Itoa(int(q))
	if q == artifacts.InfiniteQuantity {
		qStr = "∞"
	}
	return qStr
}

func ArtifactsToControlPanel(cp *ControlPanel, p *gui.Panel, i int, a *artifacts.Artifact, q uint16, top float64) {
	xI := i % IconRowMax
	yI := i / IconRowMax
	p.AddButton(&gui.ImageButton{
		ButtonGUI: gui.ButtonGUI{Icon: "artifacts/" + a.Name, X: float64(24 + xI*IconW), Y: top + float64(yI*IconH), SX: IconS, SY: IconS},
		ClickImpl: func() {
			ArtifactToHelperPanel(cp.GetHelperPanel(), a)
		},
	})
	p.AddTextLabel(ArtifactQStr(q), float64(24+xI*IconW), top+float64(yI*IconH+IconH+4))
}

func TaskToControlPanel(cp *ControlPanel, p *gui.Panel, i int, y float64, task economy.Task, w int) {
	var style uint8 = gui.ImageLabelStyleHighlight
	if task.Blocked() {
		style = gui.ImageLabelStyleDisabled
	}
	p.AddButton(&gui.ImageButton{
		ButtonGUI: gui.ButtonGUI{Icon: "tasks/" + economy.IconName(task), X: float64(24 + i*w), Y: y, SX: IconS, SY: IconS},
		Style:     style,
		ClickImpl: func() {
			TaskToHelperPanel(cp.GetHelperPanel(), task)
		},
	})
}

func VehicleToControlPanel(p *gui.Panel, i int, y float64, vehicle *vehicles.Vehicle, w int) {
	var style uint8 = gui.ImageLabelStyleHighlight
	if !vehicle.InUse {
		style = gui.ImageLabelStyleDisabled
	}
	p.AddImageLabel("vehicles/"+vehicle.T.Name, float64(24+i*w+IconW), y, IconS, IconS, style)
}

func GetHouseholdHelperSuggestions(h *social.Household) *gui.Suggestion {
	if h.TargetNumPeople < 2 {
		return &gui.Suggestion{Message: "Add people to your house.\nPeople will move over from the townhall.", Icon: "person", X: ControlPanelSX - 24, Y: PersonGUIY*ControlPanelSY + IconS/4}
	}
	return nil
}
