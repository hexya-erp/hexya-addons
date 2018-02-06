// Copyright 2017 NDP Systèmes. All Rights Reserved.
// See LICENSE file for full licensing details.

package sale

import (
	"fmt"

	"github.com/hexya-erp/hexya/hexya/models"
	"github.com/hexya-erp/hexya/hexya/models/types"
	"github.com/hexya-erp/hexya/pool/h"
)

func init() {

	h.SaleReport().DeclareManualModel()
	h.SaleReport().AddFields(map[string]models.FieldDefinition{
		"Name":              models.CharField{String: "Order Reference" /*[ readonly True]*/},
		"Date":              models.DateTimeField{String: "Date Order" /*[ readonly True]*/},
		"Product":           models.Many2OneField{RelationModel: h.ProductProduct() /* readonly=true */},
		"ProductUom":        models.Many2OneField{String: "Unit of Measure", RelationModel: h.ProductUom() /* readonly=true */},
		"ProductUomQty":     models.FloatField{String: "# of Qty" /*[ readonly True]*/},
		"QtyDelivered":      models.FloatField{String: "Qty Delivered" /*[ readonly True]*/},
		"QtyToInvoice":      models.FloatField{String: "Qty To Invoice" /*[ readonly True]*/},
		"QtyInvoiced":       models.FloatField{String: "Qty Invoiced" /*[ readonly True]*/},
		"Partner":           models.Many2OneField{RelationModel: h.Partner() /* readonly=true */},
		"Company":           models.Many2OneField{RelationModel: h.Company() /* readonly=true */},
		"User":              models.Many2OneField{String: "Salesperson", RelationModel: h.User() /* readonly=true */},
		"PriceTotal":        models.FloatField{String: "Total" /*[ readonly True]*/},
		"PriceSubtotal":     models.FloatField{String: "Untaxed Total" /*[ readonly True]*/},
		"ProductTmpl":       models.Many2OneField{String: "Product Template", RelationModel: h.ProductTemplate() /* readonly=true */},
		"Categ":             models.Many2OneField{String: "Product Category", RelationModel: h.ProductCategory() /* readonly=true */},
		"Nbr":               models.IntegerField{String: "# of Lines" /*[ readonly True]*/},
		"Pricelist":         models.Many2OneField{RelationModel: h.ProductPricelist() /* readonly=true */},
		"AnalyticAccount":   models.Many2OneField{RelationModel: h.AccountAnalyticAccount() /* readonly=true */},
		"Team":              models.Many2OneField{String: "Sales Team", RelationModel: h.CRMTeam() /* readonly=true */},
		"Country":           models.Many2OneField{String: "Partner Country", RelationModel: h.Country() /* readonly=true */},
		"CommercialPartner": models.Many2OneField{String: "Commercial Entity", RelationModel: h.Partner() /* readonly=true */},
		"State": models.SelectionField{String: "Status", Selection: types.Selection{
			"draft":  "Draft Quotation",
			"sent":   "Quotation Sent",
			"sale":   "Sales Order",
			"done":   "Sales Done",
			"cancel": "Cancelled",
		} /*[ readonly True]*/},
		"Weight": models.FloatField{String: "Gross Weight" /*[ readonly True]*/},
		"Volume": models.FloatField{ /*[ readonly True]*/ },
	})

	h.SaleReport().Methods().Select().DeclareMethod(
		`Select returns the select clause of the SQL view.`,
		func(rs h.SaleReportSet) string {
			selectStr := fmt.Sprintf(`
			      WITH cur_rate as (%s)
			       SELECT min(l.id) as id,
			              l.product_id as product_id,
			              t.uom_id as product_uom,
			              sum(l.product_uom_qty / u.factor * u2.factor) as product_uom_qty,
			              sum(l.qty_delivered / u.factor * u2.factor) as qty_delivered,
			              sum(l.qty_invoiced / u.factor * u2.factor) as qty_invoiced,
			              sum(l.qty_to_invoice / u.factor * u2.factor) as qty_to_invoice,
			              sum(l.price_total / COALESCE(cr.rate, 1.0)) as price_total,
			              sum(l.price_subtotal / COALESCE(cr.rate, 1.0)) as price_subtotal,
			              count(*) as nbr,
			              s.name as name,
			              s.date_order as date,
			              s.state as state,
			              s.partner_id as partner_id,
			              s.user_id as user_id,
			              s.company_id as company_id,
			              extract(epoch from avg(date_trunc('day',s.date_order)-date_trunc('day',s.create_date)))/(24*60*60)::decimal(16,2) as delay,
			              t.categ_id as categ_id,
			              s.pricelist_id as pricelist_id,
			              s.project_id as analytic_account_id,
			              s.team_id as team_id,
			              p.product_tmpl_id,
			              partner.country_id as country_id,
			              partner.commercial_partner_id as commercial_partner_id,
			              sum(p.weight * l.product_uom_qty / u.factor * u2.factor) as weight,
			              sum(p.volume * l.product_uom_qty / u.factor * u2.factor) as volume
			  `, h.Currency().NewSet(rs.Env()).SelectCompaniesRates())
			return selectStr
		})

	h.SaleReport().Methods().From().DeclareMethod(
		`From returns the from clause of the SQL view.`,
		func(rs h.SaleReportSet) string {
			fromStr := `
			          sale_order_line l
			                join sale_order s on (l.order_id=s.id)
			                join partner on s.partner_id = partner.id
			                  left join product_product p on (l.product_id=p.id)
			                      left join product_template t on (p.product_tmpl_id=t.id)
			              left join product_uom u on (u.id=l.product_uom_id)
			              left join product_uom u2 on (u2.id=t.uom_id)
			              left join product_pricelist pp on (s.pricelist_id = pp.id)
			              left join cur_rate cr on (cr.currency_id = pp.currency_id and
			                  cr.company_id = s.company_id and
			                  cr.date_start <= coalesce(s.date_order, now()) and
			                  (cr.date_end is null or cr.date_end > coalesce(s.date_order, now())))
			`
			return fromStr
		})

	h.SaleReport().Methods().GroupByClause().DeclareMethod(
		`GroupByClause returns the group by clause of the SQL view`,
		func(rs h.SaleReportSet) string {
			groupByStr := `
				GROUP BY l.product_id,
					l.order_id,
					t.uom_id,
					t.categ_id,
					s.name,
					s.date_order,
					s.partner_id,
					s.user_id,
					s.state,
					s.company_id,
					s.pricelist_id,
					s.project_id,
					s.team_id,
					p.product_tmpl_id,
					partner.country_id,
					partner.commercial_partner_id
			`
			return groupByStr
		})

	h.SaleReport().Methods().Init().DeclareMethod(
		`Init initializes the SaleReport view`,
		func(rs h.SaleReportSet) {
			rs.Env().Cr().Execute(fmt.Sprintf(`
				DROP VIEW IF EXISTS sale_report;
				CREATE or REPLACE VIEW sale_report as (
					%s
				FROM ( %s )
					%s
				)`, rs.Select(), rs.From(), rs.GroupByClause()))
		})

}
