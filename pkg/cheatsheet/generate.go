//go:generate go run generator.go

// This "script" generates files called Keybindings_{{.LANG}}.md
// in the docs-master/keybindings directory.
//
// The content of these generated files is a keybindings cheatsheet.
//
// To generate the cheatsheets, run:
//   go generate pkg/cheatsheet/generate.go

package cheatsheet

import (
	"cmp"
	"fmt"
	"log"
	"os"
	"slices"
	"strings"

	"github.com/jesseduffield/generics/maps"
	"github.com/jesseduffield/lazycore/pkg/utils"
	"github.com/jesseduffield/lazygit/pkg/app"
	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/jesseduffield/lazygit/pkg/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/i18n"
	"github.com/samber/lo"
)

type bindingSection struct {
	title    string
	bindings []*types.Binding
}

type header struct {
	// priority decides the order of the headers in the cheatsheet (lower means higher)
	priority int
	title    string
}

type headerWithBindings struct {
	header   header
	bindings []*types.Binding
}

func CommandToRun() string {
	return "go generate ./..."
}

func GetKeybindingsDir() string {
	return utils.GetLazyRootDirectory() + "/docs-master/keybindings"
}

func generateAtDir(cheatsheetDir string) {
	translationSetsByLang, err := i18n.GetTranslationSets()
	if err != nil {
		log.Fatal(err)
	}
	mConfig := config.NewDummyAppConfig()

	for lang := range translationSetsByLang {
		mConfig.GetUserConfig().Gui.Language = lang
		common, err := app.NewCommon(mConfig)
		if err != nil {
			log.Fatal(err)
		}
		tr, err := i18n.NewTranslationSetFromConfig(common.Log, lang)
		if err != nil {
			log.Fatal(err)
		}
		common.Tr = tr
		mApp, _ := app.NewApp(mConfig, nil, common)
		path := cheatsheetDir + "/Keybindings_" + lang + ".md"
		file, err := os.Create(path)
		if err != nil {
			panic(err)
		}

		bindings := mApp.Gui.GetCheatsheetKeybindings()
		groups := mConfig.GetUserConfig().KeybindingGroups
		bindingSections := getBindingSections(bindings, mApp.Tr, groups, defaultGroupMeta(mApp.Tr))
		content := formatSections(mApp.Tr, bindingSections)
		content = fmt.Sprintf("_This file is auto-generated. To update, make the changes in the "+
			"pkg/i18n directory and then run `%s` from the project root._\n\n%s", CommandToRun(), content)
		writeString(file, content)
	}
}

func Generate() {
	generateAtDir(GetKeybindingsDir())
}

func writeString(file *os.File, str string) {
	_, err := file.WriteString(str)
	if err != nil {
		log.Fatal(err)
	}
}

func localisedTitle(tr *i18n.TranslationSet, str string) string {
	contextTitleMap := map[string]string{
		"global":            tr.GlobalTitle,
		"navigation":        tr.NavigationTitle,
		"branches":          tr.BranchesTitle,
		"localBranches":     tr.LocalBranchesTitle,
		"files":             tr.FilesTitle,
		"status":            tr.StatusTitle,
		"submodules":        tr.SubmodulesTitle,
		"subCommits":        tr.SubCommitsTitle,
		"remoteBranches":    tr.RemoteBranchesTitle,
		"remotes":           tr.RemotesTitle,
		"reflogCommits":     tr.ReflogCommitsTitle,
		"tags":              tr.TagsTitle,
		"commitFiles":       tr.CommitFilesTitle,
		"commitMessage":     tr.CommitSummaryTitle,
		"commitDescription": tr.CommitDescriptionTitle,
		"commits":           tr.CommitsTitle,
		"confirmation":      tr.ConfirmationTitle,
		"prompt":            tr.PromptTitle,
		"information":       tr.InformationTitle,
		"main":              tr.NormalTitle,
		"patchBuilding":     tr.PatchBuildingTitle,
		"mergeConflicts":    tr.MergingTitle,
		"staging":           tr.StagingTitle,
		"menu":              tr.MenuTitle,
		"search":            tr.SearchTitle,
		"secondary":         tr.SecondaryTitle,
		"stash":             tr.StashTitle,
		"suggestions":       tr.SuggestionsCheatsheetTitle,
		"extras":            tr.ExtrasTitle,
		"worktrees":         tr.WorktreesTitle,
	}

	title, ok := contextTitleMap[str]
	if !ok {
		panic(fmt.Sprintf("title not found for %s", str))
	}

	return title
}

