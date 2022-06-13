// Copyright 2022 TrueLevel SA
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// SPDX-License-Identifier: MPL-2.0

package pages

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"go.swapbox.cash/nd300cli/model"
	"go.swapbox.cash/nd300cli/ui"
)

type customPortModel struct {
	err       error
	textInput textinput.Model
}

func newCustomPortModel(placeholder string) *customPortModel {
	ti := textinput.New()
	ti.Placeholder = placeholder
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	return &customPortModel{textInput: ti, err: nil}
}

func (cpm *customPortModel) Init() tea.Cmd {
	return textinput.Blink
}

func (cpm *customPortModel) Update(m *model.Model, tmsg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := tmsg.(type) {
	case tea.KeyMsg:
		switch msg.Type { // nolint: exhaustive // We don't care about other keys.
		case tea.KeyEnter:
			m.Port = cpm.textInput.Value()
			m.Page, cmd = NewMachinenoPage(MachinePortMin, MachinePortMax, 0)

			return m, cmd
		case tea.KeyEsc:
			m.Page, cmd = NewPortsPage()
			return m, cmd
		case tea.KeyCtrlC:
			cmd = tea.Quit
			return m, cmd
		}

	case error:
		cpm.err = msg
		return m, nil
	}

	cpm.textInput, cmd = cpm.textInput.Update(tmsg)

	return m, cmd
}

func (cpm *customPortModel) View() string {
	view := ui.Title("Enter the serial port") +
		"\n" + cpm.textInput.View() + "\n\n" +
		ui.Subtle("esc to go back")
	if cpm.err != nil {
		view += "\n" + cpm.err.Error() + "\n"
	}

	return view
}
