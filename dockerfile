FROM golang:1.16 as gobuild
ARG VERSION=latest

WORKDIR /go/src/github.com/jonleopard/dankstats-api
ADD . .
ADD vendor ./vendor

RUN CGO_ENABLED=0 GOOS=linux GO111MODULE=on go build -mod=vendor -o dankstats-api -ldflags "-X main.version=$VERSION" main.go

FROM gcr.io/distroless/base

COPY --from=gobuild /go/src/github.com/jonleopard/dankstats-api/dankstats-api /bin

EXPOSE 4000

ENTRYPOINT ["/bin/dankstats-api"]
