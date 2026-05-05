package config

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"slices"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/constants"
	"github.com/jesseduffield/lazygit/pkg/gocui"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/samber/lo"
)

func (config *UserConfig) Validate() error {
	if err := validateEnum("gui.statusPanelView", config.Gui.StatusPanelView,
		[]string{"dashboard", "allBranchesLog"}); err != nil {
		return err
	}
	if err := validateEnum("gui.showDivergenceFromBaseBranch", config.Gui.ShowDivergenceFromBaseBranch,
		[]string{"none", "onlyArrow", "arrowAndNumber"}); err != nil {
		return err
	}
	if err := validateEnum("gui.fileTreeSortOrder", config.Gui.FileTreeSortOrder,
		[]string{"mixed", "filesFirst", "foldersFirst"}); err != nil {
		return err
	}
	if err := validateEnum("git.autoForwardBranches", config.Git.AutoForwardBranches,
		[]string{"none", "onlyMainBranches", "allBranches"}); err != nil {
		return err
	}
	if err := validateEnum("git.localBranchSortOrder", config.Git.LocalBranchSortOrder,
		[]string{"date", "recency", "alphabetical"}); err != nil {
		return err
	}
	if err := validateEnum("git.remoteBranchSortOrder", config.Git.RemoteBranchSortOrder,
		[]string{"date", "alphabetical"}); err != nil {
		return err
	}
	if err := validateEnum("git.log.order", config.Git.Log.Order,
		[]string{"date-order", "author-date-order", "topo-order", "default"}); err != nil {
		return err
	}
	if err := validateEnum("git.log.showGraph", config.Git.Log.ShowGraph,
		[]string{"always", "never", "when-maximised"}); err != nil {
		return err
	}
	if err := validateKeybindings(config.Keybinding); err != nil {
		return err
	}
	if err := validateKeybindingGroups(config.KeybindingGroups, config.Keybinding); err != nil {
		return err
	}
	if err := validateCustomCommands(config.CustomCommands); err != nil {
		return err
	}
	if err := validateSpinner(config.Gui.Spinner); err != nil {
		return err
	}
	return nil
}

func validateSpinner(spinner SpinnerConfig) error {
	if len(spinner.Frames) == 0 {
		return errors.New("gui.spinner.frames must not be empty.")
	}
	firstWidth := utils.StringWidth(spinner.Frames[0])
	if lo.SomeBy(spinner.Frames, func(frame string) bool {
		return utils.StringWidth(frame) != firstWidth
	}) {
		return errors.New("All gui.spinner.frames entries must have the same width.")
	}
	return nil
}

func validateEnum(name string, value string, allowedValues []string) error {
	if slices.Contains(allowedValues, value) {
		return nil
	}
	allowedValuesStr := strings.Join(allowedValues, ", ")
	return fmt.Errorf("Unexpected value '%s' for '%s'. Allowed values: %s", value, name, allowedValuesStr)
}

// walkKeybindingStrings invokes visit on every keybinding string leaf
// in node, threading path through struct and slice nesting. Skips
// unexported fields (Value.Interface() panics on them and they never
// carry keybinding strings) and fields tagged `legacy:"alias"` (these
// are migration inputs only — finalizeKeybindings has already rewritten
// their values into the canonical fields by the time this is called).
func walkKeybindingStrings(node any, path string, visit func(path, key string) error) error {
	value := reflect.ValueOf(node)
	switch value.Kind() {
	case reflect.Struct:
		for _, f := range reflect.VisibleFields(reflect.TypeOf(node)) {
			if !f.IsExported() || f.Tag.Get("legacy") == "alias" {
				continue
			}
			childPath := f.Name
			if path != "" {
				childPath = path + "." + f.Name
			}
			if err := walkKeybindingStrings(value.FieldByName(f.Name).Interface(), childPath, visit); err != nil {
				return err
			}
		}
	case reflect.Slice:
		for i := range value.Len() {
			if err := walkKeybindingStrings(value.Index(i).Interface(), fmt.Sprintf("%s[%d]", path, i), visit); err != nil {
				return err
			}
		}
	case reflect.String:
		return visit(path, node.(string))
	default:
		log.Fatalf("walkKeybindingStrings: unexpected kind %s at %q", value.Kind(), path)
	}
	return nil
}

