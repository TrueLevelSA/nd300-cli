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

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"go.bug.st/serial"

	"go.swapbox.cash/nd300cli/model"
	"go.swapbox.cash/nd300cli/ui"
)

var help = "\n" + ui.Subtle("j/k, up/down: select") + ui.Dot +
	ui.Subtle("enter: choose") + ui.Dot +
	ui.Subtle("e: edit") + ui.Dot +
	ui.Subtle("r: refresh") + ui.Dot +
	ui.Subtle("q, esc: quit")

type portsPage struct {
	err      error
	loading  spinner.Model
	ports    []string
	selected int
}

func NewPortsPage() (*portsPage, tea.Cmd) {
	loading := spinner.New()
	loading.Spinner = spinner.MiniDot
	loading.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#e32fc5"))

	page := &portsPage{
		loading:  loading,
		selected: 0,
		ports:    nil,
		err:      nil,
	}

	return page, page.Init()
}

func (pp *portsPage) Init() tea.Cmd {
	return tea.Batch(pp.getPorts, pp.loading.Tick)
}

func (pp *portsPage) Update(m *model.Model, msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case keyUp, keyJ:
			pp.selected -= 1
			if pp.selected < 0 {
				pp.selected = len(pp.ports) - 1
			}

		case keyDown, keyK:
			pp.selected += 1
			if pp.selected >= len(pp.ports) {
				pp.selected = 0
			}

		case keyEdit:
			port := ""
			if pp.selected < len(pp.ports) {
				port = pp.ports[pp.selected]
			}

			m.Page = newCustomPortModel(port)

			return m, m.Page.Init()

		case keyReload:
			pp.ports = nil
			return m, tea.Batch(pp.getPorts, pp.loading.Tick)

		case keyEnter:
			var cmd tea.Cmd

			m.Port = pp.ports[pp.selected]
			m.Page, cmd = NewMachinenoPage(MachinePortMin, MachinePortMax, m.Machineno)

			return m, cmd
		case keyQ, keyEsc, tea.KeyCtrlC.String():
			return m, tea.Quit
		}

	case spinner.TickMsg:
		if pp.ports != nil { // Stop spinner if ports are loaded.
			break
		}

		var cmd tea.Cmd
		pp.loading, cmd = pp.loading.Update(msg)

		return m, cmd

	case []string:
		pp.ports = msg

	case error:
		pp.err = msg
	}

	return m, nil
}

func (pp *portsPage) View() string {
	switch {
	case pp.err != nil:
		return fmt.Sprintf("Error: %v\n\n", pp.err) + help
	case pp.ports == nil:
		return pp.loading.View() + " Loading serial ports..." + help
	case len(pp.ports) < 1:
		return "No serial port found.Make sure the machine is connected.\n  - Press 'r' to refresh the list of ports\n  - Press 'e' to enter a port manually." + help
	default:
		view := ui.Title("Select the machine device:") + "\n"
		for _, port := range pp.ports {
			view += ui.Checkbox(port, port == pp.ports[pp.selected]) + "\n"
		}

		view += help

		return view
	}
}

func (pp *portsPage) getPorts() tea.Msg {
	ports, err := serial.GetPortsList()
	if err != nil {
		return fmt.Errorf("failed to get ports: %w", err)
	}

	return ports
}
