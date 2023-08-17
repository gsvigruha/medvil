package social

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

func (t *TransferCategories) CollectTax(town, household *Household) {
	if int(household.Money) > t.Threshold {
		tax := (household.Money - uint32(t.Threshold)) * uint32(t.Rate) / 100
		household.Money -= tax
		town.Money += tax
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

func (t *MoneyTransfers) FundMarket(townMoney, marketMoney *uint32) {
	if int(*marketMoney) < int(*townMoney)*t.MarketFundingRate/100 {
		transfer := *townMoney*uint32(t.MarketFundingRate)/100 - *marketMoney
		if transfer <= *townMoney {
			*townMoney -= transfer
			*marketMoney += transfer
		}
	}
}
