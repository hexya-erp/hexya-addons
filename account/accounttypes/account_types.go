// Copyright 2017 NDP Systèmes. All Rights Reserved.
// See LICENSE file for full licensing details.

package accounttypes

import "github.com/hexya-erp/hexya/hexya/models/types/dates"

// A PaymentDueDates gives the amount due of an invoice at the given date
type PaymentDueDates struct {
	Date   dates.Date
	Amount float64
}

// A TaxGroup holds an amount for a given group name
type TaxGroup struct {
	GroupName string
	TaxAmount float64
	Sequence  int
}

// A DataForReconciliationWidget holds data for the reconciliation widget
type DataForReconciliationWidget struct {
	Customers []map[string]interface{} `json:"customers"`
	Suppliers []map[string]interface{} `json:"suppliers"`
	Accounts  []map[string]interface{} `json:"accounts"`
}

// An AppliedTaxData is the result of the computation of applying a tax on an amount.
type AppliedTaxData struct {
	ID              int64   `json:"id"`
	Name            string  `json:"name"`
	Amount          float64 `json:"amount"`
	Sequence        int     `json:"sequence"`
	AccountID       int64   `json:"account_id"`
	RefundAccountID int64   `json:"refund_account_id"`
	Analytic        bool    `json:"analytic"`
	Base            float64 `json:"base"`
}
