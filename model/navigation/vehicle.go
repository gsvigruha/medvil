package navigation

type Vehicle interface {
	PathType() PathType
	GetTraveller() *Traveller
	SetInUse(bool)
	SetHome(bool)
	Water() bool
}
