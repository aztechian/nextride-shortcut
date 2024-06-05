FROM gcr.io/distroless/static
# FROM alpine:3
ARG GOOS=linux GOARCH=amd64 USER=65532
ENV CGO_ENABLED=0
COPY --chmod=755 --chown=$USER:0 ./bin/nextride-shortcut-${GOOS}-${GOARCH} /server
# COPY --chmod=755 --chown=0:0 ./bin/busybox /bin/busybox

# USER $USER
EXPOSE 8080
ENTRYPOINT ["/server"]
