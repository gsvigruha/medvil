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

func (t *TransferCategories) Transfer(townMoney, householdMoney *uint32) {
	if int(*householdMoney) > t.Threshold {
		tax := (*householdMoney - uint32(t.Threshold)) * uint32(t.Rate) / 100
		*householdMoney -= tax
		*townMoney += tax
	} else if int(*householdMoney) < t.Threshold {
		subsidy := uint32(t.Threshold) - *householdMoney
		if *townMoney >= subsidy {
			*householdMoney += subsidy
			*townMoney -= subsidy
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
