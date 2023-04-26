package economy

type EquipmentType struct {
	Weapon bool
	Tool   bool
	Name   string
}

var NoEquipment = &EquipmentType{Name: "none", Tool: false, Weapon: false}
var Tool = &EquipmentType{Name: "tool", Tool: true, Weapon: false}
var Weapon = &EquipmentType{Name: "weapon", Tool: false, Weapon: true}

var EquipmentTypes = [...]*EquipmentType{
	NoEquipment,
	Tool,
	Weapon,
}

func GetEquipmentType(name string) *EquipmentType {
	for _, t := range EquipmentTypes {
		if t.Name == name {
			return t
		}
	}
	return nil
}
