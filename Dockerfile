FROM scratch

ADD assets/passwd.nobody /etc/passwd

USER nobody

ENTRYPOINT ["/m3u-filter"]

ARG ARCH=amd64

COPY build/m3u-filter_linux_${ARCH} /m3u-filter
