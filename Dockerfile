FROM golang:1.14
RUN go get github.com/google/wire/cmd/wire
#RUN go get github.com/h3isenbug/url-analytics/...
WORKDIR /src
COPY . ./
RUN make

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /srv
COPY --from=0 /src/analytics .
CMD ["./analytics"]
