# syntax=docker/dockerfile:1

FROM golang:1.21.1 AS build-stage
WORKDIR /app
COPY currency ./currency
COPY go.mod go.sum *.go *.yml ./
# ENV HTTP_RPOXY="socks5://localhost:1080"
# RUN CGO_ENABLED=0 GOOS=linux go build -o /geek
RUN go mod download && CGO_ENABLED=0 GOOS=linux go build -ldflags "-s -w" -o /geek
# CMD [ "sh","-c","while true;do sleep 10800;/geek;done" ]

FROM gcr.io/distroless/base-debian11 AS build-release-stage
WORKDIR /
COPY --from=build-stage /geek /geek
COPY --from=build-stage /app/config.yml ./
USER nonroot:nonroot
ENV HTTP_RPOXY="socks5://localhost:1080"
ENTRYPOINT ["/geek" ]
