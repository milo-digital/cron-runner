#!/bin/sh
set -e

(
    flock -x -n 200 || exit 1
    sleep 5 #Wait to make sure everything is done moving

    has_crontabs=false
    > /etc/crontabs/milo
    for i in /mnt/* ; do
      if [ -f "${i}/live/cron/${DEPLOY_ENV}.cron" ]; then
        has_crontabs=true
        echo "### $(basename \"$i\") ###" >> /etc/crontabs/milo;
        cat ${i}/live/cron/${DEPLOY_ENV}.cron >> /etc/crontabs/milo;
        echo "Added ${i}/live/cron/${DEPLOY_ENV}.cron";
      fi
    done;

    if [ "${has_crontabs}" = true ] ; then
        crontab /etc/crontabs/milo;
        echo "Reloaded Crontabs";
    else
        echo "No Crontabs Found";
    fi
) 200>/var/lock/.crongen.exclusivelock