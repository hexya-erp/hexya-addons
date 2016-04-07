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

package ir

import "sync"

type ActionType string

const (
	ACTION_ACT_WINDOW ActionType = "ir.actions.act_window"
	ACTION_SERVER     ActionType = "ir.actions.server"
)

type ActionsCollection struct {
	sync.RWMutex
	actions map[string]*BaseAction
}

// NewActionCollection returns a pointer to a new
// ActionsCollection instance
func NewActionsCollection() *ActionsCollection {
	res := ActionsCollection{
		actions: make(map[string]*BaseAction),
	}
	return &res
}

// AddAction adds the given action to our ActionsCollection
func (ar *ActionsCollection) AddAction(a *BaseAction) {
	ar.Lock()
	defer ar.Unlock()
	ar.actions[a.ID] = a
}

// GetActionById returns the Action with the given id
func (ar *ActionsCollection) GetActionById(id string) *BaseAction {
	return ar.actions[id]
}

type BaseAction struct {
	ID    string
	Type  ActionType
	Name  string
	Model string
}
