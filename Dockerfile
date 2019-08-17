FROM node:12-alpine as webbuilder

ADD . /forest
RUN cd /forest/web \
  && yarn \
  && yarn build \
  && rm -rf node_module

FROM golang:1.12-alpine as builder

COPY --from=webbuilder /forest /forest

RUN apk update \
  && apk add git make \
  && go get -u github.com/gobuffalo/packr/v2/packr2 \
  && cd /forest \
  && make build

FROM alpine 

EXPOSE 7001

RUN addgroup -g 1000 go \
  && adduser -u 1000 -G go -s /bin/sh -D go \
  && apk add --no-cache ca-certificates

COPY --from=builder /forest/forest /usr/local/bin/forest
COPY --from=webbuilder /forest/font /font

USER go

WORKDIR /home/go

CMD ["forest"]
