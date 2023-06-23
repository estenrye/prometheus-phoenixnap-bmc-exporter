# 2. Relevant Videos Used to Learn How to Write Prometheus Exporters

Date: 2023-06-23

## Status

Accepted

## Context

This ADR documents video resources used to educate the creation of this exporter.

## Decision

Review the following videos.

- [How to Build Custom Prometheus Exporter? (Step-by-Step - Real-world Example - Parse Log + HTTP)](https://www.youtube.com/watch?v=3wT0zSsQb58)

## Consequences

### Tools

To install the necessary tools for this repository:

```basn
brew bundle
```


### Creating the Exporter Directory

```bash
SOURCE_DIR="${HOME}/src"
EXPORTER_NAME="prometheus-phoenix-nap-exporter"
REPOSITORY_PATH="github.com/estenrye"

mkdir -p "${SOURCE_DIR}/${EXPORTER_NAME}/cmd/exporter"
cd "${SOURCE_DIR}/${EXPORTER_NAME}
git init
adr init
go mod init "${REPOSITORY_PATH}/${EXPORTER_NAME}"
touch cmd/exporter/main.go
go get github.com/prometheus/client_golang/prometheus
go get "k8s.io/apiserver/pkg/server/mux"
```