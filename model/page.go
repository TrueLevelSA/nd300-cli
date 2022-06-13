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

type Page interface {
	Init() tea.Cmd

	// Update is called when a message is received. Use it to inspect messages
	// and, in response, update the page, model and/or send a command.
	Update(*Model, tea.Msg) (tea.Model, tea.Cmd)

	// View renders the page of program's UI, which is just a string. The view is
	// rendered after every Update.
	View() string
}
