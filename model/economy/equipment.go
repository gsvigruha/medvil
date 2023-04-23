package economy

type EquipmentType struct {
	Weapon bool
	Tool   bool
}

var NoEquipment = &EquipmentType{Tool: false, Weapon: false}
var Tool = &EquipmentType{Tool: true, Weapon: false}
var Weapon = &EquipmentType{Tool: false, Weapon: true}
