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
		label += "c-"
	}
	if key.Mod()&gocui.ModAlt != 0 {
		label += "a-"
	}
	if key.Mod()&gocui.ModShift != 0 {
		label += "s-"
	}
	if key.Mod()&gocui.ModMeta != 0 {
		label += "m-"
	}

	if key.KeyName() == gocui.KeyName(tcell.KeyRune) {
		if key.Str() == " " {
			label += "space"
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

	if len(label) > 1 {
		label = "<" + label + ">"
	}

	return label
}

// KeyFromLabel parses a keybinding label, which may be a single token
// ("p", "<c-g>") or a sequence of tokens ("bp", "<c-g>p", "<c-g><c-p>", "abc").
// For sequences it returns a Key whose first token is the head and whose
// Rest() contains the remaining tokens.
func KeyFromLabel(label string) (gocui.Key, bool) {
	if label == "" || label == "<disabled>" {
		return gocui.Key{}, true
	}

	tokens, ok := tokenizeSequence(label)
	if !ok {
		// Tokenization can fail on inputs like a literal "<" used as a
		// single-key binding (the default for Universal.GotoTop). Fall back
		// to the single-key parser so backward compatibility is preserved.
		if k, ok := singleKeyFromLabel(label); ok {
			return k, true
		}
		return gocui.Key{}, false
	}

	// Disallow <disabled> inside a longer sequence.
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
		// Esc is reserved for chord cancel; disallow in non-first positions.
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

// tokenizeSequence splits a label into its constituent key tokens.
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
		// Consume one rune.
		_, size := utf8.DecodeRuneInString(label[i:])
		tokens = append(tokens, label[i:i+size])
		i += size
	}
	if len(tokens) == 0 {
		return nil, false
	}
	return tokens, true
}

// singleKeyFromLabel parses a single token (never a sequence) into a Key.
// This is the original KeyFromLabel logic, unchanged in semantics.
func singleKeyFromLabel(label string) (gocui.Key, bool) {
	if strings.HasPrefix(label, "<") && strings.HasSuffix(label, ">") {
		label = label[1 : len(label)-1]
	}

	if label == "-" {
		return gocui.NewKeyRune('-'), true
	}

	mod := gocui.ModNone
	for {
		modStr, remainder, ok := strings.Cut(label, "-")
		if !ok {
			break
		}

		label = remainder

		switch modStr {
		case "s":
			if (mod & gocui.ModShift) != 0 {
				return gocui.Key{}, false
			}
			mod |= gocui.ModShift
		case "c":
			if (mod & gocui.ModCtrl) != 0 {
				return gocui.Key{}, false
			}
			mod |= gocui.ModCtrl
		case "a":
			if (mod & gocui.ModAlt) != 0 {
				return gocui.Key{}, false
			}
			mod |= gocui.ModAlt
		case "m":
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

	if keyName, ok := keyByLabel[label]; ok {
		return gocui.NewKey(keyName, "", mod), true
	}

	runeCount := utf8.RuneCountInString(label)
	if runeCount != 1 {
		return gocui.Key{}, false
	}

	return gocui.NewKeyStrMod(label, mod), true
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
