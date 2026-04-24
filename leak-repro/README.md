# Pager-process leak reproduction — issue #2460

Self-contained reproducer for
<https://github.com/jesseduffield/lazygit/issues/2460>.

This branch is throwaway — it will be deleted once the issue is fixed.

## Run

```sh
./leak-repro/reproduce.sh
```

Requires `go`, `git`, `tmux`, `bash`. Exit code `2` means the leak was
reproduced; `0` means lazygit cleaned up correctly.

The script:

1. Builds lazygit from this checkout (`go build`).
2. Creates an ephemeral git repo with a tiny uncommitted diff in `$(mktemp -d)`.
3. Writes an isolated lazygit config that points the pager at
   `leak-repro/bad-pager.sh`, a shell script that installs
   `trap '' TERM HUP PIPE` (simulating a pager that doesn't honour
   `SIGTERM` — the observed state for `difft` while it's parsing a large
   diff).
4. Launches lazygit headlessly in tmux, waits for the pager to spawn,
   then sends `q` + `y` to quit.
5. After lazygit exits, checks whether the pager subtree is still alive.

On a clean shutdown the subtree should be gone. Today it survives:

```
systemd(1) ─── systemd(1153) ─── git(60998) ─── bash(60999) ─── sleep(61001)
```

The whole chain is reparented to `systemd --user` and keeps running
until `sleep 600` expires.

## Mechanism

See `pkg/gui/pty.go:92-97` and `pkg/tasks/tasks.go:335-353`. The relevant
observations:

- `pty.StartWithSize` (via `creack/pty`) sets `Setsid: true`, so the
  pager is in its own session — it won't receive `SIGHUP` when lazygit
  exits.
- On view-switch / app-exit, cleanup calls
  `TerminateProcessGracefully(cmd)` which sends `SIGTERM` to the direct
  child only (git, not the pager's process group).
- `ViewBufferManager.Close` waits up to 3 seconds for the task to stop,
  then prints `cannot kill child process` and lets lazygit exit anyway.
- There is no `SIGKILL` fallback, and the pager's process group is never
  targeted.

Under these conditions, a pager that's slow to respond to `SIGTERM`
(heavy Rust pager mid-parse, frozen difft, etc.) outlives lazygit and
stays around forever.

## Variants worth checking against any proposed fix

- **SIGKILL on lazygit** (`kill -9 $(pgrep lazygit)`) while the pager is
  running. No userspace cleanup runs at all; only OS-level cleanup can
  help here. Sending the pager `SIGHUP` via a shared process group
  before `pty.StartWithSize` sets `Setsid` would address this.
- **Terminal close** (parent shell killed while lazygit is mid-render).
  Same class as SIGKILL above for the pager; lazygit gets SIGHUP and
  may or may not have time to run defers.
- **Well-behaved pager under SIGTERM grace** — current `delta` handles
  SIGPIPE cleanly, so it will usually exit on its own when the pty
  master closes. The leak shows up only when the pager is wedged or
  ignores signals.

## Files

- `bad-pager.sh` — the misbehaving pager used by the repro.
- `reproduce.sh` — end-to-end driver (build, set up, drive tmux,
  inspect orphans, clean up).

Nothing else in this branch is touched; delete the branch once the
fix is merged.
