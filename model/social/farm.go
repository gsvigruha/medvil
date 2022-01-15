package social

import (
	"encoding/json"
	"medvil/model/navigation"
	"medvil/model/time"
)

type FarmLand struct {
	X uint16
	Y uint16
	F navigation.IField
}

type Farm struct {
	Household Household
	Land      []FarmLand
}

func (f *Farm) UnmarshalJSON(data []byte) error {
	var j map[string]json.RawMessage
	if err := json.Unmarshal(data, &j); err != nil {
		return err
	}
	if err := json.Unmarshal(j["household"], &f.Household); err != nil {
		return err
	}
	var l [][]uint16
	if err := json.Unmarshal(j["land"], &l); err != nil {
		return err
	}
	f.Land = make([]FarmLand, len(l))
	for i := range l {
		f.Land[i].X = l[i][0]
		f.Land[i].Y = l[i][1]
	}
	return nil
}

func (f *Farm) ElapseTime(Calendar *time.CalendarType) {
	f.Household.ElapseTime(Calendar)

}
