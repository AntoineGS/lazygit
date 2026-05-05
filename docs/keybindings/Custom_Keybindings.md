## Possible keybindings
| Put in        | You will get   |
|---------------|----------------|
| `<f1>`        | F1             |
| `<f2>`        | F2             |
| `<f3>`        | F3             |
| `<f4>`        | F4             |
| `<f5>`        | F5             |
| `<f6>`        | F6             |
| `<f7>`        | F7             |
| `<f8>`        | F8             |
| `<f9>`        | F9             |
| `<f10>`       | F10            |
| `<f11>`       | F11            |
| `<f12>`       | F12            |
| `<insert>`    | Insert         |
| `<delete>`    | Delete         |
| `<home>`      | Home           |
| `<end>`       | End            |
| `<pgup>`      | Pgup           |
| `<pgdown>`    | Pgdn           |
| `<up>`        | ArrowUp        |
| `<s-up>`      | ShiftArrowUp   |
| `<down>`      | ArrowDown      |
| `<s-down>`    | ShiftArrowDown |
| `<left>`      | ArrowLeft      |
| `<right>`     | ArrowRight     |
| `<tab>`       | Tab            |
| `<backtab>`   | Backtab        |
| `<enter>`     | Enter          |
| `<a-enter>`   | AltEnter       |
| `<esc>`       | Esc            |
| `<backspace>` | Backspace      |
| `<c-space>`   | CtrlSpace      |
| `<c-/>`       | CtrlSlash      |
| `<space>`     | Space          |
| `<c-a>`       | CtrlA          |
| `<c-b>`       | CtrlB          |
| `<c-c>`       | CtrlC          |
| `<c-d>`       | CtrlD          |
| `<c-e>`       | CtrlE          |
| `<c-f>`       | CtrlF          |
| `<c-g>`       | CtrlG          |
| `<c-j>`       | CtrlJ          |
| `<c-k>`       | CtrlK          |
| `<c-l>`       | CtrlL          |
| `<c-n>`       | CtrlN          |
| `<c-o>`       | CtrlO          |
| `<c-p>`       | CtrlP          |
| `<c-q>`       | CtrlQ          |
| `<c-r>`       | CtrlR          |
| `<c-s>`       | CtrlS          |
| `<c-t>`       | CtrlT          |
| `<c-u>`       | CtrlU          |
| `<c-v>`       | CtrlV          |
| `<c-w>`       | CtrlW          |
| `<c-x>`       | CtrlX          |
| `<c-y>`       | CtrlY          |
| `<c-z>`       | CtrlZ          |
| `<c-4>`       | Ctrl4          |
| `<c-5>`       | Ctrl5          |
| `<c-6>`       | Ctrl6          |
| `<c-8>`       | Ctrl8          |

## Chord Group Labels

You can label chord prefixes to make the footer more readable. Define groups
under `keybindingGroups`, nested by context name:

```yaml
keybindingGroups:
  global:
    "<b>":      { name: "Branch" }
    "<b><t>":   { name: "Pull Request" }
  localBranches:
    "<r>":      { name: "Rebase options" }
  files:
    "<s>":      { name: "Stash options" }
```

The outer key is a context name (e.g. `global`, `files`, `localBranches`,
`commits`); the inner key is the chord-prefix label (e.g. `<b>`, `<ctrl+b>`).
Lookup at popup-open time checks the originating view's context first, then
falls back to `global`.

When a chord prefix is pending, sub-bindings whose next key matches a defined
group prefix are collapsed into a single footer row using the group's `name`.
For example, with three bindings under `<b><t>` (`<b><t><o>`, `<b><t><l>`,
`<b><t><c>`), pressing `b` shows one row labeled `t: Pull Request` instead of
three rows.

Validation rules:
- The prefix must be a valid chord-key string.
- The `name` must be non-empty.
- At least one keybinding must exist under the prefix.
- The prefix must not collide with a leaf binding using the same key sequence.

`keybindingGroups` is purely a footer-labeling and mnemonic-grouping
mechanism — it doesn't change which bindings fire or where. Navigation
between panes remains the user's responsibility (use `g{n}` or `<tab>`).

### Legacy flat shape

A single map keyed by chord-prefix label (no context layer) is still accepted
and migrated under `global`:

```yaml
keybindingGroups:
  "<b>":      { name: "Branch" }
  "<s>":      { name: "Stash" }
```

Prefer the nested shape for new configs — context-scoped groups (e.g. an
`r` group on `localBranches` only) require it.

## Migrating from older versions

The chord refactor split several combined-menu keybindings into individual
bindings sharing a common chord prefix. The following deprecated keys are
still recognised; setting one rewrites the chord-HEAD prefix on every
replacement binding to the value you provide. New configs should set the
individual bindings directly.

| Deprecated key                            | Default | Now controls chord HEAD for                                                                  |
|-------------------------------------------|---------|----------------------------------------------------------------------------------------------|
| `keybinding.files.ignoreFile`             | `i`     | `ignore`, `exclude`                                                                          |
| `keybinding.files.viewStashOptions`       | `S`     | `stashAllChangesKeepIndex`, `stashIncludeUntrackedChanges`, `stashStagedChanges`, `stashUnstagedChanges` |
| `keybinding.files.viewResetOptions`       | `D`     | `nukeWorkingTree`, `discardUnstagedChanges`, `discardUntrackedFiles`, `discardStagedChanges`, `softReset`, `mixedReset`, `hardReset` |
| `keybinding.files.openStatusFilter`       | `<ctrl+b>` | `filterStaged`, `filterUnstaged`, `filterTracked`, `filterUntracked`, `noFilter`           |
| `keybinding.files.copyFileInfoToClipboard`| `y`     | `copyFileName`, `copyRelativeFilePath`, `copyAbsoluteFilePath`, `copyFileDiff`, `copyAllFilesDiff` |
| `keybinding.branches.rebaseBranch`        | `r`     | `rebaseBranchSimple`, `rebaseBranchInteractive`, `rebaseBranchOntoBase`                      |
| `keybinding.branches.mergeIntoCurrentBranch` | `M`  | `mergeRegular`, `mergeNonFFwd`, `mergeFastForward`, `mergeSquash`, `mergeSquashCommitted`    |
| `keybinding.branches.viewGitFlowOptions`  | `i`     | `gitFlowFinish`, `gitFlowStartFeature`, `gitFlowStartHotfix`, `gitFlowStartBugfix`, `gitFlowStartRelease` |
| `keybinding.universal.createRebaseOptionsMenu` | `m` | `rebaseContinue`, `rebaseAbort`, `rebaseSkip`                                                |
| `keybinding.commits.viewBisectOptions`    | `b`     | `bisectMarkBad`, `bisectMarkGood`, `bisectSkipCurrent`, `bisectSkipSelected`, `bisectReset`, `bisectStartMarkBad`, `bisectStartMarkGood`, `bisectChooseTerms` |
| `keybinding.commits.viewResetOptions`     | `g`     | `mixedResetToRef`, `softResetToRef`, `hardResetToRef`                                        |
| `keybinding.submodules.bulkMenu`          | `b`     | `bulkInit`, `bulkUpdate`, `bulkUpdateRecursive`, `bulkDeinit`                                |

Setting a deprecated key to its default value is a no-op; setting it to a
different key sequence remaps the entire chord group. The deprecated key is
not surfaced in the cheatsheet.
