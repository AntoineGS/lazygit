# Local Branch Hierarchy Display Design

## Goal

Add an optional local-branch sort mode that groups local branches by inferred
tip ancestry and shows each branch's downstream depth. This makes ordinary
stacked branches readable from the local branches panel without changing the
default display.

## Scope

In scope:

- A `hierarchy` value for `git.localBranchSortOrder`.
- A Hierarchy item in the local branch Sort Order menu.
- Best-effort parent inference from current local branch tips.
- Newest-commit-first ordering among roots and siblings.
- One ASCII `-` prefix per downstream level.
- Flat rendering while the branch list is filtered.
- Fast and legacy ancestry-loading paths.
- Unit, integration, configuration, and presentation coverage.

Out of scope:

- Persisting explicit parent relationships.
- Reconstructing branch creation history from reflogs.
- Inferring relationships from remote tracking branches.
- Changing the default local branch sort order.
- Adding collapse or expand controls.

## User Experience

The local branch Sort Order menu gains an `h: Hierarchy` option. The existing
Date default remains unchanged. Users can persist the mode with:

```yaml
git:
  localBranchSortOrder: hierarchy
```

Hierarchy mode renders roots exactly as they render today. Each downstream
level adds one neutral-colored ASCII dash immediately before the branch name:

```text
master
-feature/auth
--fix/login
--feature/2fa
-feature/billing
--fix/tax
release
```

The prefix consumes one column per level and participates in the existing name
width calculation. Branch status, worktree, pull request, commit hash, and
divergence fields retain their current alignment. Existing branch-name coloring
applies to the name, not to the neutral hierarchy prefix. If a narrow panel and
deep hierarchy cannot fit both the full prefix and the existing three-character
minimum branch name, omit leading dashes until the minimum name fits.

While a branch filter is active, only matching branches are shown and all
hierarchy prefixes are omitted. Fuzzy matches retain the existing
relevance-ranked filter order rather than being forced back into hierarchy
order. Clearing the filter restores the prefixes.

## Relationship Semantics

Git does not persist the branch from which another branch was created. The
display therefore represents relationships between current branch tips, not
historical branch-creation events.

A local branch is a parent candidate when its current tip is a strict ancestor
of the child's current tip. Branches at the same commit are not parents of one
another, although either can be a parent candidate for a later tip.

Parent selection proceeds as follows:

1. Remove candidate ancestors that have another local branch candidate between
   them and the child.
2. Prefer candidates with the smallest number of commits from candidate tip to
   child tip.
3. If candidates remain tied, prefer the candidate with the most children from
   unambiguous parent assignments elsewhere in the same load.
4. Prefer a candidate named in `git.mainBranches`.
5. Prefer the candidate with the oldest tip commit date.
6. Break any final tie alphabetically by branch name.

An assignment is unambiguous when exactly one candidate remains after the
between-branch reduction and minimum-distance comparison. Child counts include
only those assignments and are computed before resolving ambiguous children,
so the result is deterministic and does not depend on traversal order.

This produces a forest because unrelated histories and branches with no named
ancestor remain roots. Roots and siblings are sorted by descending tip commit
date, then alphabetically when dates are equal. The forest is flattened with a
depth-first traversal.

The root subtree containing the checked-out branch is moved before other root
subtrees. The checked-out branch remains at its inferred depth instead of being
detached from its parent. A detached or unborn HEAD does not participate in
ancestry inference and remains a standalone first root.

## Data Loading

The ordinary branch query remains the source of branch names, head hashes,
remote state, subjects, and tip dates. Hierarchy mode retains the parsed Unix
tip timestamp for deterministic peer ordering and tie-breaking.

### Git 2.41 And Later

Reuse the existing `for-each-ref` ahead/behind format builder and parser used by
base-branch divergence loading. Supply unique local branch tip object IDs as
bases; object IDs are safe inside Git format atoms even when a valid branch name
contains format punctuation such as `)`. Keep the mapping from each object ID
to all branch refs at that tip. For a child row and candidate base:

- `behind == 0` means the candidate is an ancestor of the child.
- `ahead` is the distance used when direct candidates remain ambiguous.
- Equal head hashes are excluded from strict parent relationships.

