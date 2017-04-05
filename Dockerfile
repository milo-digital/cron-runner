FROM alpine

RUN apk add --update \
		curl \
		tini \
		inotify-tools \
	&& rm -rf /var/cache/apk/*


COPY . /opt/cronrunner/

ENTRYPOINT ["/sbin/tini", "--"]
CMD ["/opt/cronrunner/bin/entrypoint.sh", "crond", "-f"]