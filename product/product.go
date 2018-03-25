// Copyright 2017 NDP Systèmes. All Rights Reserved.
// See LICENSE file for full licensing details.

package product

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/hexya-erp/hexya-addons/decimalPrecision"
	"github.com/hexya-erp/hexya-base/base"
	"github.com/hexya-erp/hexya/hexya/actions"
	"github.com/hexya-erp/hexya/hexya/models"
	"github.com/hexya-erp/hexya/hexya/models/operator"
	"github.com/hexya-erp/hexya/hexya/models/security"
	"github.com/hexya-erp/hexya/hexya/models/types"
	"github.com/hexya-erp/hexya/hexya/models/types/dates"
	"github.com/hexya-erp/hexya/pool/h"
	"github.com/hexya-erp/hexya/pool/q"
)

func init() {

	h.ProductCategory().DeclareModel()
	h.ProductCategory().SetDefaultOrder("Parent.Name")

	h.ProductCategory().AddFields(map[string]models.FieldDefinition{
		"Name": models.CharField{Index: true, Required: true, Translate: true},
		"Parent": models.Many2OneField{String: "Parent Category", RelationModel: h.ProductCategory(), Index: true,
			OnDelete: models.Cascade, Constraint: h.ProductCategory().Methods().CheckCategoryRecursion()},
		"Children": models.One2ManyField{String: "Child Categories", RelationModel: h.ProductCategory(),
			ReverseFK: "Parent", JSON: "child_id"},
		"Type": models.SelectionField{String: "Category Type", Selection: types.Selection{"view": "View", "normal": "Normal"},
			Default: models.DefaultValue("normal"), Help: "A category of the view type is a virtual category that can be used as the parent of another category to create a hierarchical structure."},
		"Products": models.One2ManyField{RelationModel: h.ProductTemplate(), ReverseFK: "Categ"},
		"ProductCount": models.IntegerField{String: "# Products", Compute: h.ProductCategory().Methods().ComputeProductCount(),
			Help:    "The number of products under this category (Does not consider the children categories)",
			Depends: []string{"Products"}, GoType: new(int)},
	})

	h.ProductCategory().Methods().ComputeProductCount().DeclareMethod(
		`ComputeProductCount returns the number of products within this category (not considering children categories)`,
		func(rs h.ProductCategorySet) *h.ProductCategoryData {
			return &h.ProductCategoryData{
				ProductCount: h.ProductTemplate().Search(rs.Env(), q.ProductTemplate().Categ().Equals(rs)).SearchCount(),
			}
		})

	h.ProductCategory().Methods().CheckCategoryRecursion().DeclareMethod(
		`CheckCategoryRecursion panics if there is a recursion in the category tree.`,
		func(rs h.ProductCategorySet) {
			if !rs.CheckRecursion() {
				log.Panic(rs.T("Error ! You cannot create recursive categories."))
			}
		})

	h.ProductCategory().Methods().NameGet().Extend("",
		func(rs h.ProductCategorySet) string {
			var names []string
			for current := rs; !current.IsEmpty(); current = current.Parent() {
				names = append([]string{current.Name()}, names...)
			}
			return strings.Join(names, " / ")
		})

	h.ProductCategory().Methods().SearchByName().Extend("",
		func(rs h.ProductCategorySet, name string, op operator.Operator, additionalCond q.ProductCategoryCondition, limit int) h.ProductCategorySet {
			if name == "" {
				return rs.Super().SearchByName(name, op, additionalCond, limit)
			}
			// Be sure name_search is symetric to name_get
			categoryNames := strings.Split(name, " / ")
			child := categoryNames[len(categoryNames)-1]
			cond := q.ProductCategory().Name().AddOperator(op, child)
			var categories h.ProductCategorySet
			if len(categoryNames) > 1 {
				parents := rs.SearchByName(strings.Join(categoryNames[:len(categoryNames)-1], " / "), operator.IContains, additionalCond, limit)
				if op.IsNegative() {
					categories = h.ProductCategory().Search(rs.Env(), q.ProductCategory().ID().NotIn(parents.Ids()))
					cond = cond.Or().Parent().In(categories)
				} else {
					cond = cond.And().Parent().In(parents)
				}
				for i := 1; i < len(categoryNames); i++ {
					if op.IsNegative() {
						cond = cond.AndCond(q.ProductCategory().Name().AddOperator(op, strings.Join(categoryNames[len(categoryNames)-1-i:], " / ")))
					} else {
						cond = cond.OrCond(q.ProductCategory().Name().AddOperator(op, strings.Join(categoryNames[len(categoryNames)-1-i:], " / ")))
					}
				}
			}
			return h.ProductCategory().Search(rs.Env(), cond.AndCond(additionalCond))
		})

	h.ProductPriceHistory().DeclareModel()
	h.ProductPriceHistory().SetDefaultOrder("Datetime DESC")

	h.ProductPriceHistory().AddFields(map[string]models.FieldDefinition{
		"Company": models.Many2OneField{RelationModel: h.Company(),
			Default: func(env models.Environment) interface{} {
				if env.Context().HasKey("force_company") {

					return h.Company().Browse(env, []int64{env.Context().GetInteger("force_company")})
				}
				currentUser := h.User().NewSet(env).CurrentUser()
				return currentUser.Company()
			}, Required: true},
		"Product": models.Many2OneField{RelationModel: h.ProductProduct(), JSON: "product_id",
			OnDelete: models.Cascade, Required: true},
		"Datetime": models.DateTimeField{String: "Date", Default: func(env models.Environment) interface{} {
			return dates.Now()
		}},
		"Cost": models.FloatField{String: "Cost", Digits: decimalPrecision.GetPrecision("Product Price")},
	})

	h.ProductProduct().DeclareModel()
	h.ProductProduct().SetDefaultOrder("DefaultCode", "Name", "ID")

	h.ProductProduct().AddFields(map[string]models.FieldDefinition{
		"Price": models.FloatField{Compute: h.ProductProduct().Methods().ComputeProductPrice(),
			Digits:  decimalPrecision.GetPrecision("Product Price"),
			Inverse: h.ProductProduct().Methods().InverseProductPrice()},
		"PriceExtra": models.FloatField{String: "Variant Price Extra",
			Compute: h.ProductProduct().Methods().ComputeProductPriceExtra(),
			Depends: []string{"AttributeValues", "AttributeValues.Prices", "AttributeValues.Prices.PriceExtra", "AttributeValues.Prices.ProductTmpl"},
			Digits:  decimalPrecision.GetPrecision("Product Price"),
			Help:    "This is the sum of the extra price of all attributes"},
		"LstPrice": models.FloatField{String: "Sale Price",
			Compute: h.ProductProduct().Methods().ComputeProductLstPrice(),
			Depends: []string{"ListPrice", "PriceExtra"},
			Digits:  decimalPrecision.GetPrecision("Product Price"),
			Inverse: h.ProductProduct().Methods().InverseProductLstPrice(),
			Help:    "The sale price is managed from the product template. Click on the 'Variant Prices' button to set the extra attribute prices."},
		"DefaultCode": models.CharField{String: "Internal Reference", Index: true},
		"Code": models.CharField{String: "Internal Reference",
			Compute: h.ProductProduct().Methods().ComputeProductCode(), Depends: []string{""}},
		"PartnerRef": models.CharField{String: "Customer Ref",
			Compute: h.ProductProduct().Methods().ComputePartnerRef(), Depends: []string{""}},
		"Active": models.BooleanField{String: "Active",
			Default: models.DefaultValue(true),
			Help:    "If unchecked, it will allow you to hide the product without removing it."},
		"ProductTmpl": models.Many2OneField{String: "Product Template", RelationModel: h.ProductTemplate(),
			Index: true, OnDelete: models.Cascade, Required: true, Embed: true},
		"Barcode": models.CharField{String: "Barcode", NoCopy: true, /*Unique: true,*/
			Help: "International Article Number used for product identification."},
		"AttributeValues": models.Many2ManyField{String: "Attributes", RelationModel: h.ProductAttributeValue(),
			JSON:       "attribute_value_ids", /*, OnDelete: models.Restrict*/
			Constraint: h.ProductProduct().Methods().CheckAttributeValueIds()},
		"ImageVariant": models.BinaryField{String: "Variant Image",
			Help: "This field holds the image used as image for the product variant, limited to 1024x1024px."},
		"Image": models.BinaryField{String: "Big-sized image",
			Compute: h.ProductProduct().Methods().ComputeImages(),
			Depends: []string{"ImageVariant", "ProductTmpl", "ProductTmpl.Image"},
			Inverse: h.ProductProduct().Methods().InverseImageValue(),
			Help: `Image of the product variant (Big-sized image of product template if false). It is automatically
resized as a 1024x1024px image, with aspect ratio preserved.`},
		"ImageSmall": models.BinaryField{String: "Small-sized image",
			Compute: h.ProductProduct().Methods().ComputeImages(),
			Depends: []string{"ImageVariant", "ProductTmpl", "ProductTmpl.Image"},
			Inverse: h.ProductProduct().Methods().InverseImageValue(),
			Help:    "Image of the product variant (Small-sized image of product template if false)."},
		"ImageMedium": models.BinaryField{String: "Medium-sized image",
			Compute: h.ProductProduct().Methods().ComputeImages(),
			Depends: []string{"ImageVariant", "ProductTmpl", "ProductTmpl.Image"},
			Inverse: h.ProductProduct().Methods().InverseImageValue(),
			Help:    "Image of the product variant (Medium-sized image of product template if false)."},
		"StandardPrice": models.FloatField{String: "Cost", /*, CompanyDependent : true*/
			Digits: decimalPrecision.GetPrecision("Product Price"),
			Help: `Cost of the product template used for standard stock valuation in accounting and used as a
base price on purchase orders. Expressed in the default unit of measure of the product.`},
		"Volume": models.FloatField{Help: "The volume in m3."},
		"Weight": models.FloatField{Digits: decimalPrecision.GetPrecision("Stock Weight"),
			Help: "The weight of the contents in Kg, not including any packaging, etc."},
		"PricelistItems": models.Many2ManyField{RelationModel: h.ProductPricelistItem(),
			JSON: "pricelist_item_ids", Compute: h.ProductProduct().Methods().GetPricelistItems()},
	})

	h.ProductProduct().Fields().StandardPrice().RevokeAccess(security.GroupEveryone, security.All).GrantAccess(base.GroupUser, security.All)

	h.ProductProduct().Methods().ComputeProductPrice().DeclareMethod(
		`ComputeProductPrice computes the price of this product based on the given context keys:

		- 'partner' => int64 (id of the partner)
		- 'pricelist' => int64 (id of the price list)
		- 'quantity' => float64`,
		func(rs h.ProductProductSet) *h.ProductProductData {
			if !rs.Env().Context().HasKey("pricelist") {
				return new(h.ProductProductData)
			}
			priceListID := rs.Env().Context().GetInteger("pricelist")
			priceList := h.ProductPricelist().Browse(rs.Env(), []int64{priceListID})
			if priceList.IsEmpty() {
				return new(h.ProductProductData)
			}
			quantity := rs.Env().Context().GetFloat("quantity")
			if quantity == 0 {
				quantity = 1
			}
			partnerID := rs.Env().Context().GetInteger("partner")
			partner := h.Partner().Browse(rs.Env(), []int64{partnerID})
			return &h.ProductProductData{
				Price: priceList.GetProductPrice(rs, quantity, partner, dates.Today(), h.ProductUom().NewSet(rs.Env())),
			}
		})

	h.ProductProduct().Methods().InverseProductPrice().DeclareMethod(
		`InverseProductPrice updates ListPrice from the given Price`,
		func(rs h.ProductProductSet, price float64) {
			if rs.Env().Context().HasKey("uom") {
				price = h.ProductUom().Browse(rs.Env(), []int64{rs.Env().Context().GetInteger("uom")}).ComputePrice(price, rs.Uom())
			}
			price -= rs.PriceExtra()
			rs.SetListPrice(price)
		})

	h.ProductProduct().Methods().InverseProductLstPrice().DeclareMethod(
		`InverseProductLstPrice updates ListPrice from the given LstPrice`,
		func(rs h.ProductProductSet, price float64) {
			if rs.Env().Context().HasKey("uom") {
				price = h.ProductUom().Browse(rs.Env(), []int64{rs.Env().Context().GetInteger("uom")}).ComputePrice(price, rs.Uom())
			}
			price -= rs.PriceExtra()
			rs.SetListPrice(price)
		})

	h.ProductProduct().Methods().ComputeProductPriceExtra().DeclareMethod(
		`ComputeProductPriceExtra computes the price extra of this product by suming the extras of each attribute`,
		func(rs h.ProductProductSet) *h.ProductProductData {
			var priceExtra float64
			for _, attributeValue := range rs.AttributeValues().Records() {
				for _, attributePrice := range attributeValue.Prices().Records() {
					if attributePrice.ProductTmpl().Equals(rs.ProductTmpl()) {
						priceExtra += attributePrice.PriceExtra()
					}
				}
			}
			return &h.ProductProductData{
				PriceExtra: priceExtra,
			}
		})

	h.ProductProduct().Methods().ComputeProductLstPrice().DeclareMethod(
		`ComputeProductLstPrice computes the LstPrice from the ListPrice and the extras`,
		func(rs h.ProductProductSet) *h.ProductProductData {
			listPrice := rs.ListPrice()
			if rs.Env().Context().HasKey("uom") {
				toUoM := h.ProductUom().Browse(rs.Env(), []int64{rs.Env().Context().GetInteger("uom")})
				listPrice = rs.Uom().ComputePrice(listPrice, toUoM)
			}
			return &h.ProductProductData{
				LstPrice: listPrice + rs.PriceExtra(),
			}
		})

	h.ProductProduct().Methods().ComputeProductCode().DeclareMethod(
		`ComputeProductCode computes the product code based on the context:
- 'partner_id' => int64 (id of the considered partner)`,
		func(rs h.ProductProductSet) *h.ProductProductData {
			var code string
			for _, supplierInfo := range rs.Sellers().Records() {
				if supplierInfo.Name().ID() == rs.Env().Context().GetInteger("partner_id") {
					code = supplierInfo.ProductCode()
					if code != "" {
						break
					}
				}
			}
			if code != "" {
				code = rs.DefaultCode()
			}
			return &h.ProductProductData{
				Code: code,
			}
		})

	h.ProductProduct().Methods().ComputePartnerRef().DeclareMethod(
		`ComputePartnerRef computes the product's reference (i.e. "[code] description") based on the context:
- 'partner_id' => int64 (id of the considered partner)`,
		func(rs h.ProductProductSet) *h.ProductProductData {
			var code, productName string
			for _, supplierInfo := range rs.Sellers().Records() {
				if supplierInfo.Name().ID() == rs.Env().Context().GetInteger("partner_id") {
					code = supplierInfo.ProductCode()
					if code != "" {
						break
					}
					productName = supplierInfo.ProductName()
					if productName != "" {
						break
					}
				}
			}
			if code == "" {
				code = rs.DefaultCode()
			}
			if productName == "" {
				productName = rs.Name()
			}
			return &h.ProductProductData{
				PartnerRef: rs.NameFormat(productName, code),
			}
		})

	h.ProductProduct().Methods().ComputeImages().DeclareMethod(
		`ComputeImages computes the images in different sizes.`,
		func(rs h.ProductProductSet) *h.ProductProductData {
			// TODO implement image resizing
			//@api.depends('image_variant','product_tmpl_id.image')
			/*def _compute_images(self):
			  if self._context.get('bin_size'):
			      self.image_medium = self.image_variant
			      self.image_small = self.image_variant
			      self.image = self.image_variant
			  else:
			      resized_images = tools.image_get_resized_images(self.image_variant, return_big=True, avoid_resize_medium=True)
			      self.image_medium = resized_images['image_medium']
			      self.image_small = resized_images['image_small']
			      self.image = resized_images['image']
			  if not self.image_medium:
			      self.image_medium = self.product_tmpl_id.image_medium
			  if not self.image_small:
			      self.image_small = self.product_tmpl_id.image_small
			  if not self.image:
			      self.image = self.product_tmpl_id.image

			*/
			return &h.ProductProductData{}
		})

	h.ProductProduct().Methods().InverseImageValue().DeclareMethod(
		`InverseImageValue sets all images from the given image`,
		func(rs h.ProductProductSet, image string) {
			// TODO Resize image
			//image = tools.image_resize_image_big(value)
			if rs.ProductTmpl().Image() == "" {
				rs.ProductTmpl().SetImage(image)
				return
			}
			rs.SetImageVariant(image)
		})

	h.ProductProduct().Methods().GetPricelistItems().DeclareMethod(
		`GetPricelistItems returns all price list items for this product`,
		func(rs h.ProductProductSet) *h.ProductProductData {
			rs.EnsureOne()
			priceListItems := h.ProductPricelistItem().Search(rs.Env(),
				q.ProductPricelistItem().Product().Equals(rs).Or().ProductTmpl().Equals(rs.ProductTmpl()))
			return &h.ProductProductData{
				PricelistItems: priceListItems,
			}
		})

	h.ProductProduct().Methods().CheckAttributeValueIds().DeclareMethod(
		`CheckAttributeValueIds checks that we do not have more than one value per attribute.`,
		func(rs h.ProductProductSet) {
			attributes := h.ProductAttribute().NewSet(rs.Env())
			for _, value := range rs.AttributeValues().Records() {
				if !value.Attribute().Intersect(attributes).IsEmpty() {
					log.Panic(rs.T("Error! It is not allowed to choose more than one value for a given attribute."))
				}
				attributes = attributes.Union(value.Attribute())
			}
		})

	h.ProductProduct().Methods().OnchangeUom().DeclareMethod(
		`OnchangeUom process UI triggers when changing th UoM`,
		func(rs h.ProductProductSet) (*h.ProductProductData, []models.FieldNamer) {
			if !rs.Uom().IsEmpty() && !rs.UomPo().IsEmpty() && !rs.Uom().Category().Equals(rs.UomPo().Category()) {
				return &h.ProductProductData{
					UomPo: rs.Uom(),
				}, []models.FieldNamer{h.ProductProduct().UomPo()}
			}
			return new(h.ProductProductData), []models.FieldNamer{}
		})

	h.ProductProduct().Methods().Create().Extend("",
		func(rs h.ProductProductSet, data *h.ProductProductData) h.ProductProductSet {
			product := rs.WithContext("create_product_product", true).Super().Create(data)
			// When a unique variant is created from tmpl then the standard price is set by DefineStandardPrice
			if !rs.Env().Context().HasKey("create_from_tmpl") && product.ProductTmpl().ProductVariants().Len() == 1 {
				product.DefineStandardPrice(data.StandardPrice)
			}
			return product
		})

	h.ProductProduct().Methods().Write().Extend("",
		func(rs h.ProductProductSet, data *h.ProductProductData, fieldsToReset ...models.FieldNamer) bool {
			// Store the standard price change in order to be able to retrieve the cost of a product for a given date
			res := rs.Super().Write(data, fieldsToReset...)
			if _, ok := data.Get(h.ProductProduct().StandardPrice(), fieldsToReset...); ok {
				rs.DefineStandardPrice(data.StandardPrice)
			}
			return res
		})

	h.ProductProduct().Methods().Unlink().Extend("",
		func(rs h.ProductProductSet) int64 {
			unlinkProducts := h.ProductProduct().NewSet(rs.Env())
			unlinkTemplates := h.ProductTemplate().NewSet(rs.Env())
			for _, product := range rs.Records() {
				// Check if the product is last product of this template
				otherProducts := h.ProductProduct().Search(rs.Env(),
					q.ProductProduct().ProductTmpl().Equals(product.ProductTmpl()).And().ID().NotEquals(product.ID()))
				if otherProducts.IsEmpty() {
					unlinkTemplates = unlinkTemplates.Union(product.ProductTmpl())
				}
				unlinkProducts = unlinkProducts.Union(product)
			}
			res := unlinkProducts.Super().Unlink()
			// delete templates after calling super, as deleting template could lead to deleting
			// products due to ondelete='cascade'
			unlinkTemplates.Unlink()
			return res
		})

	h.ProductProduct().Methods().Copy().Extend("",
		func(rs h.ProductProductSet, overrides *h.ProductProductData, fieldsToReset ...models.FieldNamer) h.ProductProductSet {
			if rs.Env().Context().HasKey("variant") {
				// if we copy a variant or create one, we keep the same template
				overrides.ProductTmpl = rs.ProductTmpl()
				fieldsToReset = append(fieldsToReset, h.ProductProduct().ProductTmpl())
			} else if _, ok := overrides.Get(h.ProductProduct().Name(), fieldsToReset...); !ok {
				overrides.Name = rs.Name()
				fieldsToReset = append(fieldsToReset, h.ProductProduct().Name())
			}
			return rs.Super().Copy(overrides, fieldsToReset...)
		})

	h.ProductProduct().Methods().Search().Extend("",
		func(rs h.ProductProductSet, cond q.ProductProductCondition) h.ProductProductSet {
			// FIXME: strange...
			if categID := rs.Env().Context().GetInteger("search_default_categ_id"); categID != 0 {
				categ := h.ProductCategory().Browse(rs.Env(), []int64{categID})
				cond = cond.AndCond(q.ProductProduct().Categ().ChildOf(categ))
			}
			return rs.Super().Search(cond)
		})

	h.ProductProduct().Methods().NameFormat().DeclareMethod(
		`NameFormat formats a product name string from the given arguments`,
		func(rs h.ProductProductSet, name, code string) string {
			if code == "" ||
				(rs.Env().Context().HasKey("display_default_code") && !rs.Env().Context().GetBool("display_default_code")) {
				return name
			}
			return fmt.Sprintf("[%s] %s", code, name)
		})

	h.ProductProduct().Methods().NameGet().Extend("",
		func(rs h.ProductProductSet) string {
			/*
			   def _name_get(d):
			       name = d.get('name', '')
			       code = self._context.get('display_default_code', True) and d.get('default_code', False) or False
			       if code:
			           name = '[%s] %s' % (code,name)
			       return (d['id'], name)

			   partner_id = self._context.get('partner_id')
			   if partner_id:
			       partner_ids = [partner_id, self.env['res.partner'].browse(partner_id).commercial_partner_id.id]
			   else:
			       partner_ids = []

			   # all user don't have access to seller and partner
			   # check access and use superuser
			   self.check_access_rights("read")
			   self.check_access_rule("read")

			   result = []
			   for product in self.sudo():
			       # display only the attributes with multiple possible values on the template
			       variable_attributes = product.attribute_line_ids.filtered(lambda l: len(l.value_ids) > 1).mapped('attribute_id')
			       variant = product.attribute_value_ids._variant_name(variable_attributes)

			       name = variant and "%s (%s)" % (product.name, variant) or product.name
			       sellers = []
			       if partner_ids:
			           sellers = [x for x in product.seller_ids if (x.name.id in partner_ids) and (x.product_id == product)]
			           if not sellers:
			               sellers = [x for x in product.seller_ids if (x.name.id in partner_ids) and not x.product_id]
			       if sellers:
			           for s in sellers:
			               seller_variant = s.product_name and (
			                   variant and "%s (%s)" % (s.product_name, variant) or s.product_name
			                   ) or False
			               mydict = {
			                         'id': product.id,
			                         'name': seller_variant or name,
			                         'default_code': s.product_code or product.default_code,
			                         }
			               temp = _name_get(mydict)
			               if temp not in result:
			                   result.append(temp)
			       else:
			           mydict = {
			                     'id': product.id,
			                     'name': name,
			                     'default_code': product.default_code,
			                     }
			           result.append(_name_get(mydict))
			   return result
			*/
			// display only the attributes with multiple possible values on the template
			variableAttributes := h.ProductAttribute().NewSet(rs.Env())
			for _, attrLine := range rs.AttributeLines().Records() {
				if attrLine.Values().Len() > 1 {
					variableAttributes = variableAttributes.Union(attrLine.Attribute())
				}
			}
			variant := rs.AttributeValues().VariantName(variableAttributes)
			if variant != "" {
				return fmt.Sprintf("%s (%s)", rs.PartnerRef(), variant)
			}
			return rs.PartnerRef()
		})

	h.ProductProduct().Methods().SearchByName().Extend("",
		func(rs h.ProductProductSet, name string, op operator.Operator, additionalCond q.ProductProductCondition, limit int) h.ProductProductSet {
			if name == "" {
				return rs.Super().SearchByName(name, op, additionalCond, limit)
			}
			products := h.ProductProduct().NewSet(rs.Env())
			if op.IsPositive() {
				products = rs.Search(q.ProductProduct().DefaultCode().Equals(name).AndCond(additionalCond)).Limit(limit)
				if products.IsEmpty() {
					products = rs.Search(q.ProductProduct().Barcode().Equals(name).AndCond(additionalCond)).Limit(limit)
				}
			}
			switch {
			case products.IsEmpty() && !op.IsNegative():
				// Do not merge the 2 next lines into one single search, SQL search performance would be abysmal
				// on a database with thousands of matching products, due to the huge merge+unique needed for the
				// OR operator (and given the fact that the 'name' lookup results come from the ir.translation table
				// Performing a quick memory merge of ids in Python will give much better performance
				products = h.ProductProduct().Search(rs.Env(), q.ProductProduct().DefaultCode().AddOperator(op, name)).Limit(limit)
				if limit == 0 || products.Len() < limit {
					// we may underrun the limit because of dupes in the results, that's fine
					limit2 := limit - products.Len()
					if limit2 < 0 {
						limit2 = 0
					}
					products = products.Union(h.ProductProduct().Search(rs.Env(),
						q.ProductProduct().Name().AddOperator(op, name).And().ID().NotIn(products.Ids())))
				}
			case products.IsEmpty() && op.IsNegative():
				products = h.ProductProduct().Search(rs.Env(),
					q.ProductProduct().DefaultCode().AddOperator(op, name).And().Name().AddOperator(op, name).AndCond(additionalCond))
			}
			if products.IsEmpty() && op.IsPositive() {
				ptrn, _ := regexp.Compile(`(\[(.*?)\])`)
				res := ptrn.FindAllString(name, -1)
				if len(res) > 1 {
					products = h.ProductProduct().Search(rs.Env(),
						q.ProductProduct().DefaultCode().Equals(res[1]).AndCond(additionalCond))
				}
			}
			// still no results, partner in context: search on supplier info as last hope to find something
			if products.IsEmpty() && rs.Env().Context().HasKey("partner_id") {
				partner := h.Partner().Browse(rs.Env(), []int64{rs.Env().Context().GetInteger("partner_id")})
				suppliers := h.ProductSupplierinfo().Search(rs.Env(),
					q.ProductSupplierinfo().Name().Equals(partner).
						AndCond(q.ProductSupplierinfo().ProductCode().AddOperator(op, name).Or().ProductName().AddOperator(op, name)))
				if !suppliers.IsEmpty() {
					products = h.ProductProduct().Search(rs.Env(),
						q.ProductProduct().ProductTmplFilteredOn(q.ProductTemplate().Sellers().In(suppliers)))
				}
			}
			return products
		})

	h.ProductProduct().Methods().OpenProductTemplate().DeclareMethod(
		`OpenProductTemplate is a utility method used to add an "Open Template" button in product views`,
		func(rs h.ProductProductSet) *actions.Action {
			rs.EnsureOne()
			return &actions.Action{
				Type:     actions.ActionActWindow,
				Model:    "ProductTemplate",
				ViewMode: "form",
				ResID:    rs.ProductTmpl().ID(),
				Target:   "new",
			}
		})

	h.ProductProduct().Methods().SelectSeller().DeclareMethod(
		`SelectSeller returns the ProductSupplierInfo to use for the given partner, quantity, date and UoM.
		If any of the parameters are their Go zero value, then they are not used for filtering.`,
		func(rs h.ProductProductSet, partner h.PartnerSet, quantity float64, date dates.Date, uom h.ProductUomSet) h.ProductSupplierinfoSet {
			rs.EnsureOne()
			if date.IsZero() {
				date = dates.Today()
			}
			res := h.ProductSupplierinfo().NewSet(rs.Env())
			for _, seller := range rs.Sellers().Records() {
				quantityUomSeller := quantity
				if quantityUomSeller != 0 && !uom.IsEmpty() && !uom.Equals(seller.ProductUom()) {
					quantityUomSeller = uom.ComputeQuantity(quantityUomSeller, seller.ProductUom(), true)
				}
				if !seller.DateStart().IsZero() && seller.DateStart().Greater(date) {
					continue
				}
				if !seller.DateEnd().IsZero() && seller.DateEnd().Lower(date) {
					continue
				}
				if !partner.IsEmpty() && seller.Name().Intersect(partner.Union(partner.Parent())).IsEmpty() {
					continue
				}
				if quantityUomSeller < seller.MinQty() {
					continue
				}
				if !seller.Product().IsEmpty() && !seller.Product().Equals(rs) {
					continue
				}
				res = res.Union(seller)
				break
			}
			return res
		})

	h.ProductProduct().Methods().PriceCompute().DeclareMethod(
		`PriceCompute returns the price field defined by priceType in the given uom and currency
		for the given company.`,
		func(rs h.ProductProductSet, priceType models.FieldNamer, uom h.ProductUomSet, currency h.CurrencySet, company h.CompanySet) float64 {
			rs.EnsureOne()
			product := rs
			if priceType == h.ProductProduct().StandardPrice() {
				// StandardPrice field can only be seen by users in base.group_user
				// Thus, in order to compute the sale price from the cost for users not in this group
				// We fetch the standard price as the superuser
				if company.IsEmpty() {
					company = h.User().NewSet(rs.Env()).CurrentUser().Company()
					if rs.Env().Context().HasKey("force_company") {
						company = h.Company().Browse(rs.Env(), []int64{rs.Env().Context().GetInteger("force_company")})
					}
				}
				product = rs.WithContext("force_company", company.ID()).Sudo()
			}

			price := product.Get(priceType.String()).(float64)
			if priceType == h.ProductProduct().ListPrice() {
				price += product.PriceExtra()
			}

			if !uom.IsEmpty() {
				price = product.Uom().ComputePrice(price, uom)
			}
			// Convert from current user company currency to asked one
			// This is right cause a field cannot be in more than one currency
			if !currency.IsEmpty() {
				price = product.Currency().Compute(price, currency, true)
			}
			return price
		})

	h.ProductProduct().Methods().DefineStandardPrice().DeclareMethod(
		`DefineStandardPrice stores the standard price change in order to be able to retrieve the cost of a product for
		a given date`,
		func(rs h.ProductProductSet, value float64) {
			company := h.User().NewSet(rs.Env()).CurrentUser().Company()
			if rs.Env().Context().HasKey("force_company") {
				company = h.Company().Browse(rs.Env(), []int64{rs.Env().Context().GetInteger("force_company")})
			}
			for _, product := range rs.Records() {
				h.ProductPriceHistory().Create(rs.Env(), &h.ProductPriceHistoryData{
					Product: product,
					Cost:    value,
					Company: company,
				})
			}
		})

	h.ProductProduct().Methods().GetHistoryPrice().DeclareMethod(
		`GetHistoryPrice returns the standard price of this product for the given company at the given date`,
		func(rs h.ProductProductSet, company h.CompanySet, date dates.DateTime) float64 {
			if date.IsZero() {
				date = dates.Now()
			}
			history := h.ProductPriceHistory().Search(rs.Env(),
				q.ProductPriceHistory().Company().Equals(company).
					And().Product().In(rs).
					And().Datetime().LowerOrEqual(date)).Limit(1)
			return history.Cost()
		})

	h.ProductProduct().Methods().NeedProcurement().DeclareMethod(
		`NeedProcurement`,
		func(rs h.ProductProductSet) bool {
			// When sale/product is installed alone, there is no need to create procurements. Only
			// sale_stock and sale_service need procurements
			return false
		})

	h.ProductPackaging().DeclareModel()
	h.ProductPackaging().SetDefaultOrder("Sequence")

	h.ProductPackaging().AddFields(map[string]models.FieldDefinition{
		"Name": models.CharField{String: "Packaging Type", Required: true},
		"Sequence": models.IntegerField{Default: models.DefaultValue(1),
			Help: "The first in the sequence is the default one."},
		"ProductTmpl": models.Many2OneField{String: "Product", RelationModel: h.ProductTemplate()},
		"Qty": models.FloatField{String: "Quantity per Package",
			Help: "The total number of products you can have per pallet or box."},
	})

	h.ProductSupplierinfo().DeclareModel()
	h.ProductSupplierinfo().SetDefaultOrder("Sequence", "MinQty DESC", "Price")

	h.ProductSupplierinfo().AddFields(map[string]models.FieldDefinition{
		"Name": models.Many2OneField{String: "Vendor", RelationModel: h.Partner(), JSON: "name",
			Filter: q.Partner().Supplier().Equals(true), OnDelete: models.Cascade, Required: true,
			Help: "Vendor of this product"},
		"ProductName": models.CharField{String: "Vendor Product Name",
			Help: `This vendor's product name will be used when printing a request for quotation.
Keep empty to use the internal one.`},
		"ProductCode": models.CharField{String: "Vendor Product Code",
			Help: `This vendor's product code will be used when printing a request for quotation.
Keep empty to use the internal one.`},
		"Sequence": models.IntegerField{Default: models.DefaultValue(1),
			Help: "Assigns the priority to the list of product vendor."},
		"ProductUom": models.Many2OneField{String: "Vendor Unit of Measure", RelationModel: h.ProductUom(),
			ReadOnly: true, Related: "ProductTmpl.UomPo", Help: "This comes from the product form."},
		"MinQty": models.FloatField{String: "Minimal Quantity", Default: models.DefaultValue(0), Required: true,
			Help: `The minimal quantity to purchase from this vendor, expressed in the vendor Product Unit of Measure if any,
or in the default unit of measure of the product otherwise.`},
		"Price": models.FloatField{Default: models.DefaultValue(0), Digits: decimalPrecision.GetPrecision("Product Price"),
			Required: true, Help: "The price to purchase a product"},
		"Company": models.Many2OneField{RelationModel: h.Company(), Default: func(env models.Environment) interface{} {
			return h.User().NewSet(env).CurrentUser().Company()
		}, Index: true},
		"Currency": models.Many2OneField{RelationModel: h.Currency(), Default: func(env models.Environment) interface{} {
			return h.User().NewSet(env).CurrentUser().Company().Currency()
		}, Required: true},
		"DateStart": models.DateField{String: "Start Date", Help: "Start date for this vendor price"},
		"DateEnd":   models.DateField{String: "End Date", Help: "End date for this vendor price"},
		"Product": models.Many2OneField{String: "Product Variant", RelationModel: h.ProductProduct(),
			Help: "When this field is filled in, the vendor data will only apply to the variant."},
		"ProductTmpl": models.Many2OneField{String: "Product Template", RelationModel: h.ProductTemplate(),
			Index: true, OnDelete: models.Cascade},
		"Delay": models.IntegerField{String: "Delivery Lead Time", Default: models.DefaultValue(1), Required: true,
			Help: `Lead time in days between the confirmation of the purchase order and the receipt of the
products in your warehouse. Used by the scheduler for automatic computation of the purchase order planning.`},
	})

}
