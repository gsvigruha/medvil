package economy

import (
	"medvil/model/artifacts"
	"medvil/model/navigation"
	"medvil/model/time"
	"strconv"
)

const MineFieldUseTypeNone uint8 = 0
const MineFieldUseTypeStone uint8 = 1
const MineFieldUseTypeClay uint8 = 2
const MineFieldUseTypeIron uint8 = 3
const MineFieldUseTypeGold uint8 = 4

const MineTaskDurationStone = 24 * 60
const MineTaskQuantityStone = 2
const MineTaskDurationClay = 24 * 60
const MineTaskQuantityClay = 2
const MineTaskDurationIron = 24 * 60
const MineTaskQuantityIron = 2
const MineTaskDurationGold = 24 * 60
const MineTaskQuantityGold = 2

var stone = artifacts.GetArtifact("stone")
var clay = artifacts.GetArtifact("clay")
var iron = artifacts.GetArtifact("iron_ore")
var gold = artifacts.GetArtifact("gold_ore")

type MiningTask struct {
	TaskBase
	F        *navigation.Field
	Progress uint16
	UseType  uint8
}

func (t *MiningTask) Field() *navigation.Field {
	return t.F
}

func (t *MiningTask) Complete(Calendar *time.CalendarType, tool bool) bool {
	t.Progress++
	switch t.UseType {
	case MineFieldUseTypeStone:
		if t.Progress >= MineTaskDurationStone {
			t.F.Terrain.Resources.Add(stone, MineTaskQuantityStone)
			return true
		}
	case MineFieldUseTypeClay:
		if t.Progress >= MineTaskDurationClay {
			t.F.Terrain.Resources.Add(clay, MineTaskQuantityClay)
			return true
		}
	case MineFieldUseTypeIron:
		if t.Progress >= MineTaskDurationIron {
			t.F.Terrain.Resources.Add(iron, MineTaskQuantityIron)
			return true
		}
	case MineFieldUseTypeGold:
		if t.Progress >= MineTaskDurationGold {
			t.F.Terrain.Resources.Add(gold, MineTaskQuantityGold)
			return true
		}
	}
	return false
}

func (t *MiningTask) Blocked() bool {
	return false
}

func (t *MiningTask) Name() string {
	return "mining"
}

func MiningTaskTag(f *navigation.Field, ut uint8) string {
	return strconv.Itoa(int(f.X)) + "#" + strconv.Itoa(int(f.Y)) + "#" + strconv.Itoa(int(ut))
}

func (t *MiningTask) Tag() string {
	return MiningTaskTag(t.F, t.UseType)
}

func (t *MiningTask) Expired(Calendar *time.CalendarType) bool {
	return false
}

func (t *MiningTask) Motion() uint8 {
	return navigation.MotionMine
}
