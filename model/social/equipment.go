package social

type Equipment interface {
	Weapon() bool
	Tool() bool
}

type NoEquipment struct {
}

func (e *NoEquipment) Tool() bool {
	return false
}

func (e *NoEquipment) Weapon() bool {
	return false
}

type Tool struct {
}

func (t *Tool) Tool() bool {
	return true
}

func (t *Tool) Weapon() bool {
	return false
}

type Weapon struct {
}

func (w *Weapon) Tool() bool {
	return false
}

func (w *Weapon) Weapon() bool {
	return true
}
