package materials

type Material struct {
	Name string
	// kg/m3
	Density uint16
	Liquid  bool
}

var All = [...]Material{
	Material{Name: "water", Density: 1000, Liquid: true},
	Material{Name: "stone", Density: 1800, Liquid: false},
	Material{Name: "sandstone", Density: 2400, Liquid: false},
	Material{Name: "wood", Density: 700, Liquid: false},
	Material{Name: "brick", Density: 2000, Liquid: false},
	Material{Name: "marble", Density: 2600, Liquid: false},
	Material{Name: "hay", Density: 150, Liquid: false},
	Material{Name: "reed", Density: 150, Liquid: false},
	Material{Name: "tile", Density: 2000, Liquid: false},
	Material{Name: "clay", Density: 1600, Liquid: true},
	Material{Name: "whitewash", Density: 1000, Liquid: false},
	Material{Name: "limestone", Density: 2000, Liquid: false},
	Material{Name: "iron", Density: 7800, Liquid: false},
	Material{Name: "copper", Density: 9000, Liquid: false},
	Material{Name: "silver", Density: 10000, Liquid: false},
	Material{Name: "gold", Density: 19000, Liquid: false},
	Material{Name: "organic", Density: 1000, Liquid: false},

	Material{Name: "leather", Density: 800, Liquid: false},
	Material{Name: "linen", Density: 150, Liquid: false},
	Material{Name: "wool", Density: 150, Liquid: false},
	Material{Name: "paper", Density: 150, Liquid: false},
	Material{Name: "parchment", Density: 150, Liquid: false},
}

func GetMaterial(name string) *Material {
	for i := 0; i < len(All); i++ {
		if All[i].Name == name {
			return &All[i]
		}
	}
	return nil
}
