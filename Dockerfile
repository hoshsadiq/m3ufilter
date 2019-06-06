ARG ARCH=amd64

FROM scratch

ADD assets/passwd.nobody /etc/passwd

USER nobody

COPY build/m3u-filter_linux_$ARCH /m3u-filter

ENTRYPOINT ["/m3u-filter"]