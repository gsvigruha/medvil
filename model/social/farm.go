package social

type FarmField interface {
}

type Farm struct {
	H    HouseHold
	Land []*FarmField
}
