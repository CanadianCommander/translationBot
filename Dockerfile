FROM golang:1.18-alpine AS phase-0

RUN mkdir /var/source
WORKDIR /var/source

COPY . /var/source
RUN go build ./cmd/tb

FROM node:19.3.0-alpine3.17 AS phase-1

RUN mkdir /var/source
WORKDIR /var/source

COPY --from=phase-0 /var/source/tb /var/source/tb
COPY --from=phase-0 /var/source/cmd/tsJson /var/source/cmd/tsJson

RUN cd /var/source/cmd/tsJson && \
    yarn install

ENTRYPOINT ["/var/source/tb"]