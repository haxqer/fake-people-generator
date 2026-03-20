FROM golang:1.25-bookworm AS builder

WORKDIR /src

COPY go.mod ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
  -trimpath \
  -ldflags="-s -w" \
  -o /out/fake-people-generator \
  ./cmd/server

FROM gcr.io/distroless/static-debian12:nonroot

WORKDIR /app

COPY --from=builder /out/fake-people-generator /app/server

ENV ADDR=:8080

EXPOSE 8080

ENTRYPOINT ["/app/server"]
