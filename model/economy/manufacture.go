package economy

import (
	"medvil/model/artifacts"
)

type Manufacture struct {
	Name    string
	Time    uint8
	Power   uint16
	Inputs  []artifacts.Artifacts
	Outputs []artifacts.Artifacts
}

var AllManufacture = [...]*Manufacture{
	&Manufacture{
		Name:    "sawing",
		Time:    2,
		Power:   1000,
		Inputs:  []artifacts.Artifacts{artifacts.Artifacts{A: artifacts.GetArtifact("log"), Quantity: 1}},
		Outputs: []artifacts.Artifacts{artifacts.Artifacts{A: artifacts.GetArtifact("board"), Quantity: 1}}},
	&Manufacture{
		Name:    "stonecutting",
		Time:    2,
		Power:   1000,
		Inputs:  []artifacts.Artifacts{artifacts.Artifacts{A: artifacts.GetArtifact("rock"), Quantity: 1}},
		Outputs: []artifacts.Artifacts{artifacts.Artifacts{A: artifacts.GetArtifact("cube"), Quantity: 1}}},
}

func GetAllManufactureNames() []string {
	result := make([]string, len(AllManufacture))
	for i, m := range AllManufacture {
		result[i] = m.Name
	}
	return result
}

func GetManufacture(name string) *Manufacture {
	for i := 0; i < len(AllManufacture); i++ {
		if AllManufacture[i].Name == name {
			return AllManufacture[i]
		}
	}
	return nil
}
