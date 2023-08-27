package navigation

import (
	"log"
	"math/rand"
	"os"
)

func absDist(x, y uint16) float64 {
	if x > y {
		return float64(x - y)
	}
	return float64(y - x)
}

type TargetedElement struct {
	PE   PathElement
	prev *TargetedElement
	d    uint8
	ds   float64
}

func AddNextFieldWithEstimatedDist(pe PathElement, prevE *TargetedElement, toVisit *[][]*TargetedElement, inQueue map[Location]bool, dest Destination, currentEstimatedDist int) int {
	if _, ok := inQueue[pe.GetLocation()]; ok {
		return currentEstimatedDist
	} else {
		cx, cy := pe.LocationXY()
		dx, dy, _ := dest.DestHint()
		distI := prevE.ds + 1.0/pe.GetSpeed()
		estimatedTotalDist := int(distI + absDist(cx, dx) + absDist(cy, dy))
		for len(*toVisit) <= int(estimatedTotalDist) {
			*toVisit = append(*toVisit, []*TargetedElement{})
		}
		queueForDist := &((*toVisit)[estimatedTotalDist])
		*queueForDist = append(*queueForDist, &TargetedElement{PE: pe, prev: prevE, d: prevE.d + 1, ds: distI})
		inQueue[pe.GetLocation()] = true
		if estimatedTotalDist < currentEstimatedDist {
			return estimatedTotalDist
		}
		return currentEstimatedDist
	}
}

func FindShortPathTargeted(m IMap, start Location, dest Destination, pathType PathType) []PathElement {
	var iter = 0
	var currentEstimatedDist = 0
	r := rand.New(rand.NewSource(int64(start.X*599 + start.Y)))
	visited := make(map[Location]*[]PathElement, Capacity)
	se := &TargetedElement{PE: m.GetField(start.X, start.Y).GetPathElement(start.Z), prev: nil, d: 1, ds: 0.0}
	var toVisit = [][]*TargetedElement{[]*TargetedElement{se}}
	var inQueue = make(map[Location]bool, Capacity)
	for true {
		var e *TargetedElement = nil
		for e == nil {
			if currentEstimatedDist >= len(toVisit) {
				if os.Getenv("MEDVIL_VERBOSE") == "2" {
					log.Printf("Not found path with targeted algorithm after: ", iter)
				}
				return nil
			}
			if len(toVisit[currentEstimatedDist]) > 0 {
				e = toVisit[currentEstimatedDist][0]
				toVisit[currentEstimatedDist] = toVisit[currentEstimatedDist][1:]
			} else {
				currentEstimatedDist++
			}
		}

		if dest.Check(e.PE) {
			path := make([]PathElement, e.d)
			var eI = e
			for i := range path {
				path[len(path)-1-i] = eI.PE
				eI = eI.prev
			}
			if os.Getenv("MEDVIL_VERBOSE") == "2" {
				log.Printf("Found path with targeted algorithm in: ", iter)
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
				currentEstimatedDist = AddNextFieldWithEstimatedDist(pe, e, &toVisit, inQueue, dest, currentEstimatedDist)
			}
		}
		for _, idx := range order {
			pe := neighbors[idx]
			if pe.GetSpeed() <= 1.0 {
				currentEstimatedDist = AddNextFieldWithEstimatedDist(pe, e, &toVisit, inQueue, dest, currentEstimatedDist)
			}
		}
		visited[e.PE.GetLocation()] = nil
		iter++
	}
	if os.Getenv("MEDVIL_VERBOSE") == "2" {
		log.Printf("Not found path with targeted algorithm after: ", iter)
	}
	return nil
}
