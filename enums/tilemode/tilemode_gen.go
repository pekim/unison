// Code generated from "enum.go.tmpl" - DO NOT EDIT.

// Copyright (c) 2021-2024 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package tilemode

import (
	"strings"

	"github.com/richardwilkes/toolbox/i18n"
)

// Possible values.
const (
	Clamp Enum = iota
	Repeat
	Mirror
	Decal
)

// All possible values.
var All = []Enum{
	Clamp,
	Repeat,
	Mirror,
	Decal,
}

// Enum holds the type of tiling to perform.
type Enum byte

// EnsureValid ensures this is of a known value.
func (e Enum) EnsureValid() Enum {
	if e <= Decal {
		return e
	}
	return Clamp
}

// Key returns the key used in serialization.
func (e Enum) Key() string {
	switch e {
	case Clamp:
		return "clamp"
	case Repeat:
		return "repeat"
	case Mirror:
		return "mirror"
	case Decal:
		return "decal"
	default:
		return Clamp.Key()
	}
}

// String implements fmt.Stringer.
func (e Enum) String() string {
	switch e {
	case Clamp:
		return i18n.Text("Clamp")
	case Repeat:
		return i18n.Text("Repeat")
	case Mirror:
		return i18n.Text("Mirror")
	case Decal:
		return i18n.Text("Decal")
	default:
		return Clamp.String()
	}
}

// MarshalText implements the encoding.TextMarshaler interface.
func (e Enum) MarshalText() (text []byte, err error) {
	return []byte(e.Key()), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (e *Enum) UnmarshalText(text []byte) error {
	*e = Extract(string(text))
	return nil
}

// Extract the value from a string.
func Extract(str string) Enum {
	for _, e := range All {
		if strings.EqualFold(e.Key(), str) {
			return e
		}
	}
	return Clamp
}
