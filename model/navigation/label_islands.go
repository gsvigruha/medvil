package navigation

import (
	"medvil/model/terrain"
)

func LabelIslands(m IMap, sx, sy uint16) {
	var cntr uint16 = 2
	for i := uint16(0); i < sx; i++ {
		for j := uint16(0); j < sy; j++ {
			m.GetField(i, j).IslandLabel = 0
		}
	}
	for i := uint16(0); i < sx; i++ {
		for j := uint16(0); j < sy; j++ {
			field := m.GetField(i, j)
			if field.IslandLabel == 0 {
				if field.Terrain.T == terrain.Water && field.Road == nil {
					spanFields(m, i, j, sx, sy, cntr, true)
					cntr++
				} else {
					spanFields(m, i, j, sx, sy, cntr, false)
					cntr++
				}
			}
		}
	}
}

type Coords struct {
	i uint16
	j uint16
}

func spanFields(m IMap, i, j, sx, sy uint16, c uint16, water bool) {
	var toVisit = []Coords{Coords{i: i, j: j}}
	visited := make(map[Coords]bool, Capacity)
	for len(toVisit) > 0 {
		i = toVisit[0].i
		j = toVisit[0].j

		cI := toVisit[0]
		toVisit = toVisit[1:]

		if _, ok := visited[cI]; ok {
			continue
		}

		field := m.GetField(i, j)
		field.IslandLabel = c
		for _, coords := range DirectionOrthogonalXY {
			i2 := uint16(int(i) + coords[0])
			j2 := uint16(int(j) + coords[1])
			field2 := m.GetField(i2, j2)
			coords := Coords{i: i2, j: j2}
			if field2 != nil && field2.IslandLabel == 0 {
				if water == (field2.Terrain.T == terrain.Water && field2.Road == nil) {
					toVisit = append(toVisit, coords)
				}
			}
		}
		visited[cI] = true
	}
}
