package terrain

import (
	"bytes"
	"encoding/json"
	"math"
	"medvil/model/artifacts"
	"medvil/model/time"
	"strconv"
)

type PlantHabitatType uint8

const (
	Land  PlantHabitatType = 1
	Shore PlantHabitatType = 2
)

type PlantType struct {
	Name             string
	MaturityAgeYears uint8
	TreeT            *TreeType
	Yield            artifacts.Artifacts
	Tall             bool
	Habitat          PlantHabitatType
}

func (pt *PlantType) MarshalJSON() ([]byte, error) {
	return json.Marshal(pt.Name)
}

func (pt *PlantType) UnmarshalJSON(data []byte) error {
	s := bytes.NewBuffer(data).String()
	switch s {
	case "grain":
		*pt = *AllCropTypes[0]
	case "vegetables":
		*pt = *AllCropTypes[1]
	case "reed":
		*pt = *AllCropTypes[2]
	case "herb":
		*pt = *AllCropTypes[3]
	case "oak tree":
		*pt = *AllTreeTypes[0]
	case "apple tree":
		*pt = *AllTreeTypes[1]
	}
	return nil
}

type Plant struct {
	T             *PlantType
	X             uint16
	Y             uint16
	BirthDateDays uint32
	Shape         uint8
	Ripe          bool
}

func (p *Plant) IsTree() bool {
	return p.T.TreeT != nil
}

func (p *Plant) MaturityYears(Calendar *time.CalendarType) float64 {
	return math.Min(float64(p.AgeYears(Calendar)), float64(p.T.MaturityAgeYears))
}

func (p *Plant) Maturity(Calendar *time.CalendarType) float64 {
	return p.MaturityYears(Calendar) / float64(p.T.MaturityAgeYears)
}

func (p *Plant) IsMature(Calendar *time.CalendarType) bool {
	return p.Maturity(Calendar) == 1.0
}

func (p *Plant) AgeYears(Calendar *time.CalendarType) uint32 {
	return (Calendar.DaysElapsed() - p.BirthDateDays) / (30 * 12)
}

func (p *Plant) ElapseTime(Calendar *time.CalendarType) {
	if p.T.IsAnnual() {
		p.Ripe = (Calendar.DaysElapsed() - p.BirthDateDays) > 90
	} else if p.T.TreeT != nil {
		p.Ripe = Calendar.Month >= 7
	}
}

func (p *PlantType) IsAnnual() bool {
	return p.TreeT == nil && p.MaturityAgeYears <= 1.0
}

func GetPlantType(name string) *PlantType {
	for _, t := range AllCropTypes {
		if t.Name == name {
			return t
		}
	}
	for _, t := range AllTreeTypes {
		if t.Name == name {
			return t
		}
	}
	return nil
}

func (p *Plant) CacheKey(Calendar *time.CalendarType) string {
	if p.IsTree() {
		return p.T.Name + "#" + strconv.Itoa(int(p.Shape)) + "#" + strconv.Itoa(int(p.MaturityYears(Calendar)))
	} else {
		return p.T.Name + "#" + strconv.Itoa(int(p.Shape))
	}
}
