# syntax = docker/dockerfile:1.4

ARG GO_VERSION="1.20.1-bullseye"

FROM golang:${GO_VERSION} AS build

ENV DEBIAN_FRONTEND=noninteractive

WORKDIR /src

RUN apt-get update && apt-get install -y \
    git \
    && rm -rf /var/lib/apt/lists/*

COPY ./ ./

RUN CGO_ENABLED=0 go build -o /out/hello -ldflags '-s -d -w' ./cmd/web; chmod +x /out/hello

FROM scratch
COPY --from=build /out/hello /bin/hello
CMD ["/bin/hello"]