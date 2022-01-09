package social

type JSONBuilding struct {
	Plan string
	X    uint16
	Y    uint16
}

type JSONFarm struct {
	Land       [][]uint16
	Building   JSONBuilding
	Population uint8
	Money      uint32
}

type Town struct {
	Country     *Country
	Townhall    *Townhall
	Marketplace *Marketplace
	Farms  []*Farm
}
