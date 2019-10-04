FROM golang:1.12-alpine AS build_base
RUN apk add --no-cache bash ca-certificates make gcc git libc-dev

RUN mkdir -p /build
COPY production.env /build
COPY development.env /build

WORKDIR /build

COPY go.mod .
COPY go.sum .
RUN go mod download

FROM build_base AS server_builder
COPY . .
COPY production.env /build/production.env
COPY development.env /build/development.env

RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o /build/xxxxx


FROM alpine:latest AS runtime
RUN apk add ca-certificates
COPY --from=server_builder /build/xxxxx /build/
COPY --from=server_builder /build/production.env /
COPY --from=server_builder /build/development.env /
EXPOSE 7777
ENTRYPOINT ["/build/xxxxx"]