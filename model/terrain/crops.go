package terrain

import (
	"medvil/model/artifacts"
)

var AllCropTypes = [...]*PlantType{
	&PlantType{Name: "grain", MaturityAgeYears: 1, Yield: artifacts.Artifacts{A: artifacts.GetArtifact("grain"), Quantity: 2}, Tall: true, Habitat: Land},
	&PlantType{Name: "vegetables", MaturityAgeYears: 1, Yield: artifacts.Artifacts{A: artifacts.GetArtifact("vegetable"), Quantity: 4}, Tall: false, Habitat: Land},
	&PlantType{Name: "reed", MaturityAgeYears: 3, Yield: artifacts.Artifacts{A: artifacts.GetArtifact("reed"), Quantity: 1}, Tall: true, Habitat: Shore},
	&PlantType{Name: "herb", MaturityAgeYears: 1, Yield: artifacts.Artifacts{A: artifacts.GetArtifact("herb"), Quantity: 1}, Tall: false, Habitat: Land},
}
