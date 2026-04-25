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
func BuildChordContinuations(
	allBindings []*types.Binding,
	prefix []gocui.Key,
	groups map[string]map[string]config.KeybindingGroupConfig,
	ctxName string,
) []bindingInfo {
	// Canonicalize so "<b><t>" and "bt" both resolve to the same label.
	canonicalize := func(label string) string {
		k, ok := config.KeyFromLabel(label)
		if !ok {
			return label
		}
		return config.LabelForKeySequence(k.Sequence())
	}
	normalize := func(in map[string]config.KeybindingGroupConfig) map[string]config.KeybindingGroupConfig {
		out := make(map[string]config.KeybindingGroupConfig, len(in))
		for label, g := range in {
			out[canonicalize(label)] = g
		}
		return out
	}
	contextual := normalize(groups[ctxName])
	global := normalize(groups["global"])

	lookupGroup := func(label string) (config.KeybindingGroupConfig, bool) {
		if g, ok := contextual[label]; ok && g.Name != "" {
			return g, true
		}
		if g, ok := global[label]; ok && g.Name != "" {
			return g, true
		}
		return config.KeybindingGroupConfig{}, false
	}

	result := []bindingInfo{}
	seenGroupKeys := map[string]struct{}{}

	for _, binding := range allBindings {
		seq := binding.Key.Sequence()
		if len(seq) <= len(prefix) {
			continue
		}
		matches := true
		for i, k := range prefix {
			if !seq[i].Equals(k) {
				matches = false
				break
			}
		}
		if !matches {
			continue
		}

		// Mode-mismatched bindings (AllowFurtherDispatching == true)
		// are hidden — they share keys with another mode's binding.
		if binding.GetDisabledReason != nil {
			if reason := binding.GetDisabledReason(); reason != nil && reason.AllowFurtherDispatching {
				continue
			}
		}
		isDisabled := binding.IsDisabled()

		nextKey := seq[len(prefix)]
		groupPrefix := config.LabelForKeySequence(append(append([]gocui.Key{}, prefix...), nextKey))

		displayStyle := theme.OptionsFgColor
		if binding.DisplayStyle != nil {
			displayStyle = *binding.DisplayStyle
		}

		if g, ok := lookupGroup(groupPrefix); ok {
			if _, already := seenGroupKeys[groupPrefix]; already {
				continue
			}
			seenGroupKeys[groupPrefix] = struct{}{}
			result = append(result, bindingInfo{
				key:         config.LabelForKey(nextKey),
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
				key:         config.LabelForKey(nextKey),
				description: binding.GetShortDescription(),
				style:       displayStyle,
				isGroup:     true,
			})
			continue
		}

		result = append(result, bindingInfo{
			key:         config.LabelForKey(nextKey),
			description: binding.GetShortDescription(),
			tooltip:     binding.ChordPopupExtra,
			binding:     binding,
			style:       displayStyle,
			disabled:    isDisabled,
		})
	}

	return result
}
