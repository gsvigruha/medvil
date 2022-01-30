package artifacts

import (
	"medvil/model/materials"
)

type Artifact struct {
	Name string
	M    *materials.Material
}

var All = [...]*Artifact{
	// Building
	&Artifact{Name: "log", M: materials.GetMaterial("wood")},
	&Artifact{Name: "board", M: materials.GetMaterial("wood")},
	&Artifact{Name: "rock", M: materials.GetMaterial("stone")},
	&Artifact{Name: "rock", M: materials.GetMaterial("sandstone")},
	&Artifact{Name: "rock", M: materials.GetMaterial("marble")},
	&Artifact{Name: "cube", M: materials.GetMaterial("stone")},
	&Artifact{Name: "cube", M: materials.GetMaterial("marble")},
	&Artifact{Name: "cube", M: materials.GetMaterial("sandstone")},
	&Artifact{Name: "clay", M: materials.GetMaterial("clay")},
	&Artifact{Name: "brick", M: materials.GetMaterial("brick")},
	// Metal
	&Artifact{Name: "ore", M: materials.GetMaterial("iron")},
	&Artifact{Name: "ore", M: materials.GetMaterial("gold")},
	&Artifact{Name: "ore", M: materials.GetMaterial("silver")},
	&Artifact{Name: "ore", M: materials.GetMaterial("copper")},
	// Food
	&Artifact{Name: "fruit", M: materials.GetMaterial("organic")},
	&Artifact{Name: "vegetable", M: materials.GetMaterial("organic")},
	&Artifact{Name: "grain", M: materials.GetMaterial("organic")},
	&Artifact{Name: "bread", M: materials.GetMaterial("organic")},
	&Artifact{Name: "meat", M: materials.GetMaterial("organic")},
	&Artifact{Name: "water", M: materials.GetMaterial("water")},
	&Artifact{Name: "wine", M: materials.GetMaterial("water")},
	&Artifact{Name: "beer", M: materials.GetMaterial("water")},
	// Sheets
	&Artifact{Name: "leather", M: materials.GetMaterial("leather")},
	&Artifact{Name: "linen", M: materials.GetMaterial("linen")},
	&Artifact{Name: "wool", M: materials.GetMaterial("wool")},
	&Artifact{Name: "paper", M: materials.GetMaterial("paper")},
	&Artifact{Name: "paper", M: materials.GetMaterial("parchment")},
	// Clothes
	&Artifact{Name: "clothes", M: materials.GetMaterial("leather")},
	&Artifact{Name: "clothes", M: materials.GetMaterial("linen")},
	&Artifact{Name: "clothes", M: materials.GetMaterial("wool")},
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
