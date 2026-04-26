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
under `keybindingGroups`:

```yaml
keybindingGroups:
  "<b>":      { name: "Branch" }
  "<b><t>":   { name: "Pull Request" }
  "<c>":      { name: "Commit" }
  "<s>":      { name: "Stash" }
```

When a chord prefix is pending, sub-bindings whose next key matches a defined
group prefix are collapsed into a single footer row using the group's `name`.
For example, with the config above and three bindings under `<b><t>`
(`<b><t><o>`, `<b><t><l>`, `<b><t><c>`), pressing `b` shows one row labeled
`t: Pull Request` instead of three rows.

Validation rules:
- The prefix must be a valid chord-key string.
- The `name` must be non-empty.
- At least one keybinding must exist under the prefix.
- The prefix must not collide with a leaf binding using the same key sequence.

`keybindingGroups` is purely a footer-labeling and mnemonic-grouping
mechanism — it doesn't change which bindings fire or where. Navigation
between panes remains the user's responsibility (use `g{n}` or `<tab>`).
