FROM golang:alpine AS builder

USER root

ENV GIN_MODE=release
ENV PORT=8000
ENV GOPROXY=https://goproxy.io,direct
ENV CGO_ENABLED=0

# https://stackoverflow.com/a/53590802/10475870
RUN apk add --no-cache git ca-certificates

WORKDIR /app

COPY . .
RUN go get -d entgo.io/ent/cmd/ent
RUN go generate ./ent

# Build the binary.
RUN go build -ldflags="-w -s" -o /go/bin/server


FROM scratch
WORKDIR /app
COPY --from=builder /go/bin/server ./
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
EXPOSE 8000
ENTRYPOINT ["./server"]