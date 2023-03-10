.PHONY: build run 

SERVER_PORT = 8080

build:
	docker build -t drugowick.dev/kindle-calendar-reader:latest .

run: build
	docker stop kindle-calendar-reader || true
	docker rm kindle-calendar-reader || true
	docker container run -d --name kindle-calendar-reader \
 		-p ${SERVER_PORT}:${SERVER_PORT} -e SERVER_PORT=${SERVER_PORT} \
 		drugowick.dev/kindle-calendar-reader