The existing parser must first be generalized to retain the base revision
associated with every valid ahead/behind field. It currently filters empty
fields and loses their positional identity. The refactor is behavior-preserving
for divergence loading and gives both consumers one command-construction and
parsing path. If a divergence row has no valid values, its behind count safely
defaults to zero rather than calling the minimum selector with an empty slice.

Split base object IDs into batches whose generated `--format` argument is at
most 16 KiB, then merge rows by child ref. This bounds command-line size and
peak output memory while still computing the selected mode for repositories
with many local branches. Total ancestry work remains quadratic in the number
of unique local branch tips.

### Older Git Versions

For each child, run `for-each-ref --merged=<child> refs/heads` to obtain local
ancestor candidates. Use the complete relation to remove candidates hidden
behind another local branch. Run `rev-list --count <candidate>..<child>` only
when multiple direct candidates still require distance comparison.

The older path produces the same forest semantics while accepting additional
Git processes only when hierarchy mode is selected. Synthetic detached and
unborn HEAD entries are excluded because they are not refs under `refs/heads`.

## Model And Ordering

Add the branch tip timestamp and resulting hierarchy depth to `models.Branch`.
Do not store parent pointers or child collections in the shared model. The
loader owns inference and emits the final flat order; presentation only needs
the depth.

Current code assumes `Model().Branches[0]` is HEAD in several branch and merge
helpers. Before hierarchy behavior lands, refactor those call sites to find the
branch whose `Head` field is true through the existing checked-out-ref helper.
The `SelectCheckedOutBranch` refresh behavior must likewise find and select the
actual `Head` row, ensure that nested row is visible, and update its row-zero
contract comment. Existing date, recency, and alphabetical modes continue
moving HEAD to index zero. Only hierarchy mode uses the root-subtree-first rule.

Branch presentation receives the context's filtering state in addition to each
branch's depth. This lets it suppress prefixes without mutating model depth or
changing the shared fuzzy-filter implementation.

## Configuration And Translations

Extend the local sort-order schema and validation enum with `hierarchy`, and
update the generated configuration documentation. Add English translation
fields for the Hierarchy label and its sort-menu description. Main-branch
tie-breaking compares local candidate short names directly with
`UserConfig.Git.MainBranches`; it does not use remote-resolved refs from
`MainBranches.Get()`. Do not edit Crowdin-maintained translation files or
release documentation under `docs/`.

## Failure Behavior

If ancestry loading or parsing fails, return the branches already loaded in
newest-commit-first flat order with zero hierarchy depths and log the error. As
in existing flat modes, detached or unborn HEAD remains promoted to row zero in
this fallback. The panel remains usable and a later refresh can retry. Malformed
individual matrix fields are ignored without shifting the identities of later
fields.

Concurrent refresh sequencing remains unchanged: inference and ordering finish
on the worker before the branch slice is assigned on the UI thread.

## Testing

Unit tests cover:

- Ahead/behind parsing with base-revision identity, including empty and malformed
  fields.
- Linear stacks, siblings, unrelated roots, equal head hashes, and detached
  and unborn HEAD.
- Multiple direct ancestors and each deterministic tie-break.
- Root and sibling date ordering and checked-out-subtree promotion.
- Safe OID bases, modern batched matrix loading, and the older `--merged`
  fallback with fake runners.
- Dash prefixes, width calculations, truncation, full-description mode, and
  flat relevance-ranked filtered rendering, including prefixes too deep for a
  narrow panel.
- Configuration validation for `hierarchy`.

Extend the existing `SortLocalBranches` integration test. Its fixture already
forms `master -> first -> second -> third`; selecting Hierarchy must render
`master`, `-first`, `--second`, and `---third` in that order.

Run generation and formatting, then the unit suite, lint, build, and targeted
integration test before completing the implementation.

## Commit Structure

Keep each commit compiling, formatted, lint-clean, and tested:

1. Refactor checked-out-branch lookup and refresh selection so they do not
   depend on row zero.
2. Generalize ahead/behind parsing and query reuse without changing behavior.
3. Add hierarchy inference, configuration, presentation, documentation, and
   tests as one behavior change.
