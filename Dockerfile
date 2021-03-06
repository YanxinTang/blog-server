FROM golang:alpine AS builder

USER root

ENV GIN_MODE=release
ENV PORT=8000
ENV GOPROXY=https://goproxy.io,direct
ENV CGO_ENABLED=0

RUN apk add --no-cache git

WORKDIR /app

COPY . .
RUN go get -d -v

# Build the binary.
RUN go build -ldflags="-w -s" -o /go/bin/server


FROM scratch
WORKDIR /app
COPY --from=builder /go/bin/server ./
EXPOSE 8000
ENTRYPOINT ["./server"]