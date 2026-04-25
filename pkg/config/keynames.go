package config

import (
	"log"
	"strings"
	"unicode/utf8"

	"github.com/gdamore/tcell/v3"
	"github.com/jesseduffield/lazygit/pkg/gocui"
	"github.com/samber/lo"
)

// NOTE: if you make changes to this table, be sure to update
// docs/keybindings/Custom_Keybindings.md as well

var labelByKey = map[gocui.KeyName]string{
	gocui.KeyF1:          "f1",
	gocui.KeyF2:          "f2",
	gocui.KeyF3:          "f3",
	gocui.KeyF4:          "f4",
	gocui.KeyF5:          "f5",
	gocui.KeyF6:          "f6",
	gocui.KeyF7:          "f7",
	gocui.KeyF8:          "f8",
	gocui.KeyF9:          "f9",
	gocui.KeyF10:         "f10",
	gocui.KeyF11:         "f11",
	gocui.KeyF12:         "f12",
	gocui.KeyInsert:      "insert",
	gocui.KeyDelete:      "delete",
	gocui.KeyHome:        "home",
	gocui.KeyEnd:         "end",
	gocui.KeyPgup:        "pgup",
	gocui.KeyPgdn:        "pgdown",
	gocui.KeyArrowUp:     "up",
	gocui.KeyArrowDown:   "down",
	gocui.KeyArrowLeft:   "left",
	gocui.KeyArrowRight:  "right",
	gocui.KeyTab:         "tab",
	gocui.KeyBacktab:     "backtab",
	gocui.KeyEnter:       "enter",
	gocui.KeyEsc:         "esc",
	gocui.KeyBackspace:   "backspace",
	gocui.MouseWheelUp:   "mouse wheel up",
	gocui.MouseWheelDown: "mouse wheel down",
}

var keyByLabel = lo.Invert(labelByKey)

func LabelForKey(key gocui.Key) string {
	if !key.IsSet() {
		return ""
	}

	label := ""
	if key.Mod()&gocui.ModCtrl != 0 {
		label += "ctrl+"
	}
	if key.Mod()&gocui.ModAlt != 0 {
		label += "alt+"
	}
	if key.Mod()&gocui.ModShift != 0 {
		label += "shift+"
	}
	if key.Mod()&gocui.ModMeta != 0 {
		label += "meta+"
	}

	if key.KeyName() == gocui.KeyName(tcell.KeyRune) {
		if key.Str() == " " {
			label += "space"
		} else if key.Str() == "-" && key.Mod() != gocui.ModNone {
			label += "minus"
		} else if key.Str() == "+" && key.Mod() != gocui.ModNone {
			label += "plus"
		} else {
			label += key.Str()
		}
	} else {
		value, ok := labelByKey[key.KeyName()]
		if ok {
			label += value
		} else {
			label += "unknown"
		}
	}

	if utf8.RuneCountInString(label) > 1 {
		label = "<" + label + ">"
	}

	return label
}

// KeyFromLabel accepts either a single token ("p", "<c-g>") or a
// sequence of tokens ("bp", "<c-g>p", "abc"). Sequences return a Key
// whose Rest() holds the remaining tokens.
func KeyFromLabel(label string) (gocui.Key, bool) {
	if label == "" || label == "<disabled>" {
		return gocui.Key{}, true
	}

	tokens, ok := tokenizeSequence(label)
	if !ok {
		// Tokenization fails on a literal "<" (the default for
		// Universal.GotoTop). Fall back to single-key parsing.
		if k, ok := singleKeyFromLabel(label); ok {
			return k, true
		}
		return gocui.Key{}, false
	}

	// Backward compat: an unwrapped multi-rune label that names a known
	// key ("f1", "space", "enter") parses as a single key, not a chord.
	// Arbitrary rune sequences ("abc") fall through because
	// singleKeyFromLabel rejects them.
	if len(tokens) > 1 {
		if k, ok := singleKeyFromLabel(label); ok {
			return k, true
		}
	}

	if len(tokens) > 1 {
		for _, t := range tokens {
			if t == "<disabled>" {
				return gocui.Key{}, false
			}
		}
	}

	keys := make([]gocui.Key, 0, len(tokens))
	for i, tok := range tokens {
		k, ok := singleKeyFromLabel(tok)
		if !ok {
			return gocui.Key{}, false
		}
		// Esc is reserved for chord cancel.
		if i > 0 && k.KeyName() == gocui.KeyName(tcell.KeyEscape) {
			return gocui.Key{}, false
		}
		keys = append(keys, k)
	}

	if len(keys) == 1 {
		return keys[0], true
	}
	return keys[0].WithRest(keys[1:]), true
}

