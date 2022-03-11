package terrain

import (
	"medvil/model/artifacts"
)

var AllCropTypes = [...]PlantType{
	PlantType{Name: "grain", MaturityAgeYears: 1, Yield: artifacts.Artifacts{A: artifacts.GetArtifact("grain"), Quantity: 2}},
	PlantType{Name: "vegetables", MaturityAgeYears: 1, Yield: artifacts.Artifacts{A: artifacts.GetArtifact("vegetable"), Quantity: 3}},
	PlantType{Name: "reed", MaturityAgeYears: 5, Yield: artifacts.Artifacts{A: artifacts.GetArtifact("reed"), Quantity: 1}},
}
