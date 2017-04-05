FROM alpine:3.5
LABEL description="Obtains Google Analytics RealTime API metrics, and presents them to prometheus for scraping."

RUN apk add --update ca-certificates
ENV APP_PATH /ga
RUN mkdir $APP_PATH
ADD ganalytics $APP_PATH/

WORKDIR $APP_PATH
ENTRYPOINT $APP_PATH"/ganalytics"
