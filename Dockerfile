# Build image
FROM golang:1.18-alpine AS builder

RUN apk update \
    && apk upgrade \
    && apk add --update \
    ca-certificates \
    gcc \
    git \
    libc-dev \
    make \
    && update-ca-certificates

WORKDIR $GOPATH/src/github.com/cadicallegari/user

COPY go.mod go.sum ./
RUN go mod download

ARG GIT_TAG
ARG GIT_COMMIT
ENV GIT_TAG $GIT_TAG
ENV GIT_COMMIT $GIT_COMMIT

COPY . ./
RUN make go-install

# Final image
FROM alpine:3.16

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/bin/user /usr/bin/

EXPOSE 80

CMD ["user"]
