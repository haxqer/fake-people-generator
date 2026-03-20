FROM --platform=$BUILDPLATFORM golang:1.25-bookworm AS builder

WORKDIR /src

ARG TARGETOS=linux
ARG TARGETARCH=amd64

COPY go.mod ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build \
      -trimpath \
      -ldflags="-s -w" \
      -o /out/fake-people-generator \
      ./cmd/server

FROM gcr.io/distroless/static-debian12:nonroot

WORKDIR /app

COPY --from=builder /out/fake-people-generator /app/fake-people-generator

ENV PORT=8080

EXPOSE 8080

ENTRYPOINT ["/app/fake-people-generator"]
