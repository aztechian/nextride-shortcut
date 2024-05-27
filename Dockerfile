FROM gcr.io/distroless/static
ARG GOOS=linux GOARCH=amd64
ENV CGO_ENABLED=0
COPY --chown=1001:0 ./bin/nextride-shortcut-${GOOS}-${GOARCH} /server

USER 1001
EXPOSE 8080
ENTRYPOINT ["/server"]
