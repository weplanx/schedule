FROM alpine:edge

RUN apk add tzdata

COPY ./dist/worker /app/main
WORKDIR /app

CMD [ "./main" ]
