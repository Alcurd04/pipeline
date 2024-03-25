# syntax=docker/dockerfile:1
FROM golang:latest as build
RUN mkdir -p /go/src/pipeline
WORKDIR /go/src/pipeline
ADD pipeline.go .
ADD go.mod .
#RUN go build .
RUN go get
RUN go build -o /bin/pipeline ./pipeline.go
#RUN ls

FROM ubuntu:latest
LABEL version="1.0.0"
LABEL maintainer="Test Student<test@test.ru>"
#WORKDIR /root/
COPY --from=build /bin/pipeline /bin/pipeline
#RUN ls /go/bin/pipeline/*
#RUN chmod +x /go/bin/pipeline/pipeline
CMD ["/bin/pipeline"]