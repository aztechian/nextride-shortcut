ARG GO_VERSION=1.22.3
FROM golang:${GO_VERSION} AS build
ENV CGO_ENABLED=0

WORKDIR /app
COPY . /app
RUN make all && \
  chown 1001:0 $(make printexe)

FROM gcr.io/distroless/static
COPY --from=build /app/nextride-shortcut-linux-arm64 /server

USER 1001
EXPOSE 8080
ENTRYPOINT ["/server"]

