## Build
FROM golang:1.22.1-alpine AS build

WORKDIR $GOPATH/src/shopifyx

# manage dependencies
COPY . .
RUN go mod download

RUN go build -a -o /shopifyx-server ./main.go


## Deploy
FROM alpine:latest
RUN apk add tzdata
COPY --from=build /shopifyx-server /shopifyx-server

EXPOSE 8080

ENTRYPOINT ["/shopifyx-server"]