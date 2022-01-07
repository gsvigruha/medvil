package terrain

type PlanType struct {
	Name string
}

type Plant struct {
	T     *PlanType
	X     uint16
	Y     uint16
	Age   uint8
	Shape uint8
}
