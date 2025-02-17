FROM golang:latest AS builder

WORKDIR /app

COPY ./ ./
RUN make build

FROM busybox
WORKDIR /app
COPY --from=builder /app/kubernetes-mcp-server /app/kubernetes-mcp-server
ENTRYPOINT ["/app/kubernetes-mcp-server"]
