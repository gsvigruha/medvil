package model

import (
	"math/rand"
	"medvil/model/navigation"
)

const ShortPathMaxLength = 100
const capacity = 1000

type BFSElement struct {
	F    *navigation.Field
	prev *BFSElement
	d    uint8
	z    uint8
}

func AddNextField(m *Map, x, y uint16, e *BFSElement, toVisit *[]*BFSElement, inQueue map[*navigation.Field]bool) {
	if x >= 0 && y >= 0 && x < m.SX && y < m.SY {
		field := m.GetField(x, y)
		if _, ok := inQueue[field]; ok {
			// no need to add it to the queue again
		} else {
			nextZ := uint8(int(e.z) + len(field.Building.BuildingComponents) - len(e.F.Building.BuildingComponents))
			*toVisit = append(*toVisit, &BFSElement{F: field, prev: e, d: e.d + 1, z: nextZ})
			inQueue[field] = true
		}
	}
}

func FindShortPathBFS(m *Map, start navigation.Location, dest navigation.Destination, travellerType uint8) []navigation.PathElement {
	var iter = 0
	visited := make(map[*navigation.Field]*[]*navigation.Field, capacity)
	se := &BFSElement{F: m.GetField(start.X, start.Y), prev: nil, d: 1, z: start.Z}
	var toVisit = []*BFSElement{se}
	var inQueue = make(map[*navigation.Field]bool, capacity)
	for len(toVisit) > 0 {
		e := toVisit[0]
		toVisit = toVisit[1:]

		if dest.Check(e.F) {
			path := make([]navigation.PathElement, e.d)
			var eI = e
			for i := range path {
				path[len(path)-1-i] = navigation.PathElement{F: eI.F, Z: eI.z}
				eI = eI.prev
			}
			return path
		}

		if _, ok := visited[e.F]; ok {
			continue
		}

		if e.d > ShortPathMaxLength || (e != se && !e.F.Walkable()) {
			visited[e.F] = nil
			continue
		}

		nextCoords := [][]uint16{{e.F.X + 1, e.F.Y}, {e.F.X - 1, e.F.Y}, {e.F.X, e.F.Y + 1}, {e.F.X, e.F.Y - 1}}
		order := rand.Perm(4)
		for _, idx := range order {
			x, y := nextCoords[idx][0], nextCoords[idx][1]
			if m.GetField(x, y) != nil && m.GetField(x, y).Road != nil {
				AddNextField(m, x, y, e, &toVisit, inQueue)
			}
		}
		for _, idx := range order {
			x, y := nextCoords[idx][0], nextCoords[idx][1]
			if m.GetField(x, y) != nil && m.GetField(x, y).Road == nil {
				AddNextField(m, x, y, e, &toVisit, inQueue)
			}
		}
		iter++
	}
	return nil
}
