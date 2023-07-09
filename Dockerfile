FROM golang:alpine

ENV GIN_MODE=release
ENV PORT=8080

WORKDIR /go/src/github.com/sorathank/robinhood-assignment

COPY . .

RUN apk update
RUN apk add --no-cache git
RUN go get ./...
RUN go mod vendor
RUN go build main.go

EXPOSE $PORT

CMD [ "./main" ]
