package gui

import (
	"fmt"
	"strings"

	"github.com/jesseduffield/generics/set"
	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/jesseduffield/lazygit/pkg/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/controllers/helpers"
	"github.com/jesseduffield/lazygit/pkg/gui/style"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/theme"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/samber/lo"
)

type OptionsMapMgr struct {
	c *helpers.HelperCommon
}

func (gui *Gui) renderContextOptionsMap() {
	// In demos, we render our own content to this view
	if gui.integrationTest != nil && gui.integrationTest.IsDemo() {
		return
	}
	mgr := OptionsMapMgr{c: gui.c}
	mgr.renderContextOptionsMap()
}

// Render the options available for the current context at the bottom of the screen
// STYLE GUIDE: we use the default options fg color for most keybindings. We can
// only use a different color if we're in a specific mode where the user is likely
// to want to press that key. For example, when in cherry-picking mode, we
// want to prominently show the keybinding for pasting commits.
func (self *OptionsMapMgr) renderContextOptionsMap() {
	currentContext := self.c.Context().Current()

	currentContextBindings := currentContext.GetKeybindings(self.c.KeybindingsOpts())
	globalBindings := self.c.Contexts().Global.GetKeybindings(self.c.KeybindingsOpts())

	currentContextKeys := set.NewFromSlice(
		lo.Map(currentContextBindings, func(binding *types.Binding, _ int) gocui.Key {
			return binding.Key
		}))

	allBindings := append(currentContextBindings, lo.Filter(globalBindings, func(b *types.Binding, _ int) bool {
		return !currentContextKeys.Includes(b.Key)
	})...)

	// While a chord prefix is pending, the footer is dedicated to showing
	// the available continuations rather than the regular keybindings.
	if prefix := self.c.State().GetRepoState().GetPendingChord(); len(prefix) > 0 {
		self.renderOptions(self.formatBindingInfos(self.chordContinuationBindings(allBindings, prefix)))
		return
	}

	bindingsToDisplay := lo.Filter(allBindings, func(binding *types.Binding, _ int) bool {
		return binding.DisplayOnScreen && !binding.IsDisabled()
	})

	optionsMap := lo.Map(bindingsToDisplay, func(binding *types.Binding, _ int) bindingInfo {
		displayStyle := theme.OptionsFgColor
		if binding.DisplayStyle != nil {
			displayStyle = *binding.DisplayStyle
		}

		return bindingInfo{
			key:         config.LabelForKey(binding.Key),
			description: binding.GetShortDescription(),
			style:       displayStyle,
		}
	})

	// Mode-specific local keybindings
	if currentContext.GetKey() == context.LOCAL_COMMITS_CONTEXT_KEY {
		if self.c.Modes().CherryPicking.Active() {
			optionsMap = utils.Prepend(optionsMap, bindingInfo{
				key:         self.c.KeybindingsOpts().Config.Commits.PasteCommits,
				description: self.c.Tr.PasteCommits,
				style:       style.FgCyan,
			})
		}

		if self.c.Model().BisectInfo.Started() {
			optionsMap = utils.Prepend(optionsMap, bindingInfo{
				key:         self.c.KeybindingsOpts().Config.Commits.ViewBisectOptions,
				description: self.c.Tr.ViewBisectOptions,
				style:       style.FgGreen,
			})
		}
	}

	// Mode-specific global keybindings
	if state := self.c.Model().WorkingTreeStateAtLastCommitRefresh; state.Any() {
		optionsMap = utils.Prepend(optionsMap, bindingInfo{
			key:         self.c.KeybindingsOpts().Config.Universal.CreateRebaseOptionsMenu,
			description: state.OptionsMapTitle(self.c.Tr),
			style:       style.FgYellow,
		})
	}

	if self.c.Git().Patch.PatchBuilder.Active() {
		optionsMap = utils.Prepend(optionsMap, bindingInfo{
			key:         self.c.KeybindingsOpts().Config.Universal.CreatePatchOptionsMenu,
			description: self.c.Tr.ViewPatchOptions,
			style:       style.FgYellow,
		})
	}

	self.renderOptions(self.formatBindingInfos(optionsMap))
}

