FROM golang:1.18-alpine AS phase-0

RUN mkdir /var/source
WORKDIR /var/source

COPY . /var/source
RUN go build ./cmd/tb

FROM node:18.12.1-alpine3.17 AS phase-1

RUN mkdir /var/source
WORKDIR /var/source

COPY --from=phase-0 /var/source/tb /var/source/tb
COPY --from=phase-0 /var/source/cmd/tsJson /var/source/cmd/tsJson
COPY --from=phase-0 /var/source/cmd/tsYaml /var/source/cmd/tsYaml

RUN apk add git bash

# install TS components
RUN cd /var/source/cmd/tsJson && \
    yarn install \
RUN cd /var/source/cmd/tsYaml && \
    yarn install

# setup git
RUN git config --global user.email "translation.bot@bbenetti.ca" && \
  git config --global user.name "Translation Bot"


ENTRYPOINT ["/var/source/tb"]