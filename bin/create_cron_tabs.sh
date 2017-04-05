#!/bin/sh
set -e

(
    flock -x -n 200 || exit 1

    now=`date`

    has_crontabs=false
    > /etc/crontabs/milo
    for i in /mnt/* ; do
      if [ -f "${i}/code/live/cron/${DEPLOY_ENV}.cron" ]; then
        has_crontabs=true
        echo "### $(basename $i) ###" >> /etc/crontabs/milo;
        cat ${i}/code/live/cron/${DEPLOY_ENV}.cron >> /etc/crontabs/milo;
        echo "${now} - Added ${i}/code/live/cron/${DEPLOY_ENV}.cron";
      fi
    done;

    if [ "${has_crontabs}" = true ] ; then
        if ! cmp /etc/crontabs/milo /etc/crontabs/root >/dev/null 2>&1
        then
            crontab /etc/crontabs/milo;
            echo "${now} - Reloaded Crontabs";
        fi
    fi
) 200>/var/lock/.crongen.exclusivelock