package social

type TransferCategories struct {
	TaxRate      int
	TaxThreshold int
	Subsidy      int
}

type MoneyTransfers struct {
	Farm     TransferCategories
	Workshop TransferCategories
	Mine     TransferCategories
}

func (t *TransferCategories) Transfer(townMoney, householdMoney *uint32) {
	if int(*householdMoney) > t.TaxThreshold {
		tax := (*householdMoney - uint32(t.TaxThreshold)) * uint32(t.TaxRate) / 100
		*householdMoney -= tax
		*townMoney += tax
	} else if int(*householdMoney) < t.Subsidy {
		subsidy := uint32(t.Subsidy) - *householdMoney
		if *townMoney >= subsidy {
			*householdMoney += subsidy
			*townMoney -= subsidy
		}
	}
}