func validateKeybindings(keybindingConfig KeybindingConfig) error {
	err := walkKeybindingStrings(keybindingConfig, "", func(path, key string) error {
		if !isValidKeybindingKey(key) {
			return fmt.Errorf("Unrecognized key '%s' for keybinding '%s'. For permitted values see %s",
				key, path, constants.Links.Docs.CustomKeybindings)
		}
		return rejectModifierInChordTail(key, fmt.Sprintf("keybinding '%s'", path))
	})
	if err != nil {
		return err
	}

	if len(keybindingConfig.Universal.JumpToBlock) != 5 {
		return fmt.Errorf("keybinding.universal.jumpToBlock must have 5 elements; found %d.",
			len(keybindingConfig.Universal.JumpToBlock))
	}

	return nil
}

func validateCustomCommandKey(key string) error {
	if !isValidKeybindingKey(key) {
		return fmt.Errorf("Unrecognized key '%s' for custom command. For permitted values see %s",
			key, constants.Links.Docs.CustomKeybindings)
	}
	return nil
}

func validateCustomCommands(customCommands []CustomCommand) error {
	for _, customCommand := range customCommands {
		if err := validateCustomCommandKey(customCommand.Key); err != nil {
			return err
		}

		if len(customCommand.CommandMenu) > 0 {
			if len(customCommand.Context) > 0 ||
				len(customCommand.Command) > 0 ||
				len(customCommand.Prompts) > 0 ||
				len(customCommand.LoadingText) > 0 ||
				len(customCommand.Output) > 0 ||
				len(customCommand.OutputTitle) > 0 ||
				customCommand.After != nil {
				commandRef := ""
				if len(customCommand.Key) > 0 {
					commandRef = fmt.Sprintf(" with key '%s'", customCommand.Key)
				}
				return fmt.Errorf("Error with custom command%s: it is not allowed to use both commandMenu and any of the other fields except key and description.", commandRef)
			}

			if err := validateCustomCommands(customCommand.CommandMenu); err != nil {
				return err
			}
		} else {
			for _, prompt := range customCommand.Prompts {
				if err := validateCustomCommandPrompt(prompt); err != nil {
					return err
				}
			}

			if err := validateEnum("customCommand.output", customCommand.Output,
				[]string{"", "none", "terminal", "log", "logWithPty", "popup"}); err != nil {
				return err
			}
		}
	}
	return nil
}

func validateCustomCommandPrompt(prompt CustomCommandPrompt) error {
	for _, option := range prompt.Options {
		if !isValidKeybindingKey(option.Key) {
			return fmt.Errorf("Unrecognized key '%s' for custom command prompt option. For permitted values see %s",
				option.Key, constants.Links.Docs.CustomKeybindings)
		}
	}

	return nil
}

func collectAllKeybindingStrings(node any) []string {
	var out []string
	_ = walkKeybindingStrings(node, "", func(_, key string) error {
		if key != "" && key != "<disabled>" {
			out = append(out, key)
		}
		return nil
	})
	return out
}

// rejectModifierInChordTail enforces that only the head of a chord may
// carry a modifier. `<ctrl+b>p` is fine; `b<ctrl+p>` is not — the
// non-head modifier is unusual UX and easy to type by accident.
func rejectModifierInChordTail(label, contextDesc string) error {
	if label == "" || label == "<disabled>" {
		return nil
	}
	k, ok := KeyFromLabel(label)
	if !ok {
		// isValidKeybindingKey already rejected this elsewhere.
		return nil
	}
	seq := k.Sequence()
	for i := 1; i < len(seq); i++ {
		if seq[i].Mod() != gocui.ModNone {
			return fmt.Errorf(
				"%s: chord %q has a modifier on key #%d; modifiers are only allowed on the first key of a chord",
				contextDesc, label, i+1)
		}
	}
	return nil
}

// chordContextToBindingFields maps a chord-group context name (matching
// gui ContextKey values) to the KeybindingConfig field names whose
// bindings register in that view.
//
// Universal is intentionally NOT included for non-global contexts:
// chord-dispatch wins over single-key bindings at runtime, so the
// chord-prefix-vs-Universal-leaf overlap (e.g. `Universal.Remove="d"`
// coexisting with files chord prefix `d`) is by design.
var chordContextToBindingFields = map[string][]string{
	"files":          {"Files"},
	"localBranches":  {"Branches"},
	"remoteBranches": {"Branches"},
	"tags":           {"Branches"},
	"commits":        {"Commits"},
	"reflogCommits":  {"Commits"},
	"subCommits":     {"Commits"},
	"commitFiles":    {"CommitFiles"},
	"submodules":     {"Submodules"},
	"stash":          {"Stash"},
	"worktrees":      {"Worktrees"},
	"main":           {"Main"},
	"commitMessage":  {"CommitMessage"},
	"status":         {"Status"},
}

