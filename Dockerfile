#FROM golang:1.18-alpine
FROM balenalib/raspberry-pi-debian-golang:latest

WORKDIR /godog
COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

RUN go build

CMD [ "./godog" ]