// Copyright 2017 NDP Systèmes. All Rights Reserved.
// See LICENSE file for full licensing details.

package product

func init() {

	//h.BaseConfigSettings().AddFields(map[string]models.FieldDefinition{
	//	"CompanyShareProduct": models.BooleanField{String: "Share product to all companies" /*['Share product to all companies']*/, Help: "Share your product to all companies defined in your instance.\n  * Checked : Product are visible for every company, even if a company is defined on the partner.\n  * Unchecked : Each company can see only its product (product where company is defined). Product not related to a company are visible for all companies." /*[ even if a company is defined on the partner.\n" " * Unchecked : Each company can see only its product (product where company is defined). Product not related to a company are visible for all companies."]*/},
	//	"GroupProductVariant": models.SelectionField{String: "Product Variants", Selection: types.Selection{
	//		"0": "No variants on products",
	//		"1": "Products can have several attributes defining variants (Example: size color...)",
	//	}, /*[]*/ /*["Product Variants"]*/ /*, ImpliedGroup :"product.group_product_variant"*/ Help: "Work with product variant allows you to define some variant of the same products" /*[ an ease the product management in the ecommerce for example']*/},
	//})
	//h.BaseConfigSettings().Methods().GetDefaultCompanyShareProduct().DeclareMethod(
	//	`GetDefaultCompanyShareProduct`,
	//	func(rs h.BaseConfigSettingsSet, args struct {
	//		Fields interface{}
	//	}) {
	//		//@api.model
	//		/*def get_default_company_share_product(self, fields):
	//		  product_rule = self.env.ref('product.product_comp_rule')
	//		  return {
	//		      'company_share_product': not bool(product_rule.active)
	//		  }
	//
	//		*/
	//	})
	//h.BaseConfigSettings().Methods().SetAuthCompanyShareProduct().DeclareMethod(
	//	`SetAuthCompanyShareProduct`,
	//	func(rs h.BaseConfigSettingsSet) {
	//		//@api.multi
	//		/*def set_auth_company_share_product(self):
	//		  self.ensure_one()
	//		  product_rule = self.env.ref('product.product_comp_rule')
	//		  product_rule.write({'active': not bool(self.company_share_product)})
	//		*/
	//	})

}
