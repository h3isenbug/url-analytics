FROM golang:1.14
WORKDIR /usr/local/go/src/github.com/h3isenbug/url-analytics
COPY * ./
RUN go get github.com/google/wire/cmd/wire
RUN make

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /srv
COPY --from=0 /usr/local/go/src/github.com/h3isenbug/url-analytics/analytics .
CMD ["./analytics"]
