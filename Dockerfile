FROM alpine:edge

RUN apk add tzdata

COPY dist /app
WORKDIR /app

EXPOSE 6000 8080

VOLUME [ "app/config" ]

CMD [ "./schedule-microservice" ]