// collectKeybindingsBySection returns each top-level KeybindingConfig
// field's bindings, keyed by struct field name. Legacy aliases are
// excluded by collectAllKeybindingStrings.
func collectKeybindingsBySection(keybindings KeybindingConfig) map[string][]string {
	out := map[string][]string{}
	v := reflect.ValueOf(keybindings)
	t := v.Type()
	for i := range v.NumField() {
		field := t.Field(i)
		if !field.IsExported() {
			continue
		}
		out[field.Name] = collectAllKeybindingStrings(v.Field(i).Interface())
	}
	return out
}

func validateKeybindingGroups(groups map[string]map[string]KeybindingGroupConfig, keybindings KeybindingConfig) error {
	bindingsBySection := collectKeybindingsBySection(keybindings)
	var allBindings []string
	for _, b := range bindingsBySection {
		allBindings = append(allBindings, b...)
	}

	for contextName, contextGroups := range groups {
		// Pick the binding subset to check leaf collisions against.
		// Global chord prefixes can collide with anything; per-view
		// chord prefixes can only collide with bindings registered in
		// the same view (i.e. the corresponding KeybindingConfig
		// section). Universal is excluded for non-global since
		// chord-dispatch wins at runtime.
		var collisionBindings []string
		if contextName == "global" {
			collisionBindings = allBindings
		} else if fields, known := chordContextToBindingFields[contextName]; known {
			for _, f := range fields {
				collisionBindings = append(collisionBindings, bindingsBySection[f]...)
			}
		} else {
			// Unknown context: be conservative and check against all.
			collisionBindings = allBindings
		}

		// Sort iteration to make error messages deterministic.
		canonicalSeen := map[string]string{}
		nameSeen := map[string]string{}
		sortedPrefixes := make([]string, 0, len(contextGroups))
		for p := range contextGroups {
			sortedPrefixes = append(sortedPrefixes, p)
		}
		slices.Sort(sortedPrefixes)
		for _, prefix := range sortedPrefixes {
			group := contextGroups[prefix]
			if _, ok := KeyFromLabel(prefix); !ok {
				return fmt.Errorf("Unrecognized chord prefix '%s' in keybindingGroups. For permitted values see %s",
					prefix, constants.Links.Docs.CustomKeybindings)
			}
			if err := rejectModifierInChordTail(prefix,
				fmt.Sprintf("keybindingGroups[%s][%s]", contextName, prefix)); err != nil {
				return err
			}
			if strings.TrimSpace(group.Name) == "" {
				return fmt.Errorf("keybindingGroups[%s] must have a non-empty name", prefix)
			}

			canonical := CanonicalizePrefixLabel(prefix)
			if other, dup := canonicalSeen[canonical]; dup {
				return fmt.Errorf(
					"keybindingGroups[%s][%s] resolves to the same chord prefix as keybindingGroups[%s][%s] (both canonicalize to %q); remove or rename one",
					contextName, prefix, contextName, other, canonical)
			}
			canonicalSeen[canonical] = prefix

			if other, dup := nameSeen[group.Name]; dup {
				return fmt.Errorf(
					"keybindingGroups[%s][%s] has the same name %q as keybindingGroups[%s][%s]; group names must be unique within a context so chord-prefix lookup is deterministic",
					contextName, prefix, group.Name, contextName, other)
			}
			nameSeen[group.Name] = prefix

			prefixKey, _ := KeyFromLabel(prefix)
			prefixSeq := prefixKey.Sequence()

			for _, b := range collisionBindings {
				bk, ok := KeyFromLabel(b)
				if !ok {
					continue
				}
				bseq := bk.Sequence()
				if len(bseq) == len(prefixSeq) && gocui.KeysHavePrefix(bseq, prefixSeq) {
					return fmt.Errorf("keybindingGroups[%s][%s] collides with a leaf binding using the same key sequence; a key cannot be both an action and a sub-menu", contextName, prefix)
				}
			}

			// hasChild scans all bindings: a chord tail registered via
			// a controller for view X may live in a KeybindingConfig
			// section that doesn't directly map to X (e.g. global
			// chord tails live under Universal).
			hasChild := false
			for _, b := range allBindings {
				bk, ok := KeyFromLabel(b)
				if !ok {
					continue
				}
				bseq := bk.Sequence()
				if len(bseq) > len(prefixSeq) && gocui.KeysHavePrefix(bseq, prefixSeq) {
					hasChild = true
					break
				}
			}
			if !hasChild {
				return fmt.Errorf("keybindingGroups[%s] has no bindings under it; either add a chord binding starting with %s or remove the group entry",
					prefix, prefix)
			}
		}
	}
	return nil
}
