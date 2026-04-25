package gui

import (
	"fmt"
	"sort"
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
	gui.optionsMapMgr.renderContextOptionsMap()
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

	bindingsToDisplay := lo.Filter(allBindings, func(binding *types.Binding, _ int) bool {
		return binding.DisplayOnScreen && !binding.IsDisabled()
	})

	optionsMap := lo.Map(bindingsToDisplay, func(binding *types.Binding, _ int) bindingInfo {
		displayStyle := theme.OptionsFgColor
		if binding.DisplayStyle != nil {
			displayStyle = *binding.DisplayStyle
		}

		return bindingInfo{
			key:         config.LabelForBindingKey(binding.Key),
			description: binding.GetShortDescription(),
			style:       displayStyle,
		}
	})

	// Chord-group prefixes have no underlying *types.Binding, so they
	// come from KeybindingGroups (entries with DisplayOnScreen set).
	optionsMap = append(optionsMap, self.chordGroupOptions(string(currentContext.GetKey()))...)

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
				key:         "b",
				description: self.c.Tr.ViewBisectOptions,
				style:       style.FgGreen,
			})
		}
	}

	// Mode-specific global keybindings
	if state := self.c.Model().WorkingTreeStateAtLastCommitRefresh; state.Any() {
		optionsMap = utils.Prepend(optionsMap, bindingInfo{
			key:         "m",
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

// Iterates current context first, then "global", deduplicating by
// canonical prefix label so a context override shadows the global entry.
// Inner iteration is sorted to keep bar order stable across renders.
func (self *OptionsMapMgr) chordGroupOptions(currentContextName string) []bindingInfo {
	groups := self.c.UserConfig().KeybindingGroups
	if len(groups) == 0 {
		return nil
	}

	canonicalize := func(label string) string {
		k, ok := config.KeyFromLabel(label)
		if !ok {
			return label
		}
		return config.LabelForKeySequence(k.Sequence())
	}

	seen := map[string]struct{}{}
	var result []bindingInfo

	collect := func(contextGroups map[string]config.KeybindingGroupConfig) {
		labels := make([]string, 0, len(contextGroups))
		for label := range contextGroups {
			labels = append(labels, label)
		}
		sort.Strings(labels)
		for _, label := range labels {
			group := contextGroups[label]
			if !group.DisplayOnScreen {
				continue
			}
			canonical := canonicalize(label)
			if _, dup := seen[canonical]; dup {
				continue
			}
			seen[canonical] = struct{}{}

			description := group.ShortName
			if description == "" {
				description = group.Name
			}
			result = append(result, bindingInfo{
				key:         canonical,
				description: description,
				style:       theme.OptionsFgColor,
			})
		}
	}

	collect(groups[currentContextName])
	collect(groups["global"])
	return result
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
