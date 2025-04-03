package models

type IncomeStatement struct {
	ID                                 uint    `gorm:"primaryKey;autoIncrement" json:"id"`
	StatementID                        *string `json:"statement_id,omitempty"`
	FileID                             *string `json:"file_id,omitempty"`
	NetTurnover                        *string `json:"net_turnover,omitempty"`
	ByNatureInventoryChange            *string `json:"by_nature_inventory_change,omitempty"`
	ByNatureLongTermInvestmentExpenses *string `json:"by_nature_long_term_investment_expenses,omitempty"`
	ByNatureOtherOperatingRevenues     *string `json:"by_nature_other_operating_revenues,omitempty"`
	ByNatureMaterialExpenses           *string `json:"by_nature_material_expenses,omitempty"`
	ByNatureLabourExpenses             *string `json:"by_nature_labour_expenses,omitempty"`
	ByNatureDepreciationExpenses       *string `json:"by_nature_depreciation_expenses,omitempty"`
	ByFunctionCostOfGoodsSold          *string `json:"by_function_cost_of_goods_sold,omitempty"`
	ByFunctionGrossProfit              *string `json:"by_function_gross_profit,omitempty"`
	ByFunctionSellingExpenses          *string `json:"by_function_selling_expenses,omitempty"`
	ByFunctionAdministrativeExpenses   *string `json:"by_function_administrative_expenses,omitempty"`
	ByFunctionOtherOperatingRevenues   *string `json:"by_function_other_operating_revenues,omitempty"`
	OtherOperatingExpenses             *string `json:"other_operating_expenses,omitempty"`
	EquityInvestmentEarnings           *string `json:"equity_investment_earnings,omitempty"`
	OtherLongTermInvestmentEarnings    *string `json:"other_long_term_investment_earnings,omitempty"`
	OtherInterestRevenues              *string `json:"other_interest_revenues,omitempty"`
	InvestmentFairValueAdjustments     *string `json:"investment_fair_value_adjustments,omitempty"`
	InterestExpenses                   *string `json:"interest_expenses,omitempty"`
	ExtraRevenues                      *string `json:"extra_revenues,omitempty"`
	ExtraExpenses                      *string `json:"extra_expenses,omitempty"`
	IncomeBeforeIncomeTaxes            *string `json:"income_before_income_taxes,omitempty"`
	ProvisionForIncomeTaxes            *string `json:"provision_for_income_taxes,omitempty"`
	IncomeAfterIncomeTaxes             *string `json:"income_after_income_taxes,omitempty"`
	OtherTaxes                         *string `json:"other_taxes,omitempty"`
	ExtraDividends                     *string `json:"extra_dividends,omitempty"`
	NetIncome                          *string `json:"net_income,omitempty"`
}
