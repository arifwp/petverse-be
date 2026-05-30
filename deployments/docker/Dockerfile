FROM golang:1.25-alpine AS build

WORKDIR /src
COPY go.mod ./
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/petverse-api ./cmd/api

FROM alpine:3.22

RUN adduser -D -H petverse
USER petverse

COPY --from=build /out/petverse-api /petverse-api

EXPOSE 8080
ENTRYPOINT ["/petverse-api"]
