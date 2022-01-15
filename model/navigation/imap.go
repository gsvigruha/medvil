package navigation

type IMap interface {
	GetField(x uint16, y uint16) *Field
}
