FROM golang:1.11-alpine AS build-env
RUN apk add --update ca-certificates && apk --no-cache add build-base git gcc
ENV GO111MODULE on
ENV CGO_ENABLED 0
ENV GOOS linux
WORKDIR /src
ADD . /src
RUN go fmt && go build -ldflags "-s" -a -installsuffix cgo ganalytics.go

FROM alpine:3.5
LABEL description="Obtains Google Analytics RealTime API metrics, and presents them to prometheus for scraping."

RUN apk add --update ca-certificates
WORKDIR /ga
COPY --from=build-env /src/ganalytics /ga/ganalytics
ENTRYPOINT ["/ga/ganalytics"]
