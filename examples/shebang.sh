#!/usr/bin/env bash
[ -z "$BUILDSH" ] && exec buildsh "$0" "$@"

# your code is after here...
echo "I'm in a container!"
