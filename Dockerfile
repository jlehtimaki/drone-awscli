FROM golang:alpine AS builder
ENV GOARCH=amd64 GOOS=linux CGO_ENABLED=0
WORKDIR /build
COPY . .

RUN go build

FROM ubuntu:18.04

RUN apt update && apt install less

COPY --from=builder /build/drone_awscli /bin/


ENTRYPOINT ["/bin/drone_awscli"]                                                                                                                                                                                                                                                         ~/Projects/drone_awscli     master                                                                                                                                                                                                                                                                    11:23:02 
