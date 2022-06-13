// Copyright 2022 TrueLevel SA
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// SPDX-License-Identifier: MPL-2.0

package pages

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"

	"go.swapbox.cash/nd300cli/model"
	"go.swapbox.cash/nd300cli/ui"
)

const (
	MachinePortMin = 0
	MachinePortMax = 7
)

type machinenoPage struct {
	min, max, selected byte
}

func NewMachinenoPage(min, max, selected byte) (*machinenoPage, tea.Cmd) {
	page := &machinenoPage{min: min, max: max, selected: selected}

	return page, page.Init()
}

func (mp *machinenoPage) Init() tea.Cmd {
	return nil
}

func (mp *machinenoPage) Update(m *model.Model, tmsg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := tmsg.(tea.KeyMsg); ok {
		switch msg.String() {
		case keyDown, keyJ:
			mp.selected += 1
			if mp.selected > mp.max {
				mp.selected = 0
			}
		case keyUp, keyK:
			mp.selected -= 1
			if mp.selected < mp.min {
				mp.selected = mp.max
			}
		case keyEnter:
			var cmd tea.Cmd

			m.Machineno = mp.selected
			m.Page, cmd = NewCmdPage()

			return m, cmd
		case keyEsc:
			var cmd tea.Cmd
			m.Page, cmd = NewPortsPage()

			return m, cmd
		case keyQ, tea.KeyCtrlC.String():
			return m, tea.Quit
		}
	}

	return m, nil
}

func (mp *machinenoPage) View() string {
	tpl := ui.Title("Select the machine number:") + "\n\n"
	tpl += "%s\n"
	tpl += "\n" + ui.Subtle("j/k, up/down: select") + ui.Dot + ui.Subtle("enter: choose") + ui.Dot + ui.Subtle("esq: back") + ui.Dot + ui.Subtle("q: quit")

	choices := ""
	for i := mp.min; i <= mp.max; i++ {
		choices += ui.Checkbox(fmt.Sprintf("%d", i), i == mp.selected) + "\n"
	}

	return fmt.Sprintf(tpl, choices)
}
