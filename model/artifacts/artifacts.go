package artifacts

import (
	"medvil/model/materials"
)

type Artifact struct {
	Name        string
	M           *materials.Material
	V           uint16
	Idx         uint16
	Description string
}

var All = [...]*Artifact{
	// Building
	&Artifact{Name: "log", M: materials.GetMaterial("wood"), V: 3, Idx: 1, Description: "Logs are used as firewood or wooden boards"},
	&Artifact{Name: "board", M: materials.GetMaterial("wood"), V: 3, Idx: 2, Description: "Wooden boards are for buildings or vehicles"},
	&Artifact{Name: "reed", M: materials.GetMaterial("reed"), V: 1, Idx: 3, Description: "Reeds are needed for paper or roof thatching"},
	&Artifact{Name: "stone", M: materials.GetMaterial("stone"), V: 2, Idx: 4, Description: "Raw stones are used to make stone cubes"},
	&Artifact{Name: "cube", M: materials.GetMaterial("stone"), V: 2, Idx: 5, Description: "Stones cubes are used for buildings or roads"},
	&Artifact{Name: "clay", M: materials.GetMaterial("clay"), V: 2, Idx: 6, Description: "Raw clay is needed to make brick, tiles or pots"},
	&Artifact{Name: "brick", M: materials.GetMaterial("brick"), V: 2, Idx: 7, Description: "Bricks are used as building materials"},
	&Artifact{Name: "thatch", M: materials.GetMaterial("thatch"), V: 2, Idx: 8, Description: "Thatches are used as roof materials"},
	&Artifact{Name: "tile", M: materials.GetMaterial("tile"), V: 2, Idx: 9, Description: "Tiles are used as roof materials"},
	&Artifact{Name: "pot", M: materials.GetMaterial("clay"), V: 1, Idx: 10, Description: "Pots are needed to produce beer or medicine"},
	// Metal
	&Artifact{Name: "iron_ore", M: materials.GetMaterial("iron"), V: 2, Idx: 11, Description: "Iron ore is needed to produce iron bars"},
	&Artifact{Name: "iron_bar", M: materials.GetMaterial("iron"), V: 2, Idx: 12, Description: "Iron bars are for tools, weapons or vehicles"},
	&Artifact{Name: "gold_ore", M: materials.GetMaterial("gold"), V: 2, Idx: 13, Description: "Gold ore is needed to mint coins"},
	&Artifact{Name: "gold_coin", M: materials.GetMaterial("gold"), V: 2, Idx: 14, Description: "Gold coins are used as currency"},
	// Food
	&Artifact{Name: "fruit", M: materials.GetMaterial("organic"), V: 1, Idx: 15, Description: "Fruits: food, increases health, happiness"},
	&Artifact{Name: "vegetable", M: materials.GetMaterial("organic"), V: 1, Idx: 16, Description: "Vegetables: food, increases health"},
	&Artifact{Name: "grain", M: materials.GetMaterial("organic"), V: 1, Idx: 17, Description: "Grain can be turned into flour"},
	&Artifact{Name: "flour", M: materials.GetMaterial("organic"), V: 1, Idx: 18, Description: "Flour is used to bake bread"},
	&Artifact{Name: "bread", M: materials.GetMaterial("organic"), V: 1, Idx: 19, Description: "Bread is high quality food"},
	&Artifact{Name: "meat", M: materials.GetMaterial("organic"), V: 1, Idx: 20, Description: "Meat is high quality food"},
	&Artifact{Name: "water", M: materials.GetMaterial("water"), V: 1, Idx: 21, Description: "Water is needed to reduce thirst"},
	//&Artifact{Name: "wine", M: materials.GetMaterial("water"), V: 1},
	&Artifact{Name: "beer", M: materials.GetMaterial("water"), V: 1, Idx: 22, Description: "Beer increases the happiness of villagers"},
	&Artifact{Name: "sheep", M: materials.GetMaterial("organic"), V: 3, Idx: 23, Description: "Sheep provides meat, wool and leather"},
	&Artifact{Name: "herb", M: materials.GetMaterial("organic"), V: 1, Idx: 24, Description: "Herbs are used to make medicine"},
	&Artifact{Name: "medicine", M: materials.GetMaterial("organic"), V: 1, Idx: 25, Description: "Medicine increases health"},
	// Sheets
	&Artifact{Name: "leather", M: materials.GetMaterial("leather"), V: 1, Idx: 26, Description: "Leather is needed for clothes or vehicles"},
	//&Artifact{Name: "linen", M: materials.GetMaterial("linen"), V: 1},
	&Artifact{Name: "wool", M: materials.GetMaterial("wool"), V: 1, Idx: 27, Description: "Wool is needed to produce textile"},
	&Artifact{Name: "paper", M: materials.GetMaterial("paper"), V: 1, Idx: 28, Description: "Paper makes the economy more efficient"},
	// Clothes
	&Artifact{Name: "textile", M: materials.GetMaterial("wool"), V: 2, Idx: 29, Description: "Textile is needed for clothes or vehicles"},
	&Artifact{Name: "clothes", M: materials.GetMaterial("wool"), V: 1, Idx: 30, Description: "Clothes keep villagers warm and happy"},
	// Tools
	&Artifact{Name: "tools", M: materials.GetMaterial("iron"), V: 1, Idx: 31, Description: "Tools are used to speed up all tasks"},
	// Heating
	&Artifact{Name: "firewood", M: materials.GetMaterial("wood"), V: 1, Idx: 32, Description: "Firewood	is needed to heat houses"},
	// Military
	&Artifact{Name: "sword", M: materials.GetMaterial("iron"), V: 1, Idx: 33, Description: "Swords are used by soldiers"},
	&Artifact{Name: "shield", M: materials.GetMaterial("iron"), V: 1, Idx: 34, Description: "Shields are used by soldiers"},
}

func GetArtifact(name string) *Artifact {
	for i := 0; i < len(All); i++ {
		if All[i].Name == name {
			return All[i]
		}
	}
	return nil
}

var Water = GetArtifact("water")
