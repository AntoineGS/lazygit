package config

// Stable identifiers for built-in chord groups, used by code (especially
// integration tests) that needs to look up a chord prefix without
// hard-coding the user-facing display Name. They survive translation
// of Name and user remapping of the chord prefix.
//
// IDs only need to be unique within a context: the same ID can appear
// in multiple contexts (e.g. ChordIDMerge in localBranches and
// remoteBranches; ChordIDResetToRef in commits/reflogCommits/subCommits).
const (
	ChordIDRebaseOptions          = "rebaseOptions"
	ChordIDIgnoreOptions          = "ignoreOptions"
	ChordIDStashOptions           = "stashOptions"
	ChordIDDiscardAndResetOptions = "discardAndResetOptions"
	ChordIDFilterFiles            = "filterFiles"
	ChordIDCopyToClipboard        = "copyToClipboard"
	ChordIDResetToUpstream        = "resetToUpstream"
	ChordIDDiscardChanges         = "discardChanges"
	ChordIDMerge                  = "merge"
	ChordIDDeleteBranch           = "deleteBranch"
	ChordIDGitFlowOptions         = "gitFlowOptions"
	ChordIDBranchUpstreamOptions  = "branchUpstreamOptions"
	ChordIDResetToRef             = "resetToRef"
	ChordIDDeleteRemoteBranch     = "deleteRemoteBranch"
	ChordIDDeleteTag              = "deleteTag"
	ChordIDBisectOptions          = "bisectOptions"
	ChordIDFixupCommitOptions     = "fixupCommitOptions"
	ChordIDBulkOptions            = "bulkOptions"
)

// ChordPrefixLookup resolves a chord-group ID to its current chord
// prefix label. Lookup checks the requested context first, then falls
// back to "global". Returns "" if no group with the given ID is
// registered.
//
// The unexported groups field shields this struct from
// reflection-based config walkers (e.g. the keybinding-collision
// validator), so it has no exported fields and cannot accidentally
// register a binding.
type ChordPrefixLookup struct {
	groups map[string]map[string]KeybindingGroupConfig
}

// Get returns the chord prefix label for the group with the given ID
// in the given context, or in "global" as a fallback. Returns "" if
// no matching group is registered.
func (l ChordPrefixLookup) Get(context, id string) string {
	if id == "" {
		return ""
	}
	if prefix := lookupChordID(l.groups, context, id); prefix != "" {
		return prefix
	}
	if context != "global" {
		return lookupChordID(l.groups, "global", id)
	}
	return ""
}

func lookupChordID(groups map[string]map[string]KeybindingGroupConfig, context, id string) string {
	for prefix, group := range groups[context] {
		if group.ID == id {
			return prefix
		}
	}
	return ""
}

// LookupGroup resolves a chord prefix label to its KeybindingGroupConfig
// using the standard "context first, then global" fallback. Entries with
// an empty Name are treated as missing (callers that want to display a
// title need a non-empty Name; an entry that exists only to attach an ID
// without a display name should not match).
//
// Returns the matched group and true on hit; ({}, false) on miss.
func LookupGroup(
	groups map[string]map[string]KeybindingGroupConfig,
	context, label string,
) (KeybindingGroupConfig, bool) {
	if g, ok := groups[context][label]; ok && g.Name != "" {
		return g, true
	}
	if g, ok := groups["global"][label]; ok && g.Name != "" {
		return g, true
	}
	return KeybindingGroupConfig{}, false
}

// ResolveChordPrefixes canonicalizes per-context KeybindingGroups labels
// (so equivalent spellings like "<C-b>" and "<ctrl+b>" share an entry)
// and wires the typed ChordPrefix lookup onto cfg. Callers must invoke
// this after KeybindingGroups has been finalized (defaults loaded, YAML
// decoded, legacy aliases applied). Safe to call multiple times:
// CanonicalizePrefixLabel is idempotent.
func (cfg *UserConfig) ResolveChordPrefixes() {
	for ctx, m := range cfg.KeybindingGroups {
		cfg.KeybindingGroups[ctx] = NormalizeGroupKeys(m)
	}
	cfg.Keybinding.ChordPrefix = ChordPrefixLookup{groups: cfg.KeybindingGroups}
}
