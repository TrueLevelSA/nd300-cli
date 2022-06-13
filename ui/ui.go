// Copyright 2022 TrueLevel SA
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// SPDX-License-Identifier: MPL-2.0

package ui

import (
	"github.com/muesli/termenv"
)

var (
	Term   = termenv.EnvColorProfile()
	Subtle = MakeFgStyle("241")
	Dot    = ColorFg(" • ", "236")
)

// ColorFg colors a string's foreground with the given value.
func ColorFg(val, color string) string {
	return termenv.String(val).Foreground(Term.Color(color)).String()
}

// MakeFgStyle return a function that will colorize the foreground of a given string.
func MakeFgStyle(color string) func(string) string {
	return termenv.Style{}.Foreground(Term.Color(color)).Styled
}

// Title formats a string as a title for the application.
func Title(text string) string {
	return termenv.String(text).Underline().Foreground(termenv.ANSIBlue).String()
}

func Checkbox(label string, checked bool) string {
	if checked {
		return termenv.String("▸ " + label).Foreground(termenv.ANSIYellow).String()
	}

	return "  " + label
}