// groupMeta is sourced from translation strings so the cheatsheet text
// matches the pre-chord menu-opener bindings in the user's language.
// Falls back to KeybindingGroupConfig.Name when no builtin entry exists.
type groupMeta struct {
	description string
	tooltip     string
}

func getBindingSections(
	bindings []*types.Binding,
	tr *i18n.TranslationSet,
	groups map[string]map[string]config.KeybindingGroupConfig,
	meta map[string]map[string]groupMeta,
) []*bindingSection {
	excludedViews := []string{"stagingSecondary", "patchBuildingSecondary"}
	bindingsToDisplay := lo.Filter(bindings, func(binding *types.Binding, _ int) bool {
		if lo.Contains(excludedViews, binding.ViewName) {
			return false
		}

		return (binding.Description != "" || binding.Alternative != "") && binding.Key.IsSet()
	})

	bindingsByHeader := lo.GroupBy(bindingsToDisplay, func(binding *types.Binding) header {
		return getHeader(binding, tr)
	})

	bindingGroups := maps.MapToSlice(
		bindingsByHeader,
		func(header header, hBindings []*types.Binding) headerWithBindings {
			uniqBindings := lo.UniqBy(hBindings, func(binding *types.Binding) string {
				return binding.Description + config.LabelForBindingKey(binding.Key)
			})

			return headerWithBindings{
				header:   header,
				bindings: insertChordGroupHeaders(uniqBindings, groups, meta),
			}
		},
	)

	slices.SortFunc(bindingGroups, func(a, b headerWithBindings) int {
		if a.header.priority != b.header.priority {
			return cmp.Compare(b.header.priority, a.header.priority)
		}
		return strings.Compare(a.header.title, b.header.title)
	})

	return lo.Map(bindingGroups, func(hb headerWithBindings, _ int) *bindingSection {
		return &bindingSection{
			title:    hb.header.title,
			bindings: hb.bindings,
		}
	})
}

// insertChordGroupHeaders inserts synthetic header rows in front of
// each chord prefix configured in UserConfig.KeybindingGroups.
//
// Description/Tooltip resolution precedence:
//  1. KeybindingGroupConfig.Description / .Tooltip (user override)
//  2. defaultGroupMeta entry (translated strings from the pre-chord
//     menu-opener bindings)
//  3. KeybindingGroupConfig.Name (description only)
//
// If a real binding already occupies the same key, no synthetic row is
// emitted (avoids duplicate keys in the cheatsheet).
func insertChordGroupHeaders(
	bindings []*types.Binding,
	groups map[string]map[string]config.KeybindingGroupConfig,
	meta map[string]map[string]groupMeta,
) []*types.Binding {
	if len(groups) == 0 {
		return bindings
	}

	existing := map[string]struct{}{}
	for _, b := range bindings {
		existing[b.ViewName+"|"+config.LabelForBindingKey(b.Key)] = struct{}{}
	}

	// Canonicalize lookup keys the same way chord_continuations does, so
	// labels written as "<b><t>" or "bt" both resolve.
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
	normalized := make(map[string]map[string]config.KeybindingGroupConfig, len(groups))
	for ctx, byPrefix := range groups {
		normalized[ctx] = normalize(byPrefix)
	}

	lookupGroup := func(viewName, prefixLabel string) (config.KeybindingGroupConfig, bool) {
		if g, ok := normalized[viewName][prefixLabel]; ok && g.Name != "" {
			return g, true
		}
		if g, ok := normalized["global"][prefixLabel]; ok && g.Name != "" {
			return g, true
		}
		return config.KeybindingGroupConfig{}, false
	}

	lookupMeta := func(viewName, prefixLabel string) (groupMeta, bool) {
		if m, ok := meta[viewName][prefixLabel]; ok {
			return m, true
		}
		if m, ok := meta["global"][prefixLabel]; ok {
			return m, true
		}
		return groupMeta{}, false
	}

	result := make([]*types.Binding, 0, len(bindings))
	emitted := map[string]struct{}{}

	for _, b := range bindings {
		seq := b.Key.Sequence()
		if len(seq) >= 2 {
			// Walk every proper prefix length so nested groups (e.g.
			// "u" then "ug") each get their own header row.
			for prefixLen := 1; prefixLen < len(seq); prefixLen++ {
				prefixLabel := config.LabelForKeySequence(seq[:prefixLen])
				marker := b.ViewName + "|" + prefixLabel
				if _, already := emitted[marker]; already {
					continue
				}
				if _, conflict := existing[marker]; conflict {
					emitted[marker] = struct{}{}
					continue
				}
				g, ok := lookupGroup(b.ViewName, prefixLabel)
				if !ok {
					continue
				}
				emitted[marker] = struct{}{}

				description := g.Description
				tooltip := g.Tooltip
				if description == "" || tooltip == "" {
					if m, ok := lookupMeta(b.ViewName, prefixLabel); ok {
						if description == "" {
							description = m.description
						}
						if tooltip == "" {
							tooltip = m.tooltip
						}
					}
				}
				if description == "" {
					description = g.Name
				}

				result = append(result, &types.Binding{
					ViewName:    b.ViewName,
					Key:         buildPrefixKey(seq[:prefixLen]),
					Description: description,
					Tooltip:     tooltip,
					Tag:         b.Tag,
				})
			}
		}
		result = append(result, b)
	}

	return result
}

