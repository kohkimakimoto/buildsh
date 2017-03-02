#!/bin/bash
set -e

# Workaround to use 'sudo' with arbitrary user id that is specified host machine.
# You should set '-e' docker run option like the following:
#   -e BUILDSH_USER="<user_id>:<group_id>"
if [ -n "$BUILDSH_USER" ]; then
    # split user_id and group_id
    OLD_IFS="$IFS"
    IFS=:
    arr=($BUILDSH_USER)
    IFS="$OLD_IFS"

    if [ ${#arr[@]} -ne 2 ]; then
        echo "'BUILDSH_USER' must be formated '<user_id>:<group_id>', but $BUILDSH_USER" 1>&2
        exit 1
    fi

    # Create buildbot user
    groupadd --non-unique --gid ${arr[1]} buildbot
    useradd --non-unique --uid ${arr[0]} --gid ${arr[1]} buildbot
    echo 'buildbot	ALL=(ALL) NOPASSWD:ALL' >> /etc/sudoers

    if [[ ! $@ ]]; then
        exec sudo -u buildbot -E /bin/bash
    else
        exec sudo -u buildbot -E $@
    fi
else
    if [[ ! $@ ]]; then
        exec /bin/bash
    else
        exec $@
    fi
fi
