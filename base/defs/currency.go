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
	"github.com/npiganeau/yep/pool"
	"github.com/npiganeau/yep/yep/models"
)

type ResCurrencyRate struct {
	Name     models.DateTime     `yep:"string(Date);required;index"`
	Rate     float64             `yep:"digits(12,6);help(The rate of the currency to the currency of rate 1)"`
	Currency pool.ResCurrencySet `yep:"readonly;type(many2one)"`
	Company  pool.ResCompanySet  `yep:"type(many2one)"`
}

type ResCurrency struct {
	Name          string                  `yep:"string(Currency);help(Currency Code [ISO 4217]);size(3);unique"`
	Symbol        string                  `yep:"string(Symbol);help(Currency sign, to be used when printing amounts);size(4)"`
	Rate          float64                 `yep:"string(Current Rate);help(The rate of the currency to the currency of rate 1);digits(12,6)"` //;compute(ComputeCurrentRate)"`
	Rates         pool.ResCurrencyRateSet `yep:"type(one2many);fk(Currency)"`
	Rounding      float64                 `yep:"string(Rounding Factor);digits(12,6)"`
	DecimalPlaces int                     `yep:"string(Decimal Places)"` //;compute(ComputeDecimalPlaces)"`
	Active        bool
	Position      string `yep:"type(selection);selection(after|After Amount,before|Before Amount);string(Symbol Position);help(Determines where the currency symbol should be placed after or before the amount.)"`
}

func initCurrency() {
	models.CreateModel("ResCurrencyRate")
	models.ExtendModel("ResCurrencyRate", new(ResCurrencyRate))
	models.CreateModel("ResCurrency")
	models.ExtendModel("ResCurrency", new(ResCurrency))
}
