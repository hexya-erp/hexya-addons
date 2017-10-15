package product

import (
	"github.com/hexya-erp/hexya/hexya/models"
	"github.com/hexya-erp/hexya/pool"
)

func init() {

	pool.ProductPriceListWizard().DeclareTransientModel()

	pool.ProductPriceListWizard().AddFields(map[string]models.FieldDefinition{
		"PriceList": models.Many2OneField{RelationModel: pool.ProductPricelist(), Required: true},
		"Qty1":      models.IntegerField{String: "Quantity-1", Default: models.DefaultValue(1)},
		"Qty2":      models.IntegerField{String: "Quantity-2", Default: models.DefaultValue(5)},
		"Qty3":      models.IntegerField{String: "Quantity-3", Default: models.DefaultValue(10)},
		"Qty4":      models.IntegerField{String: "Quantity-4", Default: models.DefaultValue(0)},
		"Qty5":      models.IntegerField{String: "Quantity-5", Default: models.DefaultValue(0)},
	})

	pool.ProductPriceListWizard().Methods().PrintReport().DeclareMethod(
		`PrintReport`,
		func(rs pool.ProductPriceListWizardSet) {
			//@api.multi
			/*def print_report(self):
			  """
			  To get the date and print the report
			  @return : return report
			  """
			  datas = {'ids': self.env.context.get('active_ids', [])}
			  res = self.read(['price_list', 'qty1', 'qty2', 'qty3', 'qty4', 'qty5'])
			  res = res and res[0] or {}
			  res['price_list'] = res['price_list'][0]
			  datas['form'] = res
			  return self.env['report'].get_action([], 'product.report_pricelist', data=datas)
			*/
		})

}
