build:
	docker build -t drugowick.dev/kindle-calendar-reader:latest .

run: build
	docker rm kindle-calendar-reader
	docker container run --name kindle-calendar-reader -p 8080:8080 drugowick.dev/kindle-calendar-reader
