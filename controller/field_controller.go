package controller

import (
	"medvil/model/navigation"
	"medvil/view/gui"
)


func FieldToControlPanel(p *gui.Panel, f *navigation.Field) {
	p.AddTextureLabel("terrain/"+f.Terrain.T.Name, 10, 50, 32, 32)
	var aI = 0
	for a, q := range f.Terrain.Resources.Artifacts {
		ArtifactsToControlPanel(p, aI, a, q)
		aI++
	}
}
