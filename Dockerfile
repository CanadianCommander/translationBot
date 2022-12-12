FROM golang:1.18-alpine AS phase-0

RUN mkdir /var/source
WORKDIR /var/source

COPY . /var/source
RUN go build ./cmd/tb

FROM golang:1.18-alpine AS phase-1

RUN mkdir /var/source
WORKDIR /var/source
COPY --from=phase-0 /var/source/tb /var/source/tb
ENTRYPOINT ["/var/source/tb"]