func (self *OptionsMapMgr) formatBindingInfos(bindingInfos []bindingInfo) string {
	width := self.c.Views().Options.InnerWidth() - 2 // -2 for some padding
	var builder strings.Builder
	ellipsis := "…"
	separator := " | "

	length := 0

	for i, info := range bindingInfos {
		plainText := fmt.Sprintf("%s: %s", info.description, info.key)

		// Check if adding the next formatted string exceeds the available width
		textLen := utils.StringWidth(plainText)
		if i > 0 && length+len(separator)+textLen > width {
			builder.WriteString(theme.OptionsFgColor.Sprint(separator + ellipsis))
			break
		}

		formatted := info.style.Sprintf(plainText)

		if i > 0 {
			builder.WriteString(theme.OptionsFgColor.Sprint(separator))
			length += len(separator)
		}
		builder.WriteString(formatted)
		length += textLen
	}

	return builder.String()
}

func (self *OptionsMapMgr) renderOptions(options string) {
	self.c.SetViewContent(self.c.Views().Options, options)
}

type bindingInfo struct {
	key         string
	description string
	style       style.TextStyle
}

// buildChordContinuations is the pure logic behind chordContinuationBindings:
// given the pending prefix, all currently-eligible bindings, and the configured
// groups, return the rows to show in the footer.
//
// Behavior:
//   - For each binding under the prefix, examine prefix+nextKey.
//   - If prefix+nextKey is a configured group, emit one row using group.Name
//     (deduplicated across all bindings that share the same group).
//   - Otherwise, emit one row per binding using its short description.
//   - Always append the <esc>: cancel row last.
//
// Style note: when multiple bindings under one group differ in DisplayStyle,
// the collapsed group row uses the style of whichever leaf is iterated first.
// Groups are typically homogeneous in style so this is rarely visible.
func buildChordContinuations(
	allBindings []*types.Binding,
	prefix []gocui.Key,
	groups map[string]config.KeybindingGroupConfig,
) []bindingInfo {
	// Normalize group keys to canonical form so that "<b><t>" and "bt" both
	// resolve to the same entry (LabelForKeySequence always produces the
	// canonical form without redundant angle brackets for single chars).
	normalizedGroups := make(map[string]config.KeybindingGroupConfig, len(groups))
	for label, g := range groups {
		if k, ok := config.KeyFromLabel(label); ok {
			normalizedGroups[config.LabelForKeySequence(k.Sequence())] = g
		} else {
			normalizedGroups[label] = g
		}
	}

	result := []bindingInfo{}
	seenGroupKeys := map[string]struct{}{}

	for _, binding := range allBindings {
		if binding.IsDisabled() {
			continue
		}
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

		nextKey := seq[len(prefix)]
		groupPrefix := config.LabelForKeySequence(append(append([]gocui.Key{}, prefix...), nextKey))

		displayStyle := theme.OptionsFgColor
		if binding.DisplayStyle != nil {
			displayStyle = *binding.DisplayStyle
		}

		if g, ok := normalizedGroups[groupPrefix]; ok {
			if _, already := seenGroupKeys[groupPrefix]; already {
				continue
			}
			seenGroupKeys[groupPrefix] = struct{}{}
			result = append(result, bindingInfo{
				key:         config.LabelForKey(nextKey),
				description: g.Name,
				style:       displayStyle,
			})
			continue
		}

		result = append(result, bindingInfo{
			key:         config.LabelForKey(nextKey),
			description: binding.GetShortDescription(),
			style:       displayStyle,
		})
	}

	result = append(result, bindingInfo{
		key:         "<esc>",
		description: "cancel",
		style:       theme.OptionsFgColor,
	})
	return result
}

// chordContinuationBindings reads keybindingGroups from config and delegates
// the collapse logic to buildChordContinuations.
func (self *OptionsMapMgr) chordContinuationBindings(allBindings []*types.Binding, prefix []gocui.Key) []bindingInfo {
	groups := self.c.UserConfig().KeybindingGroups
	return buildChordContinuations(allBindings, prefix, groups)
}
