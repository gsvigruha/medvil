package navigation

type Vehicle interface {
	TravellerType() uint8
	GetTraveller() *Traveller
	SetInUse(bool)
	SetHome(bool)
}
