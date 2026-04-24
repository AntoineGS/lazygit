#!/usr/bin/env bash
# Reproduces the pager-process leak reported in
# https://github.com/jesseduffield/lazygit/issues/2460
#
# Requirements: go, git, tmux, bash.
#
# Usage:
#   ./leak-repro/reproduce.sh
#
# Exit codes:
#   0  no orphan pager processes (bug not reproduced / fixed)
#   1  environmental problem (missing binary, build failure, etc.)
#   2  orphan pager survived lazygit exit (bug reproduced)

set -eu

HERE="$(cd "$(dirname "$0")" && pwd)"
LG_SRC="$(cd "$HERE/.." && pwd)"
TMP="$(mktemp -d -t lg-leak-repro-XXXX)"

cleanup() {
    # Stash a copy of any logs the user might want to inspect, then clean up.
    if [ -f "$TMP/lazygit.out" ]; then
        cp "$TMP/lazygit.out" "${TMPDIR:-/tmp}/lazygit-leak-repro.out" 2>/dev/null || true
    fi
    tmux -S "$TMP/tmux.sock" kill-server 2>/dev/null || true
    # Reap any orphans we created so the host isn't left with zombies
    # after the repro runs.
    pkill -9 -f "$TMP/bad-pager.sh" 2>/dev/null || true
    pgrep -f "$TMP/" -a 2>/dev/null | awk '{print $1}' | xargs -r kill -9 2>/dev/null || true
    rm -rf "$TMP"
}
trap cleanup EXIT

echo "=> workspace: $TMP"

# ---- 1. prerequisites --------------------------------------------------
for bin in go git tmux bash awk sed pgrep; do
    if ! command -v "$bin" >/dev/null 2>&1; then
        echo "error: '$bin' is required but not found in PATH" >&2
        exit 1
    fi
done

# ---- 2. build lazygit from this checkout ------------------------------
echo "=> building lazygit from $LG_SRC"
( cd "$LG_SRC" && go build -o "$TMP/lazygit" . ) || {
    echo "error: build failed" >&2
    exit 1
}

# ---- 3. install a SIGTERM-ignoring pager ------------------------------
install -m 0755 "$HERE/bad-pager.sh" "$TMP/bad-pager.sh"

# ---- 4. isolated lazygit config ---------------------------------------
mkdir -p "$TMP/lg-config"
cat > "$TMP/lg-config/config.yml" <<CONFIG
git:
  pagers:
    - colorArg: never
      pager: $TMP/bad-pager.sh
CONFIG

# ---- 5. throwaway repo with a diff for lazygit to render --------------
REPO="$TMP/repo"
mkdir -p "$REPO"
git -C "$REPO" init -q -b main
git -C "$REPO" -c user.email=test@test -c user.name=test commit -q --allow-empty -m bootstrap
seq 1 100 > "$REPO/file.txt"
git -C "$REPO" add file.txt
git -C "$REPO" -c user.email=test@test -c user.name=test commit -q -m initial
sed -i 's/$/-changed/' "$REPO/file.txt"

# ---- 6. drive lazygit headlessly in tmux ------------------------------
SOCK="$TMP/tmux.sock"
tmux -S "$SOCK" new-session -d -s s -x 200 -y 50 \
    "cd '$REPO' && '$TMP/lazygit' \
        --use-config-dir='$TMP/lg-config' \
        --use-config-file='$TMP/lg-config/config.yml' \
        > '$TMP/lazygit.out' 2>&1"

echo "=> lazygit launched; waiting 2s for pager to spawn"
sleep 2

# pgrep -f also matches the tmux server (its argv contains the full
# command we passed it); use -x to match the lazygit binary exactly.
# We copied lazygit into $TMP, so the basename is just "lazygit".
LG_PID="$(pgrep -x lazygit | while read -r p; do
    [ "$(readlink -f /proc/$p/exe 2>/dev/null)" = "$TMP/lazygit" ] && echo "$p" && break
done)"
if [ -z "$LG_PID" ]; then
    echo "error: lazygit did not start; its output was:" >&2
    sed 's/^/   /' "$TMP/lazygit.out" >&2
    exit 1
fi

echo "=> process tree while lazygit is running:"
if command -v pstree >/dev/null 2>&1; then
    pstree -ps "$LG_PID" | sed 's/^/   /'
else
    ps -eo pid,ppid,pgid,sid,stat,cmd --no-headers | \
        awk -v root="$LG_PID" '
            BEGIN { p[root]=1 }
            { pid=$1; ppid=$2; if (ppid in p) p[pid]=1; if (pid in p) print "   " $0 }'
fi

echo "=> sending quit ('q', then 'y' in case a confirm prompt appears)"
tmux -S "$SOCK" send-keys -t s 'q' || true
sleep 0.3
# The tmux session usually dies on the first 'q' because lazygit closes
# its window; send the confirmation only if the session is still up.
if tmux -S "$SOCK" has-session -t s 2>/dev/null; then
    tmux -S "$SOCK" send-keys -t s 'y' || true
fi

# lazygit's ViewBufferManager.Close has a 3-second deadline; wait longer
# so we can distinguish "still shutting down" from "truly leaked".
sleep 6

echo
echo "=> after lazygit exited:"
if kill -0 "$LG_PID" 2>/dev/null; then
    echo "   (unexpected: lazygit is still running as PID $LG_PID)"
    exit 1
else
    echo "   lazygit: gone"
fi

orphans="$(pgrep -f "$TMP/bad-pager.sh" || true)"
if [ -n "$orphans" ]; then
    echo "   *** ORPHAN PAGER SURVIVED ***"
    for p in $orphans; do
        ps -o pid,ppid,pgid,sid,stat,etime,cmd -p "$p" | tail -n +2 | sed 's/^/     /'
    done
    if command -v pstree >/dev/null 2>&1; then
        echo "     tree:"
        pstree -ps "$(echo "$orphans" | head -1)" | sed 's/^/     /'
    fi
    exit 2
fi

echo "   no orphan pager processes"
