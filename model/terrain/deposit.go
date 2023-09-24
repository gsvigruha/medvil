package terrain

import (
	"medvil/model/artifacts"
)

type DepositType struct {
	Name string
	A    *artifacts.Artifact
}

var Rock = &DepositType{A: artifacts.GetArtifact("stone"), Name: "rock"}
var Mud = &DepositType{A: artifacts.GetArtifact("clay"), Name: "mud"}
var IronBog = &DepositType{A: artifacts.GetArtifact("iron_ore"), Name: "iron_bog"}
var Gold = &DepositType{A: artifacts.GetArtifact("gold_ore"), Name: "gold"}

var DepositTypes = [...]*DepositType{
	Rock,
	Mud,
	IronBog,
	Gold,
}

func GetDepositType(name string) *DepositType {
	for _, t := range DepositTypes {
		if t.Name == name {
			return t
		}
	}
	return nil
}

type Deposit struct {
	T *DepositType
	Q uint16
}
