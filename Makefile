all: build

build:
	wire /src/cmd/analytics
	go build -o analytics /src/cmd/analytics/
