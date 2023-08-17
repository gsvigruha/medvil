package social

import (
	"math"
)

type TransferCategories struct {
	Rate      int
	Threshold int
}

type MoneyTransfers struct {
	Farm              TransferCategories
	Workshop          TransferCategories
	Mine              TransferCategories
	Factory           TransferCategories
	Tower             TransferCategories
	MarketFundingRate int
}

func CollectTax[H House](houses []H, town *Town, t TransferCategories) {
	for _, h := range houses {
		household := h.GetHousehold()
		if int(household.Money) > t.Threshold {
			tax := (household.Money - uint32(t.Threshold)) * uint32(t.Rate) / 100
			household.Money -= tax
			town.Townhall.Household.Money += tax
			dHappiness := uint8(t.Rate * 2)
			for _, person := range household.People {
				if person.Happiness >= dHappiness {
					person.Happiness = -dHappiness
				} else {
					person.Happiness = 0
				}
			}
		}
	}
}

func (t *MoneyTransfers) FundMarket(townMoney, marketMoney *uint32) {
	if int(*marketMoney) < int(*townMoney)*t.MarketFundingRate/100 {
		transfer := *townMoney*uint32(t.MarketFundingRate)/100 - *marketMoney
		if transfer <= *townMoney {
			*townMoney -= transfer
			*marketMoney += transfer
		}
	}
}

func SumSubsidyNeeded[H House](houses []H, transfer TransferCategories) uint32 {
	var subsidy uint32 = 0
	for _, h := range houses {
		if h.GetHousehold().Money < uint32(transfer.Threshold) {
			subsidy += uint32(transfer.Threshold) - h.GetHousehold().Money
		}
	}
	return subsidy
}

func SendSubsidy[H House](houses []H, t *Town, transfer TransferCategories, ratio float64) {
	for _, h := range houses {
		if h.GetHousehold().Money < uint32(transfer.Threshold) {
			q := uint32(math.Floor(float64(transfer.Threshold)-float64(h.GetHousehold().Money)) * ratio)
			h.GetHousehold().Money += q
			t.Townhall.Household.Money -= q
		}
	}
}
