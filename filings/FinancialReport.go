package filings

import "errors"

// TODO this will need a walk script that'll go through raw reports and parse these.
// there should be an option to only create new and to override all
type FinancialReport struct {
	CIK, Year, Quarter                                               int64
	Revenue, OperatingExpense, NetIncome                             int64
	CurrentAssets, TotalAssets, CurrentLiabilities, TotalLiabilities int64
	OperatingCash, CapitalExpenditures                               int64
}

func (fr *FinancialReport) IsValid() error {
	missingFields := ""

	// not looping through every struct field because some may be optional
	if fr.Revenue == 0 {
		missingFields += "Revenue,"
	}

	if fr.OperatingExpense == 0 {
		missingFields += "OperatingExpense,"
	}

	if fr.NetIncome == 0 {
		missingFields += "NetIncome,"
	}

	if fr.TotalAssets == 0 {
		missingFields += "TotalAssets,"
	}

	if fr.TotalLiabilities == 0 {
		missingFields += "TotalLiabilities,"
	}

	if fr.CurrentAssets == 0 {
		missingFields += "CurrentAssets,"
	}

	if fr.CurrentLiabilities == 0 {
		missingFields += "CurrentLiabilities,"
	}

	if fr.OperatingCash == 0 {
		missingFields += "OperatingCash,"
	}

	if fr.CapitalExpenditures == 0 {
		missingFields += "CapitalExpenditures,"
	}

	if len(missingFields) > 0 {
		return errors.New(missingFields)
	} else {
		return nil
	}
}

type FinancialReportRaw struct {
	CIK, Year, Quarter int64
	RawFields          map[string]int64
}

func (frr *FinancialReportRaw) GetPreviousQuarter() (int64, int64) {
	if frr.Quarter == 1 {
		return frr.Year - 1, 4
	} else {
		return frr.Year, frr.Quarter - 1
	}
}

type RawFieldNameList interface {
	GetInt64RawFieldNames() []string
	GetVariablePeriodFieldNames() []string
}

// This could be done with a db table, but I like the idea of having something so critical in source control...
type BasicRawFieldNameList struct{}

func (brfnl *BasicRawFieldNameList) GetInt64RawFieldNames() []string {
	return []string{
		"Revenues",
		"SalesRevenueNet",
		"CostsAndExpenses",
		"OperatingExpenses",
		"NetIncomeLoss",

		"Assets",
		"LiabilitiesCurrent",
		"LongTermDebtNoncurrent",
		"DeferredTaxLiabilitiesNoncurrent",
		"AssetsCurrent",

		"NetCashProvidedByUsedInOperatingActivities",
		"PaymentsToAcquirePropertyPlantAndEquipment",
	}
}

func (brfnl *BasicRawFieldNameList) GetVariablePeriodFieldNames() []string {
	return []string{
		"NetCashProvidedByUsedInOperatingActivities",
		"PaymentsToAcquirePropertyPlantAndEquipment",
	}
}

type RawToScreenableMapping interface {
	GetRawToScreenableMapping(fr *FinancialReport) map[*int64][][]string
}

type BasicRawToScreenableMapping struct{}

func (brtsm *BasicRawToScreenableMapping) GetRawToScreenableMapping(fr *FinancialReport) map[*int64][][]string {
	var mapping map[*int64][][]string = make(map[*int64][][]string)

	mapping[&fr.Revenue] = [][]string{{"Revenues"}, {"SalesRevenueNet"}}
	mapping[&fr.OperatingExpense] =
		[][]string{{"CostsAndExpenses"},
			{"OperatingExpenses"},
			{"CostOfGoodsAndServicesSold"}}
	mapping[&fr.NetIncome] = [][]string{{"NetIncomeLoss"}}

	mapping[&fr.CurrentAssets] = [][]string{{"AssetsCurrent"}}
	mapping[&fr.TotalAssets] = [][]string{{"Assets"}}
	mapping[&fr.CurrentLiabilities] = [][]string{{"LiabilitiesCurrent"}}
	mapping[&fr.TotalLiabilities] = [][]string{{
		"LiabilitiesCurrent",
		"DeferredTaxLiabilitiesNoncurrent",
		"LongTermDebtNoncurrent",
	}, {"Liabilities"}}

	mapping[&fr.OperatingCash] = [][]string{{"NetCashProvidedByUsedInOperatingActivities"}}
	mapping[&fr.CapitalExpenditures] = [][]string{{"PaymentsToAcquirePropertyPlantAndEquipment"}}

	return mapping
}
