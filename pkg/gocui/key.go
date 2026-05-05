// Copyright 2026 The gocui Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gocui

import "github.com/gdamore/tcell/v3"

type Key struct {
	keyName KeyName
	str     string

	mod Modifier

	// Stored as a pointer so Key stays comparable (slices aren't); two
	// independently-parsed chord keys therefore have distinct identities.
	rest *[]Key
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

func (k Key) Rest() []Key {
	if k.rest == nil {
		return nil
	}
	return *k.rest
}

func (k Key) HasRest() bool {
	return k.rest != nil && len(*k.rest) > 0
}

// Sequence returns the chord sequence starting with k. Returned
// elements always have empty rest; callers must iterate the slice to
// see the chord shape, not call HasRest on individual elements.
func (k Key) Sequence() []Key {
	head := k
	head.rest = nil
	if k.rest == nil {
		return []Key{head}
	}
	return append([]Key{head}, (*k.rest)...)
}

func (k Key) WithRest(rest []Key) Key {
	cp := append([]Key(nil), rest...)
	k.rest = &cp
	return k
}

// KeysHavePrefix reports whether seq begins with prefix. Returns true
// when the two are equal length and equal element-wise; callers that
// need to distinguish "equals" from "proper prefix" should compare
// lengths separately.
//
// This is the canonical chord-prefix-match function — chord_continuations
// and the validator share it so they can't drift out of sync.
func KeysHavePrefix(seq, prefix []Key) bool {
	if len(prefix) > len(seq) {
		return false
	}
	for i := range prefix {
		if !seq[i].Equals(prefix[i]) {
			return false
		}
	}
	return true
}
