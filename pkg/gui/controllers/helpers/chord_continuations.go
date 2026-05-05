package helpers

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/jesseduffield/lazygit/pkg/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/style"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/theme"
)

// A bindingInfo is either a leaf binding (binding != nil) or a group
// placeholder (isGroup == true).
type bindingInfo struct {
	key         string
	description string
	tooltip     string
	binding     *types.Binding
	style       style.TextStyle
	disabled    bool
	isGroup     bool
}

// BuildChordContinuations emits one row per matching binding's next key,
// collapsing bindings under a configured KeybindingGroups entry into a
// single group row. Group lookup tries groups[ctxName] then
// groups["global"].
//
// `groups` must already be normalized (every per-context key
// canonicalized); BuildChordContinuations does not call NormalizeGroupKeys
// itself. ChordMenuHelper memoizes the normalization for us — open and
// refresh hit this on every keystroke and rebuilding the maps each call
// shows up under -benchmem.
func BuildChordContinuations(
	allBindings []*types.Binding,
	prefix []gocui.Key,
	groups map[string]map[string]config.KeybindingGroupConfig,
	ctxName string,
) []bindingInfo {
	result := []bindingInfo{}
	seenGroupKeys := map[string]struct{}{}
	prefixLabel := config.LabelForKeySequence(prefix)

	for _, binding := range allBindings {
		seq := binding.Key.Sequence()
		// "Strict" prefix: a binding whose sequence equals prefix is the
		// completion of the chord, not a continuation. The chord popup
		// only ever shows the *next* key after `prefix`.
		if len(seq) <= len(prefix) {
			continue
		}
		if !gocui.KeysHavePrefix(seq, prefix) {
			continue
		}

		// Mode-mismatched bindings (AllowFurtherDispatching == true)
		// are hidden — they share keys with another mode's binding.
		if binding.GetDisabledReason != nil {
			if reason := binding.GetDisabledReason(); reason != nil && reason.AllowFurtherDispatching {
				continue
			}
		}

		if binding.IsHiddenInChordPopup() {
			continue
		}

		isDisabled := binding.IsDisabled()

		nextKey := seq[len(prefix)]
		nextKeyLabel := config.LabelForKey(nextKey)
		groupPrefix := prefixLabel + nextKeyLabel

		displayStyle := theme.OptionsFgColor
		if binding.DisplayStyle != nil {
			displayStyle = *binding.DisplayStyle
		}

		if g, ok := config.LookupGroup(groups, ctxName, groupPrefix); ok {
			if _, already := seenGroupKeys[groupPrefix]; already {
				continue
			}
			seenGroupKeys[groupPrefix] = struct{}{}
			result = append(result, bindingInfo{
				key:         nextKeyLabel,
				description: g.Name,
				style:       displayStyle,
				isGroup:     true,
			})
			continue
		}

		// Implicit group: the binding's sequence is longer than
		// prefix+1, so pressing this row's key must extend the chord
		// rather than fire the binding.
		if len(seq) > len(prefix)+1 {
			if _, already := seenGroupKeys[groupPrefix]; already {
				continue
			}
			seenGroupKeys[groupPrefix] = struct{}{}
			result = append(result, bindingInfo{
				key:         nextKeyLabel,
				description: binding.GetShortDescription(),
				style:       displayStyle,
				isGroup:     true,
			})
			continue
		}

		result = append(result, bindingInfo{
			key:         nextKeyLabel,
			description: binding.GetShortDescription(),
			tooltip:     binding.ChordPopupExtra,
			binding:     binding,
			style:       displayStyle,
			disabled:    isDisabled,
		})
	}

	return result
}
