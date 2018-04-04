// Copyright 2017 NDP Systèmes. All Rights Reserved.
// See LICENSE file for full licensing details.

package sale

import (
	"github.com/hexya-erp/hexya/hexya/models"
	"github.com/hexya-erp/hexya/pool/h"
)

func init() {

	h.AccountInvoice().AddFields(map[string]models.FieldDefinition{
		"Team": models.Many2OneField{String: "Sales Team", RelationModel: h.CRMTeam(),
			Default: func(env models.Environment) interface{} {
				return h.CRMTeam().NewSet(env).GetDefaultTeam(h.User().NewSet(env))
			}},
		"PartnerShipping": models.Many2OneField{String: "Delivery Address", RelationModel: h.Partner(),
			OnChange: h.AccountInvoice().Methods().OnchangePartnerShipping(),
			Help:     "Delivery address for current invoice."}, /* readonly=true */ /*[ states {'draft': [('readonly']*/ /*[ False)]}]*/
	})

	h.AccountInvoice().Fields().Comment().SetDefault(
		func(env models.Environment) interface{} {
			invoiceType := "out_invoice"
			if env.Context().HasKey("type") {
				invoiceType = env.Context().GetString("type")
			}
			if invoiceType == "out_invoice" {
				return h.User().NewSet(env).CurrentUser().Company().SaleNote()
			}
			return ""
		})

	h.AccountInvoice().Methods().OnchangePartnerShipping().DeclareMethod(
		`OnchangePartnerShipping triggers the change of fiscal position
		when the shipping address is modified.`,
		func(rs h.AccountInvoiceSet) (*h.AccountInvoiceData, []models.FieldNamer) {
			fiscalPosition := h.AccountFiscalPosition().NewSet(rs.Env()).GetFiscalPosition(rs.Partner(), rs.PartnerShipping())
			return &h.AccountInvoiceData{
				FiscalPosition: fiscalPosition,
			}, []models.FieldNamer{h.AccountInvoice().FiscalPosition()}
		})

	h.AccountInvoice().Methods().OnchangePartner().Extend("",
		func(rs h.AccountInvoiceSet) (*h.AccountInvoiceData, []models.FieldNamer) {
			data, fields := rs.Super().OnchangePartner()
			data.PartnerShipping = rs.Partner().AddressGet([]string{"delivery"})["delivery"]
			fields = append(fields, h.AccountInvoice().PartnerShipping())
			return data, fields
		})

	//h.AccountInvoice().Methods().ActionInvoicePaid().Extend("",
	//	func(rs h.AccountInvoiceSet) bool {
	//		res := rs.Super().ActionInvoicePaid()
	//		todo := make(map[struct {
	//			order h.SaleOrderSet
	//			name  string
	//		}]bool)
	//		for _, invoice := range rs.Records() {
	//			for _, line := range invoice.InvoiceLines().Records() {
	//				for _, saleLine := range line.SaleLines {
	//					todo[struct {
	//						order h.SaleOrderSet
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

	//h.AccountInvoice().Methods().OrderLinesLayouted().DeclareMethod(
	//	`OrderLinesLayouted returns this sale order lines ordered by sale_layout_category sequence.
	//	Used to render the report.`,
	//	func(rs h.AccountInvoiceSet) {
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

	h.AccountInvoice().Methods().GetDeliveryPartner().Extend("",
		func(rs h.AccountInvoiceSet) h.PartnerSet {
			rs.EnsureOne()
			if !rs.PartnerShipping().IsEmpty() {
				return rs.PartnerShipping()
			}
			return rs.Super().GetDeliveryPartner()
		})

	h.AccountInvoice().Methods().GetRefundCommonFields().Extend("",
		func(rs h.AccountInvoiceSet) []models.FieldNamer {
			return append(rs.Super().GetRefundCommonFields(),
				h.AccountInvoice().Team(), h.AccountInvoice().PartnerShipping())
		})

	h.AccountInvoiceLine().SetDefaultOrder("Invoice", "LayoutCategory", "Sequence", "ID")

	h.AccountInvoiceLine().AddFields(map[string]models.FieldDefinition{
		"SaleLines": models.Many2ManyField{String: "Sale Order Lines", RelationModel: h.SaleOrderLine(),
			JSON: "sale_line_ids", NoCopy: true, ReadOnly: true},
		"LayoutCategory":         models.Many2OneField{String: "Section", RelationModel: h.SaleLayoutCategory()},
		"LayoutCategorySequence": models.IntegerField{String: "Layout Sequence"},
	})

}
