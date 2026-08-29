# Build a small static Go binary, cross-compiled to the target platform so multi-arch
# builds compile natively on the runner instead of under slow QEMU emulation. The frontend
# is embedded via go:embed, so there is no Node build stage.
FROM --platform=$BUILDPLATFORM golang:1.27-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
ARG TARGETOS TARGETARCH TARGETVARIANT
ARG VERSION=dev
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH \
    GOARM=$(printf '%s' "$TARGETVARIANT" | tr -d v) \
    go build -ldflags="-w -s -X main.version=$VERSION" -o /who-dat ./cmd/server

FROM alpine:latest
RUN apk --no-cache add ca-certificates && adduser -D -u 1000 app
COPY --from=build /who-dat /usr/local/bin/who-dat
USER app
EXPOSE 8080
HEALTHCHECK --interval=30s --timeout=3s --retries=3 \
  CMD wget -q -O- http://localhost:8080/health || exit 1
ENTRYPOINT ["who-dat"]
