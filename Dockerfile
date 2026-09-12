# download go dependencies for source code
FROM golang:1.27-alpine3.24 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN apk add --no-cache \
        aspell-en=2026.02.25-r0 \
        aspell=~0.60.8.2-r0 \
        make=~4.4.1-r4 \
    && go mod download

# build the server
COPY . ./
RUN make build/bin/server \
    GO_ARGS="CGO_ENABLED=0" \
    && go clean -cache

# copy the server to a minimal build image
FROM scratch
WORKDIR /app
COPY --from=builder /app/build/bin/server .
ENTRYPOINT [ "/app/server" ]
