package economy

type Tag []uint16

type Tags struct {
	T []Tag
}

func SingleTag(e ...uint16) Tag {
	return e
}

func MakeTags(t Tag) Tags {
	return Tags{T: []Tag{t}}
}

func AppendTags(ts Tags, t Tag) Tags {
	return Tags{T: append(ts.T, t)}
}

func (ts Tags) Count(t Tag) int {
	if len(t) == 0 { // Empty Tag
		return 1
	}
	var cnt = 0
	for _, tI := range ts.T {
		var same = true
		for i := range t {
			if i >= len(tI) || t[i] != tI[i] {
				same = false
			}
		}
		if same {
			cnt++
		}
	}
	return cnt
}

var EmptyTags = Tags{T: []Tag{}}
var EmptyTag = Tag{}

const TagFoodShopping uint16 = 1
const TagToolPurchase uint16 = 2
const TagHeatingFuelShopping uint16 = 3
const TagBeerShopping uint16 = 4
const TagMedicineShopping uint16 = 5
const TagTextileShopping uint16 = 6
const TagRepairShopping uint16 = 7
const TagOrderInput uint16 = 8
const TagWeaponBuying uint16 = 9
const TagPaperPurchase uint16 = 10
const TagManufactureInput uint16 = 11

const TagSellArtifacts uint16 = 21

const TagStorageTarget uint16 = 31

const TagMarket uint16 = 101
