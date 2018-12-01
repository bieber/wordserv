FROM golang:1.10-alpine AS builder
WORKDIR /go/src/github.com/bieber/wordserv
RUN apk add git
COPY ./ /go/src/github.com/bieber/wordserv
RUN go get -d github.com/bieber/wordserv
RUN mkdir /app
RUN go build -o /app/server github.com/bieber/wordserv

FROM alpine:latest
LABEL maintainer="docker@biebersprojects.com"
EXPOSE 80
ENV PORT 80
ENV BOOK_DIR /app/books

WORKDIR /app
ENV PORT 80
ENV BOOK_DIR /app/books
COPY --from=builder /app/server /app/server
ENTRYPOINT /app/server
