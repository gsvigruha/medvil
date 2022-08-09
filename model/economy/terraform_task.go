package economy

import (
	"medvil/model/navigation"
	"medvil/model/time"
	"strconv"
)

const TerraformTaskTypeLevelForBuilding uint8 = 1
const TerraformTaskTypeLevelForRoad uint8 = 2

const TerraformTaskMaxProgress = 30 * 24

type TerraformTask struct {
	TaskBase
	M        navigation.IMap
	F        *navigation.Field
	T        uint8
	progress uint16
}

func (t *TerraformTask) Field() *navigation.Field {
	return t.F
}

func (t *TerraformTask) Complete(Calendar *time.CalendarType, tool bool) bool {
	if t.progress < TerraformTaskMaxProgress {
		t.progress++
		return false
	} else {
		switch t.T {
		case TerraformTaskTypeLevelForBuilding:
			navigation.LevelFieldForBuilding(t.F, t.M)
		case TerraformTaskTypeLevelForRoad:
			navigation.LevelFieldForRoad(t.F, t.M)
		}
		t.F.Construction = false
		return true
	}
}

func (t *TerraformTask) Blocked() bool {
	return !t.F.Empty()
}

func (t *TerraformTask) Name() string {
	return "terraform"
}

func (t *TerraformTask) Tag() string {
	return TerraformTaskTag(t.F)
}

func TerraformTaskTag(f *navigation.Field) string {
	return strconv.Itoa(int(f.X)) + "#" + strconv.Itoa(int(f.X))
}

func (t *TerraformTask) Expired(Calendar *time.CalendarType) bool {
	return false
}

func (t *TerraformTask) Motion() uint8 {
	return navigation.MotionStand
}
