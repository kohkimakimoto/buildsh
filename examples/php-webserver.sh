#!/usr/bin/env bash
[ -z "$BUILDSH" ] && (sleep 1 && python -m webbrowser http://localhost:8888) &
[ -z "$BUILDSH" ] && export BUILDSH_DOCKER_OPTIONS="--rm -e TZ=Asia/Tokyo -i -t -p 8888:8888" && exec buildsh "$0" "$@"

# start php builtin server
php -S 0.0.0.0:8888
