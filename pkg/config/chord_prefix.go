package config

import (
	"fmt"
	"reflect"
	"strings"
)

// ChordPrefixConfig is a derived view over KeybindingGroups, exposing
// each chord-prefix key as a Go-typed slot for integration tests:
//
//	t.GlobalPress(t.keys.ChordPrefix.Global.RebaseOptions)
//
// Each leaf field's `chordGroup:"<context>:<name>"` tag is matched
// against KeybindingGroups by ResolveChordPrefixes.
type ChordPrefixConfig struct {
	Global         GlobalChordPrefixes         `yaml:"-"`
	Files          FilesChordPrefixes          `yaml:"-"`
	LocalBranches  LocalBranchesChordPrefixes  `yaml:"-"`
	RemoteBranches RemoteBranchesChordPrefixes `yaml:"-"`
	Tags           TagsChordPrefixes           `yaml:"-"`
	Commits        CommitsChordPrefixes        `yaml:"-"`
	ReflogCommits  ReflogCommitsChordPrefixes  `yaml:"-"`
	SubCommits     SubCommitsChordPrefixes     `yaml:"-"`
	CommitFiles    CommitFilesChordPrefixes    `yaml:"-"`
	Submodules     SubmodulesChordPrefixes     `yaml:"-"`
}

type GlobalChordPrefixes struct {
	RebaseOptions string `chordGroup:"global:Rebase options"`
}

type FilesChordPrefixes struct {
	IgnoreOptions          string `chordGroup:"files:Ignore options"`
	StashOptions           string `chordGroup:"files:Stash options"`
	DiscardAndResetOptions string `chordGroup:"files:Discard / reset options"`
	FilterFiles            string `chordGroup:"files:Filter files"`
	CopyToClipboard        string `chordGroup:"files:Copy to clipboard"`
	ResetToUpstream        string `chordGroup:"files:Reset to upstream"`
	DiscardChanges         string `chordGroup:"files:Discard changes"`
}

type LocalBranchesChordPrefixes struct {
	Merge                 string `chordGroup:"localBranches:Merge"`
	DeleteBranch          string `chordGroup:"localBranches:Delete branch"`
	GitFlowOptions        string `chordGroup:"localBranches:Git flow options"`
	RebaseOptions         string `chordGroup:"localBranches:Rebase options"`
	BranchUpstreamOptions string `chordGroup:"localBranches:Branch upstream options"`
	ResetToRef            string `chordGroup:"localBranches:Reset to ref"`
}

type RemoteBranchesChordPrefixes struct {
	DeleteRemoteBranch string `chordGroup:"remoteBranches:Delete remote branch"`
	Merge              string `chordGroup:"remoteBranches:Merge"`
	RebaseOptions      string `chordGroup:"remoteBranches:Rebase options"`
	ResetToRef         string `chordGroup:"remoteBranches:Reset to ref"`
}

type TagsChordPrefixes struct {
	DeleteTag  string `chordGroup:"tags:Delete tag"`
	ResetToRef string `chordGroup:"tags:Reset to ref"`
}

type CommitsChordPrefixes struct {
	BisectOptions      string `chordGroup:"commits:Bisect options"`
	FixupCommitOptions string `chordGroup:"commits:Fixup commit options"`
	ResetToRef         string `chordGroup:"commits:Reset to ref"`
}

type ReflogCommitsChordPrefixes struct {
	ResetToRef string `chordGroup:"reflogCommits:Reset to ref"`
}

type SubCommitsChordPrefixes struct {
	ResetToRef string `chordGroup:"subCommits:Reset to ref"`
}

type CommitFilesChordPrefixes struct {
	CopyToClipboard string `chordGroup:"commitFiles:Copy to clipboard"`
}

type SubmodulesChordPrefixes struct {
	BulkOptions string `chordGroup:"submodules:Bulk options"`
}

// ResolveChordPrefixes populates cfg.Keybinding.ChordPrefix from
// KeybindingGroups. Unmatched fields stay empty so callers can detect
// configuration mismatches. Safe to call multiple times.
func (cfg *UserConfig) ResolveChordPrefixes() {
	populateChordPrefixStruct(
		reflect.ValueOf(&cfg.Keybinding.ChordPrefix).Elem(),
		cfg.KeybindingGroups,
	)
}

func populateChordPrefixStruct(v reflect.Value, groups map[string]map[string]KeybindingGroupConfig) {
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		fv := v.Field(i)
		if fv.Kind() == reflect.Struct {
			populateChordPrefixStruct(fv, groups)
			continue
		}
		if fv.Kind() != reflect.String {
			continue
		}
		tag := t.Field(i).Tag.Get("chordGroup")
		if tag == "" {
			continue
		}
		ctx, name, ok := splitChordGroupTag(tag)
		if !ok {
			continue
		}
		for prefix, group := range groups[ctx] {
			if group.Name == name {
				fv.SetString(prefix)
				break
			}
		}
	}
}

func splitChordGroupTag(tag string) (context, name string, ok bool) {
	idx := strings.Index(tag, ":")
	if idx < 0 {
		return "", "", false
	}
	return tag[:idx], tag[idx+1:], true
}

// MustChordPrefix is for test setup paths where a missing group is a
// programming error.
func (cfg *UserConfig) MustChordPrefix(context, name string) string {
	for prefix, group := range cfg.KeybindingGroups[context] {
		if group.Name == name {
			return prefix
		}
	}
	for prefix, group := range cfg.KeybindingGroups["global"] {
		if group.Name == name {
			return prefix
		}
	}
	panic(fmt.Sprintf("MustChordPrefix: no chord group %q in context %q (or global)", name, context))
}
