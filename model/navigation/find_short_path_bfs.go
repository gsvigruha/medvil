package navigation

import (
	"log"
	"math/rand"
	"os"
)

const ShortPathMaxLength = 100
const Capacity = 100

type BFSElement struct {
	PE   PathElement
	prev *BFSElement
	d    uint8
}

func AddNextField(pe PathElement, prevE *BFSElement, toVisit *[]*BFSElement, inQueue map[Location]bool) {
	if _, ok := inQueue[pe.GetLocation()]; ok {
		// no need to add it to the queue again
	} else {
		*toVisit = append(*toVisit, &BFSElement{PE: pe, prev: prevE, d: prevE.d + 1})
		inQueue[pe.GetLocation()] = true
	}
}

func CheckField(pe PathElement, pathType PathType) bool {
	if pathType == PathTypePedestrian {
		return pe.Walkable() && !pe.Crowded()
	} else if pathType == PathTypeCart {
		return pe.Driveble()
	} else if pathType == PathTypeBoat {
		return pe.Sailable()
	}
	return false
}

func FindShortPathBFS(m IMap, start Location, dest Destination, pathType PathType) []PathElement {
	var iter = 0
	r := rand.New(rand.NewSource(int64(start.X*599 + start.Y)))
	visited := make(map[Location]*[]PathElement, Capacity)
	se := &BFSElement{PE: m.GetField(start.X, start.Y).GetPathElement(start.Z), prev: nil, d: 1}
	var toVisit = make([]*BFSElement, 1, Capacity)
	toVisit[0] = se
	var inQueue = make(map[Location]bool, Capacity)
	for len(toVisit) > 0 {
		e := toVisit[0]
		toVisit = toVisit[1:]

		if dest.Check(e.PE) {
			path := make([]PathElement, e.d)
			var eI = e
			for i := range path {
				path[len(path)-1-i] = eI.PE
				eI = eI.prev
			}
			if os.Getenv("MEDVIL_VERBOSE") == "2" {
				log.Printf("Found path with BFS in: ", iter)
			}
			return path
		}

		if _, ok := visited[e.PE.GetLocation()]; ok {
			continue
		}

		if e.d > ShortPathMaxLength || (e != se && !CheckField(e.PE, pathType)) {
			visited[e.PE.GetLocation()] = nil
			continue
		}

		neighbors := e.PE.GetNeighbors(m)
		order := r.Perm(len(neighbors))
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
		visited[e.PE.GetLocation()] = nil
		iter++
	}
	if os.Getenv("MEDVIL_VERBOSE") == "2" {
		log.Printf("Not found path with BFS after: ", iter)
	}
	return nil
}
