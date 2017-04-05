#!/bin/sh

set -e

WATCH_FOLDER=/mnt/
# Abort, if already running.
if [[ -n "$(ps | grep inotifywait | grep -v grep)" ]]; then
	echo "Already watching directory: ${WATCH_FOLDER}" >&2
	exit 1
fi


# Watch the live certificates directory. When changes are detected, install
# combined certificates and reload HAproxy.
echo "Watching directory: ${WATCH_FOLDER}"
inotifywait \
	--event create \
	--event delete \
	--event modify \
	--event move \
	--format "%e %w%f" \
	--monitor \
	--quiet \
	-r \
	"${WATCH_FOLDER}" |
while read CHANGED
do
	echo "$CHANGED"
    /opt/cronrunner/bin/create_cron_tabs.sh
done