// Copyright 2017 NDP Systèmes. All Rights Reserved.
// See LICENSE file for full licensing details.

package sale

import (
	"github.com/hexya-erp/hexya/hexya/models"
	"github.com/hexya-erp/hexya/pool"
)

func init() {

	pool.AccountInvoice().AddFields(map[string]models.FieldDefinition{
		"Team": models.Many2OneField{String: "Sales Team", RelationModel: pool.CRMTeam(),
			Default: func(env models.Environment, vals models.FieldMap) interface{} {
				return pool.CRMTeam().NewSet(env).GetDefaultTeam(pool.User().NewSet(env))
			}},
		"PartnerShipping": models.Many2OneField{String: "Delivery Address", RelationModel: pool.Partner(),
			OnChange: pool.AccountInvoice().Methods().OnchangePartnerShipping(),
			Help:     "Delivery address for current invoice."}, /* readonly=true */ /*[ states {'draft': [('readonly']*/ /*[ False)]}]*/
	})

	pool.AccountInvoice().Fields().Comment().SetDefault(
		func(env models.Environment, vals models.FieldMap) interface{} {
			invoiceType := "out_invoice"
			if env.Context().HasKey("type") {
				invoiceType = env.Context().GetString("type")
			}
			if invoiceType == "out_invoice" {
				return pool.User().NewSet(env).CurrentUser().Company().SaleNote()
			}
			return ""
		})

	pool.AccountInvoice().Methods().OnchangePartnerShipping().DeclareMethod(
		`OnchangePartnerShipping triggers the change of fiscal position
		when the shipping address is modified.`,
		func(rs pool.AccountInvoiceSet) (*pool.AccountInvoiceData, []models.FieldNamer) {
			fiscalPosition := pool.AccountFiscalPosition().NewSet(rs.Env()).GetFiscalPosition(rs.Partner(), rs.PartnerShipping())
			return &pool.AccountInvoiceData{
				FiscalPosition: fiscalPosition,
			}, []models.FieldNamer{pool.AccountInvoice().FiscalPosition()}
		})

	pool.AccountInvoice().Methods().OnchangePartner().Extend("",
		func(rs pool.AccountInvoiceSet) (*pool.AccountInvoiceData, []models.FieldNamer) {
			data, fields := rs.Super().OnchangePartner()
			data.PartnerShipping = rs.Partner().AddressGet([]string{"delivery"})["delivery"]
			fields = append(fields, pool.AccountInvoice().PartnerShipping())
			return data, fields
		})

	//pool.AccountInvoice().Methods().ActionInvoicePaid().Extend("",
	//	func(rs pool.AccountInvoiceSet) bool {
	//		res := rs.Super().ActionInvoicePaid()
	//		todo := make(map[struct {
	//			order pool.SaleOrderSet
	//			name  string
	//		}]bool)
	//		for _, invoice := range rs.Records() {
	//			for _, line := range invoice.InvoiceLines().Records() {
	//				for _, saleLine := range line.SaleLines {
	//					todo[struct {
	//						order pool.SaleOrderSet
	//						name  string
	//					}{
	//						order: saleLine.Order(), name: invoice.Number()}] = true
	//				}
	//			}
	//		}
	//		for key := range todo {
	//			key.order.MessagePost(rs.T("Invoice %s paid", key.name))
	//		}
	//		return res
	//	})

	//pool.AccountInvoice().Methods().OrderLinesLayouted().DeclareMethod(
	//	`OrderLinesLayouted returns this sale order lines ordered by sale_layout_category sequence.
	//	Used to render the report.`,
	//	func(rs pool.AccountInvoiceSet) {
	//		//@api.multi
	//		/*
	//		  self.ensure_one()
	//		  report_pages = [[]]
	//		  for category, lines in groupby(self.invoice_line_ids, lambda l: l.layout_category_id):
	//		      # If last added category induced a pagebreak, this one will be on a new page
	//		      if report_pages[-1] and report_pages[-1][-1]['pagebreak']:
	//		          report_pages.append([])
	//		      # Append category to current report page
	//		      report_pages[-1].append({
	//		          'name': category and category.name or 'Uncategorized',
	//		          'subtotal': category and category.subtotal,
	//		          'pagebreak': category and category.pagebreak,
	//		          'lines': list(lines)
	//		      })
	//
	//		  return report_pages
	//
	//		*/
	//	})

	pool.AccountInvoice().Methods().GetDeliveryPartner().Extend("",
		func(rs pool.AccountInvoiceSet) pool.PartnerSet {
			rs.EnsureOne()
			if !rs.PartnerShipping().IsEmpty() {
				return rs.PartnerShipping()
			}
			return rs.Super().GetDeliveryPartner()
		})

	pool.AccountInvoice().Methods().GetRefundCommonFields().Extend("",
		func(rs pool.AccountInvoiceSet) []models.FieldNamer {
			return append(rs.Super().GetRefundCommonFields(),
				pool.AccountInvoice().Team(), pool.AccountInvoice().PartnerShipping())
		})

	pool.AccountInvoiceLine().SetDefaultOrder("Invoice", "LayoutCategory", "Sequence", "ID")

	pool.AccountInvoiceLine().AddFields(map[string]models.FieldDefinition{
		"SaleLines": models.Many2ManyField{String: "Sale Order Lines", RelationModel: pool.SaleOrderLine(),
			JSON: "sale_line_ids", NoCopy: true /*[ readonly True]*/},
		"LayoutCategory":         models.Many2OneField{String: "Section", RelationModel: pool.SaleLayoutCategory()},
		"LayoutCategorySequence": models.IntegerField{String: "Layout Sequence"},
	})

}