// defaultGroupMeta returns cheatsheet meta inherited from the pre-chord
// menu-opener bindings, keyed by (viewName, prefix label).
func defaultGroupMeta(tr *i18n.TranslationSet) map[string]map[string]groupMeta {
	return map[string]map[string]groupMeta{
		"global": {
			"m": {description: tr.ViewMergeRebaseOptions, tooltip: tr.ViewMergeRebaseOptionsTooltip},
		},
		"files": {
			"i":        {description: tr.Actions.IgnoreExcludeFile},
			"S":        {description: tr.ViewStashOptions, tooltip: tr.ViewStashOptionsTooltip},
			"D":        {description: tr.Reset, tooltip: tr.FileResetOptionsTooltip},
			"<ctrl+b>": {description: tr.FileFilter},
			"y":        {description: tr.CopyToClipboardMenu},
		},
		"localBranches": {
			"d": {description: tr.Delete, tooltip: tr.BranchDeleteTooltip},
			"r": {description: tr.RebaseBranch, tooltip: tr.RebaseBranchTooltip},
			"M": {description: tr.Merge, tooltip: tr.MergeBranchTooltip},
			"i": {description: tr.GitFlowOptions},
			"u": {description: tr.ViewBranchUpstreamOptions, tooltip: tr.ViewBranchUpstreamOptionsTooltip},
		},
		"remoteBranches": {
			"d": {description: tr.Delete, tooltip: tr.DeleteRemoteBranchTooltip},
		},
		"tags": {
			"d": {description: tr.Delete},
		},
		"commits": {
			"b": {description: tr.ViewBisectOptions},
			"f": {description: tr.Fixup, tooltip: tr.FixupTooltip},
			"g": {description: tr.ViewResetOptions, tooltip: tr.ResetTooltip},
		},
		"submodules": {
			"b": {description: tr.ViewBulkSubmoduleOptions},
		},
	}
}

func buildPrefixKey(seq []gocui.Key) gocui.Key {
	if len(seq) == 1 {
		return seq[0]
	}
	return seq[0].WithRest(seq[1:])
}

func getHeader(binding *types.Binding, tr *i18n.TranslationSet) header {
	if binding.Tag == "navigation" {
		return header{priority: 2, title: localisedTitle(tr, "navigation")}
	}

	if binding.ViewName == "" {
		return header{priority: 3, title: localisedTitle(tr, "global")}
	}

	return header{priority: 1, title: localisedTitle(tr, binding.ViewName)}
}

func formatSections(tr *i18n.TranslationSet, bindingSections []*bindingSection) string {
	var content strings.Builder
	content.WriteString(fmt.Sprintf("# Lazygit %s\n", tr.Keybindings))

	for _, section := range bindingSections {
		content.WriteString(formatTitle(section.title))
		content.WriteString("| Key | Action | Info |\n")
		content.WriteString("|-----|--------|-------------|\n")
		for _, binding := range section.bindings {
			content.WriteString(formatBinding(binding))
		}
	}

	return content.String()
}

func formatTitle(title string) string {
	return fmt.Sprintf("\n## %s\n\n", title)
}

func formatBinding(binding *types.Binding) string {
	action := config.LabelForBindingKey(binding.Key)
	description := binding.Description
	if binding.Alternative != "" {
		action += fmt.Sprintf(" (%s)", binding.Alternative)
	}

	// Replace newlines with <br> tags for proper markdown table formatting
	tooltip := strings.ReplaceAll(binding.Tooltip, "\n", "<br>")

	// Escape pipe characters to avoid breaking the table format
	action = strings.ReplaceAll(action, `|`, `\|`)
	description = strings.ReplaceAll(description, `|`, `\|`)
	tooltip = strings.ReplaceAll(tooltip, `|`, `\|`)

	// Use backticks for keyboard keys. Two backticks are needed with an inner space
	//  to escape a key that is itself a backtick.
	return fmt.Sprintf("| `` %s `` | %s | %s |\n", action, description, tooltip)
}
