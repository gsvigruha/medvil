package model

import (
	"medvil/model/navigation"
)

const ShortPathMaxLength = 10

type BFSElement struct {
	F    *navigation.Field
	prev *BFSElement
	d    uint8
}

func FindShortPathBFS(m *Map, sx, sy, ex, ey uint16, travellerType uint8) []*navigation.Field {
	visited := make(map[*navigation.Field]*[]*navigation.Field)
	var toVisit = []*BFSElement{&BFSElement{F: m.GetField(sx, sy), prev: nil, d: 1}}
	for len(toVisit) > 0 {
		e := toVisit[0]
		toVisit = toVisit[1:]

		if e.F.X == ex && e.F.Y == ey {
			path := make([]*navigation.Field, e.d)
			var eI = e
			for i := range path {
				path[len(path)-1-i] = eI.F
				eI = eI.prev
			}
			return path
		}

		if _, ok := visited[e.F]; ok {
			continue
		}

		if e.d > ShortPathMaxLength || !e.F.Walkable() {
			visited[e.F] = nil
			continue
		}

		toVisit = append(
			toVisit,
			&BFSElement{F: m.GetField(e.F.X+1, e.F.Y), prev: e, d: e.d + 1},
			&BFSElement{F: m.GetField(e.F.X-1, e.F.Y), prev: e, d: e.d + 1},
			&BFSElement{F: m.GetField(e.F.X, e.F.Y+1), prev: e, d: e.d + 1},
			&BFSElement{F: m.GetField(e.F.X, e.F.Y-1), prev: e, d: e.d + 1},
		)
	}
	return nil
}
