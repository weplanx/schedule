FROM alpine:edge

RUN apk add tzdata

COPY dist /app
WORKDIR /app

EXPOSE 6000

VOLUME [ "app/config" ]

CMD [ "./schedule-microservice" ]