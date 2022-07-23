build:
	docker build -t sse_server:1.0 .

tag:
	docker tag sse_server:1.0 sse_server:latest

start:
	docker run --rm -d --name sse_server -p 3003:3003 sse_server:latest

stop:
	docker stop sse_server

run:
	PORT=3003 go run main.go

bin:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o sse_server .

install:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go install -ldflags="-s -w"



.PHONY: build tag start stop run bin install