// A token is either a single rune or a <...>-bracketed group.
func tokenizeSequence(label string) ([]string, bool) {
	var tokens []string
	i := 0
	for i < len(label) {
		if label[i] == '<' {
			end := strings.Index(label[i:], ">")
			if end == -1 {
				return nil, false
			}
			tokens = append(tokens, label[i:i+end+1])
			i += end + 1
			continue
		}
		_, size := utf8.DecodeRuneInString(label[i:])
		tokens = append(tokens, label[i:i+size])
		i += size
	}
	if len(tokens) == 0 {
		return nil, false
	}
	return tokens, true
}

func singleKeyFromLabel(label string) (gocui.Key, bool) {
	if strings.HasPrefix(label, "<") && strings.HasSuffix(label, ">") {
		label = label[1 : len(label)-1]
	}

	mod := gocui.ModNone
	for {
		// A bare "-" or "+" with any (or no) modifiers is a literal rune
		// key; this also covers lenient forms like `<c-->` and `<c++>`,
		// neither of which we emit (we use `<c-minus>` and `<c-+>`).
		if label == "-" || label == "+" {
			return gocui.NewKeyStrMod(label, mod), true
		}

		sepIdx := strings.IndexAny(label, "-+")
		if sepIdx == -1 {
			break
		}
		modStr, remainder := label[:sepIdx], label[sepIdx+1:]

		label = remainder

		switch modStr {
		case "s", "shift":
			if (mod & gocui.ModShift) != 0 {
				return gocui.Key{}, false
			}
			mod |= gocui.ModShift
		case "c", "ctrl":
			if (mod & gocui.ModCtrl) != 0 {
				return gocui.Key{}, false
			}
			mod |= gocui.ModCtrl
		case "a", "alt":
			if (mod & gocui.ModAlt) != 0 {
				return gocui.Key{}, false
			}
			mod |= gocui.ModAlt
		case "m", "meta":
			if (mod & gocui.ModMeta) != 0 {
				return gocui.Key{}, false
			}
			mod |= gocui.ModMeta
		default:
			return gocui.Key{}, false
		}
	}

	if label == "space" {
		return gocui.NewKeyStrMod(" ", mod), true
	}

	if label == "minus" {
		if mod == gocui.ModShift {
			return gocui.Key{}, false
		}
		return gocui.NewKeyStrMod("-", mod), true
	}

	if label == "plus" {
		if mod == gocui.ModShift {
			return gocui.Key{}, false
		}
		return gocui.NewKeyStrMod("+", mod), true
	}

	if keyName, ok := keyByLabel[label]; ok {
		return gocui.NewKey(keyName, "", mod), true
	}

	runeCount := utf8.RuneCountInString(label)
	if runeCount != 1 {
		return gocui.Key{}, false
	}

	// Shift on a bare rune is invalid: terminals fold shift into the rune
	// itself (shift+a arrives as "A"), so the binding could never fire.
	// Space is exempt and handled above; combined with other modifiers,
	// shift is fine because the terminal can't fold it into the rune then.
	if mod == gocui.ModShift {
		return gocui.Key{}, false
	}

	// An ASCII uppercase letter with any modifier is invalid. Ctrl+letter
	// events always arrive with a lowercase rune — control codes have no
	// case distinction (the terminal sends the same byte for ctrl+a and
	// ctrl+A), and CSI-u protocols report the unshifted codepoint with
	// shift as a separate modifier (alt+shift+a → rune='a' mod=Alt|Shift).
	// Users should write <c-s-a> rather than <c-A>.
	if mod != gocui.ModNone && len(label) == 1 && label[0] >= 'A' && label[0] <= 'Z' {
		return gocui.Key{}, false
	}

	return gocui.NewKeyStrMod(label, mod), true
}

// LabelForKeySequence produces a canonical chord label like "<b><p>"
// from a key slice. Inverse of tokenizeSequence.
func LabelForKeySequence(keys []gocui.Key) string {
	var b strings.Builder
	for _, k := range keys {
		b.WriteString(LabelForKey(k))
	}
	return b.String()
}

// LabelForBindingKey renders the full label of a binding key, including
// any chord tail. Use this for UI display; LabelForKey returns only the
// head key.
func LabelForBindingKey(key gocui.Key) string {
	return LabelForKeySequence(key.Sequence())
}

func isValidKeybindingKey(key string) bool {
	_, ok := KeyFromLabel(key)
	return ok
}

func GetValidatedKeyBindingKey(label string) gocui.Key {
	key, ok := KeyFromLabel(label)
	if !ok {
		log.Fatalf("Unrecognized key %s, this should have been caught by user config validation", label)
	}

	return key
}
