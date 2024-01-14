package economy

import (
	"medvil/model/navigation"
	"medvil/model/time"
)

const TerraformTaskTypeLevelForBuilding uint8 = 1
const TerraformTaskTypeLevelForRoad uint8 = 2

const TerraformTaskMaxProgress = 30 * 24

type TerraformTask struct {
	TaskBase
	M        navigation.IMap
	F        *navigation.Field
	T        uint8
	Progress uint16
}

func (t *TerraformTask) Destination() navigation.Destination {
	return t.F
}

func (t *TerraformTask) Complete(m navigation.IMap, tool bool) bool {
	if t.Progress < TerraformTaskMaxProgress {
		t.Progress++
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

func (t *TerraformTask) Tags() Tags {
	return MakeTags(TerraformTaskTag(t.F))
}

func TerraformTaskTag(f *navigation.Field) Tag {
	return SingleTag(f.X, f.Y)
}

func (t *TerraformTask) Expired(Calendar *time.CalendarType) bool {
	return false
}

func (t *TerraformTask) Motion() uint8 {
	return navigation.MotionFieldWork
}

func (t *TerraformTask) Description() string {
	return "Terraform land in order to build on it"
}
