FROM gcr.io/distroless/static-debian12:nonroot
# FROM alpine:3
ARG GOOS=linux GOARCH=amd64 USER=65532
ENV CGO_ENABLED=0
COPY --chmod=755 --chown=$USER:0 ./bin/nextride-shortcut-${GOOS}-${GOARCH} /server

# USER $USER
EXPOSE 8080
ENTRYPOINT ["/server"]
