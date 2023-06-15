# pull official base image
FROM golang:1.20 as builder

WORKDIR /workspace
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download && go mod verify

# Copy the go source
COPY . .

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o gophkeeper-server ./cmd/server/main.go

## Use distroless as minimal base image to package the manager binary
## Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM gcr.io/distroless/static:nonroot
WORKDIR /app

COPY --from=builder /workspace/gophkeeper-server .
COPY --from=builder /workspace/example.server.json server.json

EXPOSE 3001
USER 65532:65532

ENTRYPOINT ["/app/gophkeeper-server"]