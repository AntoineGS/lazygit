# Pager-process leak reproduction — issue #2460

Self-contained, interactive reproducer for
<https://github.com/jesseduffield/lazygit/issues/2460>.

This branch is throwaway — delete it once the issue is fixed.

## Run

```sh
./leak-repro/reproduce.sh
```

Requires `go`, `git`, `bash`. The script:

1. Builds lazygit from this checkout (`go build`).
2. Creates a throwaway repo in `$(mktemp -d)` with **5 uncommitted
   files**, so you have something to click through.
3. Installs `leak-repro/bad-pager.sh` as lazygit's pager — a shell
   script that installs `trap '' TERM HUP PIPE` and holds ~50 MiB of
   RSS. This mimics the state real pagers (`difft` in particular) end
   up in during heavy parsing: not responsive to `SIGTERM`, not exiting
   on `SIGPIPE`, holding memory.
4. Launches lazygit **in your current terminal** so you can interact
   with it.
5. When you quit lazygit, the script resumes and reports how many
   pagers survived and how much RSS they're still holding.

## How to observe the leak

Before you press ENTER to launch lazygit, open another terminal and run

```sh
watch -n1 'free -m; echo; pgrep -af /tmp/lg-leak-repro | wc -l'
```

(the exact path is printed by the script). Then inside lazygit:

- Use `j` / `k` (or arrow keys) to move through `a.txt … e.txt` in the
  Files panel. **Each selection spawns a new pager.**
- Old pagers don't die — they ignore `SIGTERM` and aren't in lazygit's
  process group, so they keep running and keep their ~50 MiB.
- Watch the other terminal: orphan count climbs, `free -m` drops.

When you're ready, press `q` to quit lazygit. The script prints a
table of PIDs + RSS and asks whether to clean them up.

Expected (today): something like
```
orphan count : 5
orphan RSS   : 285 MB
MemAvailable : before=28314 MB, after=28029 MB (delta=285 MB)
```

Expected after a fix: `no orphan pagers`.

## Mechanism

See `pkg/gui/pty.go:92-97` and `pkg/tasks/tasks.go:335-353`.

- `pty.StartWithSize` (via `creack/pty`) sets `Setsid: true`, so the
  pager is the leader of its own session — no `SIGHUP` cascade from
  lazygit exit.
- On view-switch / app-exit, cleanup calls
  `TerminateProcessGracefully(cmd)` which sends `SIGTERM` to the direct
  child only (git, not the pager's process group).
- `ViewBufferManager.Close` waits up to 3 seconds for the task to stop,
  then logs `cannot kill child process` and lets lazygit exit anyway.
- There is no `SIGKILL` fallback, and the pager's process group is
  never targeted.

Under these conditions a pager that's slow to respond to `SIGTERM`
(heavy Rust pager mid-parse, frozen difft, etc.) outlives lazygit
and stays around forever, holding its RSS until something else reaps
it.

## Variants worth testing against any proposed fix

- **SIGKILL on lazygit** (`kill -9 $(pgrep lazygit)`) while the pager
  is running. No userspace cleanup runs at all — only OS-level cleanup
  can help here. Dropping `Setsid` on the pty (or sending `SIGHUP` to
  the pager's session before `pty.StartWithSize` detaches it) would
  close this.
- **Terminal close** (parent shell killed while lazygit is mid-render).
  Same class as SIGKILL for the pager.
- **Well-behaved pager under SIGTERM grace** — `delta` handles
  `SIGPIPE` cleanly, so it exits on its own when the pty master closes.
  The leak only manifests when the pager is wedged or ignores signals,
  which is exactly the `difft`-under-load case users report.

## Files

- `bad-pager.sh` — the misbehaving pager used by the repro
  (~50 MiB RSS per instance, ignores TERM/HUP/PIPE).
- `reproduce.sh` — interactive driver: build, set up, hand off to you,
  inspect orphans, offer cleanup.

Nothing else in this branch is touched; delete it once the fix lands.
