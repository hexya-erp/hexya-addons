/*   Copyright (C) 2008-2016 by Nicolas Piganeau
 *   (See AUTHORS file)
 *
 *   This program is free software; you can redistribute it and/or modify
 *   it under the terms of the GNU General Public License as published by
 *   the Free Software Foundation; either version 2 of the License, or
 *   (at your option) any later version.
 *
 *   This program is distributed in the hope that it will be useful,
 *   but WITHOUT ANY WARRANTY; without even the implied warranty of
 *   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *   GNU General Public License for more details.
 *
 *   You should have received a copy of the GNU General Public License
 *   along with this program; if not, write to the
 *   Free Software Foundation, Inc.,
 *   59 Temple Place - Suite 330, Boston, MA  02111-1307, USA.
 */

package res

import (
	"github.com/npiganeau/yep/yep/models"
	"time"
)

/*
import (
	"fmt"
	"github.com/jinzhu/gorm"
	"image"
	"image/color"
	"time"
)

*/
type ResPartner struct {
	ID   int64
	Name string
	Date time.Time `orm:"null"`
	//Title            *PartnerTitle
	Parent   *ResPartner   `orm:"rel(fk);null"`
	Children []*ResPartner `orm:"reverse(many);null"`
	Ref      string
	Lang     string
	TZ       string
	TZOffset string
	User     *ResUsers `orm:"reverse(one);null"`
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
	Birthdate        time.Time `orm:"null"`
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
