# build app
FROM golang:1.20-alpine3.16 AS app-builder

ARG VERSION=dev
ARG REVISION=dev
ARG BUILDTIME

# Install necessary packages for the build
RUN apk add --no-cache git tzdata sox flac ffmpeg

ENV SERVICE=redactedhook

WORKDIR /src

# Cache go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy rest of the source code
COPY . ./

#RUN go build -ldflags "-s -w -X main.version=${VERSION} -X main.commit=${REVISION} -X main.date=${BUILDTIME}" -o bin/propolis cmd/propolis/main.go
RUN go build -trimpath -ldflags "-s -w -X main.version=${VERSION} -X main.commit=${REVISION} -X main.date=${BUILDTIME}" -o bin/propolis cmd/propolis/main.go

# build runner
FROM alpine:latest

LABEL org.opencontainers.image.source="https://github.com/dsrtusr88/propolis"

ENV HOME="/propolis" \
    XDG_CONFIG_HOME="/propolis" \
    XDG_DATA_HOME="/propolis"

# Install runtime dependencies
RUN apk --no-cache add ca-certificates curl tzdata sox flac ffmpeg

WORKDIR /propolis

VOLUME /propolis

COPY --from=app-builder /src/bin/propolis /usr/local/bin/

EXPOSE 42135

ENTRYPOINT ["/usr/local/bin/propolis", "--config", "config.toml"]
