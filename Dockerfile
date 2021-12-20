FROM golang:1.17-alpine AS build

ENV APP=./cmd/app
ENV BIN=/bin/faceitsubscriber
ENV PATH_ROJECT=${GOPATH}/src/github.com/nakiner/faceit-subscriber
ENV GO111MODULE=on
ENV GOSUMDB=off
ENV GOFLAGS=-mod=vendor

WORKDIR ${PATH_ROJECT}
COPY . ${PATH_ROJECT}

RUN CGO_ENABLED=0 GOOS=linux go build \
    -a -o ${BIN} ${APP}

FROM alpine:3.14 as production

RUN apk add --update --no-cache tzdata
ENV TZ Europe/Moscow

COPY --from=build /bin/faceitsubscriber /bin/faceitsubscriber
ENTRYPOINT ["/bin/faceitsubscriber"]