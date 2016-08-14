// Copyright 2016 NDP Systèmes. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package defs

import (
	"time"

	"github.com/npiganeau/yep/yep/models"
)

type ResPartner struct {
	ID   int64
	Name string
	Date time.Time `yep:"type(date)"`
	//Title            *PartnerTitle
	Parent   *ResPartner
	Children []*ResPartner `yep:"type(one2many)"`
	Ref      string
	Lang     string
	TZ       string
	TzOffset string
	User     *ResUsers
	VAT      string
	//Banks            []*PartnerBank
	Website string
	Comment string
	//Categories       []*PartnerCategory
	CreditLimit float64
	EAN13       string
	Active      bool
	Customer    bool
	Supplier    bool
	Employee    bool
	Function    string
	Type        string
	Street      string
	Street2     string
	Zip         string
	City        string
	//State            *CountryState
	//Country          *Country
	Email            string
	Phone            string
	Fax              string
	Mobile           string
	Birthdate        models.Date
	IsCompany        bool
	UseParentAddress bool
	//Image            image.Image
	//Company          *Company
	//Color            color.Color
	//Users []*ResUsers `orm:"reverse(many)"`

	//'has_image': fields.function(_has_image, type="boolean"),
	//'company_id': fields.many2one('res.company', 'Company', select=1),
	//'color': fields.integer('Color Index'),
	//'user_ids': fields.one2many('res.users', 'partner_id', 'Users'),
	//'contact_address': fields.function(_address_display,  type='char', string='Complete Address'),
	//
	//# technical field used for managing commercial fields
	//'commercial_partner_id': fields.function(_commercial_partner_id, type='many2one', relation='res.partner', string='Commercial Entity', store=_commercial_partner_store_triggers)

}

func initPartner() {
	models.CreateModel("ResPartner")
	models.ExtendModel("ResPartner", new(ResPartner))
}
