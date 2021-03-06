package model

import (
	"math/rand"
	"medvil/model/navigation"
)

const ShortPathMaxLength = 100
const capacity = 1000

type BFSElement struct {
	PE   navigation.PathElement
	prev *BFSElement
	d    uint8
}

func AddNextField(pe navigation.PathElement, prevE *BFSElement, toVisit *[]*BFSElement, inQueue map[navigation.Location]bool) {
	if _, ok := inQueue[pe.GetLocation()]; ok {
		// no need to add it to the queue again
	} else {
		*toVisit = append(*toVisit, &BFSElement{PE: pe, prev: prevE, d: prevE.d + 1})
		inQueue[pe.GetLocation()] = true
	}
}

func CheckField(pe navigation.PathElement, travellerType uint8) bool {
	if travellerType == navigation.TravellerTypePedestrian {
		return pe.Walkable()
	} else if travellerType == navigation.TravellerTypeBoat {
		return pe.Sailable()
	}
	return false
}

func FindShortPathBFS(m *Map, start navigation.Location, dest navigation.Destination, travellerType uint8) []navigation.PathElement {
	var iter = 0
	visited := make(map[navigation.Location]*[]navigation.PathElement, capacity)
	se := &BFSElement{PE: m.GetField(start.X, start.Y), prev: nil, d: 1}
	var toVisit = []*BFSElement{se}
	var inQueue = make(map[navigation.Location]bool, capacity)
	for len(toVisit) > 0 {
		e := toVisit[0]
		toVisit = toVisit[1:]

		if dest.Check(e.PE) {
			path := make([]navigation.PathElement, e.d)
			var eI = e
			for i := range path {
				path[len(path)-1-i] = eI.PE
				eI = eI.prev
			}
			return path
		}

		if _, ok := visited[e.PE.GetLocation()]; ok {
			continue
		}

		if e.d > ShortPathMaxLength || (e != se && !CheckField(e.PE, travellerType)) {
			visited[e.PE.GetLocation()] = nil
			continue
		}

		neighbors := e.PE.GetNeighbors(m)
		order := rand.Perm(len(neighbors))
		for _, idx := range order {
			pe := neighbors[idx]
			if pe.GetSpeed() > 1.0 {
				AddNextField(pe, e, &toVisit, inQueue)
			}
		}
		for _, idx := range order {
			pe := neighbors[idx]
			if pe.GetSpeed() <= 1.0 {
				AddNextField(pe, e, &toVisit, inQueue)
			}
		}
		iter++
	}
	return nil
}
