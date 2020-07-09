all: build

build:
	wire ./cmd/analytics
	go build -o analytics ./cmd/analytics/main.go ./cmd/analytics/inject_{http,message_queue,repo,services}.go ./cmd/analytics/wire_gen.go
