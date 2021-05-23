FROM scratch
copy dist/ocd /
expose 8080
VOLUME ["/tmp"]
ENTRYPOINT ["/ocd"]