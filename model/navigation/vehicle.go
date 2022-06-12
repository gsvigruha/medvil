package navigation

type Vehicle interface {
	TravellerType() uint8
	GetTraveller() *Traveller
}
