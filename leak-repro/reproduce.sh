#!/usr/bin/env bash
# Interactive reproducer for the pager-process leak
# (https://github.com/jesseduffield/lazygit/issues/2460).
#
# Requirements: go, git, bash.
#
# What you'll see:
#   - lazygit launches in this terminal against a throwaway repo with 5
#     modified files.
#   - Every diff you open spawns a pager that ignores SIGTERM. Flicking
#     up/down the file list spawns new ones without reaping the old, so
#     RSS climbs with every switch.
#   - When you quit lazygit (press `q`), this script resumes and shows
#     all the orphan pagers that survived, with their RSS.
#
# Tip: open another terminal and run `watch -n1 free -m` (or `htop`) so
# you can see total memory change as you navigate.

set -eu

HERE="$(cd "$(dirname "$0")" && pwd)"
LG_SRC="$(cd "$HERE/.." && pwd)"
TMP="$(mktemp -d -t lg-leak-repro-XXXX)"

cleanup() {
    if [ "${ORPHANS_FOUND:-0}" = "1" ] && [ "${AUTO_CLEANUP:-1}" = "1" ]; then
        pkill -9 -f "$TMP/bad-pager.sh" 2>/dev/null || true
        # The orphan chain is `git -> bash(bad-pager) -> sleep`. Killing
        # bad-pager also reaps its sleep child; reap dangling gits too.
        pgrep -f "$TMP/" 2>/dev/null | xargs -r kill -9 2>/dev/null || true
    fi
    rm -rf "$TMP"
}
trap cleanup EXIT

if ! [ -t 0 ] || ! [ -t 1 ]; then
    echo "error: this reproducer needs an interactive terminal" >&2
    exit 1
fi

# ---- prerequisites -----------------------------------------------------
for bin in go git bash sed awk pgrep ps; do
    command -v "$bin" >/dev/null 2>&1 || {
        echo "error: '$bin' is required but not found in PATH" >&2
        exit 1
    }
done

# ---- build lazygit from this checkout ---------------------------------
echo "=> building lazygit from $LG_SRC"
( cd "$LG_SRC" && go build -o "$TMP/lazygit" . ) || {
    echo "error: build failed" >&2
    exit 1
}

# ---- install the misbehaving pager ------------------------------------
install -m 0755 "$HERE/bad-pager.sh" "$TMP/bad-pager.sh"

# ---- isolated lazygit config that points at our pager -----------------
mkdir -p "$TMP/lg-config"
cat > "$TMP/lg-config/config.yml" <<CONFIG
git:
  pagers:
    - colorArg: never
      pager: $TMP/bad-pager.sh
CONFIG

# ---- throwaway repo with 5 modified files -----------------------------
REPO="$TMP/repo"
mkdir -p "$REPO"
git -C "$REPO" init -q -b main
git -C "$REPO" -c user.email=test@test -c user.name=test \
    commit -q --allow-empty -m bootstrap
for name in a b c d e; do
    seq 1 100 > "$REPO/$name.txt"
done
git -C "$REPO" add .
git -C "$REPO" -c user.email=test@test -c user.name=test \
    commit -q -m initial
for name in a b c d e; do
    sed -i "s/\$/-$name/" "$REPO/$name.txt"
done

# ---- record RSS baseline so we can compare after quitting -------------
RSS_BEFORE_KB=$(awk '/^MemAvailable:/ {print $2}' /proc/meminfo 2>/dev/null || echo 0)

# ---- instructions -------------------------------------------------------
cat <<INFO

============================================================
Ready to launch lazygit against $REPO
============================================================

  - 5 files (a..e.txt) are uncommitted. Use j/k or arrows in the
    Files panel to move between them. Each selection spawns a
    new pager.
  - The pager holds ~50 MiB of RSS and ignores SIGTERM/HUP/PIPE,
    so switching files leaks memory.
  - When you're done, quit lazygit normally (press q).

Optional: open another terminal and run

    watch -n1 'free -m; echo; pgrep -af "$TMP/bad-pager" | wc -l'

to watch free memory drop and orphan count climb.

Press ENTER to launch lazygit.
INFO
# shellcheck disable=SC2034
read -r _

"$TMP/lazygit" \
    --use-config-dir="$TMP/lg-config" \
    --use-config-file="$TMP/lg-config/config.yml" \
    -p "$REPO" \
    || true

# ---- inspect orphans --------------------------------------------------
echo
echo "=> lazygit exited; inspecting orphan pagers..."
sleep 1

mapfile -t ORPHANS < <(pgrep -f "$TMP/bad-pager.sh" || true)
ORPHAN_COUNT=${#ORPHANS[@]}

RSS_AFTER_KB=$(awk '/^MemAvailable:/ {print $2}' /proc/meminfo 2>/dev/null || echo 0)

if [ "$ORPHAN_COUNT" -eq 0 ]; then
    echo "   no orphan pagers — either the leak is fixed or no diffs opened"
    exit 0
fi

ORPHANS_FOUND=1
TOTAL_RSS_KB=0
printf '\n   %-7s %-8s %-10s %s\n' "PID" "RSS(MB)" "ELAPSED" "CMD"
printf '   %s\n' "-----------------------------------------------------------"
for p in "${ORPHANS[@]}"; do
    read -r rss etime cmd < <(ps -o rss=,etime=,cmd= -p "$p" 2>/dev/null) || continue
    TOTAL_RSS_KB=$(( TOTAL_RSS_KB + rss ))
    printf '   %-7s %-8s %-10s %s\n' "$p" "$(( rss / 1024 ))" "$etime" "$cmd"
done

echo
echo "   orphan count : $ORPHAN_COUNT"
echo "   orphan RSS   : $(( TOTAL_RSS_KB / 1024 )) MB"
if [ "$RSS_BEFORE_KB" -gt 0 ] && [ "$RSS_AFTER_KB" -gt 0 ]; then
    DELTA_KB=$(( RSS_BEFORE_KB - RSS_AFTER_KB ))
    echo "   MemAvailable : before=$(( RSS_BEFORE_KB / 1024 )) MB, after=$(( RSS_AFTER_KB / 1024 )) MB (delta=$(( DELTA_KB / 1024 )) MB)"
fi

echo
echo "These processes will keep running (and holding RSS) until killed."
printf "Clean them up now? [Y/n] "
read -r answer
case "${answer:-Y}" in
    [nN]*)
        echo "Leaving orphans alive. PIDs: ${ORPHANS[*]}"
        AUTO_CLEANUP=0
        ;;
    *)
        echo "Killing orphans..."
        ;;
esac

exit 2
