FROM golang:1.18-alpine AS builder
RUN mkdir /app
WORKDIR /app
RUN apk add git
COPY ./ ./
RUN go build -o server

FROM alpine:latest
LABEL maintainer="docker@biebersprojects.com"
EXPOSE 80

WORKDIR /app
ENV PORT 80
ENV BOOK_DIR /app/books
COPY --from=builder /app/server /app/server
ENTRYPOINT /app/server
