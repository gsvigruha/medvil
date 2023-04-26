package terrain

import (
	"medvil/model/artifacts"
)

var AllCropTypes = [...]*PlantType{
	&PlantType{Name: "grain", MaturityAgeYears: 1, Yield: artifacts.Artifacts{A: artifacts.GetArtifact("grain"), Quantity: 2}, Tall: true, Habitat: Cultivated},
	&PlantType{Name: "vegetables", MaturityAgeYears: 1, Yield: artifacts.Artifacts{A: artifacts.GetArtifact("vegetable"), Quantity: 3}, Tall: false, Habitat: Cultivated},
	&PlantType{Name: "reed", MaturityAgeYears: 3, Yield: artifacts.Artifacts{A: artifacts.GetArtifact("reed"), Quantity: 1}, Tall: true, Habitat: Shore},
	&PlantType{Name: "herb", MaturityAgeYears: 1, Yield: artifacts.Artifacts{A: artifacts.GetArtifact("herb"), Quantity: 1}, Tall: false, Habitat: Cultivated},
}
