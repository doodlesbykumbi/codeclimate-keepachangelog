ARG BASE=1.13.1-alpine3.10
ARG PARSE_A_CHANGELOG_VERSION=0.2.3
FROM golang:${BASE} as build

WORKDIR /usr/src/app

COPY engine.json ./engine.json.template
RUN apk add --no-cache jq
RUN export go_version=$(go version | cut -d ' ' -f 3) && \
    cat engine.json.template | jq '.version = .version + "/" + env.go_version' > ./engine.json

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY codeclimate-keepachangelog.go ./
RUN go build -o codeclimate-keepachangelog .

FROM golang:${BASE}
LABEL maintainer="Code Climate <hello@codeclimate.com>"

RUN apk update && apk upgrade && apk add bash curl-dev ruby-dev build-base ruby ruby-io-console ruby-bundler
RUN gem install --no-document parse_a_changelog -v "${PARSE_A_CHANGELOG_VERSION}"

RUN adduser -u 9000 -D app

WORKDIR /usr/src/app

COPY --from=build /usr/src/app/engine.json /
COPY --from=build /usr/src/app/codeclimate-keepachangelog ./

USER app

VOLUME /code

CMD ["/usr/src/app/codeclimate-keepachangelog"]
