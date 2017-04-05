#!/bin/sh

set -e

/app/cron-runner &

/usr/sbin/crond -f
