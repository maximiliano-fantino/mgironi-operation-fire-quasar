FROM golang:1.17-alpine

MAINTAINER mgioni@gmail.com

ENV GIN_MODE=release

WORKDIR $GOPATH/src

COPY . .

RUN go get -d -v ./... 

RUN go install -v ./...

EXPOSE $PORT

CMD ["operation-fire-quasar", "-profile=server"]
