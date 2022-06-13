// Copyright 2022 TrueLevel SA
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"go.swapbox.cash/nd300cli/pages"
)

func main() {
	var (
		port      string
		machineno int
	)

	flag.StringVar(&port, "p", "", "The machine serial port (optional)")
	flag.IntVar(&machineno, "m", -1, "The machine number (must be a number between 0 and 7 included, -1 to ignore.)")
	flag.Parse()

	model, cmd := pages.Start(port, machineno)
	p := tea.NewProgram(model, tea.WithAltScreen())

	if cmd != nil {
		go func() { p.Send(cmd()) }()
	}

	if err := p.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "unexpected error: %s\n", err)
		os.Exit(1)
	}
}
