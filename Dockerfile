# syntax=docker/dockerfile:1

FROM golang:1.21.1
WORKDIR /app
COPY currency ./currency
COPY go.mod go.sum *.go *.yml ./
ENV HTTP_RPOXY="socks5://localhost:1080"
# RUN CGO_ENABLED=0 GOOS=linux go build -o /geek
RUN go mod download && CGO_ENABLED=0 GOOS=linux go build -ldflags "-s -w" -o /geek
CMD [ "/geek" ]