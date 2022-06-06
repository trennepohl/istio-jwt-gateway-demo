FROM golang:alpine3.13 as builder

WORKDIR /go/src/github.com/istio-oauth-poc/authorization

ENV GO111MODULE=on

RUN apk add git --no-cache

COPY . .

RUN go mod tidy && go mod vendor

RUN CGO_ENABLED=0 GOOS=linux go build -o authorization

FROM alpine:3.7

RUN addgroup -S authorization && adduser -S -g authorization authorization

USER authorization

COPY --from=builder /go/src/github.com/istio-oauth-poc/authorization/authorization /usr/bin/authorization

ADD --chown=authorization:authorization certs/ /data

EXPOSE 4000

CMD ["/usr/bin/authorization"]