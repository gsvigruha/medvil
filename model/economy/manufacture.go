package economy

import (
	"medvil/model/artifacts"
	"medvil/model/building"
)

type Manufacture struct {
	Name                  string
	Time                  uint16
	Power                 uint16
	BuildingExtensionType *building.BuildingExtensionType
	Inputs                []artifacts.Artifacts
	Outputs               []artifacts.Artifacts
	Description           string
}

var AllManufacture = [...]*Manufacture{
	&Manufacture{
		Name:                  "sawing",
		Time:                  30 * 24,
		Power:                 1000,
		BuildingExtensionType: building.Workshop,
		Inputs:                []artifacts.Artifacts{artifacts.Artifacts{A: artifacts.GetArtifact("log"), Quantity: 1}},
		Outputs:               []artifacts.Artifacts{artifacts.Artifacts{A: artifacts.GetArtifact("board"), Quantity: 1}}},
	&Manufacture{
		Name:                  "sawmill",
		Time:                  10 * 24,
		Power:                 1000,
		BuildingExtensionType: building.WaterMillWheel,
		Inputs:                []artifacts.Artifacts{artifacts.Artifacts{A: artifacts.GetArtifact("log"), Quantity: 1}},
		Outputs:               []artifacts.Artifacts{artifacts.Artifacts{A: artifacts.GetArtifact("board"), Quantity: 1}}},
	&Manufacture{
		Name:                  "stonecutting",
		Time:                  10 * 24,
		Power:                 1000,
		BuildingExtensionType: building.Workshop,
		Inputs:                []artifacts.Artifacts{artifacts.Artifacts{A: artifacts.GetArtifact("stone"), Quantity: 1}},
		Outputs:               []artifacts.Artifacts{artifacts.Artifacts{A: artifacts.GetArtifact("cube"), Quantity: 1}}},
	&Manufacture{
		Name:                  "tiling",
		Time:                  10 * 24,
		Power:                 1000,
		BuildingExtensionType: building.Kiln,
		Inputs: []artifacts.Artifacts{
			artifacts.Artifacts{A: artifacts.GetArtifact("clay"), Quantity: 2},
			artifacts.Artifacts{A: artifacts.GetArtifact("log"), Quantity: 1}},
		Outputs: []artifacts.Artifacts{artifacts.Artifacts{A: artifacts.GetArtifact("tile"), Quantity: 2}}},
	&Manufacture{
		Name:                  "brickmaking",
		Time:                  10 * 24,
		Power:                 1000,
		BuildingExtensionType: building.Kiln,
		Inputs: []artifacts.Artifacts{
			artifacts.Artifacts{A: artifacts.GetArtifact("clay"), Quantity: 2},
			artifacts.Artifacts{A: artifacts.GetArtifact("log"), Quantity: 1}},
		Outputs: []artifacts.Artifacts{artifacts.Artifacts{A: artifacts.GetArtifact("brick"), Quantity: 2}}},
	&Manufacture{
		Name:                  "pottery",
		Time:                  10 * 24,
		Power:                 1000,
		BuildingExtensionType: building.Kiln,
		Inputs: []artifacts.Artifacts{
			artifacts.Artifacts{A: artifacts.GetArtifact("clay"), Quantity: 2},
			artifacts.Artifacts{A: artifacts.GetArtifact("log"), Quantity: 1}},
		Outputs: []artifacts.Artifacts{artifacts.Artifacts{A: artifacts.GetArtifact("pot"), Quantity: 3}}},
	&Manufacture{
		Name:                  "thatching",
		Time:                  10 * 24,
		Power:                 1000,
		BuildingExtensionType: building.Workshop,
		Inputs:                []artifacts.Artifacts{artifacts.Artifacts{A: artifacts.GetArtifact("reed"), Quantity: 1}},
		Outputs:               []artifacts.Artifacts{artifacts.Artifacts{A: artifacts.GetArtifact("thatch"), Quantity: 1}}},
	&Manufacture{
		Name:                  "iron_smelting",
		Time:                  10 * 24,
		Power:                 1000,
		BuildingExtensionType: building.Forge,
		Inputs: []artifacts.Artifacts{
			artifacts.Artifacts{A: artifacts.GetArtifact("iron_ore"), Quantity: 1},
			artifacts.Artifacts{A: artifacts.GetArtifact("log"), Quantity: 1}},
		Outputs: []artifacts.Artifacts{artifacts.Artifacts{A: artifacts.GetArtifact("iron_bar"), Quantity: 1}}},
	&Manufacture{
		Name:                  "goldsmith",
		Time:                  10 * 24,
		Power:                 1000,
		BuildingExtensionType: building.Forge,
		Inputs: []artifacts.Artifacts{
			artifacts.Artifacts{A: artifacts.GetArtifact("gold_ore"), Quantity: 1},
			artifacts.Artifacts{A: artifacts.GetArtifact("log"), Quantity: 1}},
		Outputs: []artifacts.Artifacts{artifacts.Artifacts{A: artifacts.GetArtifact("gold_coin"), Quantity: 1}}},
	&Manufacture{
		Name:                  "milling",
		Time:                  10 * 24,
		Power:                 1000,
		BuildingExtensionType: building.WaterMillWheel,
		Inputs:                []artifacts.Artifacts{artifacts.Artifacts{A: artifacts.GetArtifact("grain"), Quantity: 1}},
		Outputs:               []artifacts.Artifacts{artifacts.Artifacts{A: artifacts.GetArtifact("flour"), Quantity: 1}}},
	&Manufacture{
		Name:                  "baking",
		Time:                  10 * 24,
		Power:                 1000,
		BuildingExtensionType: building.Cooker,
		Inputs: []artifacts.Artifacts{
			artifacts.Artifacts{A: artifacts.GetArtifact("flour"), Quantity: 1},
			artifacts.Artifacts{A: artifacts.GetArtifact("water"), Quantity: 1}},
		Outputs: []artifacts.Artifacts{artifacts.Artifacts{A: artifacts.GetArtifact("bread"), Quantity: 2}}},
	&Manufacture{
		Name:                  "brewing",
		Time:                  60 * 24,
		Power:                 1000,
		BuildingExtensionType: building.Cooker,
		Inputs: []artifacts.Artifacts{
			artifacts.Artifacts{A: artifacts.GetArtifact("grain"), Quantity: 1},
			artifacts.Artifacts{A: artifacts.GetArtifact("pot"), Quantity: 2},
			artifacts.Artifacts{A: artifacts.GetArtifact("water"), Quantity: 2}},
		Outputs: []artifacts.Artifacts{artifacts.Artifacts{A: artifacts.GetArtifact("beer"), Quantity: 2}}},
	&Manufacture{
		Name:                  "butchering",
		Time:                  90 * 24,
		Power:                 1000,
		BuildingExtensionType: building.Workshop,
		Inputs:                []artifacts.Artifacts{artifacts.Artifacts{A: artifacts.GetArtifact("sheep"), Quantity: 1}},
		Outputs: []artifacts.Artifacts{
			artifacts.Artifacts{A: artifacts.GetArtifact("meat"), Quantity: 10},
			artifacts.Artifacts{A: artifacts.GetArtifact("leather"), Quantity: 3}}},
	&Manufacture{
		Name:                  "toolsmith",
		Time:                  30 * 24,
		Power:                 1000,
		BuildingExtensionType: building.Forge,
		Inputs: []artifacts.Artifacts{
			artifacts.Artifacts{A: artifacts.GetArtifact("iron_bar"), Quantity: 1},
			artifacts.Artifacts{A: artifacts.GetArtifact("board"), Quantity: 1},
			artifacts.Artifacts{A: artifacts.GetArtifact("leather"), Quantity: 1}},
		Outputs: []artifacts.Artifacts{artifacts.Artifacts{A: artifacts.GetArtifact("tools"), Quantity: 1}}},
	&Manufacture{
		Name:                  "papermill",
		Time:                  30 * 24,
		Power:                 1000,
		BuildingExtensionType: building.WaterMillWheel,
		Inputs: []artifacts.Artifacts{
			artifacts.Artifacts{A: artifacts.GetArtifact("reed"), Quantity: 1},
			artifacts.Artifacts{A: artifacts.GetArtifact("water"), Quantity: 1}},
		Outputs: []artifacts.Artifacts{artifacts.Artifacts{A: artifacts.GetArtifact("paper"), Quantity: 3}}},
	&Manufacture{
		Name:                  "sewing",
		Time:                  30 * 24,
		Power:                 1000,
		BuildingExtensionType: building.Workshop,
		Inputs:                []artifacts.Artifacts{artifacts.Artifacts{A: artifacts.GetArtifact("wool"), Quantity: 2}},
		Outputs:               []artifacts.Artifacts{artifacts.Artifacts{A: artifacts.GetArtifact("textile"), Quantity: 1}}},
	&Manufacture{
		Name:                  "medicine",
		Time:                  120 * 24,
		Power:                 1000,
		BuildingExtensionType: building.Cooker,
		Inputs: []artifacts.Artifacts{
			artifacts.Artifacts{A: artifacts.GetArtifact("herb"), Quantity: 1},
			artifacts.Artifacts{A: artifacts.GetArtifact("pot"), Quantity: 1}},
		Outputs: []artifacts.Artifacts{artifacts.Artifacts{A: artifacts.GetArtifact("medicine"), Quantity: 1}}},
	&Manufacture{
		Name:                  "swordsmith",
		Time:                  30 * 24,
		Power:                 1000,
		BuildingExtensionType: building.Forge,
		Inputs: []artifacts.Artifacts{
			artifacts.Artifacts{A: artifacts.GetArtifact("iron_bar"), Quantity: 2},
			artifacts.Artifacts{A: artifacts.GetArtifact("board"), Quantity: 1},
			artifacts.Artifacts{A: artifacts.GetArtifact("leather"), Quantity: 1}},
		Outputs: []artifacts.Artifacts{
			artifacts.Artifacts{A: artifacts.GetArtifact("sword"), Quantity: 1},
			artifacts.Artifacts{A: artifacts.GetArtifact("shield"), Quantity: 1}}},
}

func GetAllManufactureNames() []string {
	result := make([]string, len(AllManufacture))
	for i, m := range AllManufacture {
		result[i] = m.Name
	}
	return result
}

func GetManufactureNames(extensions []*building.BuildingExtension) []string {
	var filteredExtensions []*building.BuildingExtension
	for _, extension := range extensions {
		if extension.T != building.Deck {
			filteredExtensions = append(filteredExtensions, extension)
		}
	}
	result := make([]string, 0, len(AllManufacture))
	for _, m := range AllManufacture {
		if m.BuildingExtensionType == nil && len(filteredExtensions) == 0 {
			result = append(result, m.Name)
		} else {
			for _, extension := range filteredExtensions {
				if extension != nil && extension.T == m.BuildingExtensionType {
					result = append(result, m.Name)
					break
				}
			}
		}
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

func (m *Manufacture) IsInput(a *artifacts.Artifact) bool {
	for _, a2 := range m.Inputs {
		if a == a2.A {
			return true
		}
	}
	return false
}

func (m *Manufacture) IsOutput(a *artifacts.Artifact) bool {
	for _, a2 := range m.Outputs {
		if a == a2.A {
			return true
		}
	}
	return false
}
