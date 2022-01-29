package model

import (
	"medvil/model/navigation"
)

const ShortPathMaxLength = 20

type BFSElement struct {
	F    *navigation.Field
	prev *BFSElement
	d    uint8
}

func AddNextField(m *Map, x, y uint16, e *BFSElement, toVisit *[]*BFSElement, inQueue map[*navigation.Field]bool) {
	if x >= 0 && y >= 0 && x < m.SX && y < m.SY {
		field := m.GetField(x, y)
		if _, ok := inQueue[field]; ok {
			// no need to add it to the queue again
		} else {
			*toVisit = append(*toVisit, &BFSElement{F: field, prev: e, d: e.d + 1})
			inQueue[field] = true
		}
	}
}

func FindShortPathBFS(m *Map, sx, sy, ex, ey uint16, travellerType uint8) []*navigation.Field {
	var iter = 0
	visited := make(map[*navigation.Field]*[]*navigation.Field)
	var toVisit = []*BFSElement{&BFSElement{F: m.GetField(sx, sy), prev: nil, d: 1}}
	var inQueue = make(map[*navigation.Field]bool)
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

		AddNextField(m, e.F.X+1, e.F.Y, e, &toVisit, inQueue)
		AddNextField(m, e.F.X-1, e.F.Y, e, &toVisit, inQueue)
		AddNextField(m, e.F.X, e.F.Y+1, e, &toVisit, inQueue)
		AddNextField(m, e.F.X, e.F.Y-1, e, &toVisit, inQueue)
		iter++
	}
	return nil
}
