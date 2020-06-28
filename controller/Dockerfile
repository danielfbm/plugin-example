# Build the manager binary
FROM golang:1.13 as builder

WORKDIR /workspace
# Copy the Go Modules manifests
# ENV GOPROXY=https://goproxy.cn,direct
COPY ./ .
# COPY go.sum go.sum
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

# Copy the go source
# COPY main.go main.go
# COPY api/ api/
# COPY controllers/ controllers/
# COPY plugins/ plugins/

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -a -o manager main.go
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -a -o foo.po plugins/foo/main.go
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -a -o bar.po plugins/bar/main.go
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -a -o bar_grpc.po plugins/bargrpc/main.go
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -a -o foobar.po plugins/foobar/main.go


# Use distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM alpine
WORKDIR /
COPY --from=builder /workspace/manager .
COPY --from=builder /workspace/foo.po /plugins/foo.po
COPY --from=builder /workspace/bar.po /plugins/bar.po
COPY --from=builder /workspace/bar_grpc.po /plugins/bar-grpc.po
COPY --from=builder /workspace/foobar.po /plugins/foobar.po
# USER nonroot:nonroot

ENTRYPOINT ["/manager"]
