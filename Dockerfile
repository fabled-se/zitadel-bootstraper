FROM golang:1.21 as build

WORKDIR /app

COPY . .

RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build ./cmd/zitadel-bootstrapper

FROM alpine:latest

COPY --from=build /app/zitadel-bootstrapper ./

CMD ["./zitadel-bootstrapper"]
