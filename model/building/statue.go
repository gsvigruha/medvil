package building

import (
	"medvil/model/artifacts"
)

type StatueType struct {
	Name      string
	Happiness uint8
	Cost      []artifacts.Artifacts
}

var FountainType = &StatueType{
	Name: "fountain", Happiness: 10, Cost: []artifacts.Artifacts{artifacts.Artifacts{artifacts.GetArtifact("cube"), 2}},
}

var StatueTypes = [...]*StatueType{
	FountainType,
}

func GetStatueType(name string) *StatueType {
	for _, t := range StatueTypes {
		if t.Name == name {
			return t
		}
	}
	return nil
}

type Statue struct {
	T            *StatueType
	Construction bool
}
