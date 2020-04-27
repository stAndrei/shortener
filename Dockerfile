# Build
FROM golang:alpine AS build

ARG TAG
ARG BUILD

ENV APP shortener
ENV REPO prologic/$APP

RUN apk add --update git make build-base && \
    rm -rf /var/cache/apk/*

WORKDIR /go/src/github.com/$REPO
COPY . /go/src/github.com/$REPO
RUN make TAG=$TAG BUILD=$BUILD build

# Runtime
FROM scratch

ENV APP shortener
ENV REPO prologic/$APP

LABEL shortener.app main

COPY --from=build /go/src/github.com/${REPO}/${APP} /${APP}

EXPOSE 8000/tcp

ENTRYPOINT ["/shortener"]
CMD []
