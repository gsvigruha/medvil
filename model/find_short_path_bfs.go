package model

import (
	"medvil/model/navigation"
)

const ShortPathMaxLength = 10

type BFSElement struct {
	L    navigation.Location
	prev *BFSElement
	d    uint8
}

func FindShortPathBFS(m *Map, sx, sy, ex, ey uint16, travellerType uint8) []navigation.Location {
	visited := make(map[navigation.Location]*[]navigation.Location)
	var toVisit = []*BFSElement{&BFSElement{L: navigation.Location{X: sx, Y: sy, F: m.GetField(sx, sy)}, prev: nil, d: 1}}
	for len(toVisit) > 0 {
		e := toVisit[0]
		toVisit = toVisit[1:]
	
		if e.L.X == ex && e.L.Y == ey {
			path := make([]navigation.Location, e.d)
			var eI = e
			for i := range path {
				path[len(path)-1-i] = eI.L
				eI = eI.prev
			}
			return path
		}

		if _, ok := visited[e.L]; ok {
			continue
		}

		if e.d > ShortPathMaxLength || !e.L.F.Walkable() {
			visited[e.L] = nil
			continue
		}

		toVisit = append(
			toVisit,
			&BFSElement{L: navigation.Location{X: e.L.X + 1, Y: e.L.Y, F: m.GetField(e.L.X+1, e.L.Y)}, prev: e, d: e.d + 1},
			&BFSElement{L: navigation.Location{X: e.L.X - 1, Y: e.L.Y, F: m.GetField(e.L.X-1, e.L.Y)}, prev: e, d: e.d + 1},
			&BFSElement{L: navigation.Location{X: e.L.X, Y: e.L.Y + 1, F: m.GetField(e.L.X, e.L.Y+1)}, prev: e, d: e.d + 1},
			&BFSElement{L: navigation.Location{X: e.L.X, Y: e.L.Y - 1, F: m.GetField(e.L.X, e.L.Y-1)}, prev: e, d: e.d + 1},
		)
	}
	return nil
}
