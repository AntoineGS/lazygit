// Copyright 2026 The gocui Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gocui

import "github.com/gdamore/tcell/v3"

type Key struct {
	keyName KeyName
	str     string

	mod Modifier

	// rest is the tail of a chord sequence. For single-key bindings rest is empty.
	rest []Key
}

func NewKey(keyName KeyName, str string, mod Modifier) Key {
	return Key{
		keyName: keyName,
		str:     str,
		mod:     mod,
	}
}

func NewKeyName(keyName KeyName) Key {
	return Key{
		keyName: keyName,
		str:     "",
		mod:     ModNone,
	}
}

func NewKeyRune(ch rune) Key {
	return Key{
		keyName: KeyName(tcell.KeyRune),
		str:     string(ch),
		mod:     ModNone,
	}
}

func NewKeyStrMod(str string, mod Modifier) Key {
	return Key{
		keyName: KeyName(tcell.KeyRune),
		str:     str,
		mod:     mod,
	}
}

func (k Key) KeyName() KeyName {
	return k.keyName
}

func (k Key) Str() string {
	return k.str
}

func (k Key) Mod() Modifier {
	return k.mod
}

func (k Key) IsSet() bool {
	return k.keyName != 0
}

func (k Key) Equals(otherKey Key) bool {
	return k.keyName == otherKey.keyName && k.str == otherKey.str && k.mod == otherKey.mod
}

// Rest returns the tail of a chord sequence, or nil for single-key bindings.
func (k Key) Rest() []Key {
	return k.rest
}

// HasRest reports whether this Key carries a chord tail.
func (k Key) HasRest() bool {
	return len(k.rest) > 0
}

// Sequence returns the full chord sequence starting with k itself.
// For single-key bindings it returns a one-element slice containing k
// with its rest cleared.
func (k Key) Sequence() []Key {
	head := k
	head.rest = nil
	return append([]Key{head}, k.rest...)
}

// WithRest returns a copy of k with the given rest appended (replacing any
// existing rest). Used by the parser to build chord keys.
func (k Key) WithRest(rest []Key) Key {
	k.rest = rest
	return k
}
