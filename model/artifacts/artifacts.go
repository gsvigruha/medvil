package artifacts

import (
	"medvil/model/materials"
)

type Artifact struct {
	Name string
	M    *materials.Material
	V    uint16
}

var All = [...]*Artifact{
	// Building
	&Artifact{Name: "log", M: materials.GetMaterial("wood"), V: 3},
	&Artifact{Name: "board", M: materials.GetMaterial("wood"), V: 3},
	&Artifact{Name: "reed", M: materials.GetMaterial("reed"), V: 1},
	&Artifact{Name: "stone", M: materials.GetMaterial("stone"), V: 2},
	&Artifact{Name: "cube", M: materials.GetMaterial("stone"), V: 2},
	&Artifact{Name: "clay", M: materials.GetMaterial("clay"), V: 2},
	&Artifact{Name: "brick", M: materials.GetMaterial("brick"), V: 2},
	// Metal
	&Artifact{Name: "iron_ore", M: materials.GetMaterial("iron"), V: 2},
	&Artifact{Name: "iron_bar", M: materials.GetMaterial("iron"), V: 2},
	&Artifact{Name: "gold_ore", M: materials.GetMaterial("gold"), V: 2},
	&Artifact{Name: "gold_coin", M: materials.GetMaterial("gold"), V: 2},
	// Food
	&Artifact{Name: "fruit", M: materials.GetMaterial("organic"), V: 1},
	&Artifact{Name: "vegetable", M: materials.GetMaterial("organic"), V: 1},
	&Artifact{Name: "grain", M: materials.GetMaterial("organic"), V: 1},
	&Artifact{Name: "flour", M: materials.GetMaterial("organic"), V: 1},
	&Artifact{Name: "bread", M: materials.GetMaterial("organic"), V: 1},
	&Artifact{Name: "meat", M: materials.GetMaterial("organic"), V: 1},
	&Artifact{Name: "water", M: materials.GetMaterial("water"), V: 1},
	&Artifact{Name: "wine", M: materials.GetMaterial("water"), V: 1},
	&Artifact{Name: "beer", M: materials.GetMaterial("water"), V: 1},
	// Sheets
	&Artifact{Name: "leather", M: materials.GetMaterial("leather"), V: 1},
	&Artifact{Name: "linen", M: materials.GetMaterial("linen"), V: 1},
	&Artifact{Name: "wool", M: materials.GetMaterial("wool"), V: 1},
	&Artifact{Name: "paper", M: materials.GetMaterial("paper"), V: 1},
	&Artifact{Name: "paper", M: materials.GetMaterial("parchment"), V: 1},
	// Clothes
	&Artifact{Name: "clothes", M: materials.GetMaterial("leather"), V: 1},
	&Artifact{Name: "clothes", M: materials.GetMaterial("linen"), V: 1},
	&Artifact{Name: "clothes", M: materials.GetMaterial("wool"), V: 1},
	// Tools
}

func GetArtifact(name string) *Artifact {
	for i := 0; i < len(All); i++ {
		if All[i].Name == name {
			return All[i]
		}
	}
	return nil
}
