# syntax=docker/dockerfile:1.6
# Build stage
FROM --platform=$BUILDPLATFORM golang:1.24 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

ARG TARGETOS
ARG TARGETARCH
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -o /eso-filesecret-provider

# Final stage
FROM gcr.io/distroless/static:nonroot
COPY --from=builder /eso-filesecret-provider /eso-filesecret-provider
ENTRYPOINT ["/eso-filesecret-provider"]
EXPOSE 8080