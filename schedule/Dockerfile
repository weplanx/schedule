FROM alpine:edge

RUN apk add tzdata

COPY ./dist/main /app/schedule
WORKDIR /app

CMD [ "./main" ]
