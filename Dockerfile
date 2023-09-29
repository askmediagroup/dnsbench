FROM gcr.io/distroless/static-debian12
COPY dnsbench /dnsbench
ENTRYPOINT ["/dnsbench"]
