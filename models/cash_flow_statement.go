package models

type CashFlowStatement struct {
	ID                                                     uint    `gorm:"primaryKey;autoIncrement" json:"id"`
	StatementID                                            *string `json:"statement_id,omitempty"`
	FileID                                                 *string `json:"file_id,omitempty"`
	CfoDmCashReceivedFromCustomers                         *string `json:"cfo_dm_cash_received_from_customers,omitempty"`
	CfoDmCashPaidToSuppliersEmployees                      *string `json:"cfo_dm_cash_paid_to_suppliers_employees,omitempty"`
	CfoDmOtherCashReceivedPaid                             *string `json:"cfo_dm_other_cash_received_paid,omitempty"`
	CfoDmOperatingCashFlow                                 *string `json:"cfo_dm_operating_cash_flow,omitempty"`
	CfoDmInterestPaid                                      *string `json:"cfo_dm_interest_paid,omitempty"`
	CfoDmIncomeTaxesPaid                                   *string `json:"cfo_dm_income_taxes_paid,omitempty"`
	CfoDmExtraItemsCashFlow                                *string `json:"cfo_dm_extra_items_cash_flow,omitempty"`
	CfoDmNetOperatingCashFlow                              *string `json:"cfo_dm_net_operating_cash_flow,omitempty"`
	CfoImIncomeBeforeIncomeTaxes                           *string `json:"cfo_im_income_before_income_taxes,omitempty"`
	CfoImIncomeBeforeChangesInWorkingCapital               *string `json:"cfo_im_income_before_changes_in_working_capital,omitempty"`
	CfoImOperatingCashFlow                                 *string `json:"cfo_im_operating_cash_flow,omitempty"`
	CfoImInterestPaid                                      *string `json:"cfo_im_interest_paid,omitempty"`
	CfoImIncomeTaxesPaid                                   *string `json:"cfo_im_income_taxes_paid,omitempty"`
	CfoImExtraItemsCashFlow                                *string `json:"cfo_im_extra_items_cash_flow,omitempty"`
	CfoImNetOperatingCashFlow                              *string `json:"cfo_im_net_operating_cash_flow,omitempty"`
	CfiAcquisitionOfStocksShares                           *string `json:"cfi_acquisition_of_stocks_shares,omitempty"`
	CfiSaleProceedsFromStocksShares                        *string `json:"cfi_sale_proceeds_from_stocks_shares,omitempty"`
	CfiAcquisitionOfFixedAssetsIntangibleAssets            *string `json:"cfi_acquisition_of_fixed_assets_intangible_assets,omitempty"`
	CfiSaleProceedsFromFixedAssetsIntangibleAssets         *string `json:"cfi_sale_proceeds_from_fixed_assets_intangible_assets,omitempty"`
	CfiLoansMade                                           *string `json:"cfi_loans_made,omitempty"`
	CfiRepaymentsOfLoansReceived                           *string `json:"cfi_repayments_of_loans_received,omitempty"`
	CfiInterestReceived                                    *string `json:"cfi_interest_received,omitempty"`
	CfiDividendsReceived                                   *string `json:"cfi_dividends_received,omitempty"`
	CfiNetInvestingCashFlow                                *string `json:"cfi_net_investing_cash_flow,omitempty"`
	CffProceedsFromStocksBondsIssuanceOrContributedCapital *string `json:"cff_proceeds_from_stocks_bonds_issuance_or_contributed_capital,omitempty"`
	CffLoansReceived                                       *string `json:"cff_loans_received,omitempty"`
	CffSubsidiesGrantsDonationsReceived                    *string `json:"cff_subsidies_grants_donations_received,omitempty"`
	CffRepaymentsOfLoansMade                               *string `json:"cff_repayments_of_loans_made,omitempty"`
	CffRepaymentsOfLeaseObligations                        *string `json:"cff_repayments_of_lease_obligations,omitempty"`
	CffDividendsPaid                                       *string `json:"cff_dividends_paid,omitempty"`
	CffNetFinancingCashFlow                                *string `json:"cff_net_financing_cash_flow,omitempty"`
	EffectOfExchangeRateChange                             *string `json:"effect_of_exchange_rate_change,omitempty"`
	NetIncrease                                            *string `json:"net_increase,omitempty"`
	AtBeginningOfYear                                      *string `json:"at_beginning_of_year,omitempty"`
	AtEndOfYear                                            *string `json:"at_end_of_year,omitempty"`
}
