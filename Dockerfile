# FROM gcr.io/distroless/static
FROM alpine:3
ARG GOOS=linux GOARCH=amd64
ENV CGO_ENABLED=0
COPY --chmod=755 --chown=1001:0 ./bin/nextride-shortcut-${GOOS}-${GOARCH} /server
# COPY --chmod=755 --chown=0:0 ./bin/busybox /bin/busybox

USER 1001
EXPOSE 8080
ENTRYPOINT ["/server"]
