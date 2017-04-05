#!/bin/sh

set -e

/opt/cronrunner/bin/create_cron_tabs.sh

/opt/cronrunner/bin/watch.sh &

exec "$@"