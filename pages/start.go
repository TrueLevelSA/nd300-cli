// Copyright 2022 TrueLevel SA
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// SPDX-License-Identifier: MPL-2.0

package pages

import (
	tea "github.com/charmbracelet/bubbletea"

	"go.swapbox.cash/nd300cli/model"
)

func Start(port string, machineno int) (*model.Model, tea.Cmd) {
	switch {
	case port != "" && 0 <= machineno && machineno <= 7:
		page, cmd := NewCmdPage()
		m := model.New(page)
		m.Port = port
		m.Machineno = byte(machineno)

		return m, cmd
	case port != "":
		page, cmd := NewMachinenoPage(MachinePortMin, MachinePortMax, 0)
		m := model.New(page)
		m.Port = port

		return m, cmd
	case 0 <= machineno && machineno <= 7:
		page, cmd := NewPortsPage()
		m := model.New(page)
		m.Machineno = byte(machineno)

		return m, cmd
	default:
		page, cmd := NewPortsPage()
		m := model.New(page)

		return m, cmd
	}
}
