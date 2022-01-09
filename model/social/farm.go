package social

type FarmField interface {
}

type Farm struct {
	Household   Household
	Land []*FarmField
}
