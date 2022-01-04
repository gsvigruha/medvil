package materials

type Material struct {
	Name string
	// kg/m3
	Density uint16
}

var All = [...]Material{
	Material{Name: "stone", Density: 1800},
	Material{Name: "sandstone", Density: 2400},
	Material{Name: "wood", Density: 700},
	Material{Name: "brick", Density: 2000},
	Material{Name: "marble", Density: 2600},
	Material{Name: "hay", Density: 150},
	Material{Name: "tile", Density: 2000},
	Material{Name: "clay", Density: 1600},
	Material{Name: "limestone", Density: 2000},
	Material{Name: "iron", Density: 7800},
	Material{Name: "copper", Density: 9000},
	Material{Name: "silver", Density: 10000},
	Material{Name: "gold", Density: 19000}}

func GetMaterial(name string) *Material {
	for i := 0; i < len(All); i++ {
		if All[i].Name == name {
			return &All[i]
		}
	}
	return nil
}
