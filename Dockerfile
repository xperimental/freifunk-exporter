FROM golang:1.19-alpine AS builder

RUN apk add make git

WORKDIR /build

COPY go.mod go.sum /build/
RUN go mod download
RUN go mod verify

COPY . /build/
RUN make

FROM busybox
LABEL maintainer="Robert Jacob <xperimental@solidproject.de>"

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /build/freifunk-exporter /bin/freifunk-exporter

USER nobody
EXPOSE 9295

ENTRYPOINT ["/bin/freifunk-exporter"]
