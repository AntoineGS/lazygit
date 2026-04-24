#!/usr/bin/env bash
# Reproducer pager. Two jobs:
#   1. Ignore the signals lazygit relies on for cleanup (TERM, HUP, PIPE).
#      This mimics pagers like difft when they're stuck in their parse /
#      render phase and don't notice their pty going away.
#   2. Hold a chunk of RSS so each orphan is obviously visible in
#      `ps`/`free`. Each instance parks ~50 MiB, so switching through
#      a handful of diffs in lazygit adds up fast.

trap '' TERM HUP PIPE

# Drain stdin (the diff) so git doesn't block writing to us.
DIFF=$(cat)

# Pad RSS so the orphan is unmistakable in ps/free. ~50 MiB once base64
# has expanded the raw zeros.
PADDING=$(head -c 40000000 /dev/zero | base64)

# Pretend to be busy forever. lazygit will SIGTERM us on quit; we ignore
# it. Nothing else can make us exit until the host reaps us.
sleep 3600
