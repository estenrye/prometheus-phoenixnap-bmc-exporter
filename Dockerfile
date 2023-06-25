# Build binary
FROM --platform=$BUILDPLATFORM golang:1.19-alpine AS build-env
ADD . /app
WORKDIR /app
ARG TARGETOS
ARG TARGETARCH
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
    go build -ldflags="-s -w" -o prometheus-phoenixnap-bmc-exporter ./cmd/exporter/main.go

# Create image
FROM scratch
COPY --from=build-env /app/prometheus-phoenixnap-bmc-exporter /
COPY --from=build-env /etc/ssl /etc/ssl
ENTRYPOINT ["/prometheus-phoenixnap-bmc-exporter"]