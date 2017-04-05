#!/bin/sh

set -e

while [ true ]
do
    /opt/cronrunner/bin/create_cron_tabs.sh
    sleep 60
done