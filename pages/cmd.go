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
	"strings"
	"sync"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"go.swapbox.cash/nd300"
	"go.swapbox.cash/nd300cli/model"
	"go.swapbox.cash/nd300cli/ui"
)

const (
	historyDefaultSize = 128
	serialTimeout      = 5 * time.Second
	logBufferSize      = 32
)

var availableCmds = []nd300.Cmd{
	nd300.RequestMachineStatus,
	nd300.SingleMachinePayout,
	nd300.ResetDispenser,
}

type (
	connect struct{}
	stop    struct{}
)

type cmdPage struct {
	conn       *nd300.Conn
	loadingMsg string
	loading    spinner.Model
	history    []string
	selected   int
	hlock      sync.RWMutex
	notes      byte
	busy       bool
}

func NewCmdPage() (*cmdPage, tea.Cmd) {
	page := &cmdPage{
		selected:   0,
		history:    make([]string, 0, historyDefaultSize),
		loading:    spinner.New(),
		loadingMsg: "executing command...",
		busy:       false,
		conn:       nil,
		notes:      1,
	}

	page.loading.Spinner = spinner.MiniDot
	page.loading.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#e32fc5"))

	return page, page.Init()
}

func pprintSerialMsg(msg nd300.SerialMsg) string {
	typ := "??"

	switch msg.Type {
	case nd300.RX:
		typ = "←"
	case nd300.TX:
		typ = "→"
	}

	return fmt.Sprintf("%s % x (%s)", typ, msg.Data, msg.Data.CmdOrStatus())
}

func (c *cmdPage) Init() tea.Cmd {
	return func() tea.Msg { return connect{} }
}

func (c *cmdPage) connect(port string, machineno byte) tea.Cmd {
	c.conn = nd300.NewConn(
		port,
		nd300.SerialMode,
		nd300.STX,
		machineno,
		nd300.MsgLen,
		serialTimeout,
		logBufferSize,
	)

	return c.readSerial()
}

func (c *cmdPage) selectUp() {
	c.selected -= 1
	if c.selected < 0 {
		c.selected = len(availableCmds) - 1
	}
}

func (c *cmdPage) selectDown() {
	c.selected += 1
	if c.selected >= len(availableCmds) {
		c.selected = 0
	}
}

func (c *cmdPage) increaseNotes() {
	if availableCmds[c.selected] == nd300.SingleMachinePayout {
		c.notes++
	}
}

func (c *cmdPage) decreaseNotes() {
	if availableCmds[c.selected] == nd300.SingleMachinePayout {
		c.notes--
	}
}

func (c *cmdPage) execCmd() (cmd tea.Cmd) {
	c.AddToHistory("$ " + availableCmds[c.selected].String())

	switch availableCmds[c.selected] {
	case nd300.RequestMachineStatus:
		cmd = c.RequestMachineStatus()
	case nd300.ResetDispenser:
		cmd = c.ResetMachine()
	case nd300.SingleMachinePayout:
		cmd = c.Payout()
	case nd300.MultipleMachinesPayout:
		panic("Multiple machine payout is not supported.")
	}

	c.busy = true
	cmd = tea.Batch(spinner.Tick, cmd)

	return
}

func (c *cmdPage) Update(model *model.Model, m tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := m.(type) {
	case connect:
		cmd = c.connect(model.Port, model.Machineno)
	case spinner.TickMsg:
		if c.busy {
			c.loading, cmd = c.loading.Update(msg)
		}
	case tea.KeyMsg:
		if c.busy {
			if msg.String() == "q" || msg.String() == tea.KeyCtrlC.String() {
				return model, tea.Quit
			}

			return model, nil
		}

		switch msg.String() {
		case keyUp, keyJ:
			c.selectUp()
		case keyDown, keyK:
			c.selectDown()
		case keyRight, keySemiColon:
			c.increaseNotes()
		case keyLeft, keyL:
			c.decreaseNotes()
		// TODO use "b" for custom byte commands.
		case keyEnter:
			cmd = c.execCmd()
		case tea.KeyCtrlL.String():
			c.clearHistory()
		case keyEsc:
			model.Page, cmd = NewPortsPage()
		case keyQ, tea.KeyCtrlC.String():
			cmd = tea.Quit
		}
	case nd300.SerialMsg:
		c.AddToHistory(pprintSerialMsg(msg))
		cmd = c.readSerial()
	case string:
		c.AddToHistory(msg)
	case stop:
	}

	return model, cmd
}

func (c *cmdPage) readSerial() tea.Cmd {
	return func() tea.Msg {
		return <-c.conn.Log
	}
}

func (c *cmdPage) View() string {
	view := c.GetHistory()
	if view != "" {
		view += "\n"
	}

	if c.busy {
		view += c.loading.View() + " " + c.loadingMsg + "\n"
	} else {
		view += ui.Title("Select a command:") + "\n"
		for i, cmd := range availableCmds {
			view += ui.Checkbox(cmd.String(), i == c.selected)
			if availableCmds[i] == nd300.SingleMachinePayout {
				view += fmt.Sprintf(" %3d notes", c.notes)
				if i == c.selected {
					view += ui.Dot + ui.Subtle("l/;, left/right: decrease/increase number of notes")
				}
			}
			view += "\n"
		}
		view += ui.Subtle("j/k, up/down: select") + ui.Dot +
			ui.Subtle("enter: choose") + ui.Dot +
			// ui.Subtle("b: enter custom byte msg") + ui.Dot + // TODO Enable custom byte msg.
			ui.Subtle("Ctrl+L: clear history") + ui.Dot +
			ui.Subtle("q: quit") + ui.Dot +
			ui.Subtle("esq: back")
	}

	return view
}

func (c *cmdPage) AddToHistory(entry string) {
	c.hlock.Lock()
	c.history = append(c.history, "["+time.Now().Format("15:04:05")+"] "+entry)
	c.hlock.Unlock()
}

func (c *cmdPage) GetHistory() string {
	c.hlock.RLock()
	out := strings.Join(c.history, "\n")
	c.hlock.RUnlock()

	return out
}

func (c *cmdPage) clearHistory() {
	c.hlock.Lock()
	c.history = nil
	c.hlock.Unlock()
}

func (c *cmdPage) RequestMachineStatus() tea.Cmd {
	return func() tea.Msg {
		var msg string
		if status, count, err := c.conn.RequestStatus(); err != nil {
			msg = "status: " + status.String() + " error: " + err.Error()
		} else {
			msg = fmt.Sprintf("status: %s, notes distributed: %d", status, count)
		}

		c.busy = false

		return msg
	}
}

func (c *cmdPage) ResetMachine() tea.Cmd {
	return func() tea.Msg {
		var msg string
		if err := c.conn.ResetMachine(); err != nil {
			msg = "error: " + err.Error()
		} else {
			msg = "machine reset"
		}

		c.busy = false

		return msg
	}
}

func (c *cmdPage) Payout() tea.Cmd {
	return func() tea.Msg {
		var msg string
		if count, err := c.conn.Payout(c.notes); err != nil {
			msg = fmt.Sprintf("error: %s, notes paid: %d", err.Error(), count)
		} else {
			msg = fmt.Sprintf("notes paid: %d", count)
		}

		c.busy = false

		return msg
	}
}
