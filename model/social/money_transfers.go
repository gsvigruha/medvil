package social

import (
	"math"
)

type HomeProvider interface {
	GetHome() Home
}

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
	Trader            TransferCategories
	MarketFundingRate int
}

func CollectTax[H HomeProvider](homes []H, town *Town, t TransferCategories) {
	for _, home := range homes {
		if int(home.GetHome().GetMoney()) > t.Threshold {
			tax := (home.GetHome().GetMoney() - uint32(t.Threshold)) * uint32(t.Rate) / 100
			home.GetHome().Spend(tax)
			town.Townhall.Household.Money += tax
			dHappiness := uint8(t.Rate * 2)
			for _, person := range home.GetHome().GetPeople() {
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

func SumSubsidyNeeded[H HomeProvider](homes []H, transfer TransferCategories) uint32 {
	var subsidy uint32 = 0
	for _, h := range homes {
		if h.GetHome().GetMoney() < uint32(transfer.Threshold) {
			subsidy += uint32(transfer.Threshold) - h.GetHome().GetMoney()
		}
	}
	return subsidy
}

func SendSubsidy[H HomeProvider](homes []H, t *Town, transfer TransferCategories, ratio float64) {
	for _, h := range homes {
		if h.GetHome().GetMoney() < uint32(transfer.Threshold) {
			q := uint32(math.Floor(float64(transfer.Threshold)-float64(h.GetHome().GetMoney())) * ratio)
			h.GetHome().Earn(q)
			t.Townhall.Household.Money -= q
		}
	}
}
