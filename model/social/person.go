package social

type Person struct {
	FX        uint16
	FY        uint16
	FZ        uint8
	PX        uint8
	PY        uint8
	Hunger    uint8
	Thirst    uint8
	Happiness uint8
	Health    uint8
	HouseHold *HouseHold
}
