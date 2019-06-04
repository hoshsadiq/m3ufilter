FROM scratch

ADD assets/passwd.nobody /etc/passwd

USER nobody

COPY build/m3u-filter_linux_amd64 /m3u-filter

ENTRYPOINT ["/m3u-filter"]