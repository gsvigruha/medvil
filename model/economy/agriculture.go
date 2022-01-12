package economy

type AgriculturalTaskType struct {
	Duration uint16
}

type AgriculturalTask struct {
	T        *AgriculturalTaskType
	L        Location
	Progress uint16
}

func (t *AgriculturalTask) Tick() {
}

func (t *AgriculturalTask) Location() Location {
	return t.L
}
