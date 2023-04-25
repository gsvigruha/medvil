package economy

type EquipmentType struct {
	Weapon bool
	Tool   bool
	Name   string
}

var NoEquipment = &EquipmentType{Name: "none", Tool: false, Weapon: false}
var Tool = &EquipmentType{Name: "tool", Tool: true, Weapon: false}
var Weapon = &EquipmentType{Name: "weapon", Tool: false, Weapon: true}
