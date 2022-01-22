package navigation

type IMap interface {
	GetField(x uint16, y uint16) *Field
	ShortPath(sx, sy, ex, ey uint16, travellerType uint8) *Path
}
