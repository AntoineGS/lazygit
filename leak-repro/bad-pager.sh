#!/usr/bin/env bash
# Minimal reproducer pager that ignores the signals lazygit uses to
# terminate children. Mirrors the bad state real-world pagers (notably
# difft during tree-sitter parsing) can end up in.

trap '' TERM HUP PIPE

# Drain stdin so `git diff` doesn't block on write, but never produce
# output: this keeps the pager "alive" and in a state where it won't
# notice the pty master closing.
cat > /dev/null &
sleep 600 &
wait
