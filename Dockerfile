FROM golang:1.20 AS build
WORKDIR /make
COPY . .

RUN CGO_ENABLED=0 go build -o /bin/calendar-events ./cmd/main.go

FROM alpine:latest
RUN apk add --no-cache bash
COPY credentials.json credentials.json
COPY outlookCredentials.json outlookCredentials.json
COPY --from=build /bin/calendar-events /app/calendar-events
ENTRYPOINT ["/app/calendar-events"]
