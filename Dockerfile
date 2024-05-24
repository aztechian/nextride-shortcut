FROM golang:1.22.3

WORKDIR /app
COPY . /app
RUN make all && \
  chown 1001:0 $(make printexe)

# FROM scratch
FROM busybox
COPY --from=0 /app/nextride-shortcut-linux-arm64 /server

EXPOSE 8080
ENTRYPOINT ["/server"]

