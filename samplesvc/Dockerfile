FROM golang:alpine3.13 as builder

WORKDIR /go/src/github.com/istio-oauth-poc/samplesvc

ENV GO111MODULE=on

RUN apk add git --no-cache

COPY . .

RUN go mod tidy && go mod vendor

RUN CGO_ENABLED=0 GOOS=linux go build -o samplesvc

FROM alpine:3.7

RUN addgroup -S samplesvc && adduser -S -g samplesvc samplesvc

USER samplesvc

COPY --from=builder /go/src/github.com/istio-oauth-poc/samplesvc/samplesvc /usr/bin/samplesvc

EXPOSE 4000

CMD ["/usr/bin/samplesvc"]