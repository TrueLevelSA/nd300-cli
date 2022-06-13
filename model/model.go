// Copyright 2022 TrueLevel SA
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// SPDX-License-Identifier: MPL-2.0

package model

import (
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	Page      Page
	Port      string
	Machineno byte
}

func (m *Model) Init() tea.Cmd {
	return m.Page.Init()
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m.Page.Update(m, msg)
}

func (m *Model) View() string {
	return m.Page.View()
}

func New(page Page) *Model {
	return &Model{
		Page:      page,
		Machineno: 0,
		Port:      "",
	}
}
