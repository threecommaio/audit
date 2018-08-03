FROM golang:1.11beta2 as builder
ENV WORKDIR /go/src/github.com/threecommaio/audit

WORKDIR ${WORKDIR}
COPY go.mod go.sum ${WORKDIR}/
ENV GO111MODULE on
RUN go get

COPY . ${WORKDIR}

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app

FROM alpine:latest
ENV WORKDIR /go/src/github.com/threecommaio/audit
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder ${WORKDIR}/app /go/bin/app
ENTRYPOINT [ "/go/bin/app" ]
